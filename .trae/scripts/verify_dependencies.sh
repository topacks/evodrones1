#!/bin/bash

# EVO DRONES - Script de Verificação de Dependências
# Atualizado com base na análise do projeto em 2024

echo "🔍 Verificando dependências do projeto EVO DRONES..."
echo "================================================"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para verificar se estamos no diretório correto
check_directory() {
    if [ ! -d "backend" ] || [ ! -d "frontend" ] || [ ! -f "docker-compose.yml" ]; then
        echo -e "${RED}❌ Erro: Execute este script no diretório raiz do projeto EVO DRONES${NC}"
        exit 1
    fi
}

# Função para verificar dependências de uma struct
check_struct_dependencies() {
    local struct_name=$1
    echo -e "${BLUE}🔍 Verificando dependências da struct: ${struct_name}${NC}"
    
    echo "📁 Arquivos que usam a struct:"
    grep -r "$struct_name" backend/ --include="*.go" | head -10
    
    echo "📝 Campos da struct em uso:"
    grep -r "\.$struct_name" backend/ --include="*.go" | head -5
    
    echo ""
}

# Função para verificar dependências de uma função
check_function_dependencies() {
    local function_name=$1
    echo -e "${BLUE}🔍 Verificando dependências da função: ${function_name}${NC}"
    
    echo "📁 Arquivos que chamam a função:"
    grep -r "$function_name" backend/ --include="*.go" | head -10
    
    echo "🧪 Testes relacionados:"
    grep -r "Test.*$function_name\|$function_name.*Test" backend/ --include="*_test.go" | head -5
    
    echo ""
}

# Função para verificar rotas
check_route_dependencies() {
    local route_path=$1
    echo -e "${BLUE}🔍 Verificando dependências da rota: ${route_path}${NC}"
    
    echo "🛣️ Definições de rota:"
    grep -r "$route_path" backend/routes/ --include="*.go"
    
    echo "🔒 Middleware aplicado:"
    grep -r "middleware" backend/routes/ --include="*.go" | grep -i "$route_path" || echo "Nenhum middleware específico encontrado"
    
    echo ""
}

# Função para executar testes
run_tests() {
    echo -e "${YELLOW}🧪 Executando testes...${NC}"
    
    cd backend
    
    echo "📋 Executando go vet..."
    if go vet ./...; then
        echo -e "${GREEN}✅ go vet passou${NC}"
    else
        echo -e "${RED}❌ go vet falhou${NC}"
        return 1
    fi
    
    echo "📋 Executando go fmt..."
    if [ -z "$(go fmt ./...)" ]; then
        echo -e "${GREEN}✅ go fmt passou${NC}"
    else
        echo -e "${RED}❌ go fmt encontrou problemas de formatação${NC}"
        return 1
    fi
    
    echo "📋 Executando testes unitários..."
    if go test ./...; then
        echo -e "${GREEN}✅ Todos os testes passaram${NC}"
    else
        echo -e "${RED}❌ Alguns testes falharam${NC}"
        return 1
    fi
    
    echo "📋 Verificando build..."
    if go build ./cmd/api; then
        echo -e "${GREEN}✅ Build bem-sucedido${NC}"
        rm -f api # Remove o binário gerado
    else
        echo -e "${RED}❌ Build falhou${NC}"
        return 1
    fi
    
    cd ..
}

# Função para verificar dependências de módulos Go
check_go_modules() {
    echo -e "${YELLOW}📦 Verificando módulos Go...${NC}"
    
    cd backend
    
    echo "📋 Executando go mod tidy..."
    go mod tidy
    
    echo "📋 Verificando módulos..."
    if go mod verify; then
        echo -e "${GREEN}✅ Módulos verificados com sucesso${NC}"
    else
        echo -e "${RED}❌ Problemas encontrados nos módulos${NC}"
        return 1
    fi
    
    echo "📋 Verificando dependências circulares..."
    if go list -deps ./... > /dev/null; then
        echo -e "${GREEN}✅ Nenhuma dependência circular encontrada${NC}"
    else
        echo -e "${RED}❌ Dependências circulares detectadas${NC}"
        return 1
    fi
    
    cd ..
}

# Função para verificar estruturas críticas
check_critical_structures() {
    echo -e "${YELLOW}🏗️ Verificando estruturas críticas...${NC}"
    
    local critical_structs=("User" "Drone" "MaintenanceOrder" "TechnicalReport" "DroneItem")
    
    for struct in "${critical_structs[@]}"; do
        check_struct_dependencies "$struct"
    done
}

# Função para verificar endpoints críticos
check_critical_endpoints() {
    echo -e "${YELLOW}🌐 Verificando endpoints críticos...${NC}"
    
    local critical_routes=("/api/auth" "/api/drones" "/api/maintenance-orders" "/api/technical-reports")
    
    for route in "${critical_routes[@]}"; do
        check_route_dependencies "$route"
    done
}

# Função para gerar relatório de dependências
generate_dependency_report() {
    echo -e "${YELLOW}📊 Gerando relatório de dependências...${NC}"
    
    local report_file=".trae/reports/dependency_report_$(date +%Y%m%d_%H%M%S).md"
    mkdir -p .trae/reports
    
    {
        echo "# Relatório de Dependências - $(date)"
        echo ""
        echo "## Estruturas Críticas"
        echo ""
        
        local critical_structs=("User" "Drone" "MaintenanceOrder" "TechnicalReport" "DroneItem")
        for struct in "${critical_structs[@]}"; do
            echo "### $struct"
            echo "\`\`\`"
            grep -r "$struct" backend/ --include="*.go" | wc -l | xargs echo "Referências encontradas:"
            echo "\`\`\`"
            echo ""
        done
        
        echo "## Endpoints Críticos"
        echo ""
        
        local critical_routes=("/api/auth" "/api/drones" "/api/maintenance-orders" "/api/technical-reports")
        for route in "${critical_routes[@]}"; do
            echo "### $route"
            echo "\`\`\`"
            grep -r "$route" backend/routes/ --include="*.go" | wc -l | xargs echo "Definições encontradas:"
            echo "\`\`\`"
            echo ""
        done
        
    } > "$report_file"
    
    echo -e "${GREEN}📄 Relatório gerado: $report_file${NC}"
}

# Função principal
main() {
    check_directory
    
    case "${1:-all}" in
        "struct")
            if [ -z "$2" ]; then
                echo -e "${RED}❌ Uso: $0 struct <nome_da_struct>${NC}"
                exit 1
            fi
            check_struct_dependencies "$2"
            ;;
        "function")
            if [ -z "$2" ]; then
                echo -e "${RED}❌ Uso: $0 function <nome_da_funcao>${NC}"
                exit 1
            fi
            check_function_dependencies "$2"
            ;;
        "route")
            if [ -z "$2" ]; then
                echo -e "${RED}❌ Uso: $0 route <caminho_da_rota>${NC}"
                exit 1
            fi
            check_route_dependencies "$2"
            ;;
        "test")
            run_tests
            ;;
        "modules")
            check_go_modules
            ;;
        "critical")
            check_critical_structures
            check_critical_endpoints
            ;;
        "report")
            generate_dependency_report
            ;;
        "all")
            echo -e "${GREEN}🚀 Executando verificação completa...${NC}"
            check_critical_structures
            check_critical_endpoints
            check_go_modules
            run_tests
            generate_dependency_report
            echo -e "${GREEN}✅ Verificação completa finalizada!${NC}"
            ;;
        "help")
            echo "Uso: $0 [comando] [parâmetro]"
            echo ""
            echo "Comandos disponíveis:"
            echo "  struct <nome>     - Verificar dependências de uma struct"
            echo "  function <nome>   - Verificar dependências de uma função"
            echo "  route <caminho>   - Verificar dependências de uma rota"
            echo "  test             - Executar todos os testes"
            echo "  modules          - Verificar módulos Go"
            echo "  critical         - Verificar estruturas e endpoints críticos"
            echo "  report           - Gerar relatório de dependências"
            echo "  all              - Executar verificação completa (padrão)"
            echo "  help             - Mostrar esta ajuda"
            ;;
        *)
            echo -e "${RED}❌ Comando inválido: $1${NC}"
            echo "Use '$0 help' para ver os comandos disponíveis"
            exit 1
            ;;
    esac
}

# Executar função principal com todos os argumentos
main "$@"