#!/bin/bash

# Fechar processos nas portas 8080, 8081, 8082
echo "Fechando processos nas portas 8080, 8081, 8082..."
lsof -ti:8080 | xargs -r kill -9
lsof -ti:8081 | xargs -r kill -9
lsof -ti:8082 | xargs -r kill -9

# Aguardar um momento para garantir que as portas sejam liberadas
sleep 2

# Iniciar o backend
cd backend
echo "Iniciando servidor backend..."
go run cmd/api/main.go