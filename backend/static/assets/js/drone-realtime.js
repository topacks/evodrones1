/**
 * EVO DRONES - Visualização em Tempo Real
 * Script para controle e visualização de drones em tempo real
 */

class DroneRealtimeViewer {
  constructor(options = {}) {
    // Configurações padrão
    this.config = {
      mapContainerId: 'map-container',
      droneDataEndpoint: '/api/drones/realtime',
      refreshInterval: 2000, // 2 segundos
      animationSpeed: 300,   // 300ms para animações
      maxHistoryPoints: 50,  // Máximo de pontos no histórico
      ...options
    };
    
    // Estado da aplicação
    this.state = {
      drones: new Map(),      // Mapa de drones ativos
      selectedDroneId: null, // Drone selecionado
      isConnected: false,    // Status da conexão
      lastUpdate: null,      // Timestamp da última atualização
      mapInitialized: false, // Status de inicialização do mapa
      viewMode: '3d',        // Modo de visualização (2d/3d)
      filters: {             // Filtros ativos
        status: 'all',
        type: 'all'
      },
      history: {},           // Histórico de posições
      socket: null           // Conexão WebSocket
    };
    
    // Elementos da UI
    this.elements = {};
    
    // Inicialização
    this.init();
  }
  
  /**
   * Inicializa o visualizador
   */
  init() {
    // Capturar elementos da UI
    this.cacheElements();
    
    // Configurar eventos
    this.setupEventListeners();
    
    // Inicializar mapa
    this.initMap();
    
    // Iniciar conexão com servidor
    this.connectToServer();
    
    // Inicializar gráficos
    this.initCharts();
    
    // Inicializar filtros
    this.initFilters();
    
    console.log('DroneRealtimeViewer inicializado');
  }
  
  /**
   * Armazena referências aos elementos da UI
   */
  cacheElements() {
    // Elementos principais
    this.elements.mapContainer = document.getElementById(this.config.mapContainerId);
    this.elements.droneList = document.getElementById('drone-list');
    this.elements.statusIndicator = document.getElementById('connection-status');
    this.elements.lastUpdateTime = document.getElementById('last-update-time');
    this.elements.viewModeToggle = document.getElementById('view-mode-toggle');
    this.elements.filterControls = document.getElementById('filter-controls');
    this.elements.droneDetails = document.getElementById('drone-details');
    this.elements.menuToggle = document.getElementById('menu-toggle');
    this.elements.sidebar = document.getElementById('realtime-sidebar');
    this.elements.sidebarClose = document.getElementById('sidebar-close');
    
    // Gráficos
    this.elements.altitudeChart = document.getElementById('altitude-chart');
    this.elements.batteryChart = document.getElementById('battery-chart');
    this.elements.speedChart = document.getElementById('speed-chart');
    
    // Estatísticas
    this.elements.totalDrones = document.getElementById('total-drones');
    this.elements.activeDrones = document.getElementById('active-drones');
    this.elements.alertCount = document.getElementById('alert-count');
  }
  
  /**
   * Configura os listeners de eventos
   */
  setupEventListeners() {
    // Toggle do modo de visualização
    if (this.elements.viewModeToggle) {
      this.elements.viewModeToggle.addEventListener('click', () => this.toggleViewMode());
    }
    
    // Filtros
    if (this.elements.filterControls) {
      const filterInputs = this.elements.filterControls.querySelectorAll('select, input');
      filterInputs.forEach(input => {
        input.addEventListener('change', () => this.applyFilters());
      });
    }
    
    // Toggle do menu em dispositivos móveis
    if (this.elements.menuToggle) {
      this.elements.menuToggle.addEventListener('click', () => {
        this.elements.sidebar.classList.add('show');
      });
    }
    
    // Fechar sidebar
    if (this.elements.sidebarClose) {
      this.elements.sidebarClose.addEventListener('click', () => {
        this.elements.sidebar.classList.remove('show');
      });
    }
    
    // Eventos de redimensionamento
    window.addEventListener('resize', () => this.handleResize());
  }
  
  /**
   * Inicializa o mapa (usando Leaflet ou outra biblioteca)
   */
  initMap() {
    // Verificar se o elemento do mapa existe
    if (!this.elements.mapContainer) {
      console.error('Elemento do mapa não encontrado');
      return;
    }
    
    try {
      // Inicializar mapa com Leaflet (exemplo)
      this.map = L.map(this.elements.mapContainer.id).setView([-23.5505, -46.6333], 13);
      
      // Adicionar camada de mapa
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; OpenStreetMap contributors'
      }).addTo(this.map);
      
      // Camada para os drones
      this.dronesLayer = L.layerGroup().addTo(this.map);
      
      // Camada para as rotas
      this.routesLayer = L.layerGroup().addTo(this.map);
      
      this.state.mapInitialized = true;
      console.log('Mapa inicializado com sucesso');
    } catch (error) {
      console.error('Erro ao inicializar o mapa:', error);
      
      // Fallback para modo de visualização simples
      this.elements.mapContainer.innerHTML = `
        <div class="d-flex align-items-center justify-content-center h-100">
          <div class="text-center">
            <div class="mb-3">
              <i class="bi bi-map text-muted" style="font-size: 3rem;"></i>
            </div>
            <h4>Mapa não disponível</h4>
            <p class="text-muted">Não foi possível carregar o mapa. Verifique se a biblioteca Leaflet está disponível.</p>
          </div>
        </div>
      `;
    }
  }
  
  /**
   * Conecta ao servidor para receber atualizações em tempo real
   */
  connectToServer() {
    // Tentar usar WebSockets primeiro
    this.tryWebSocketConnection();
    
    // Fallback para polling se WebSocket falhar
    if (!this.state.isConnected) {
      this.startPolling();
    }
  }
  
  /**
   * Tenta estabelecer uma conexão WebSocket
   */
  tryWebSocketConnection() {
    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/api/drones/ws`;
      
      this.state.socket = new WebSocket(wsUrl);
      
      this.state.socket.onopen = () => {
        console.log('Conexão WebSocket estabelecida');
        this.state.isConnected = true;
        this.updateConnectionStatus(true);
      };
      
      this.state.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        this.processData(data);
      };
      
      this.state.socket.onerror = (error) => {
        console.error('Erro na conexão WebSocket:', error);
        this.state.isConnected = false;
        this.updateConnectionStatus(false);
        this.startPolling(); // Fallback para polling
      };
      
      this.state.socket.onclose = () => {
        console.log('Conexão WebSocket fechada');
        this.state.isConnected = false;
        this.updateConnectionStatus(false);
        
        // Tentar reconectar após 5 segundos
        setTimeout(() => this.tryWebSocketConnection(), 5000);
      };
    } catch (error) {
      console.error('Falha ao iniciar WebSocket:', error);
      this.startPolling(); // Fallback para polling
    }
  }
  
  /**
   * Inicia o polling como fallback
   */
  startPolling() {
    console.log('Iniciando polling como fallback');
    
    // Limpar intervalo existente se houver
    if (this.pollingInterval) {
      clearInterval(this.pollingInterval);
    }
    
    // Fazer a primeira requisição imediatamente
    this.fetchDroneData();
    
    // Configurar intervalo para atualizações
    this.pollingInterval = setInterval(() => {
      this.fetchDroneData();
    }, this.config.refreshInterval);
  }
  
  /**
   * Busca dados dos drones via API REST
   */
  async fetchDroneData() {
    try {
      const response = await fetch(this.config.droneDataEndpoint);
      
      if (!response.ok) {
        throw new Error(`Erro HTTP: ${response.status}`);
      }
      
      const data = await response.json();
      this.processData(data);
      this.updateConnectionStatus(true);
    } catch (error) {
      console.error('Erro ao buscar dados dos drones:', error);
      this.updateConnectionStatus(false);
    }
  }
  
  /**
   * Processa os dados recebidos do servidor
   */
  processData(data) {
    // Atualizar timestamp
    this.state.lastUpdate = new Date();
    this.updateLastUpdateTime();
    
    // Processar dados de cada drone
    if (Array.isArray(data.drones)) {
      // Criar conjunto de IDs atuais para detectar drones removidos
      const currentIds = new Set();
      
      // Atualizar drones existentes e adicionar novos
      data.drones.forEach(drone => {
        currentIds.add(drone.id);
        
        // Adicionar ao histórico
        this.addToHistory(drone);
        
        // Atualizar ou adicionar drone
        if (this.state.drones.has(drone.id)) {
          this.updateDrone(drone);
        } else {
          this.addDrone(drone);
        }
      });
      
      // Remover drones que não estão mais presentes
      this.state.drones.forEach((drone, id) => {
        if (!currentIds.has(id)) {
          this.removeDrone(id);
        }
      });
      
      // Atualizar estatísticas
      this.updateStats(data.stats || {});
      
      // Atualizar gráficos
      this.updateCharts();
      
      // Atualizar lista de drones
      this.renderDroneList();
      
      // Atualizar detalhes do drone selecionado
      if (this.state.selectedDroneId) {
        this.updateDroneDetails(this.state.selectedDroneId);
      }
    }
  }
  
  /**
   * Adiciona um novo drone ao mapa e à lista
   */
  addDrone(drone) {
    // Adicionar ao estado
    this.state.drones.set(drone.id, drone);
    
    // Adicionar ao mapa se inicializado
    if (this.state.mapInitialized) {
      const marker = this.createDroneMarker(drone);
      drone.marker = marker;
      this.dronesLayer.addLayer(marker);
    }
    
    console.log(`Drone adicionado: ${drone.id}`);
  }
  
  /**
   * Atualiza um drone existente
   */
  updateDrone(drone) {
    const existingDrone = this.state.drones.get(drone.id);
    
    // Atualizar dados
    Object.assign(existingDrone, drone);
    
    // Atualizar posição no mapa
    if (this.state.mapInitialized && existingDrone.marker) {
      const newLatLng = [drone.latitude, drone.longitude];
      
      // Animação suave da posição
      existingDrone.marker.setLatLng(newLatLng);
      
      // Atualizar popup se existir
      if (existingDrone.marker.getPopup()) {
        existingDrone.marker.setPopupContent(this.createPopupContent(drone));
      }
      
      // Atualizar ícone baseado no status
      existingDrone.marker.setIcon(this.getDroneIcon(drone.status));
    }
  }
  
  /**
   * Remove um drone do mapa e da lista
   */
  removeDrone(droneId) {
    const drone = this.state.drones.get(droneId);
    
    // Remover do mapa
    if (this.state.mapInitialized && drone && drone.marker) {
      this.dronesLayer.removeLayer(drone.marker);
    }
    
    // Remover do estado
    this.state.drones.delete(droneId);
    
    // Se era o drone selecionado, limpar seleção
    if (this.state.selectedDroneId === droneId) {
      this.state.selectedDroneId = null;
      this.clearDroneDetails();
    }
    
    console.log(`Drone removido: ${droneId}`);
  }
  
  /**
   * Cria um marcador para o drone no mapa
   */
  createDroneMarker(drone) {
    // Criar ícone baseado no status
    const icon = this.getDroneIcon(drone.status);
    
    // Criar marcador
    const marker = L.marker([drone.latitude, drone.longitude], { icon });
    
    // Adicionar popup
    marker.bindPopup(this.createPopupContent(drone));
    
    // Adicionar evento de clique
    marker.on('click', () => {
      this.selectDrone(drone.id);
    });
    
    return marker;
  }
  
  /**
   * Retorna o ícone apropriado baseado no status do drone
   */
  getDroneIcon(status) {
    // Cores baseadas no status
    const colors = {
      online: '#10B981',    // Verde
      offline: '#EF4444',   // Vermelho
      warning: '#FBBF24',   // Amarelo
      maintenance: '#6366F1' // Roxo
    };
    
    const color = colors[status] || colors.offline;
    
    // Criar ícone personalizado
    return L.divIcon({
      className: 'drone-icon',
      html: `
        <div class="drone-marker" style="color: ${color};">
          <i class="bi bi-drone-fill"></i>
          <div class="drone-pulse" style="background-color: ${color};"></div>
        </div>
      `,
      iconSize: [40, 40],
      iconAnchor: [20, 20]
    });
  }
  
  /**
   * Cria o conteúdo do popup para um drone
   */
  createPopupContent(drone) {
    return `
      <div class="drone-popup">
        <h5>${drone.name || `Drone ${drone.id.substring(0, 8)}`}</h5>
        <div class="drone-popup-info">
          <p><strong>Status:</strong> <span class="status-${drone.status}">${drone.status}</span></p>
          <p><strong>Bateria:</strong> ${drone.battery}%</p>
          <p><strong>Altitude:</strong> ${drone.altitude}m</p>
          <p><strong>Velocidade:</strong> ${drone.speed}km/h</p>
        </div>
        <button class="btn-realtime btn-primary btn-sm w-100" onclick="droneViewer.selectDrone('${drone.id}')">
          Ver detalhes
        </button>
      </div>
    `;
  }
  
  /**
   * Seleciona um drone para exibir detalhes
   */
  selectDrone(droneId) {
    this.state.selectedDroneId = droneId;
    this.updateDroneDetails(droneId);
    
    // Atualizar seleção visual na lista
    if (this.elements.droneList) {
      const items = this.elements.droneList.querySelectorAll('.drone-list-item');
      items.forEach(item => {
        if (item.dataset.droneId === droneId) {
          item.classList.add('active');
        } else {
          item.classList.remove('active');
        }
      });
    }
    
    // Centralizar no mapa
    if (this.state.mapInitialized) {
      const drone = this.state.drones.get(droneId);
      if (drone && drone.marker) {
        this.map.setView(drone.marker.getLatLng(), 16);
        drone.marker.openPopup();
      }
    }
  }
  
  /**
   * Atualiza os detalhes do drone selecionado
   */
  updateDroneDetails(droneId) {
    if (!this.elements.droneDetails) return;
    
    const drone = this.state.drones.get(droneId);
    if (!drone) {
      this.clearDroneDetails();
      return;
    }
    
    // Formatar dados para exibição
    const lastUpdate = new Date(drone.timestamp).toLocaleTimeString();
    const batteryClass = this.getBatteryStatusClass(drone.battery);
    const signalClass = this.getSignalStatusClass(drone.signal);
    
    // Atualizar conteúdo
    this.elements.droneDetails.innerHTML = `
      <div class="realtime-card">
        <div class="card-header">
          <h3 class="card-title">${drone.name || `Drone ${drone.id.substring(0, 8)}`}</h3>
          <span class="status-badge status-${drone.status}">
            <i class="bi bi-circle-fill"></i>
            ${drone.status}
          </span>
        </div>
        
        <div class="drone-details-grid">
          <div class="detail-item">
            <div class="detail-icon">
              <i class="bi bi-battery-full"></i>
            </div>
            <div class="detail-content">
              <h4>Bateria</h4>
              <div class="progress">
                <div class="progress-bar ${batteryClass}" style="width: ${drone.battery}%">
                  ${drone.battery}%
                </div>
              </div>
            </div>
          </div>
          
          <div class="detail-item">
            <div class="detail-icon">
              <i class="bi bi-reception-4"></i>
            </div>
            <div class="detail-content">
              <h4>Sinal</h4>
              <div class="progress">
                <div class="progress-bar ${signalClass}" style="width: ${drone.signal}%">
                  ${drone.signal}%
                </div>
              </div>
            </div>
          </div>
          
          <div class="detail-item">
            <div class="detail-icon">
              <i class="bi bi-speedometer"></i>
            </div>
            <div class="detail-content">
              <h4>Velocidade</h4>
              <p>${drone.speed} km/h</p>
            </div>
          </div>
          
          <div class="detail-item">
            <div class="detail-icon">
              <i class="bi bi-arrow-up-circle"></i>
            </div>
            <div class="detail-content">
              <h4>Altitude</h4>
              <p>${drone.altitude} metros</p>
            </div>
          </div>
          
          <div class="detail-item">
            <div class="detail-icon">
              <i class="bi bi-geo-alt"></i>
            </div>
            <div class="detail-content">
              <h4>Coordenadas</h4>
              <p>${drone.latitude.toFixed(6)}, ${drone.longitude.toFixed(6)}</p>
            </div>
          </div>
          
          <div class="detail-item">
            <div class="detail-icon">
              <i class="bi bi-clock"></i>
            </div>
            <div class="detail-content">
              <h4>Última atualização</h4>
              <p>${lastUpdate}</p>
            </div>
          </div>
        </div>
        
        <div class="drone-actions mt-4">
          <button class="btn-realtime btn-primary">
            <i class="bi bi-camera-video"></i> Visualizar câmera
          </button>
          <button class="btn-realtime btn-outline">
            <i class="bi bi-clipboard-data"></i> Relatório completo
          </button>
        </div>
      </div>
      
      <div class="realtime-card mt-4">
        <div class="card-header">
          <h3 class="card-title">Histórico de voo</h3>
        </div>
        <div id="flight-history-chart" style="height: 200px;"></div>
      </div>
    `;
    
    // Inicializar gráfico de histórico de voo
    this.renderFlightHistoryChart(droneId);
  }
  
  /**
   * Limpa os detalhes do drone
   */
  clearDroneDetails() {
    if (this.elements.droneDetails) {
      this.elements.droneDetails.innerHTML = `
        <div class="realtime-card">
          <div class="text-center py-5">
            <i class="bi bi-drone text-muted" style="font-size: 3rem;"></i>
            <h4 class="mt-3">Nenhum drone selecionado</h4>
            <p class="text-muted">Selecione um drone no mapa ou na lista para ver detalhes.</p>
          </div>
        </div>
      `;
    }
  }
  
  /**
   * Renderiza o gráfico de histórico de voo
   */
  renderFlightHistoryChart(droneId) {
    const historyContainer = document.getElementById('flight-history-chart');
    if (!historyContainer) return;
    
    const history = this.state.history[droneId] || [];
    if (history.length < 2) {
      historyContainer.innerHTML = `
        <div class="text-center py-4">
          <p class="text-muted">Dados insuficientes para exibir o histórico.</p>
        </div>
      `;
      return;
    }
    
    // Implementar gráfico com Chart.js ou outra biblioteca
    // Exemplo simplificado:
    try {
      const ctx = document.createElement('canvas');
      historyContainer.innerHTML = '';
      historyContainer.appendChild(ctx);
      
      const labels = history.map(point => {
        const date = new Date(point.timestamp);
        return date.toLocaleTimeString();
      });
      
      const altitudeData = history.map(point => point.altitude);
      const speedData = history.map(point => point.speed);
      
      new Chart(ctx, {
        type: 'line',
        data: {
          labels: labels,
          datasets: [
            {
              label: 'Altitude (m)',
              data: altitudeData,
              borderColor: '#30B4E8',
              backgroundColor: 'rgba(48, 180, 232, 0.1)',
              tension: 0.4,
              fill: true
            },
            {
              label: 'Velocidade (km/h)',
              data: speedData,
              borderColor: '#FF7A00',
              backgroundColor: 'rgba(255, 122, 0, 0.1)',
              tension: 0.4,
              fill: true
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          interaction: {
            intersect: false,
            mode: 'index'
          },
          scales: {
            y: {
              beginAtZero: true
            }
          }
        }
      });
    } catch (error) {
      console.error('Erro ao renderizar gráfico de histórico:', error);
      historyContainer.innerHTML = `
        <div class="text-center py-4">
          <p class="text-muted">Erro ao carregar o gráfico de histórico.</p>
        </div>
      `;
    }
  }
  
  /**
   * Adiciona um ponto ao histórico do drone
   */
  addToHistory(drone) {
    if (!this.state.history[drone.id]) {
      this.state.history[drone.id] = [];
    }
    
    // Adicionar ponto ao histórico
    this.state.history[drone.id].push({
      timestamp: drone.timestamp || new Date().toISOString(),
      latitude: drone.latitude,
      longitude: drone.longitude,
      altitude: drone.altitude,
      speed: drone.speed,
      battery: drone.battery
    });
    
    // Limitar tamanho do histórico
    if (this.state.history[drone.id].length > this.config.maxHistoryPoints) {
      this.state.history[drone.id].shift();
    }
  }
  
  /**
   * Renderiza a lista de drones
   */
  renderDroneList() {
    if (!this.elements.droneList) return;
    
    // Ordenar drones por status (online primeiro)
    const sortedDrones = Array.from(this.state.drones.values()).sort((a, b) => {
      if (a.status === 'online' && b.status !== 'online') return -1;
      if (a.status !== 'online' && b.status === 'online') return 1;
      return a.name?.localeCompare(b.name || '') || 0;
    });
    
    // Aplicar filtros
    const filteredDrones = this.applyFiltersToList(sortedDrones);
    
    // Renderizar lista
    this.elements.droneList.innerHTML = filteredDrones.length > 0 ? 
      filteredDrones.map(drone => this.renderDroneListItem(drone)).join('') :
      `<div class="text-center py-4"><p class="text-muted">Nenhum drone encontrado.</p></div>`;
    
    // Adicionar eventos de clique
    const items = this.elements.droneList.querySelectorAll('.drone-list-item');
    items.forEach(item => {
      item.addEventListener('click', () => {
        const droneId = item.dataset.droneId;
        this.selectDrone(droneId);
      });
    });
  }
  
  /**
   * Renderiza um item da lista de drones
   */
  renderDroneListItem(drone) {
    const isSelected = this.state.selectedDroneId === drone.id;
    const batteryClass = this.getBatteryStatusClass(drone.battery);
    
    return `
      <div class="drone-list-item ${isSelected ? 'active' : ''}" data-drone-id="${drone.id}">
        <div class="drone-list-icon">
          <i class="bi bi-drone-fill"></i>
          <span class="status-indicator status-${drone.status}"></span>
        </div>
        <div class="drone-list-info">
          <h4>${drone.name || `Drone ${drone.id.substring(0, 8)}`}</h4>
          <div class="drone-list-details">
            <span class="battery-indicator ${batteryClass}">
              <i class="bi bi-battery"></i> ${drone.battery}%
            </span>
            <span class="altitude-indicator">
              <i class="bi bi-arrow-up-circle"></i> ${drone.altitude}m
            </span>
          </div>
        </div>
      </div>
    `;
  }
  
  /**
   * Aplica filtros à lista de drones
   */
  applyFiltersToList(drones) {
    return drones.filter(drone => {
      // Filtrar por status
      if (this.state.filters.status !== 'all' && drone.status !== this.state.filters.status) {
        return false;
      }
      
      // Filtrar por tipo
      if (this.state.filters.type !== 'all' && drone.type !== this.state.filters.type) {
        return false;
      }
      
      return true;
    });
  }
  
  /**
   * Aplica os filtros selecionados
   */
  applyFilters() {
    if (!this.elements.filterControls) return;
    
    // Atualizar estado dos filtros
    const statusFilter = this.elements.filterControls.querySelector('#status-filter');
    const typeFilter = this.elements.filterControls.querySelector('#type-filter');
    
    if (statusFilter) this.state.filters.status = statusFilter.value;
    if (typeFilter) this.state.filters.type = typeFilter.value;
    
    // Atualizar visualização
    this.renderDroneList();
    this.updateMapMarkers();
  }
  
  /**
   * Atualiza os marcadores no mapa de acordo com os filtros
   */
  updateMapMarkers() {
    if (!this.state.mapInitialized) return;
    
    this.state.drones.forEach((drone, id) => {
      // Verificar se o drone deve ser exibido
      const shouldShow = 
        (this.state.filters.status === 'all' || drone.status === this.state.filters.status) &&
        (this.state.filters.type === 'all' || drone.type === this.state.filters.type);
      
      // Mostrar ou ocultar marcador
      if (drone.marker) {
        if (shouldShow) {
          this.dronesLayer.addLayer(drone.marker);
        } else {
          this.dronesLayer.removeLayer(drone.marker);
        }
      }
    });
  }
  
  /**
   * Alterna entre os modos de visualização 2D e 3D
   */
  toggleViewMode() {
    this.state.viewMode = this.state.viewMode === '2d' ? '3d' : '2d';
    
    // Atualizar botão
    if (this.elements.viewModeToggle) {
      this.elements.viewModeToggle.innerHTML = this.state.viewMode === '3d' ? 
        '<i class="bi bi-map"></i> Modo 2D' : 
        '<i class="bi bi-badge-3d"></i> Modo 3D';
    }
    
    // Implementar mudança de visualização
    if (this.state.mapInitialized) {
      // Exemplo: mudar para uma camada 3D ou 2D
      console.log(`Modo de visualização alterado para: ${this.state.viewMode}`);
    }
  }
  
  /**
   * Inicializa os gráficos
   */
  initCharts() {
    // Implementar inicialização de gráficos
    // Exemplo com Chart.js
    try {
      if (this.elements.altitudeChart) {
        this.charts = {};
        
        // Gráfico de altitude
        const altitudeCtx = this.elements.altitudeChart.getContext('2d');
        this.charts.altitude = new Chart(altitudeCtx, {
          type: 'line',
          data: {
            labels: [],
            datasets: [{
              label: 'Altitude (m)',
              data: [],
              borderColor: '#30B4E8',
              backgroundColor: 'rgba(48, 180, 232, 0.1)',
              tension: 0.4,
              fill: true
            }]
          },
          options: {
            responsive: true,
            maintainAspectRatio: false
          }
        });
        
        // Gráfico de bateria
        if (this.elements.batteryChart) {
          const batteryCtx = this.elements.batteryChart.getContext('2d');
          this.charts.battery = new Chart(batteryCtx, {
            type: 'doughnut',
            data: {
              labels: ['Restante', 'Utilizado'],
              datasets: [{
                data: [75, 25],
                backgroundColor: ['#30B4E8', '#E2E8F0']
              }]
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              cutout: '70%'
            }
          });
        }
        
        // Gráfico de velocidade
        if (this.elements.speedChart) {
          const speedCtx = this.elements.speedChart.getContext('2d');
          this.charts.speed = new Chart(speedCtx, {
            type: 'bar',
            data: {
              labels: [],
              datasets: [{
                label: 'Velocidade (km/h)',
                data: [],
                backgroundColor: '#FF7A00'
              }]
            },
            options: {
              responsive: true,
              maintainAspectRatio: false
            }
          });
        }
      }
    } catch (error) {
      console.error('Erro ao inicializar gráficos:', error);
    }
  }
  
  /**
   * Atualiza os gráficos com novos dados
   */
  updateCharts() {
    if (!this.charts) return;
    
    try {
      // Atualizar gráfico de altitude
      if (this.charts.altitude && this.state.selectedDroneId) {
        const history = this.state.history[this.state.selectedDroneId] || [];
        
        if (history.length > 0) {
          const labels = history.map(point => {
            const date = new Date(point.timestamp);
            return date.toLocaleTimeString();
          });
          
          const data = history.map(point => point.altitude);
          
          this.charts.altitude.data.labels = labels;
          this.charts.altitude.data.datasets[0].data = data;
          this.charts.altitude.update();
        }
      }
      
      // Atualizar gráfico de bateria
      if (this.charts.battery && this.state.selectedDroneId) {
        const drone = this.state.drones.get(this.state.selectedDroneId);
        
        if (drone) {
          this.charts.battery.data.datasets[0].data = [drone.battery, 100 - drone.battery];
          this.charts.battery.update();
        }
      }
      
      // Atualizar gráfico de velocidade
      if (this.charts.speed && this.state.selectedDroneId) {
        const history = this.state.history[this.state.selectedDroneId] || [];
        
        if (history.length > 0) {
          const labels = history.map(point => {
            const date = new Date(point.timestamp);
            return date.toLocaleTimeString();
          });
          
          const data = history.map(point => point.speed);
          
          this.charts.speed.data.labels = labels;
          this.charts.speed.data.datasets[0].data = data;
          this.charts.speed.update();
        }
      }
    } catch (error) {
      console.error('Erro ao atualizar gráficos:', error);
    }
  }
  
  /**
   * Inicializa os filtros
   */
  initFilters() {
    if (!this.elements.filterControls) return;
    
    // Exemplo de inicialização de filtros
    this.elements.filterControls.innerHTML = `
      <div class="filter-group">
        <label for="status-filter">Status</label>
        <select id="status-filter" class="form-select">
          <option value="all">Todos</option>
          <option value="online">Online</option>
          <option value="offline">Offline</option>
          <option value="warning">Alerta</option>
          <option value="maintenance">Manutenção</option>
        </select>
      </div>
      
      <div class="filter-group">
        <label for="type-filter">Tipo</label>
        <select id="type-filter" class="form-select">
          <option value="all">Todos</option>
          <option value="quadcopter">Quadricóptero</option>
          <option value="hexacopter">Hexacóptero</option>
          <option value="octocopter">Octocóptero</option>
        </select>
      </div>
    `;
  }
  
  /**
   * Atualiza o indicador de status da conexão
   */
  updateConnectionStatus(isConnected) {
    if (!this.elements.statusIndicator) return;
    
    if (isConnected) {
      this.elements.statusIndicator.innerHTML = `
        <i class="bi bi-wifi text-success"></i>
        <span>Conectado</span>
      `;
      this.elements.statusIndicator.classList.remove('status-offline');
      this.elements.statusIndicator.classList.add('status-online');
    } else {
      this.elements.statusIndicator.innerHTML = `
        <i class="bi bi-wifi-off text-danger"></i>
        <span>Desconectado</span>
      `;
      this.elements.statusIndicator.classList.remove('status-online');
      this.elements.statusIndicator.classList.add('status-offline');
    }
  }
  
  /**
   * Atualiza o horário da última atualização
   */
  updateLastUpdateTime() {
    if (!this.elements.lastUpdateTime || !this.state.lastUpdate) return;
    
    const timeString = this.state.lastUpdate.toLocaleTimeString();
    this.elements.lastUpdateTime.textContent = `Última atualização: ${timeString}`;
  }
  
  /**
   * Atualiza as estatísticas gerais
   */
  updateStats(stats) {
    // Total de drones
    if (this.elements.totalDrones) {
      this.elements.totalDrones.textContent = stats.total || this.state.drones.size;
    }
    
    // Drones ativos
    if (this.elements.activeDrones) {
      const activeCount = stats.active || 
        Array.from(this.state.drones.values()).filter(d => d.status === 'online').length;
      this.elements.activeDrones.textContent = activeCount;
    }
    
    // Alertas
    if (this.elements.alertCount) {
      const alertCount = stats.alerts || 
        Array.from(this.state.drones.values()).filter(d => d.status === 'warning').length;
      this.elements.alertCount.textContent = alertCount;
    }
  }
  
  /**
   * Retorna a classe CSS para o status da bateria
   */
  getBatteryStatusClass(batteryLevel) {
    if (batteryLevel <= 20) return 'bg-danger';
    if (batteryLevel <= 50) return 'bg-warning';
    return 'bg-success';
  }
  
  /**
   * Retorna a classe CSS para o status do sinal
   */
  getSignalStatusClass(signalLevel) {
    if (signalLevel <= 20) return 'bg-danger';
    if (signalLevel <= 60) return 'bg-warning';
    return 'bg-success';
  }
  
  /**
   * Manipula o evento de redimensionamento da janela
   */
  handleResize() {
    // Atualizar layout responsivo
    if (window.innerWidth >= 992) {
      // Desktop
      if (this.elements.sidebar) {
        this.elements.sidebar.classList.remove('show');
      }
    }
    
    // Atualizar tamanho do mapa
    if (this.state.mapInitialized && this.map) {
      this.map.invalidateSize();
    }
  }
}

// Inicializar quando o DOM estiver pronto
document.addEventListener('DOMContentLoaded', () => {
  // Criar instância do visualizador
  window.droneViewer = new DroneRealtimeViewer({
    mapContainerId: 'map-container',
    refreshInterval: 2000
  });
});

// Dados de exemplo para desenvolvimento
function generateMockData() {
  return {
    drones: [
      {
        id: 'drone-001',
        name: 'EVO Phantom Pro',
        status: 'online',
        latitude: -23.5505,
        longitude: -46.6333,
        altitude: 120,
        speed: 35,
        battery: 78,
        signal: 92,
        type: 'quadcopter',
        timestamp: new Date().toISOString()
      },
      {
        id: 'drone-002',
        name: 'EVO Explorer X',
        status: 'warning',
        latitude: -23.5605,
        longitude: -46.6433,
        altitude: 85,
        speed: 22,
        battery: 42,
        signal: 65,
        type: 'hexacopter',
        timestamp: new Date().toISOString()
      },
      {
        id: 'drone-003',
        name: 'EVO Surveyor',
        status: 'offline',
        latitude: -23.5705,
        longitude: -46.6233,
        altitude: 0,
        speed: 0,
        battery: 15,
        signal: 0,
        type: 'octocopter',
        timestamp: new Date().toISOString()
      }
    ],
    stats: {
      total: 3,
      active: 1,
      alerts: 1
    }
  };
}