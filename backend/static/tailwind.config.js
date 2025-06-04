/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./login.html",
    "./pages/**/*.html",
    "./components/**/*.html",
    "./assets/js/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        'evo-primary': '#30B4E8', // Azul principal EVO
        'evo-secondary': '#FFC107', // Amarelo secundário
        'evo-dark': '#1E293B',    // Azul escuro/cinza para texto e fundos
        'evo-light': '#F8F9FA',   // Cinza claro para fundos
        'evo-success': '#198754',  // Verde para sucesso
        'evo-danger': '#DC3545',   // Vermelho para perigo
        'evo-warning': '#FFC107',  // Amarelo para aviso
        'evo-info': '#0DCAF0',     // Azul claro para informação
        'evo-gray': '#6C757D',     // Cinza para texto secundário
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
      backgroundImage: {
        'gradient-evo-tech': 'linear-gradient(135deg, #1E293B 0%, #30B4E8 100%)',
        'gradient-evo-primary': 'linear-gradient(to right, #30B4E8, #2596be)',
      },
      boxShadow: {
        'card': '0 10px 30px rgba(0, 0, 0, 0.08)',
        'input': '0 2px 10px rgba(48, 180, 232, 0.1)',
      },
      borderRadius: {
        'xl': '1rem', // 16px
        '2xl': '1.5rem', // 24px
      }
    },
  },
  plugins: [],
}