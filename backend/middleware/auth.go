package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/root/evodrones/backend/config"
)

// TokenBlacklist mantém o registro de tokens inválidos
type TokenBlacklist struct {
	sync.RWMutex
	entries map[string]time.Time
}

var (
	blacklist     = &TokenBlacklist{entries: make(map[string]time.Time)}
	tokenLifespan = 24 * time.Hour
)

// Add adiciona um token à blacklist
func (b *TokenBlacklist) Add(token string) {
	b.Lock()
	defer b.Unlock()
	b.entries[token] = time.Now().Add(tokenLifespan)
}

// Contains verifica se um token está na blacklist
func (b *TokenBlacklist) Contains(token string) bool {
	b.RLock()
	defer b.RUnlock()
	expiry, exists := b.entries[token]
	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		delete(b.entries, token)
		return false
	}
	return true
}

// cleanExpired remove os tokens expirados da blacklist
func (b *TokenBlacklist) cleanExpired() {
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	for token, expiry := range b.entries {
		if now.After(expiry) {
			delete(b.entries, token)
		}
	}
}

// Iniciar a limpeza periódica da blacklist
func init() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			blacklist.cleanExpired()
		}
	}()
}

// generateSecureToken gera um token seguro para CSRF
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, role string, cfg *config.JWTConfig) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

func GenerateRefreshToken(userID string, cfg *config.JWTConfig) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

func ValidateToken(tokenString string, cfg *config.JWTConfig) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "token inválido")
	}

	return claims, nil
}

// AuthMiddleware verifica a autenticidade do token JWT
func AuthMiddleware(cfg *config.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Verifica o token CSRF para métodos não-GET
		if c.Method() != fiber.MethodGet {
			csrfToken := c.Get("X-CSRF-Token")
			if csrfToken == "" {
				log.Printf("Aviso: CSRF token não fornecido - IP: %s, Path: %s", c.IP(), c.Path())
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "token CSRF não fornecido",
				})
			}
		}

		// Verifica o token de autenticação
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf("Aviso: Tentativa de acesso sem token - IP: %s, Path: %s", c.IP(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token de autenticação não fornecido",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Aviso: Formato de token inválido - IP: %s, Header: %s", c.IP(), authHeader)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "formato de token inválido",
			})
		}

		token := parts[1]

		// Verifica se o token está na blacklist
		if blacklist.Contains(token) {
			log.Printf("Aviso: Tentativa de uso de token revogado - IP: %s, Path: %s", c.IP(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token revogado",
			})
		}

		// Valida o token
		claims, err := ValidateToken(token, cfg)
		if err != nil {
			log.Printf("Aviso: Falha na validação do token - IP: %s, Path: %s, Erro: %v", c.IP(), c.Path(), err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token inválido ou expirado",
			})
		}

		// Adiciona informações do usuário ao contexto
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		// Adiciona um token CSRF para a próxima requisição
		csrfToken, err := generateSecureToken(32)
		if err != nil {
			log.Printf("Erro: Falha ao gerar token CSRF - %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "erro interno do servidor",
			})
		}

		c.Locals("csrf_token", csrfToken)

		// Registra o acesso bem-sucedido
		log.Printf("Info: Acesso autorizado - UserID: %s, Role: %s, Path: %s", claims.UserID, claims.Role, c.Path())

		return c.Next()
	}
}

// LogoutMiddleware invalida o token atual
func LogoutMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				blacklist.Add(parts[1])
			}
		}
		return c.Next()
	}
}

func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "role não encontrada")
		}

		for _, role := range roles {
			if role == userRole.(string) {
				return c.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "acesso negado")
	}
}
