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
