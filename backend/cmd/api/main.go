package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/root/evodrones/backend/config"
	"github.com/root/evodrones/backend/controllers"
	"github.com/root/evodrones/backend/db"
	"github.com/root/evodrones/backend/internal/permissions"
	customMiddleware "github.com/root/evodrones/backend/middleware"
	"github.com/root/evodrones/backend/repositories"
	"github.com/root/evodrones/backend/routes"
	"github.com/root/evodrones/backend/services"
)

func main() {
	// Carregar configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Inicializar banco de dados
	if err := db.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	// Criar instância do Fiber
	app := fiber.New(fiber.Config{
		AppName:      "EvoDrones API",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New())
	
	// Configuração CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.AllowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour), // 12 horas
	}))

	// Middleware de segurança
	app.Use(customMiddleware.SecurityHeaders())
	app.Use(customMiddleware.RateLimiterMiddleware())
	app.Use(customMiddleware.RequestIDMiddleware())

	// Rate limiting para APIs públicas
	apiLimiter := limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests, please try again later.",
			})
		},
	})

	// Servir arquivos estáticos do frontend
	app.Static("/assets", "./static/assets")
	app.Static("/", "./static")

	// Rotas específicas para páginas principais
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/public/home/index.html") // Certificando que aponta para o novo arquivo
	})

	// Rota para a home antiga, caso necessário manter por algum motivo ou para redirecionamento futuro
	app.Get("/home_antiga", func(c *fiber.Ctx) error {
		// Se a home antiga foi movida ou renomeada, ajustar o caminho aqui
		// Exemplo: return c.SendFile("./static/pages/public/home_antiga/index.html")
		// Por enquanto, pode retornar um 404 ou redirecionar para a nova home
		return c.Status(fiber.StatusNotFound).SendString("Página não encontrada ou movida.")
	})

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/public/home/index.html")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/auth/login/index.html")
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/auth/register/index.html")
	})

	// Rotas dos dashboards
	app.Get("/dashboard/client", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/dashboard/client/index.html")
	})

	app.Get("/dashboard/technician", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/dashboard/technician/index.html")
	})

	// Rota para visualização 3D
	app.Get("/drone/3d-view", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pages/client/drone-3d-view.html")
	})

	// Rota para demo layouts
	app.Get("/demo-layouts/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		return c.SendFile("./static/demo-layouts/" + path)
	})

	// Fallback para outras páginas
	app.Get("/pages/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		return c.SendFile("./static/pages/" + path)
	})

	// Rota de health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API está funcionando",
		})
	})

	// Repositórios
	dbInstance := db.GetDB()
	droneRepo := repositories.NewDroneRepository(dbInstance)
	userRepo := repositories.NewUserRepository(dbInstance)
	maintenanceRepo := repositories.NewMaintenanceOrderRepository(dbInstance.DB)
	reportRepo := repositories.NewTechnicalReportRepository(dbInstance.DB)

	// Serviços
	authService := services.NewAuthService(userRepo, cfg)
	droneService := services.NewDroneService(droneRepo, userRepo)
	maintenanceService := services.NewMaintenanceOrderService(maintenanceRepo, droneRepo, userRepo)
	reportService := services.NewTechnicalReportService(reportRepo, maintenanceRepo, userRepo)

	// Inicializar sistema RBAC
	permissionManager, err := permissions.InitializePermissionSystemComplete(dbInstance.DB)
	if err != nil {
		log.Fatalf("Erro ao inicializar sistema RBAC: %v", err)
	}

	// Serviço de autorização centralizado (mantido para compatibilidade)
	_ = services.NewAuthorizationService(userRepo, droneRepo, maintenanceRepo, reportRepo)

	// Controladores
	authController := controllers.NewAuthController(authService)
	droneController := controllers.NewDroneController(droneService)
	maintenanceController := controllers.NewMaintenanceOrderController(maintenanceService)
	reportController := controllers.NewTechnicalReportController(reportService)

	// Grupo de rotas da API com rate limiting
	api := app.Group("/api", apiLimiter)

	// Configurar rotas públicas de autenticação
	routes.SetupAuthRoutes(api, authController)

	// Configurar rotas protegidas com autenticação JWT
	api.Use(customMiddleware.AuthMiddleware(&cfg.JWT))

	// Configurar o serviço de autorização
	authorizationService := services.NewAuthorizationService(userRepo, droneRepo, maintenanceRepo, reportRepo)

	// Configurar rotas protegidas
	routes.SetupRoutes(api, cfg, droneController, maintenanceController, authorizationService)
	routes.SetupTechnicalReportRoutes(api, cfg, reportController, authorizationService)

	// Configurar rotas de permissões (se necessário)
	// permissions.SetupPermissionRoutes(api, cfg, droneController, maintenanceController, reportController, permissionManager.Service)

	// Middleware para rotas não encontradas
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error": "Endpoint não encontrado",
		})
	})

	// Log do status do sistema RBAC
	status := permissions.GetSystemStatus(permissionManager)
	log.Printf("Sistema RBAC inicializado: %+v", status)

	// Iniciar servidor
	log.Printf("Servidor iniciando na porta %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
