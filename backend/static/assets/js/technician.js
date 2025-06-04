/**
 * Technician Dashboard JavaScript - EVO DRONES
 * Este arquivo contém todas as funcionalidades específicas para o painel do técnico
 */

class TechnicianDashboard {
    constructor() {
        // Propriedades
        this.charts = {};
        this.filters = {
            maintenanceStatus: 'all',
            maintenanceType: 'all',
            dateRange: 'month'
        };
        this.notifications = [];
        this.maintenanceOrders = [];
        this.inventoryItems = [];
        this.activities = [];
        this.schedules = [];
        
        // Inicialização
        this.init();
    }
    
    /**
     * Inicializa o dashboard
     */
    init() {
        this.cacheElements();
        this.setupEventListeners();
        this.loadData();
        this.initCharts();
        this.initCalendar();
        this.setupResponsiveness();
        this.startAutoRefresh();
    }
    
    /**
     * Armazena referências aos elementos DOM
     */
    cacheElements() {
        // Elementos do sidebar
        this.sidebar = document.getElementById('sidebar');
        this.sidebarToggle = document.getElementById('sidebarToggle');
        this.sidebarToggleMobile = document.getElementById('sidebarToggleMobile');
        
        // Elementos de notificação
        this.notificationsDropdown = document.getElementById('notificationsDropdown');
        
        // Elementos de gráficos
        this.maintenanceChartEl = document.getElementById('maintenanceChart');
        this.maintenanceTypeChartEl = document.getElementById('maintenanceTypeChart');
        
        // Elementos de data
        this.currentDateEl = document.getElementById('currentDate');
        
        // Elementos de filtro
        this.chartDropdown = document.getElementById('chartDropdown');
    }
    
    /**
     * Configura os event listeners
     */
    setupEventListeners() {
        // Toggle do sidebar
        if (this.sidebarToggle) {
            this.sidebarToggle.addEventListener('click', () => this.toggleSidebar());
        }
        
        if (this.sidebarToggleMobile) {
            this.sidebarToggleMobile.addEventListener('click', () => this.toggleSidebar());
        }
        
        // Dropdown de notificações
        if (this.notificationsDropdown) {
            this.notificationsDropdown.addEventListener('click', (e) => {
                e.preventDefault();
                this.loadNotifications();
            });
        }
        
        // Filtros de gráficos
        document.querySelectorAll('.dropdown-menu .dropdown-item').forEach(item => {
            item.addEventListener('click', (e) => {
                e.preventDefault();
                const filter = e.target.textContent.trim();
                const dropdown = e.target.closest('.dropdown');
                if (dropdown) {
                    const button = dropdown.querySelector('.dropdown-toggle');
                    if (button) {
                        button.textContent = filter;
                        this.updateFilters(dropdown.id, filter);
                    }
                }
            });
        });
        
        // Eventos de clique em itens da tabela
        document.querySelectorAll('table tbody tr').forEach(row => {
            row.addEventListener('click', () => {
                const orderId = row.cells[0].textContent;
                this.viewMaintenanceDetails(orderId);
            });
        });
        
        // Eventos de clique em eventos do calendário
        document.querySelectorAll('.tech-calendar-event').forEach(event => {
            event.addEventListener('click', (e) => {
                const eventText = e.target.textContent;
                this.viewScheduleDetails(eventText);
            });
        });
        
        // Eventos de redimensionamento da janela
        window.addEventListener('resize', () => this.handleResize());
    }
    
    /**
     * Carrega dados do servidor
     */
    loadData() {
        // Simulação de carregamento de dados
        // Em um ambiente real, isso seria substituído por chamadas de API
        this.loadMaintenanceOrders();
        this.loadInventoryItems();
        this.loadActivities();
        this.loadSchedules();
        this.updateDateTime();
    }
    
    /**
     * Carrega ordens de manutenção
     */
    loadMaintenanceOrders() {
        // Simulação - Em um ambiente real, isso seria uma chamada de API
        console.log('Carregando ordens de manutenção...');
        
        // Dados simulados
        this.maintenanceOrders = [
            {
                id: '#OM-1234',
                drone: 'Phantom Pro',
                client: 'Empresa ABC',
                type: 'Preventiva',
                status: 'Em Progresso',
                date: '10/04/2023',
                details: 'Manutenção preventiva programada para o drone Phantom Pro.'
            },
            {
                id: '#OM-1233',
                drone: 'EVO X3',
                client: 'Tech Solutions',
                type: 'Emergencial',
                status: 'Concluída',
                date: '09/04/2023',
                details: 'Reparo emergencial após queda durante voo.'
            },
            {
                id: '#OM-1232',
                drone: 'Mavic Air',
                client: 'João Silva',
                type: 'Corretiva',
                status: 'Aguardando Peças',
                date: '08/04/2023',
                details: 'Substituição de componentes da câmera.'
            },
            {
                id: '#OM-1231',
                drone: 'Inspire 2',
                client: 'Aero Filmagens',
                type: 'Preventiva',
                status: 'Agendada',
                date: '12/04/2023',
                details: 'Verificação de rotina e calibração de sensores.'
            }
        ];
    }
    
    /**
     * Carrega itens do inventário
     */
    loadInventoryItems() {
        // Simulação - Em um ambiente real, isso seria uma chamada de API
        console.log('Carregando itens do inventário...');
        
        // Dados simulados
        this.inventoryItems = [
            {
                id: 1,
                name: 'Bateria LiPo 5200mAh',
                category: 'Baterias',
                compatibility: 'Phantom Pro',
                quantity: 2,
                status: 'Crítico'
            },
            {
                id: 2,
                name: 'Controlador de Voo F7',
                category: 'Eletrônicos',
                compatibility: 'EVO X3',
                quantity: 5,
                status: 'Baixo'
            },
            {
                id: 3,
                name: 'Gimbal 3-Eixos',
                category: 'Câmeras',
                compatibility: 'Mavic Air',
                quantity: 4,
                status: 'Baixo'
            },
            {
                id: 4,
                name: 'Hélices 9450',
                category: 'Peças',
                compatibility: 'Phantom Pro',
                quantity: 12,
                status: 'Normal'
            },
            {
                id: 5,
                name: 'Antena 5.8GHz',
                category: 'Comunicação',
                compatibility: 'Transmissores',
                quantity: 8,
                status: 'Normal'
            }
        ];
    }
    
    /**
     * Carrega atividades recentes
     */
    loadActivities() {
        // Simulação - Em um ambiente real, isso seria uma chamada de API
        console.log('Carregando atividades recentes...');
        
        // Dados simulados
        this.activities = [
            {
                id: 1,
                type: 'maintenance_completed',
                title: 'Manutenção Concluída',
                description: 'Você concluiu a manutenção do drone Phantom Pro #1233.',
                time: '2h atrás',
                icon: 'check2-circle',
                color: 'success'
            },
            {
                id: 2,
                type: 'new_order',
                title: 'Nova Ordem Atribuída',
                description: 'Nova ordem de manutenção #1234 foi atribuída a você.',
                time: '4h atrás',
                icon: 'clipboard-check',
                color: 'primary'
            },
            {
                id: 3,
                type: 'parts_requested',
                title: 'Peças Solicitadas',
                description: 'Sua solicitação de peças foi aprovada pelo gerente.',
                time: '8h atrás',
                icon: 'box-seam',
                color: 'warning'
            },
            {
                id: 4,
                type: 'report_sent',
                title: 'Relatório Enviado',
                description: 'Você enviou o relatório técnico para o cliente Empresa ABC.',
                time: '1d atrás',
                icon: 'file-earmark-text',
                color: 'info'
            },
            {
                id: 5,
                type: 'login',
                title: 'Login no Sistema',
                description: 'Você fez login no sistema às 08:45.',
                time: '1d atrás',
                icon: 'box-arrow-in-right',
                color: 'secondary'
            }
        ];
    }
    
    /**
     * Carrega agendamentos
     */
    loadSchedules() {
        // Simulação - Em um ambiente real, isso seria uma chamada de API
        console.log('Carregando agendamentos...');
        
        // Dados simulados
        this.schedules = [
            {
                id: 1,
                title: 'Manutenção #1234',
                date: '2023-04-10',
                time: '09:00',
                type: 'maintenance',
                color: 'primary'
            },
            {
                id: 2,
                title: 'Manutenção #1235',
                date: '2023-04-10',
                time: '14:00',
                type: 'maintenance',
                color: 'warning'
            },
            {
                id: 3,
                title: 'Treinamento',
                date: '2023-04-11',
                time: '10:00',
                type: 'training',
                color: 'success'
            },
            {
                id: 4,
                title: 'Inventário',
                date: '2023-04-12',
                time: '13:00',
                type: 'inventory',
                color: 'info'
            },
            {
                id: 5,
                title: 'Manutenção #1236',
                date: '2023-04-13',
                time: '09:00',
                type: 'maintenance',
                color: 'primary'
            },
            {
                id: 6,
                title: 'Manutenção #1237',
                date: '2023-04-13',
                time: '15:00',
                type: 'maintenance',
                color: 'primary'
            },
            {
                id: 7,
                title: 'Emergência #E22',
                date: '2023-04-14',
                time: '11:00',
                type: 'emergency',
                color: 'danger'
            }
        ];
    }
    
    /**
     * Carrega notificações
     */
    loadNotifications() {
        // Simulação - Em um ambiente real, isso seria uma chamada de API
        console.log('Carregando notificações...');
        
        // Dados simulados
        this.notifications = [
            {
                id: 1,
                title: 'Nova ordem de manutenção #1234',
                time: '30 minutos atrás',
                read: false,
                type: 'maintenance'
            },
            {
                id: 2,
                title: 'Peça solicitada aprovada',
                time: '2 horas atrás',
                read: false,
                type: 'inventory'
            },
            {
                id: 3,
                title: 'Lembrete: Manutenção agendada',
                time: '5 horas atrás',
                read: false,
                type: 'reminder'
            },
            {
                id: 4,
                title: 'Relatório técnico aprovado',
                time: '1 dia atrás',
                read: true,
                type: 'report'
            },
            {
                id: 5,
                title: 'Novo treinamento disponível',
                time: '2 dias atrás',
                read: true,
                type: 'training'
            }
        ];
    }
    
    /**
     * Inicializa os gráficos
     */
    initCharts() {
        // Os gráficos já estão sendo inicializados no HTML inline
        // Esta função pode ser usada para inicializar gráficos adicionais ou
        // para atualizar os gráficos existentes com novos dados
        console.log('Gráficos inicializados');
    }
    
    /**
     * Inicializa o calendário
     */
    initCalendar() {
        console.log('Calendário inicializado');
        // Aqui poderia ser implementada uma biblioteca de calendário como FullCalendar
    }
    
    /**
     * Atualiza a data e hora atual
     */
    updateDateTime() {
        if (this.currentDateEl) {
            const options = { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' };
            const currentDate = new Date().toLocaleDateString('pt-BR', options);
            this.currentDateEl.textContent = currentDate;
        }
    }
    
    /**
     * Alterna a visibilidade do sidebar
     */
    toggleSidebar() {
        if (this.sidebar) {
            this.sidebar.classList.toggle('show');
        }
    }
    
    /**
     * Atualiza os filtros e recarrega os dados
     */
    updateFilters(filterId, value) {
        switch (filterId) {
            case 'chartDropdown':
                if (value === 'Esta Semana') {
                    this.filters.dateRange = 'week';
                } else if (value === 'Este Mês') {
                    this.filters.dateRange = 'month';
                } else if (value === 'Este Ano') {
                    this.filters.dateRange = 'year';
                }
                break;
            // Adicione mais casos conforme necessário
        }
        
        // Recarregar dados com base nos novos filtros
        this.refreshData();
    }
    
    /**
     * Atualiza os dados com base nos filtros
     */
    refreshData() {
        // Aqui você atualizaria os dados com base nos filtros
        console.log('Atualizando dados com filtros:', this.filters);
        
        // Recarregar dados
        this.loadData();
        
        // Atualizar gráficos
        this.updateCharts();
    }
    
    /**
     * Atualiza os gráficos com novos dados
     */
    updateCharts() {
        // Aqui você atualizaria os gráficos com novos dados
        console.log('Atualizando gráficos');
    }
    
    /**
     * Exibe detalhes de uma ordem de manutenção
     */
    viewMaintenanceDetails(orderId) {
        // Encontrar a ordem pelo ID
        const order = this.maintenanceOrders.find(order => order.id === orderId);
        if (order) {
            console.log('Visualizando detalhes da ordem:', order);
            // Aqui você poderia abrir um modal ou navegar para uma página de detalhes
            alert(`Detalhes da Ordem ${orderId}:\n${order.details}`);
        }
    }
    
    /**
     * Exibe detalhes de um agendamento
     */
    viewScheduleDetails(eventText) {
        // Encontrar o agendamento pelo título
        const schedule = this.schedules.find(schedule => schedule.title === eventText);
        if (schedule) {
            console.log('Visualizando detalhes do agendamento:', schedule);
            // Aqui você poderia abrir um modal ou navegar para uma página de detalhes
            alert(`Agendamento: ${schedule.title}\nData: ${schedule.date} às ${schedule.time}`);
        }
    }
    
    /**
     * Configura a responsividade
     */
    setupResponsiveness() {
        this.handleResize();
    }
    
    /**
     * Manipula o redimensionamento da janela
     */
    handleResize() {
        const width = window.innerWidth;
        
        // Ajustar layout com base na largura da janela
        if (width < 768) {
            // Layout móvel
            if (this.sidebar && this.sidebar.classList.contains('show')) {
                this.sidebar.classList.remove('show');
            }
        }
    }
    
    /**
     * Inicia a atualização automática dos dados
     */
    startAutoRefresh() {
        // Atualizar dados a cada 5 minutos
        setInterval(() => {
            this.loadData();
        }, 5 * 60 * 1000);
        
        // Atualizar data e hora a cada minuto
        setInterval(() => {
            this.updateDateTime();
        }, 60 * 1000);
    }
}

// Inicializar o dashboard quando o DOM estiver carregado
document.addEventListener('DOMContentLoaded', () => {
    const dashboard = new TechnicianDashboard();
});