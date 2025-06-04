/**
 * EVO DRONES - Painel Administrativo
 * Script principal para o painel administrativo
 */

class AdminDashboard {
    constructor() {
        // Inicializar propriedades
        this.charts = {};
        this.dataCache = {};
        this.refreshInterval = null;
        this.sidebarVisible = true;
        this.darkMode = false;
        
        // Inicializar o dashboard
        this.init();
    }
    
    /**
     * Inicializa o dashboard administrativo
     */
    init() {
        // Cachear elementos DOM frequentemente utilizados
        this.cacheElements();
        
        // Configurar event listeners
        this.setupEventListeners();
        
        // Carregar dados iniciais
        this.loadData();
        
        // Inicializar componentes
        this.initCharts();
        this.initDataTables();
        this.initTooltips();
        this.initDatePicker();
        
        // Atualizar data/hora
        this.updateDateTime();
        setInterval(() => this.updateDateTime(), 60000);
        
        // Verificar tamanho da tela para ajustar o sidebar
        this.handleResponsiveLayout();
        
        // Configurar atualização automática de dados
        this.setupAutoRefresh();
        
        console.log('AdminDashboard inicializado com sucesso!');
    }
    
    /**
     * Cacheia elementos DOM frequentemente utilizados
     */
    cacheElements() {
        this.elements = {
            sidebar: document.getElementById('sidebar'),
            main: document.getElementById('main'),
            sidebarToggle: document.getElementById('sidebarToggleMobile'),
            sidebarToggleDesktop: document.getElementById('sidebarToggle'),
            currentDate: document.getElementById('currentDate'),
            userGrowthChart: document.getElementById('userGrowthChart'),
            userDistributionChart: document.getElementById('userDistributionChart'),
            usersTable: document.getElementById('usersTable'),
            notificationBadge: document.querySelector('.badge.bg-danger'),
            filterDropdowns: document.querySelectorAll('.dropdown-menu[aria-labelledby="chartDropdown"] .dropdown-item'),
            themeToggle: document.getElementById('themeToggle'),
            refreshButton: document.querySelector('.btn-refresh'),
            exportButton: document.querySelector('.btn-outline-primary'),
            taskCheckboxes: document.querySelectorAll('.form-check-input[type="checkbox"]')
        };
    }
    
    /**
     * Configura os event listeners
     */
    setupEventListeners() {
        // Toggle do sidebar em dispositivos móveis
        if (this.elements.sidebarToggle) {
            this.elements.sidebarToggle.addEventListener('click', () => this.toggleSidebar());
        }
        
        // Toggle do sidebar em desktop
        if (this.elements.sidebarToggleDesktop) {
            this.elements.sidebarToggleDesktop.addEventListener('click', () => this.toggleSidebar());
        }
        
        // Filtros de gráficos
        if (this.elements.filterDropdowns) {
            this.elements.filterDropdowns.forEach(item => {
                item.addEventListener('click', (e) => {
                    e.preventDefault();
                    this.updateChartFilter(e.target);
                });
            });
        }
        
        // Toggle de tema claro/escuro
        if (this.elements.themeToggle) {
            this.elements.themeToggle.addEventListener('click', () => this.toggleDarkMode());
        }
        
        // Botão de atualizar dados
        if (this.elements.refreshButton) {
            this.elements.refreshButton.addEventListener('click', () => this.refreshData());
        }
        
        // Botão de exportar
        if (this.elements.exportButton) {
            this.elements.exportButton.addEventListener('click', () => this.exportData());
        }
        
        // Checkboxes de tarefas
        if (this.elements.taskCheckboxes) {
            this.elements.taskCheckboxes.forEach(checkbox => {
                checkbox.addEventListener('change', (e) => this.handleTaskCompletion(e.target));
            });
        }
        
        // Listener para redimensionamento da janela
        window.addEventListener('resize', () => this.handleResponsiveLayout());
        
        // Listener para cliques em notificações
        document.querySelectorAll('.dropdown-item').forEach(item => {
            item.addEventListener('click', (e) => {
                if (e.currentTarget.closest('[aria-labelledby="notificationsDropdown"]')) {
                    this.handleNotificationClick(e.currentTarget);
                }
            });
        });
        
        // Listener para ações de usuários na tabela
        document.querySelectorAll('.dropdown-menu[aria-labelledby^="dropdownMenuButton"] .dropdown-item').forEach(item => {
            item.addEventListener('click', (e) => {
                e.preventDefault();
                const action = e.currentTarget.textContent.trim();
                const userId = e.currentTarget.closest('tr').querySelector('td:first-child').textContent;
                this.handleUserAction(action, userId);
            });
        });
    }
    
    /**
     * Carrega dados iniciais (simulados)
     */
    loadData() {
        // Simular carregamento de dados
        console.log('Carregando dados do dashboard...');
        
        // Dados simulados para o dashboard
        this.dataCache = {
            users: this.generateMockUsers(),
            drones: this.generateMockDrones(),
            maintenance: this.generateMockMaintenance(),
            revenue: this.generateMockRevenue(),
            activities: this.generateMockActivities(),
            tasks: this.generateMockTasks(),
            notifications: this.generateMockNotifications()
        };
        
        console.log('Dados carregados com sucesso!');
    }
    
    /**
     * Inicializa os gráficos
     */
    initCharts() {
        try {
            // Verificar se Chart.js está disponível
            if (typeof Chart === 'undefined') {
                console.error('Chart.js não está disponível!');
                return;
            }
            
            // Configurar gráfico de crescimento de usuários
            if (this.elements.userGrowthChart) {
                const ctx = this.elements.userGrowthChart.getContext('2d');
                this.charts.userGrowth = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: ['Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez', 'Jan', 'Fev', 'Mar'],
                        datasets: [
                            {
                                label: 'Clientes',
                                data: [65, 78, 90, 105, 112, 120, 135, 142, 150, 162, 175, 185],
                                borderColor: '#36b9cc',
                                backgroundColor: 'rgba(54, 185, 204, 0.1)',
                                borderWidth: 2,
                                pointBackgroundColor: '#36b9cc',
                                pointBorderColor: '#fff',
                                pointBorderWidth: 2,
                                pointRadius: 4,
                                tension: 0.3,
                                fill: true
                            },
                            {
                                label: 'Técnicos',
                                data: [25, 28, 30, 32, 35, 38, 40, 42, 45, 48, 50, 52],
                                borderColor: '#1cc88a',
                                backgroundColor: 'rgba(28, 200, 138, 0.1)',
                                borderWidth: 2,
                                pointBackgroundColor: '#1cc88a',
                                pointBorderColor: '#fff',
                                pointBorderWidth: 2,
                                pointRadius: 4,
                                tension: 0.3,
                                fill: true
                            },
                            {
                                label: 'Administradores',
                                data: [5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11],
                                borderColor: '#4e73df',
                                backgroundColor: 'rgba(78, 115, 223, 0.1)',
                                borderWidth: 2,
                                pointBackgroundColor: '#4e73df',
                                pointBorderColor: '#fff',
                                pointBorderWidth: 2,
                                pointRadius: 4,
                                tension: 0.3,
                                fill: true
                            }
                        ]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                position: 'top',
                                align: 'end',
                                labels: {
                                    boxWidth: 12,
                                    usePointStyle: true,
                                    pointStyle: 'circle'
                                }
                            },
                            tooltip: {
                                backgroundColor: 'rgba(0, 0, 0, 0.8)',
                                padding: 12,
                                usePointStyle: true,
                                callbacks: {
                                    label: function(context) {
                                        return context.dataset.label + ': ' + context.raw + ' usuários';
                                    }
                                }
                            }
                        },
                        scales: {
                            x: {
                                grid: {
                                    display: false
                                }
                            },
                            y: {
                                beginAtZero: true,
                                grid: {
                                    color: 'rgba(0, 0, 0, 0.05)'
                                },
                                ticks: {
                                    precision: 0
                                }
                            }
                        }
                    }
                });
            }
            
            // Configurar gráfico de distribuição de usuários
            if (this.elements.userDistributionChart) {
                const ctx = this.elements.userDistributionChart.getContext('2d');
                this.charts.userDistribution = new Chart(ctx, {
                    type: 'doughnut',
                    data: {
                        labels: ['Administradores', 'Técnicos', 'Clientes'],
                        datasets: [{
                            data: [11, 52, 185],
                            backgroundColor: ['#4e73df', '#1cc88a', '#36b9cc'],
                            borderColor: ['#ffffff', '#ffffff', '#ffffff'],
                            borderWidth: 2,
                            hoverOffset: 5
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        cutout: '70%',
                        plugins: {
                            legend: {
                                display: false
                            },
                            tooltip: {
                                backgroundColor: 'rgba(0, 0, 0, 0.8)',
                                padding: 12,
                                callbacks: {
                                    label: function(context) {
                                        const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                        const percentage = Math.round((context.raw / total) * 100);
                                        return context.label + ': ' + context.raw + ' (' + percentage + '%)';
                                    }
                                }
                            }
                        }
                    }
                });
            }
            
            console.log('Gráficos inicializados com sucesso!');
        } catch (error) {
            console.error('Erro ao inicializar gráficos:', error);
        }
    }
    
    /**
     * Inicializa as DataTables
     */
    initDataTables() {
        try {
            // Verificar se DataTables está disponível
            if (typeof $.fn.DataTable === 'undefined') {
                console.warn('DataTables não está disponível!');
                return;
            }
            
            // Inicializar tabela de usuários
            if (this.elements.usersTable) {
                $(this.elements.usersTable).DataTable({
                    language: {
                        url: '//cdn.datatables.net/plug-ins/1.13.6/i18n/pt-BR.json'
                    },
                    pageLength: 5,
                    lengthMenu: [5, 10, 25, 50],
                    responsive: true,
                    dom: '<"d-flex justify-content-between align-items-center mb-3"<"d-flex align-items-center"l><"d-flex"f>>t<"d-flex justify-content-between align-items-center mt-3"<"text-muted"i><"d-flex"p>>',
                    initComplete: function() {
                        // Personalizar elementos da DataTable
                        $('.dataTables_length select').addClass('form-select form-select-sm');
                        $('.dataTables_filter input').addClass('form-control form-control-sm ms-2').attr('placeholder', 'Buscar...');
                        $('.dataTables_info').addClass('small');
                        $('.dataTables_paginate').addClass('pagination-sm');
                    }
                });
            }
            
            console.log('DataTables inicializadas com sucesso!');
        } catch (error) {
            console.error('Erro ao inicializar DataTables:', error);
        }
    }
    
    /**
     * Inicializa tooltips
     */
    initTooltips() {
        try {
            // Verificar se Bootstrap está disponível
            if (typeof bootstrap === 'undefined' || typeof bootstrap.Tooltip === 'undefined') {
                console.warn('Bootstrap Tooltips não estão disponíveis!');
                return;
            }
            
            // Inicializar todos os tooltips
            const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
            tooltipTriggerList.map(function (tooltipTriggerEl) {
                return new bootstrap.Tooltip(tooltipTriggerEl);
            });
            
            console.log('Tooltips inicializados com sucesso!');
        } catch (error) {
            console.error('Erro ao inicializar tooltips:', error);
        }
    }
    
    /**
     * Inicializa o seletor de data
     */
    initDatePicker() {
        // Implementação futura se necessário
        // Requer biblioteca adicional como Flatpickr ou DateRangePicker
    }
    
    /**
     * Atualiza a data e hora exibidas
     */
    updateDateTime() {
        if (this.elements.currentDate) {
            const options = { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' };
            const currentDate = new Date().toLocaleDateString('pt-BR', options);
            this.elements.currentDate.textContent = currentDate;
        }
    }
    
    /**
     * Alterna a visibilidade da sidebar
     */
    toggleSidebar() {
        if (this.elements.sidebar && this.elements.main) {
            this.sidebarVisible = !this.sidebarVisible;
            
            if (window.innerWidth < 768) {
                // Em dispositivos móveis
                this.elements.sidebar.classList.toggle('show');
            } else {
                // Em desktop
                if (this.sidebarVisible) {
                    this.elements.sidebar.style.transform = 'translateX(0)';
                    this.elements.main.style.marginLeft = '250px';
                } else {
                    this.elements.sidebar.style.transform = 'translateX(-100%)';
                    this.elements.main.style.marginLeft = '0';
                }
            }
            
            // Atualizar ícone do botão
            const icon = this.elements.sidebarToggle.querySelector('i');
            if (icon) {
                icon.className = this.sidebarVisible ? 'bi bi-x-lg' : 'bi bi-list';
            }
            
            // Disparar evento de redimensionamento para ajustar gráficos
            window.dispatchEvent(new Event('resize'));
        }
    }
    
    /**
     * Atualiza o filtro do gráfico
     * @param {HTMLElement} filterElement - Elemento do filtro clicado
     */
    updateChartFilter(filterElement) {
        // Remover classe ativa de todos os itens
        const dropdownItems = filterElement.closest('.dropdown-menu').querySelectorAll('.dropdown-item');
        dropdownItems.forEach(item => item.classList.remove('active'));
        
        // Adicionar classe ativa ao item selecionado
        filterElement.classList.add('active');
        
        // Atualizar texto do botão dropdown
        const dropdownButton = document.getElementById('chartDropdown');
        if (dropdownButton) {
            dropdownButton.textContent = filterElement.textContent;
        }
        
        // Atualizar dados do gráfico com base no filtro selecionado
        const filterValue = filterElement.textContent.trim();
        this.updateChartData(filterValue);
    }
    
    /**
     * Atualiza os dados do gráfico com base no filtro
     * @param {string} filterValue - Valor do filtro selecionado
     */
    updateChartData(filterValue) {
        // Implementação de exemplo - atualizar dados do gráfico com base no filtro
        console.log(`Atualizando gráfico para o período: ${filterValue}`);
        
        // Dados simulados para diferentes períodos
        let data;
        
        switch (filterValue) {
            case 'Últimos 30 Dias':
                data = {
                    clients: [150, 155, 160, 165, 170, 175, 180, 185],
                    technicians: [48, 49, 50, 51, 52, 52, 52, 52],
                    admins: [10, 10, 10, 11, 11, 11, 11, 11]
                };
                break;
                
            case 'Desde o Início':
                data = {
                    clients: [10, 25, 45, 65, 85, 105, 125, 145, 165, 185],
                    technicians: [5, 10, 15, 20, 25, 30, 35, 40, 45, 52],
                    admins: [2, 3, 4, 5, 6, 7, 8, 9, 10, 11]
                };
                break;
                
            default: // Últimos 12 Meses
                data = {
                    clients: [65, 78, 90, 105, 112, 120, 135, 142, 150, 162, 175, 185],
                    technicians: [25, 28, 30, 32, 35, 38, 40, 42, 45, 48, 50, 52],
                    admins: [5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11]
                };
        }
        
        // Atualizar dados do gráfico
        if (this.charts.userGrowth) {
            this.charts.userGrowth.data.datasets[0].data = data.clients;
            this.charts.userGrowth.data.datasets[1].data = data.technicians;
            this.charts.userGrowth.data.datasets[2].data = data.admins;
            this.charts.userGrowth.update();
        }
    }
    
    /**
     * Alterna entre modo claro e escuro
     */
    toggleDarkMode() {
        this.darkMode = !this.darkMode;
        document.body.setAttribute('data-theme', this.darkMode ? 'dark' : 'light');
        
        // Atualizar ícone do botão
        if (this.elements.themeToggle) {
            const icon = this.elements.themeToggle.querySelector('i');
            if (icon) {
                icon.className = this.darkMode ? 'bi bi-sun' : 'bi bi-moon';
            }
        }
        
        // Atualizar gráficos para o tema atual
        this.updateChartsTheme();
        
        console.log(`Modo ${this.darkMode ? 'escuro' : 'claro'} ativado`);
    }
    
    /**
     * Atualiza o tema dos gráficos
     */
    updateChartsTheme() {
        // Implementação futura - atualizar cores dos gráficos com base no tema
    }
    
    /**
     * Atualiza os dados do dashboard
     */
    refreshData() {
        console.log('Atualizando dados do dashboard...');
        
        // Simular carregamento de dados
        const loadingToast = this.showToast('Atualizando dados...', 'info');
        
        // Simular delay de rede
        setTimeout(() => {
            // Atualizar dados simulados
            this.loadData();
            
            // Atualizar gráficos
            this.updateCharts();
            
            // Atualizar estatísticas
            this.updateStats();
            
            // Esconder toast de carregamento
            this.hideToast(loadingToast);
            
            // Mostrar toast de sucesso
            this.showToast('Dados atualizados com sucesso!', 'success');
            
            console.log('Dados atualizados com sucesso!');
        }, 1500);
    }
    
    /**
     * Atualiza os gráficos com novos dados
     */
    updateCharts() {
        // Implementação futura - atualizar gráficos com novos dados
        if (this.charts.userGrowth) {
            this.charts.userGrowth.update();
        }
        
        if (this.charts.userDistribution) {
            this.charts.userDistribution.update();
        }
    }
    
    /**
     * Atualiza as estatísticas do dashboard
     */
    updateStats() {
        // Implementação futura - atualizar estatísticas com novos dados
    }
    
    /**
     * Exporta dados para CSV
     */
    exportData() {
        console.log('Exportando dados...');
        
        // Implementação de exemplo - exportar dados da tabela de usuários
        try {
            // Verificar se há dados para exportar
            if (!this.dataCache.users || this.dataCache.users.length === 0) {
                this.showToast('Não há dados para exportar!', 'warning');
                return;
            }
            
            // Criar cabeçalho CSV
            let csv = 'ID,Nome,Email,Tipo,Status,Registrado em\n';
            
            // Adicionar linhas de dados
            this.dataCache.users.forEach(user => {
                csv += `${user.id},${user.name},${user.email},${user.type},${user.status},${user.registered}\n`;
            });
            
            // Criar blob e link para download
            const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.setAttribute('href', url);
            link.setAttribute('download', 'usuarios_evo_drones.csv');
            link.style.visibility = 'hidden';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            
            this.showToast('Dados exportados com sucesso!', 'success');
        } catch (error) {
            console.error('Erro ao exportar dados:', error);
            this.showToast('Erro ao exportar dados!', 'danger');
        }
    }
    
    /**
     * Manipula a conclusão de uma tarefa
     * @param {HTMLElement} checkbox - Elemento checkbox da tarefa
     */
    handleTaskCompletion(checkbox) {
        const taskLabel = checkbox.nextElementSibling;
        if (taskLabel) {
            if (checkbox.checked) {
                taskLabel.style.textDecoration = 'line-through';
                taskLabel.style.opacity = '0.6';
                this.showToast(`Tarefa "${taskLabel.textContent.trim()}" concluída!`, 'success');
            } else {
                taskLabel.style.textDecoration = 'none';
                taskLabel.style.opacity = '1';
            }
        }
        
        // Atualizar progresso das tarefas
        this.updateTaskProgress();
    }
    
    /**
     * Atualiza o progresso das tarefas
     */
    updateTaskProgress() {
        const totalTasks = this.elements.taskCheckboxes.length;
        const completedTasks = Array.from(this.elements.taskCheckboxes).filter(checkbox => checkbox.checked).length;
        const progressPercentage = (completedTasks / totalTasks) * 100;
        
        // Atualizar texto de progresso
        const progressText = document.querySelector('.card-footer .text-muted.small');
        if (progressText) {
            progressText.textContent = `${completedTasks} de ${totalTasks} tarefas concluídas`;
        }
        
        // Atualizar barra de progresso
        const progressBar = document.querySelector('.card-footer .progress-bar');
        if (progressBar) {
            progressBar.style.width = `${progressPercentage}%`;
            progressBar.setAttribute('aria-valuenow', progressPercentage);
        }
    }
    
    /**
     * Manipula o clique em uma notificação
     * @param {HTMLElement} notificationElement - Elemento da notificação clicada
     */
    handleNotificationClick(notificationElement) {
        const notificationText = notificationElement.textContent.trim();
        console.log(`Notificação clicada: ${notificationText}`);
        
        // Implementação de exemplo - marcar notificação como lida
        if (notificationElement.classList.contains('unread')) {
            notificationElement.classList.remove('unread');
            
            // Atualizar contador de notificações
            const unreadCount = document.querySelectorAll('.dropdown-item.unread').length;
            if (this.elements.notificationBadge) {
                this.elements.notificationBadge.textContent = unreadCount;
                
                // Esconder badge se não houver notificações não lidas
                if (unreadCount === 0) {
                    this.elements.notificationBadge.style.display = 'none';
                }
            }
        }
    }
    
    /**
     * Manipula ações de usuário na tabela
     * @param {string} action - Ação selecionada
     * @param {string} userId - ID do usuário
     */
    handleUserAction(action, userId) {
        console.log(`Ação "${action}" para o usuário ${userId}`);
        
        // Implementação de exemplo - diferentes ações
        switch (action) {
            case 'Ver Detalhes':
                this.showUserDetails(userId);
                break;
                
            case 'Editar':
                this.editUser(userId);
                break;
                
            case 'Excluir':
                this.deleteUser(userId);
                break;
                
            default:
                console.warn(`Ação desconhecida: ${action}`);
        }
    }
    
    /**
     * Mostra detalhes de um usuário
     * @param {string} userId - ID do usuário
     */
    showUserDetails(userId) {
        // Implementação de exemplo - mostrar modal com detalhes do usuário
        console.log(`Mostrando detalhes do usuário ${userId}`);
        
        // Simular busca de dados do usuário
        const user = this.dataCache.users.find(u => u.id === userId);
        
        if (!user) {
            this.showToast(`Usuário ${userId} não encontrado!`, 'warning');
            return;
        }
        
        // Aqui seria implementada a lógica para exibir um modal com os detalhes
        alert(`Detalhes do usuário ${userId}:\n\nNome: ${user.name}\nEmail: ${user.email}\nTipo: ${user.type}\nStatus: ${user.status}\nRegistrado em: ${user.registered}`);
    }
    
    /**
     * Edita um usuário
     * @param {string} userId - ID do usuário
     */
    editUser(userId) {
        // Implementação de exemplo - redirecionar para página de edição
        console.log(`Editando usuário ${userId}`);
        
        // Aqui seria implementada a lógica para redirecionar ou abrir modal de edição
        alert(`Editar usuário ${userId} - Funcionalidade a ser implementada`);
    }
    
    /**
     * Exclui um usuário
     * @param {string} userId - ID do usuário
     */
    deleteUser(userId) {
        // Implementação de exemplo - confirmar exclusão
        console.log(`Solicitação para excluir usuário ${userId}`);
        
        if (confirm(`Tem certeza que deseja excluir o usuário ${userId}?`)) {
            // Simular exclusão
            console.log(`Excluindo usuário ${userId}...`);
            
            // Remover da lista de dados em cache
            this.dataCache.users = this.dataCache.users.filter(u => u.id !== userId);
            
            // Remover da tabela
            const row = document.querySelector(`tr td:first-child:contains('${userId}')`).closest('tr');
            if (row && typeof $.fn.DataTable !== 'undefined') {
                $(this.elements.usersTable).DataTable().row(row).remove().draw();
            }
            
            this.showToast(`Usuário ${userId} excluído com sucesso!`, 'success');
        }
    }
    
    /**
     * Manipula o layout responsivo
     */
    handleResponsiveLayout() {
        const isMobile = window.innerWidth < 768;
        
        if (isMobile) {
            // Em dispositivos móveis
            if (this.elements.sidebar) {
                this.elements.sidebar.classList.remove('show');
                this.sidebarVisible = false;
            }
            
            if (this.elements.main) {
                this.elements.main.style.marginLeft = '0';
            }
        } else {
            // Em desktop
            if (this.elements.sidebar) {
                this.elements.sidebar.style.transform = this.sidebarVisible ? 'translateX(0)' : 'translateX(-100%)';
            }
            
            if (this.elements.main) {
                this.elements.main.style.marginLeft = this.sidebarVisible ? '250px' : '0';
            }
        }
        
        // Atualizar gráficos para o novo tamanho
        if (this.charts.userGrowth) {
            this.charts.userGrowth.resize();
        }
        
        if (this.charts.userDistribution) {
            this.charts.userDistribution.resize();
        }
    }
    
    /**
     * Configura a atualização automática de dados
     */
    setupAutoRefresh() {
        // Limpar intervalo existente, se houver
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
        }
        
        // Configurar novo intervalo (a cada 5 minutos)
        this.refreshInterval = setInterval(() => {
            console.log('Atualizando dados automaticamente...');
            this.refreshData();
        }, 5 * 60 * 1000); // 5 minutos
    }
    
    /**
     * Mostra um toast de notificação
     * @param {string} message - Mensagem a ser exibida
     * @param {string} type - Tipo de toast (success, info, warning, danger)
     * @returns {HTMLElement} - Elemento do toast criado
     */
    showToast(message, type = 'info') {
        // Verificar se o container de toasts existe, senão criar
        let toastContainer = document.querySelector('.toast-container');
        if (!toastContainer) {
            toastContainer = document.createElement('div');
            toastContainer.className = 'toast-container position-fixed bottom-0 end-0 p-3';
            document.body.appendChild(toastContainer);
        }
        
        // Criar elemento de toast
        const toastId = `toast-${Date.now()}`;
        const toast = document.createElement('div');
        toast.className = `toast align-items-center text-white bg-${type} border-0`;
        toast.id = toastId;
        toast.setAttribute('role', 'alert');
        toast.setAttribute('aria-live', 'assertive');
        toast.setAttribute('aria-atomic', 'true');
        
        // Conteúdo do toast
        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">
                    ${message}
                </div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Fechar"></button>
            </div>
        `;
        
        // Adicionar ao container
        toastContainer.appendChild(toast);
        
        // Inicializar e mostrar o toast
        if (typeof bootstrap !== 'undefined' && typeof bootstrap.Toast !== 'undefined') {
            const bsToast = new bootstrap.Toast(toast, {
                autohide: true,
                delay: 3000
            });
            bsToast.show();
        } else {
            // Fallback se Bootstrap não estiver disponível
            toast.style.display = 'block';
            setTimeout(() => {
                toast.remove();
            }, 3000);
        }
        
        return toast;
    }
    
    /**
     * Esconde um toast
     * @param {HTMLElement} toast - Elemento do toast a ser escondido
     */
    hideToast(toast) {
        if (!toast) return;
        
        if (typeof bootstrap !== 'undefined' && typeof bootstrap.Toast !== 'undefined') {
            const bsToast = bootstrap.Toast.getInstance(toast);
            if (bsToast) {
                bsToast.hide();
            }
        } else {
            // Fallback se Bootstrap não estiver disponível
            toast.style.display = 'none';
            toast.remove();
        }
    }
    
    /**
     * Gera dados simulados de usuários
     * @returns {Array} - Array de objetos de usuário
     */
    generateMockUsers() {
        return [
            { id: '#U-1234', name: 'João Silva', email: 'joao.silva@email.com', type: 'Cliente', status: 'Ativo', registered: '10/04/2023' },
            { id: '#U-1233', name: 'Maria Santos', email: 'maria.santos@email.com', type: 'Técnico', status: 'Ativo', registered: '09/04/2023' },
            { id: '#U-1232', name: 'Carlos Oliveira', email: 'carlos.oliveira@email.com', type: 'Admin', status: 'Ativo', registered: '08/04/2023' },
            { id: '#U-1231', name: 'Ana Pereira', email: 'ana.pereira@email.com', type: 'Cliente', status: 'Inativo', registered: '05/04/2023' },
            { id: '#U-1230', name: 'Roberto Almeida', email: 'roberto.almeida@email.com', type: 'Cliente', status: 'Ativo', registered: '01/04/2023' }
        ];
    }
    
    /**
     * Gera dados simulados de drones
     * @returns {Array} - Array de objetos de drone
     */
    generateMockDrones() {
        return [
            { id: 'DRN-001', model: 'EVO-X1', status: 'Ativo', battery: 85, lastMaintenance: '15/03/2023' },
            { id: 'DRN-002', model: 'EVO-X2', status: 'Em Manutenção', battery: 0, lastMaintenance: '10/04/2023' },
            { id: 'DRN-003', model: 'EVO-X1', status: 'Ativo', battery: 72, lastMaintenance: '01/04/2023' },
            { id: 'DRN-004', model: 'EVO-X3', status: 'Inativo', battery: 15, lastMaintenance: '20/02/2023' },
            { id: 'DRN-005', model: 'EVO-X2', status: 'Ativo', battery: 93, lastMaintenance: '05/04/2023' }
        ];
    }
    
    /**
     * Gera dados simulados de manutenções
     * @returns {Array} - Array de objetos de manutenção
     */
    generateMockMaintenance() {
        return [
            { id: 'MNT-001', droneId: 'DRN-002', type: 'Preventiva', status: 'Em Andamento', technician: 'Maria Santos', startDate: '10/04/2023', estimatedCompletion: '12/04/2023' },
            { id: 'MNT-002', droneId: 'DRN-004', type: 'Corretiva', status: 'Agendada', technician: 'Pedro Costa', startDate: '15/04/2023', estimatedCompletion: '18/04/2023' },
            { id: 'MNT-003', droneId: 'DRN-001', type: 'Preventiva', status: 'Concluída', technician: 'Maria Santos', startDate: '15/03/2023', estimatedCompletion: '15/03/2023' }
        ];
    }
    
    /**
     * Gera dados simulados de receita
     * @returns {Object} - Objeto com dados de receita
     */
    generateMockRevenue() {
        return {
            monthly: 58500,
            target: 90000,
            growth: 5,
            breakdown: {
                services: 35000,
                sales: 15000,
                maintenance: 8500
            },
            history: [45000, 48000, 52000, 49000, 53000, 56000, 58500]
        };
    }
    
    /**
     * Gera dados simulados de atividades
     * @returns {Array} - Array de objetos de atividade
     */
    generateMockActivities() {
        return [
            { type: 'user', title: 'Novo Usuário', description: 'Novo usuário João Silva foi registrado no sistema.', time: '2h atrás' },
            { type: 'system', title: 'Atualização de Sistema', description: 'Sistema atualizado para a versão 2.5.0.', time: '4h atrás' },
            { type: 'alert', title: 'Alerta de Segurança', description: 'Tentativa de login não autorizada detectada.', time: '8h atrás' },
            { type: 'backup', title: 'Backup Concluído', description: 'Backup automático do banco de dados concluído com sucesso.', time: '1d atrás' }
        ];
    }
    
    /**
     * Gera dados simulados de tarefas
     * @returns {Array} - Array de objetos de tarefa
     */
    generateMockTasks() {
        return [
            { id: 'task1', description: 'Revisar relatórios mensais', priority: 'Hoje', completed: false },
            { id: 'task2', description: 'Aprovar solicitações de manutenção', priority: 'Urgente', completed: false },
            { id: 'task3', description: 'Atualizar políticas de segurança', priority: 'Esta Semana', completed: false },
            { id: 'task4', description: 'Reunião com equipe técnica', priority: 'Amanhã', completed: false },
            { id: 'task5', description: 'Planejar expansão da frota', priority: 'Próxima Semana', completed: false }
        ];
    }
    
    /**
     * Gera dados simulados de notificações
     * @returns {Array} - Array de objetos de notificação
     */
    generateMockNotifications() {
        return [
            { id: 'notif1', title: 'Novo usuário registrado', read: false, time: '10 minutos atrás' },
            { id: 'notif2', title: 'Alerta de segurança', read: false, time: '30 minutos atrás' },
            { id: 'notif3', title: 'Atualização disponível', read: false, time: '1 hora atrás' },
            { id: 'notif4', title: 'Manutenção concluída', read: true, time: '3 horas atrás' },
            { id: 'notif5', title: 'Novo relatório disponível', read: true, time: '1 dia atrás' }
        ];
    }
}

// Inicializar o dashboard quando o DOM estiver carregado
document.addEventListener('DOMContentLoaded', function() {
    // Verificar se estamos na página do dashboard admin
    if (document.querySelector('.admin-dashboard')) {
        window.adminDashboard = new AdminDashboard();
    }
});