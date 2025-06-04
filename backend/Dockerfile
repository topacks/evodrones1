FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências do sistema
RUN apk add --no-cache git

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o evodrones-api ./cmd/api

# Imagem final
FROM alpine:3.18

WORKDIR /app

# Instalar dependências necessárias
RUN apk add --no-cache ca-certificates tzdata

# Criar diretório para uploads
RUN mkdir -p /app/uploads

# Copiar o binário compilado
COPY --from=builder /app/evodrones-api .

# Expor a porta
EXPOSE 8080

# Executar a aplicação
CMD ["./evodrones-api"]