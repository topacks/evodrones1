// EVO DRONES - Dashboard do Cliente

class ClientDashboard {
    constructor() {
        this.currentSection = 'dashboard';
        this.drones = [];
        this.maintenanceOrders = [];
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.initializeAOS();
        this.loadUserData();
        this.updateDashboardStats();
        this.loadRecentActivities();
    }

    setupEventListeners() {
        // Toggle sidebar mobile
        document.getElementById('menu-toggle')?.addEventListener('click', () => {
            document.getElementById('sidebar-wrapper').classList.toggle('show');
        });

        // Close sidebar mobile
        document.getElementById('close-sidebar')?.addEventListener('click', () => {
            document.getElementById('sidebar-wrapper').classList.remove('show');
        });

        // Update time
        setInterval(() => this.updateLastUpdate(), 60000);
    }

    initializeAOS() {
        if (typeof AOS !== 'undefined') {
            AOS.init({
                duration: 600,
                easing: 'ease-out',
                once: true
            });
        }
    }

    showSection(sectionName) {
        // Hide all sections
        document.querySelectorAll('.section-content').forEach(section => {
            section.classList.add('d-none');
        });

        // Show selected section
        const targetSection = document.getElementById(`${sectionName}-section`);
        if (targetSection) {
            targetSection.classList.remove('d-none');
            this.currentSection = sectionName;
        }

        // Update sidebar active state
        document.querySelectorAll('.sidebar .list-group-item').forEach(item => {
            item.classList.remove('active');
        });

        // Add active class to current section
        const activeItem = document.querySelector(`[onclick="showSection('${sectionName}')"]`);
        if (activeItem) {
            activeItem.classList.add('active');
        }

        // Load section specific content
        this.loadSectionContent(sectionName);
    }

    loadSectionContent(sectionName) {
        switch (sectionName) {
            case 'drones':
                this.loadDronesSection();
                break;
            case 'maintenance':
                this.loadMaintenanceSection();
                break;
            case 'orders':
                this.loadOrdersSection();
                break;
            case 'history':
                this.loadHistorySection();
                break;
            case 'profile':
                this.loadProfileSection();
                break;
        }
    }

    loadUserData() {
        // Simular carregamento de dados do usuário
        const userData = {
            name: 'João Silva',
            email: 'joao@email.com',
            phone: '(11) 99999-9999'
        };

        document.getElementById('user-name').textContent = userData.name;
    }

    updateDashboardStats() {
        // Atualizar estatísticas da dashboard
        document.getElementById('total-drones').textContent = this.drones.length;
        document.getElementById('active-orders').textContent = this.maintenanceOrders.filter(o => o.status === 'active').length;
        document.getElementById('completed-orders').textContent = this.maintenanceOrders.filter(o => o.status === 'completed').length;
        document.getElementById('drones-count').textContent = this.drones.length;
        document.getElementById('orders-count').textContent = this.maintenanceOrders.filter(o => o.status === 'active').length;
    }

    loadRecentActivities() {
        const activitiesContainer = document.getElementById('recent-activities');
        
        if (this.maintenanceOrders.length === 0) {
            activitiesContainer.innerHTML = `
                <div class="text-center py-4 text-muted">
                    <i class="bi bi-clock-history fs-1 mb-2 d-block"></i>
                    Nenhuma atividade recente
                </div>
            `;
            return;
        }

        // Implementar timeline de atividades
        const activities = this.maintenanceOrders.slice(0, 5);
        activitiesContainer.innerHTML = activities.map(activity => `
            <div class="timeline-item">
                <div class="d-flex justify-content-between align-items-start">
                    <div>
                        <h6 class="mb-1">${activity.title}</h6>
                        <p class="text-muted mb-1">${activity.description}</p>
                        <small class="text-muted">${activity.date}</small>
                    </div>
                    <span class="badge bg-${activity.statusColor}">${activity.status}</span>
                </div>
            </div>
        `).join('');
    }

    updateLastUpdate() {
        const now = new Date();
        const timeString = now.toLocaleTimeString('pt-BR', { 
            hour: '2-digit', 
            minute: '2-digit' 
        });
        document.getElementById('last-update').textContent = `Hoje, ${timeString}`;
    }

    // Métodos para carregar seções específicas
    loadDronesSection() {
        const dronesContainer = document.getElementById('drones-list');
        
        if (this.drones.length === 0) {
            dronesContainer.innerHTML = `
                <div class="col-12">
                    <div class="card card-tech text-center py-5">
                        <div class="card-body">
                            <i class="bi bi-drone fs-1 text-muted mb-3 d-block"></i>
                            <h5 class="text-muted mb-3">Nenhum drone cadastrado</h5>
                            <p class="text-muted mb-4">Cadastre seu primeiro drone para começar a usar nossos serviços</p>
                            <button class="btn btn-evo-primary" onclick="openDroneRegistrationModal()">
                                <i class="bi bi-plus-circle me-2"></i>Cadastrar Primeiro Drone
                            </button>
                        </div>
                    </div>
                </div>
            `;
        } else {
            dronesContainer.innerHTML = this.drones.map(drone => `
                <div class="col-md-6 col-lg-4">
                    <div class="card card-tech h-100">
                        <div class="card-body">
                            <div class="d-flex justify-content-between align-items-start mb-3">
                                <div class="icon-circle bg-evo-primary text-white">
                                    <i class="bi bi-drone"></i>
                                </div>
                                <span class="badge bg-success">Ativo</span>
                            </div>
                            <h5 class="card-title">${drone.brand} ${drone.model}</h5>
                            <p class="text-muted mb-2">Serial: ${drone.serial}</p>
                            <p class="text-muted mb-3">Cor: ${drone.color}</p>
                            <div class="d-flex gap-2">
                                <button class="btn btn-evo-accent btn-sm flex-fill" onclick="createMaintenanceForDrone('${drone.id}')">
                                    <i class="bi bi-tools me-1"></i>Manutenção
                                </button>
                                <button class="btn btn-outline-secondary btn-sm" onclick="editDrone('${drone.id}')">
                                    <i class="bi bi-pencil"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            `).join('');
        }
    }

    loadMaintenanceSection() {
        const container = document.getElementById('maintenance-form-container');
        
        if (this.drones.length === 0) {
            container.innerHTML = `
                <div class="card card-tech text-center py-5">
                    <div class="card-body">
                        <i class="bi bi-exclamation-triangle fs-1 text-warning mb-3 d-block"></i>
                        <h5 class="text-muted mb-3">Nenhum drone cadastrado</h5>
                        <p class="text-muted mb-4">Você precisa cadastrar pelo menos um drone antes de solicitar manutenção</p>
                        <button class="btn btn-evo-primary" onclick="showSection('drones')">
                            <i class="bi bi-drone me-2"></i>Cadastrar Drone
                        </button>
                    </div>
                </div>
            `;
        } else {
            container.innerHTML = `
                <div class="row">
                    <div class="col-lg-8">
                        <div class="card card-tech">
                            <div class="card-header bg-transparent">
                                <h5 class="mb-0">Nova Solicitação de Manutenção</h5>
                            </div>
                            <div class="card-body">
                                <form id="maintenance-form">
                                    <div class="mb-4">
                                        <label class="form-label fw-medium">Selecione o Drone</label>
                                        <select class="form-select" id="selected-drone" required>
                                            <option value="">Escolha um drone...</option>
                                            ${this.drones.map(drone => `
                                                <option value="${drone.id}">${drone.brand} ${drone.model} - ${drone.serial}</option>
                                            `).join('')}
                                        </select>
                                    </div>
                                    
                                    <div class="mb-4">
                                        <label class="form-label fw-medium">Tipo de Problema</label>
                                        <div class="row g-3">
                                            <div class="col-md-6">
                                                <div class="form-check">
                                                    <input class="form-check-input" type="checkbox" id="motor-problem">
                                                    <label class="form-check-label" for="motor-problem">Problema no Motor</label>
                                                </div>
                                            </div>
                                            <div class="col-md-6">
                                                <div class="form-check">
                                                    <input class="form-check-input" type="checkbox" id="battery-problem">
                                                    <label class="form-check-label" for="battery-problem">Problema na Bateria</label>
                                                </div>
                                            </div>
                                            <div class="col-md-6">
                                                <div class="form-check">
                                                    <input class="form-check-input" type="checkbox" id="camera-problem">
                                                    <label class="form-check-label" for="camera-problem">Problema na Câmera</label>
                                                </div>
                                            </div>
                                            <div class="col-md-6">
                                                <div class="form-check">
                                                    <input class="form-check-input" type="checkbox" id="control-problem">
                                                    <label class="form-check-label" for="control-problem">Problema no Controle</label>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    
                                    <div class="mb-4">
                                        <label class="form-label fw-medium">Descrição do Problema</label>
                                        <textarea class="form-control" id="problem-description" rows="4" placeholder="Descreva detalhadamente o problema encontrado..."></textarea>
                                    </div>
                                    
                                    <div class="mb-4">
                                        <label class="form-label fw-medium">Fotos do Drone (mínimo 5)</label>
                                        <input type="file" class="form-control" id="drone-photos" multiple accept="image/*" required>
                                        <div class="form-text">Adicione pelo menos 5 fotos do drone mostrando o problema</div>
                                    </div>
                                    
                                    <div class="d-flex gap-3">
                                        <button type="button" class="btn btn-evo-primary" onclick="submitMaintenanceRequest()">
                                            <i class="bi bi-send me-2"></i>Enviar Solicitação
                                        </button>
                                        <button type="button" class="btn btn-outline-secondary" onclick="showSection('dashboard')">
                                            Cancelar
                                        </button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="col-lg-4">
                        <div class="card card-tech">
                            <div class="card-header bg-transparent">
                                <h6 class="mb-0">Informações Importantes</h6>
                            </div>
                            <div class="card-body">
                                <div class="mb-3">
                                    <i class="bi bi-clock text-evo-primary me-2"></i>
                                    <small>Tempo médio de resposta: 24h</small>
                                </div>
                                <div class="mb-3">
                                    <i class="bi bi-shield-check text-success me-2"></i>
                                    <small>Garantia de qualidade</small>
                                </div>
                                <div class="mb-3">
                                    <i class="bi bi-truck text-info me-2"></i>
                                    <small>Coleta e entrega gratuita</small>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        }
    }

    loadOrdersSection() {
        const container = document.getElementById('orders-list');
        
        if (this.maintenanceOrders.length === 0) {
            container.innerHTML = `
                <div class="card card-tech text-center py-5">
                    <div class="card-body">
                        <i class="bi bi-clipboard-x fs-1 text-muted mb-3 d-block"></i>
                        <h5 class="text-muted mb-3">Nenhuma ordem encontrada</h5>
                        <p class="text-muted mb-4">Você ainda não possui solicitações de manutenção</p>
                        <button class="btn btn-evo-primary" onclick="showSection('maintenance')">
                            <i class="bi bi-plus-circle me-2"></i>Criar Primeira Solicitação
                        </button>
                    </div>
                </div>
            `;
        } else {
            container.innerHTML = `
                <div class="row g-4">
                    ${this.maintenanceOrders.map(order => `
                        <div class="col-lg-6">
                            <div class="card card-tech">
                                <div class="card-body">
                                    <div class="d-flex justify-content-between align-items-start mb-3">
                                        <h6 class="mb-0">${order.title}</h6>
                                        <span class="badge bg-${order.statusColor}">${order.status}</span>
                                    </div>
                                    <p class="text-muted mb-2">${order.description}</p>
                                    <small class="text-muted">Data: ${order.date}</small>
                                    <div class="mt-3">
                                        <button class="btn btn-outline-primary btn-sm me-2">
                                            <i class="bi bi-eye me-1"></i>Ver Detalhes
                                        </button>
                                        ${order.status === 'active' ? '<button class="btn btn-outline-danger btn-sm"><i class="bi bi-x-circle me-1"></i>Cancelar</button>' : ''}
                                    </div>
                                </div>
                            </div>
                        </div>
                    `).join('')}
                </div>
            `;
        }
    }

    loadHistorySection() {
        const container = document.getElementById('history-content');
        const completedOrders = this.maintenanceOrders.filter(order => order.status === 'completed');
        
        if (completedOrders.length === 0) {
            container.innerHTML = `
                <div class="card card-tech text-center py-5">
                    <div class="card-body">
                        <i class="bi bi-clock-history fs-1 text-muted mb-3 d-block"></i>
                        <h5 class="text-muted mb-3">Nenhum histórico encontrado</h5>
                        <p class="text-muted mb-0">Suas manutenções concluídas aparecerão aqui</p>
                    </div>
                </div>
            `;
        } else {
            container.innerHTML = `
                <div class="card card-tech">
                    <div class="card-body">
                        <div class="timeline">
                            ${completedOrders.map(order => `
                                <div class="timeline-item">
                                    <div class="d-flex justify-content-between align-items-start">
                                        <div>
                                            <h6 class="mb-1">${order.title}</h6>
                                            <p class="text-muted mb-1">${order.description}</p>
                                            <small class="text-muted">${order.date}</small>
                                        </div>
                                        <span class="badge bg-success">Concluído</span>
                                    </div>
                                </div>
                            `).join('')}
                        </div>
                    </div>
                </div>
            `;
        }
    }

    loadProfileSection() {
        const container = document.getElementById('profile-content');
        
        container.innerHTML = `
            <div class="row">
                <div class="col-lg-8">
                    <div class="card card-tech">
                        <div class="card-header bg-transparent">
                            <h5 class="mb-0">Informações Pessoais</h5>
                        </div>
                        <div class="card-body">
                            <form id="profile-form">
                                <div class="row g-3">
                                    <div class="col-md-6">
                                        <label class="form-label">Nome Completo</label>
                                        <input type="text" class="form-control" id="profile-name" value="João Silva">
                                    </div>
                                    <div class="col-md-6">
                                        <label class="form-label">Email</label>
                                        <input type="email" class="form-control" id="profile-email" value="joao@email.com">
                                    </div>
                                    <div class="col-md-6">
                                        <label class="form-label">Telefone</label>
                                        <input type="tel" class="form-control" id="profile-phone" value="(11) 99999-9999">
                                    </div>
                                    <div class="col-md-6">
                                        <label class="form-label">CPF</label>
                                        <input type="text" class="form-control" id="profile-cpf" value="000.000.000-00">
                                    </div>
                                    <div class="col-12">
                                        <label class="form-label">Endereço</label>
                                        <input type="text" class="form-control" id="profile-address" value="Rua das Flores, 123">
                                    </div>
                                </div>
                                <div class="mt-4">
                                    <button type="button" class="btn btn-evo-primary me-2">
                                        <i class="bi bi-check-circle me-2"></i>Salvar Alterações
                                    </button>
                                    <button type="button" class="btn btn-outline-secondary">
                                        Cancelar
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
                <div class="col-lg-4">
                    <div class="card card-tech">
                        <div class="card-header bg-transparent">
                            <h6 class="mb-0">Estatísticas</h6>
                        </div>
                        <div class="card-body">
                            <div class="mb-3">
                                <div class="d-flex justify-content-between">
                                    <span>Drones Cadastrados</span>
                                    <strong>${this.drones.length}</strong>
                                </div>
                            </div>
                            <div class="mb-3">
                                <div class="d-flex justify-content-between">
                                    <span>Manutenções Realizadas</span>
                                    <strong>${this.maintenanceOrders.filter(o => o.status === 'completed').length}</strong>
                                </div>
                            </div>
                            <div class="mb-3">
                                <div class="d-flex justify-content-between">
                                    <span>Cliente desde</span>
                                    <strong>Jan 2024</strong>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    logout() {
        if (confirm('Tem certeza que deseja sair?')) {
            window.location.href = '../auth/login.html';
        }
    }

    // Adicionar dados de exemplo para demonstração
    loadSampleData() {
        // Dados de exemplo de drones
        this.drones = [
            {
                id: '1',
                brand: 'DJI',
                model: 'Mavic Air 2',
                serial: 'DJI001234',
                color: 'Cinza'
            },
            {
                id: '2',
                brand: 'DJI',
                model: 'Mini 3 Pro',
                serial: 'DJI005678',
                color: 'Branco'
            }
        ];

        // Dados de exemplo de ordens de manutenção
        this.maintenanceOrders = [
            {
                id: '1',
                title: 'Manutenção DJI Mavic Air 2',
                description: 'Problema na câmera - imagens tremidas',
                status: 'active',
                statusColor: 'warning',
                date: '15/01/2024'
            },
            {
                id: '2',
                title: 'Reparo DJI Mini 3 Pro',
                description: 'Substituição de hélice danificada',
                status: 'completed',
                statusColor: 'success',
                date: '10/01/2024'
            }
        ];
    }
}

// Funções globais para compatibilidade
function showSection(sectionName) {
    window.dashboard.showSection(sectionName);
}

function logout() {
    window.dashboard.logout();
}

// Funções auxiliares para modais e ações
function openDroneRegistrationModal() {
    // Criar modal de cadastro de drone
    const modalHtml = `
        <div class="modal fade" id="droneRegistrationModal" tabindex="-1">
            <div class="modal-dialog modal-lg">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">Cadastrar Novo Drone</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <form id="drone-registration-form">
                            <div class="row g-3">
                                <div class="col-md-6">
                                    <label class="form-label">Marca *</label>
                                    <select class="form-select" id="drone-brand" required>
                                        <option value="">Selecione a marca...</option>
                                        <option value="DJI">DJI</option>
                                        <option value="Autel">Autel</option>
                                        <option value="Parrot">Parrot</option>
                                        <option value="Skydio">Skydio</option>
                                        <option value="Outro">Outro</option>
                                    </select>
                                </div>
                                <div class="col-md-6">
                                    <label class="form-label">Modelo *</label>
                                    <input type="text" class="form-control" id="drone-model" placeholder="Ex: Mavic Air 2" required>
                                </div>
                                <div class="col-md-6">
                                    <label class="form-label">Número de Série *</label>
                                    <input type="text" class="form-control" id="drone-serial" placeholder="Ex: DJI001234" required>
                                </div>
                                <div class="col-md-6">
                                    <label class="form-label">Cor</label>
                                    <input type="text" class="form-control" id="drone-color" placeholder="Ex: Cinza">
                                </div>
                                <div class="col-12">
                                    <label class="form-label">Observações</label>
                                    <textarea class="form-control" id="drone-notes" rows="3" placeholder="Informações adicionais sobre o drone..."></textarea>
                                </div>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
                        <button type="button" class="btn btn-evo-primary" onclick="saveDroneRegistration()">Cadastrar Drone</button>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Remover modal existente se houver
    const existingModal = document.getElementById('droneRegistrationModal');
    if (existingModal) {
        existingModal.remove();
    }
    
    // Adicionar modal ao DOM
    document.body.insertAdjacentHTML('beforeend', modalHtml);
    
    // Mostrar modal
    const modal = new bootstrap.Modal(document.getElementById('droneRegistrationModal'));
    modal.show();
}

function saveDroneRegistration() {
    const form = document.getElementById('drone-registration-form');
    const formData = new FormData(form);
    
    const droneData = {
        id: Date.now().toString(),
        brand: document.getElementById('drone-brand').value,
        model: document.getElementById('drone-model').value,
        serial: document.getElementById('drone-serial').value,
        color: document.getElementById('drone-color').value || 'Não informado',
        notes: document.getElementById('drone-notes').value
    };
    
    // Validar campos obrigatórios
    if (!droneData.brand || !droneData.model || !droneData.serial) {
        alert('Por favor, preencha todos os campos obrigatórios.');
        return;
    }
    
    // Adicionar drone à lista
    window.dashboard.drones.push(droneData);
    
    // Atualizar estatísticas
    window.dashboard.updateDashboardStats();
    
    // Fechar modal
    const modal = bootstrap.Modal.getInstance(document.getElementById('droneRegistrationModal'));
    modal.hide();
    
    // Mostrar mensagem de sucesso
    showSuccessMessage('Drone cadastrado com sucesso!');
    
    // Se estiver na seção de drones, recarregar
    if (window.dashboard.currentSection === 'drones') {
        window.dashboard.loadDronesSection();
    }
}

function createMaintenanceForDrone(droneId) {
    // Ir para seção de manutenção e pré-selecionar o drone
    showSection('maintenance');
    
    setTimeout(() => {
        const droneSelect = document.getElementById('selected-drone');
        if (droneSelect) {
            droneSelect.value = droneId;
        }
    }, 100);
}

function editDrone(droneId) {
    const drone = window.dashboard.drones.find(d => d.id === droneId);
    if (!drone) return;
    
    // Abrir modal de edição (similar ao de cadastro)
    openDroneRegistrationModal();
    
    // Preencher campos com dados existentes
    setTimeout(() => {
        document.getElementById('drone-brand').value = drone.brand;
        document.getElementById('drone-model').value = drone.model;
        document.getElementById('drone-serial').value = drone.serial;
        document.getElementById('drone-color').value = drone.color;
        document.getElementById('drone-notes').value = drone.notes || '';
        
        // Alterar título e botão
        document.querySelector('#droneRegistrationModal .modal-title').textContent = 'Editar Drone';
        document.querySelector('#droneRegistrationModal .btn-evo-primary').textContent = 'Salvar Alterações';
        document.querySelector('#droneRegistrationModal .btn-evo-primary').setAttribute('onclick', `updateDroneRegistration('${droneId}')`);
    }, 100);
}

function updateDroneRegistration(droneId) {
    const droneIndex = window.dashboard.drones.findIndex(d => d.id === droneId);
    if (droneIndex === -1) return;
    
    const updatedData = {
        id: droneId,
        brand: document.getElementById('drone-brand').value,
        model: document.getElementById('drone-model').value,
        serial: document.getElementById('drone-serial').value,
        color: document.getElementById('drone-color').value || 'Não informado',
        notes: document.getElementById('drone-notes').value
    };
    
    // Validar campos obrigatórios
    if (!updatedData.brand || !updatedData.model || !updatedData.serial) {
        alert('Por favor, preencha todos os campos obrigatórios.');
        return;
    }
    
    // Atualizar drone na lista
    window.dashboard.drones[droneIndex] = updatedData;
    
    // Fechar modal
    const modal = bootstrap.Modal.getInstance(document.getElementById('droneRegistrationModal'));
    modal.hide();
    
    // Mostrar mensagem de sucesso
    showSuccessMessage('Drone atualizado com sucesso!');
    
    // Recarregar seção de drones
    if (window.dashboard.currentSection === 'drones') {
        window.dashboard.loadDronesSection();
    }
}

function submitMaintenanceRequest() {
    const selectedDrone = document.getElementById('selected-drone').value;
    const description = document.getElementById('problem-description').value;
    const photos = document.getElementById('drone-photos').files;
    
    // Validações
    if (!selectedDrone) {
        alert('Por favor, selecione um drone.');
        return;
    }
    
    if (!description.trim()) {
        alert('Por favor, descreva o problema.');
        return;
    }
    
    if (photos.length < 5) {
        alert('Por favor, adicione pelo menos 5 fotos do drone.');
        return;
    }
    
    // Simular envio da solicitação
    const drone = window.dashboard.drones.find(d => d.id === selectedDrone);
    const newOrder = {
        id: Date.now().toString(),
        title: `Manutenção ${drone.brand} ${drone.model}`,
        description: description,
        status: 'active',
        statusColor: 'warning',
        date: new Date().toLocaleDateString('pt-BR')
    };
    
    window.dashboard.maintenanceOrders.push(newOrder);
    window.dashboard.updateDashboardStats();
    
    showSuccessMessage('Solicitação de manutenção enviada com sucesso!');
    
    // Limpar formulário
    document.getElementById('maintenance-form').reset();
    
    // Voltar para dashboard
    setTimeout(() => {
        showSection('dashboard');
    }, 2000);
}

function showSuccessMessage(message) {
    // Criar toast de sucesso
    const toastHtml = `
        <div class="toast-container position-fixed top-0 end-0 p-3">
            <div class="toast show" role="alert">
                <div class="toast-header bg-success text-white">
                    <i class="bi bi-check-circle me-2"></i>
                    <strong class="me-auto">Sucesso</strong>
                    <button type="button" class="btn-close btn-close-white" data-bs-dismiss="toast"></button>
                </div>
                <div class="toast-body">
                    ${message}
                </div>
            </div>
        </div>
    `;
    
    document.body.insertAdjacentHTML('beforeend', toastHtml);
    
    // Remover toast após 3 segundos
    setTimeout(() => {
        const toastContainer = document.querySelector('.toast-container');
        if (toastContainer) {
            toastContainer.remove();
        }
    }, 3000);
}

// Inicializar dashboard quando DOM estiver pronto
document.addEventListener('DOMContentLoaded', () => {
    window.dashboard = new ClientDashboard();
    
    // Carregar dados de exemplo para demonstração
    window.dashboard.loadSampleData();
    window.dashboard.updateDashboardStats();
});