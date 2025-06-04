document.addEventListener('DOMContentLoaded', function() {
    // Inicializa AOS se estiver disponível
    if (typeof AOS !== 'undefined') {
        AOS.init({
            duration: 800,
            once: true
        });
    }
    
    // Toggle password visibility
    document.querySelectorAll('.toggle-password').forEach(button => {
        button.addEventListener('click', function() {
            const input = this.previousElementSibling;
            const icon = this.querySelector('i');
            
            if (input.type === 'password') {
                input.type = 'text';
                icon.classList.replace('bi-eye', 'bi-eye-slash');
            } else {
                input.type = 'password';
                icon.classList.replace('bi-eye-slash', 'bi-eye');
            }
        });
    });
    
    // Form validation
    const forms = document.querySelectorAll('.needs-validation');
    Array.from(forms).forEach(form => {
        form.addEventListener('submit', event => {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            }
            form.classList.add('was-validated');
        }, false);
    });
    
    // Navegação entre etapas do cadastro
    document.querySelectorAll('.next-step').forEach(button => {
        button.addEventListener('click', function() {
            // Encontra a aba atual e a próxima
            const currentTab = document.querySelector('.tab-pane.active');
            const currentTabId = currentTab.getAttribute('id');
            const nextTabId = currentTabId === 'dados-pessoais' ? 'dados-drone' : 'confirmacao';
            
            // Atualiza o resumo se estiver indo para a etapa de confirmação
            if (nextTabId === 'confirmacao') {
                updateSummary();
            }
            
            // Ativa a próxima aba
            const nextTab = new bootstrap.Tab(document.querySelector(`#${nextTabId}-tab`));
            nextTab.show();
        });
    });
    
    document.querySelectorAll('.prev-step').forEach(button => {
        button.addEventListener('click', function() {
            // Encontra a aba atual e a anterior
            const currentTab = document.querySelector('.tab-pane.active');
            const currentTabId = currentTab.getAttribute('id');
            const prevTabId = currentTabId === 'confirmacao' ? 'dados-drone' : 'dados-pessoais';
            
            // Ativa a aba anterior
            const prevTab = new bootstrap.Tab(document.querySelector(`#${prevTabId}-tab`));
            prevTab.show();
        });
    });
    
    // Atualiza o resumo das informações no cadastro
    function updateSummary() {
        if (document.getElementById('resumoNome')) {
            document.getElementById('resumoNome').textContent = document.getElementById('nome').value;
            document.getElementById('resumoEmail').textContent = document.getElementById('email').value;
            document.getElementById('resumoCPF').textContent = document.getElementById('cpf').value;
            document.getElementById('resumoTelefone').textContent = document.getElementById('telefone').value;
            
            document.getElementById('resumoModelo').textContent = document.getElementById('modeloDrone').value;
            document.getElementById('resumoSerie').textContent = document.getElementById('numeroSerie').value;
            document.getElementById('resumoData').textContent = document.getElementById('dataCompra').value;
            document.getElementById('resumoLocal').textContent = document.getElementById('localCompra').value;
        }
    }
    
    // Alternar entre abas de login (Cliente/Técnico)
    const loginTabs = document.querySelectorAll('#loginTabs .nav-link');
    loginTabs.forEach(tab => {
        tab.addEventListener('click', function(e) {
            e.preventDefault();
            loginTabs.forEach(t => t.classList.remove('active'));
            this.classList.add('active');
            
            const target = this.getAttribute('data-bs-target');
            document.querySelectorAll('.login-tab-content').forEach(content => {
                content.classList.remove('active');
            });
            document.querySelector(target).classList.add('active');
        });
    });
});