#!/bin/bash

# Script para migrar o projeto para uma estrutura simplificada (versão não interativa)
# Elimina a duplicidade entre frontend e backend/static

echo "Iniciando migração para estrutura simplificada..."

# Diretórios principais
FRONTEND_DIR="/root/evodrones/frontend"
BACKEND_STATIC_DIR="/root/evodrones/backend/static"
BACKUP_DIR="/root/evodrones/scripts/backup"

# 1. Garantir que todos os arquivos do frontend estejam sincronizados com backend/static
echo "Sincronizando arquivos do frontend para backend/static..."
/root/evodrones/scripts/sync-frontend.sh

# 2. Fazer backup dos scripts que serão modificados ou removidos
echo "Fazendo backup dos scripts..."
mkdir -p "$BACKUP_DIR"
cp -f /root/evodrones/scripts/sync-frontend.sh "$BACKUP_DIR/sync-frontend.sh.bak" 2>/dev/null || true
cp -f /root/evodrones/scripts/watch-frontend.sh "$BACKUP_DIR/watch-frontend.sh.bak" 2>/dev/null || true
cp -f /root/evodrones/scripts/live-preview.sh "$BACKUP_DIR/live-preview.sh.bak" 2>/dev/null || true

# 3. Atualizar o script live-preview.sh para apontar diretamente para backend/static
echo "Atualizando script live-preview.sh..."
cat > /root/evodrones/scripts/live-preview.sh << 'EOF'
#!/bin/bash

# Script para iniciar um servidor de desenvolvimento com atualização em tempo real

echo "Iniciando servidor de desenvolvimento com visualização em tempo real..."

# Verificar se o browser-sync está instalado
if ! command -v browser-sync &> /dev/null; then
    echo "Instalando browser-sync..."
    npm install -g browser-sync
fi

# Diretório dos arquivos estáticos
STATIC_DIR="/root/evodrones/backend/static"

# Iniciar o browser-sync
echo "Iniciando browser-sync para o diretório $STATIC_DIR"
browser-sync start --server "$STATIC_DIR" --files "$STATIC_DIR/**/*" --port 3000 --ui-port 3001 --no-notify
EOF
chmod +x /root/evodrones/scripts/live-preview.sh

# 4. Criar um README explicando a nova estrutura
echo "Criando documentação da nova estrutura..."
cat > /root/evodrones/README-ESTRUTURA.md << 'EOF'
# EVO DRONES - Estrutura Simplificada

## Visão Geral

O projeto EVO DRONES foi reestruturado para seguir uma estrutura mais simples e alinhada com os padrões de projetos Go. A principal mudança foi a eliminação da duplicidade entre os diretórios `frontend` e `backend/static`.

## Nova Estrutura

- **Backend**: Contém toda a aplicação Go e os arquivos estáticos
  - `backend/static/`: Diretório único para todos os arquivos estáticos (HTML, CSS, JS, imagens)
  - `backend/cmd/api/`: Ponto de entrada da aplicação
  - `backend/config/`: Configurações
  - `backend/controllers/`: Handlers HTTP
  - `backend/middleware/`: Middlewares
  - `backend/models/`: Estruturas de dados
  - `backend/repositories/`: Camada de dados
  - `backend/services/`: Lógica de negócio
  - `backend/routes/`: Definição de rotas
  - `backend/db/`: Conexão e migrações

- **Database**: Contém os scripts de inicialização do banco de dados
  - `database/init.sql`: Schema inicial e dados

## Fluxo de Trabalho

1. **Desenvolvimento Frontend**:
   - Edite os arquivos diretamente em `backend/static/`
   - Use o script `scripts/live-preview.sh` para visualização em tempo real

2. **Execução do Projeto**:
   - Use o script `start.sh` para iniciar o servidor Go
   - Acesse o frontend em `http://localhost:8080`

## Scripts Disponíveis

- `start.sh`: Inicia o servidor backend
- `scripts/live-preview.sh`: Inicia um servidor de desenvolvimento com atualização em tempo real

## Observações

- O servidor Go é responsável por servir todos os arquivos estáticos
- Não é necessário sincronizar arquivos, pois agora existe apenas um diretório para os arquivos estáticos
- A estrutura simplificada facilita a manutenção e o desenvolvimento
EOF

# 5. Remover o diretório frontend e scripts desnecessários
echo "Removendo diretório frontend e scripts desnecessários..."
rm -rf "$FRONTEND_DIR"
rm -f /root/evodrones/scripts/sync-frontend.sh
rm -f /root/evodrones/scripts/watch-frontend.sh
rm -f /root/evodrones/dev-preview.sh
rm -f /root/evodrones/README-DEV-PREVIEW.md

echo "Migração concluída com sucesso!"
echo "A nova estrutura está documentada em README-ESTRUTURA.md"
echo "Use o script start.sh para iniciar o servidor backend."
echo "Use o script scripts/live-preview.sh para visualização em tempo real durante o desenvolvimento."