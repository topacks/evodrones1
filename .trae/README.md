# EVO DRONES - Sistema de Memórias Estruturais

## 📋 Visão Geral

Este diretório contém o sistema de memórias estruturais do projeto EVO DRONES, projetado para manter a consistência, qualidade e integridade do código durante o desenvolvimento.

## 📁 Estrutura de Diretórios

```
.trae/
├── README.md                    # Este arquivo
├── memories/                    # Memórias estruturais do projeto
│   ├── project_structure.md     # Arquitetura e organização
│   ├── coding_patterns.md       # Padrões e convenções
│   └── dependencies_map.md      # Mapeamento de dependências
├── scripts/                     # Scripts de automação
│   └── verify_dependencies.sh   # Verificação de dependências
├── rules/                       # Regras do projeto
│   └── project_rules.md         # Regras específicas
└── reports/                     # Relatórios gerados (criado automaticamente)
    └── dependency_report_*.md   # Relatórios de dependências
```

## 🧠 Memórias Estruturais

### 📖 project_structure.md
Contém a documentação completa da arquitetura do projeto:
- Estrutura de diretórios
- Tecnologias utilizadas
- Padrões de arquitetura
- Configurações importantes
- Modelos de dados

### 🎯 coding_patterns.md
Define os padrões de codificação obrigatórios:
- Convenções de nomenclatura
- Estruturas de controllers, services e repositories
- Padrões de tratamento de erros
- Validação de entrada
- Padrões de teste

### 🔗 dependencies_map.md
Mapeia todas as dependências críticas:
- Relações entre controllers, services e repositories
- Modelos compartilhados
- Dependências de middleware
- Regras de modificação
- Fluxos de dados críticos

## 🛠️ Scripts de Automação

### verify_dependencies.sh
Script principal para verificação de dependências e validação de código.

#### Uso Básico
```bash
# Verificação completa (recomendado antes de qualquer modificação)
./.trae/scripts/verify_dependencies.sh all

# Verificar dependências de uma struct específica
./.trae/scripts/verify_dependencies.sh struct User

# Verificar dependências de uma função
./.trae/scripts/verify_dependencies.sh function getUserByID

# Verificar dependências de uma rota
./.trae/scripts/verify_dependencies.sh route "/api/auth"

# Executar apenas testes
./.trae/scripts/verify_dependencies.sh test

# Verificar módulos Go
./.trae/scripts/verify_dependencies.sh modules

# Verificar estruturas críticas
./.trae/scripts/verify_dependencies.sh critical

# Gerar relatório de dependências
./.trae/scripts/verify_dependencies.sh report

# Mostrar ajuda
./.trae/scripts/verify_dependencies.sh help
```

## 📋 Workflow de Desenvolvimento

### Antes de Modificar Código
1. **Consultar memórias estruturais** para entender o contexto
2. **Executar verificação de dependências**:
   ```bash
   ./.trae/scripts/verify_dependencies.sh all
   ```
3. **Analisar impacto** das mudanças propostas
4. **Planejar implementação** seguindo padrões estabelecidos

### Durante a Implementação
1. **Seguir padrões** documentados em `coding_patterns.md`
2. **Manter consistência** com arquitetura em `project_structure.md`
3. **Verificar dependências** conforme `dependencies_map.md`
4. **Executar testes** frequentemente

### Após Modificações
1. **Executar verificação completa**:
   ```bash
   ./.trae/scripts/verify_dependencies.sh all
   ```
2. **Validar testes** e build
3. **Atualizar memórias** se necessário
4. **Gerar relatório** de dependências

## 🚨 Regras Críticas

### ❌ NUNCA FAÇA
- Alterar nomes de structs existentes sem verificação completa
- Modificar assinaturas de funções públicas
- Remover campos de modelos sem análise de impacto
- Alterar esquema de banco sem migração
- Ignorar falhas nos testes

### ✅ SEMPRE FAÇA
- Verificar dependências antes de modificar
- Executar testes após mudanças
- Seguir padrões de nomenclatura
- Usar prepared statements
- Implementar validação de entrada
- Documentar mudanças significativas

## 📊 Relatórios

O sistema gera automaticamente relatórios de dependências em `.trae/reports/`:
- **dependency_report_YYYYMMDD_HHMMSS.md** - Relatório completo de dependências
- Contém análise de estruturas críticas
- Mapeia endpoints importantes
- Identifica pontos de atenção

## 🔧 Manutenção do Sistema

### Atualização das Memórias
As memórias estruturais devem ser atualizadas quando:
- Houver mudanças significativas na arquitetura
- Novos módulos ou serviços forem adicionados
- Padrões de codificação forem alterados
- Durante refatorações importantes

### Validação Periódica
Execute verificação completa regularmente:
```bash
# Verificação semanal recomendada
./.trae/scripts/verify_dependencies.sh all
```

### Backup das Memórias
As memórias estruturais são críticas para o projeto:
- Mantenha-as versionadas no Git
- Faça backup antes de grandes mudanças
- Documente alterações no changelog

## 🆘 Troubleshooting

### Script não executa
```bash
# Verificar permissões
ls -la .trae/scripts/verify_dependencies.sh

# Tornar executável se necessário
chmod +x .trae/scripts/verify_dependencies.sh
```

### Testes falhando
```bash
# Verificar dependências Go
cd backend && go mod tidy

# Executar testes com verbose
cd backend && go test -v ./...
```

### Build falhando
```bash
# Verificar sintaxe
cd backend && go vet ./...

# Verificar formatação
cd backend && go fmt ./...
```

## 📞 Suporte

Para dúvidas sobre o sistema de memórias estruturais:
1. Consulte este README
2. Verifique as memórias em `/memories/`
3. Execute o script de verificação
4. Analise os relatórios gerados

---

**Lembre-se**: Este sistema foi projetado para manter a qualidade e consistência do código. Use-o sempre antes de fazer modificações significativas no projeto!