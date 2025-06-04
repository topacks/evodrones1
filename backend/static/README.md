# 🚁 EVO DRONES - Frontend

Frontend moderno e responsivo da plataforma EVO DRONES, especializada em manutenção e reparo de drones.

## 🎨 Design System

### Identidade Visual
- **Cores Principais:**
  - `#30B4E8` - Azul ciano (primária)
  - `#0E8BC0` - Azul escuro
  - `#FF7A00` - Laranja (destaque)
  - `#0F172A` - Azul escuro quase preto
  - `#F8FAFC` - Branco suave

### Tecnologias
- **Bootstrap 5.3.3** - Framework CSS
- **Bootstrap Icons** - Ícones
- **AOS (Animate On Scroll)** - Animações
- **Google Fonts (Inter)** - Tipografia
- **CSS Custom Properties** - Variáveis CSS
- **Vanilla JavaScript** - Funcionalidades

## 📁 Estrutura do Projeto

```
frontend/
├── assets/
│   ├── css/
│   │   └── main.css              # CSS principal com design system
│   ├── js/
│   │   └── main.js               # JavaScript principal
│   └── images/                   # Imagens e assets
├── pages/
│   ├── auth/
│   │   ├── login/
│   │   │   └── index.html        # Página de login (Cliente/Técnico)
│   │   ├── register/
│   │   │   └── index.html        # Cadastro multi-step
│   │   └── recover/              # Recuperação de senha
│   ├── dashboard/
│   │   ├── client/               # Dashboard do cliente
│   │   └── technician/           # Dashboard do técnico
│   ├── drones/                   # Gestão de drones
│   ├── maintenance/              # Ordens de manutenção
│   └── public/
│       └── home/
│           └── index.html        # Página inicial pública
├── components/                   # Componentes reutilizáveis
├── package.json
├── index.html                    # Redirecionamento principal
└── README.md
```

## 🚀 Funcionalidades Implementadas

### ✅ Páginas Criadas
- [x] **Página Inicial** - Landing page com hero section, serviços, sobre, diferenciais
- [x] **Login** - Autenticação para clientes e técnicos com abas
- [x] **Cadastro** - Formulário multi-step com validação
- [x] **Redirecionamento** - Página principal com loader animado

### 🎯 Características
- **Design Responsivo** - Mobile-first approach
- **Animações Suaves** - AOS e CSS animations
- **Validação de Formulários** - JavaScript nativo
- **Loading States** - Feedback visual para ações
- **Toast Notifications** - Sistema de notificações
- **Máscaras de Input** - CPF, telefone, CEP
- **Busca de CEP** - Integração com ViaCEP
- **Força de Senha** - Indicador visual
- **Efeitos Visuais** - Glow, hover, partículas

## 🎨 Componentes Visuais

### Cards Tecnológicos
```css
.card-tech {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(48, 180, 232, 0.1);
  border-radius: 16px;
}
```

### Botões Customizados
- `btn-evo-primary` - Botão principal com gradiente
- `btn-evo-accent` - Botão de destaque laranja
- `btn-outline-evo-primary` - Botão outline

### Animações
- `animate-float` - Flutuação suave
- `animate-pulse-glow` - Pulsação com brilho
- `animate-spin-slow` - Rotação lenta
- `animate-bounce` - Efeito bounce

## 🛠️ Como Usar

### Instalação
```bash
cd frontend
npm install
```

### Desenvolvimento
```bash
npm run dev
# Abre servidor local na porta 3000
```

### Build
```bash
npm run build
# Prepara arquivos para produção
```

### Preview
```bash
npm run preview
# Visualiza build na porta 3001
```

## 📱 Responsividade

### Breakpoints
- **Mobile:** < 576px
- **Tablet:** 576px - 992px
- **Desktop:** > 992px

### Adaptações Mobile
- Navbar colapsível
- Cards empilhados
- Botões full-width
- Texto redimensionado
- Espaçamentos otimizados

## 🎯 Próximas Implementações

### 📋 Páginas Pendentes
- [ ] Dashboard do Cliente
- [ ] Dashboard do Técnico
- [ ] Gestão de Drones (CRUD)
- [ ] Ordens de Manutenção
- [ ] Relatórios Técnicos
- [ ] Perfil do Usuário
- [ ] Configurações
- [ ] Recuperação de Senha

### 🔧 Funcionalidades
- [ ] Integração com API Backend
- [ ] Sistema de autenticação JWT
- [ ] Upload de arquivos/imagens
- [ ] Chat/Mensagens
- [ ] Notificações push
- [ ] PWA (Progressive Web App)
- [ ] Dark mode
- [ ] Internacionalização (i18n)

## 🎨 Guia de Estilo

### Cores
```css
:root {
  --evo-primary: #30B4E8;
  --evo-primary-dark: #0E8BC0;
  --evo-accent: #FF7A00;
  --evo-dark: #0F172A;
  --evo-light: #F8FAFC;
}
```

### Tipografia
- **Família:** Inter (Google Fonts)
- **Pesos:** 300, 400, 500, 600, 700, 800
- **Tamanhos:** Sistema de escala Bootstrap

### Espaçamentos
- **Padding:** 1rem, 1.5rem, 2rem, 3rem, 5rem
- **Margin:** Sistema Bootstrap (0-5)
- **Border Radius:** 8px, 12px, 16px, 20px, 24px

## 🔧 Utilitários JavaScript

### Classe Principal
```javascript
class EvoDrones {
  static showToast(message, type)
  static formatCurrency(value)
  static formatDate(date)
  static debounce(func, wait)
}
```

### Funcionalidades
- Navbar responsiva
- Animações on scroll
- Sistema de partículas
- Validação de formulários
- Toggle de senha
- Tooltips e modais

## 📞 Suporte

Para dúvidas ou sugestões sobre o frontend:
- **Email:** dev@evodrones.com.br
- **Documentação:** [Link para docs]
- **Issues:** [Link para GitHub issues]

---

**EVO DRONES** - Especialistas em Manutenção de Drones 🚁
