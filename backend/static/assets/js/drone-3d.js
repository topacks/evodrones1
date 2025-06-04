/**
 * EVO DRONES - Visualização 3D
 * 
 * Este arquivo contém a implementação da visualização 3D dos drones
 * utilizando Three.js para renderização e controle da cena 3D.
 */

class Drone3DViewer {
    constructor(options = {}) {
        // Configurações padrão
        this.config = {
            canvasId: 'scene-canvas',
            apiEndpoint: '/api/drones/realtime',
            updateInterval: 1000,
            maxHistoryPoints: 100,
            reconnectInterval: 3000,
            ...options
        };
        
        // Propriedades da cena
        this.scene = null;
        this.camera = null;
        this.renderer = null;
        this.controls = null;
        this.clock = null;
        
        // Dados dos drones
        this.drones = new Map();
        this.selectedDroneId = null;
        
        // Estado da conexão
        this.websocket = null;
        this.isConnected = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        
        // Modos de câmera
        this.cameraMode = 'follow'; // follow, free, firstPerson
        
        // Elementos da interface
        this.elements = {};
        
        // Inicializar
        this.init();
    }
    
    /**
     * Inicializa a visualização 3D
     */
    init() {
        this.cacheElements();
        this.setupEventListeners();
        this.initThreeJS();
        this.connectToServer();
        this.initUI();
        
        // Iniciar loop de renderização
        this.animate();
    }
    
    /**
     * Armazena referências aos elementos do DOM
     */
    cacheElements() {
        this.elements = {
            canvas: document.getElementById(this.config.canvasId),
            sidebar: document.getElementById('scene-sidebar'),
            sidebarToggle: document.getElementById('sidebar-toggle'),
            droneControls: document.querySelector('.drone-controls'),
            controlButtons: document.querySelectorAll('.control-btn'),
            cameraModeButtons: document.querySelectorAll('.camera-mode-btn'),
            altitudeBar: document.querySelector('.altitude-bar'),
            altitudeValue: document.querySelector('.altitude-value'),
            speedValue: document.querySelector('.speed-value'),
            batteryValue: document.querySelector('.battery-value'),
            batteryTime: document.querySelector('.battery-time'),
            compassNeedle: document.querySelector('.compass-needle'),
            coordinates: document.querySelector('.coordinates'),
            locationInfo: document.querySelector('.scene-info p')
        };
    }
    
    /**
     * Configura os event listeners
     */
    setupEventListeners() {
        // Responsividade
        window.addEventListener('resize', this.handleResize.bind(this));
        
        // Toggle da sidebar
        if (this.elements.sidebarToggle) {
            this.elements.sidebarToggle.addEventListener('click', this.toggleSidebar.bind(this));
        }
        
        // Botões de controle do drone
        this.elements.controlButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                this.handleDroneControl(e.currentTarget);
            });
        });
        
        // Botões de modo de câmera
        this.elements.cameraModeButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                this.setCameraMode(e.currentTarget);
            });
        });
    }
    
    /**
     * Inicializa o Three.js
     */
    initThreeJS() {
        // Criar cena
        this.scene = new THREE.Scene();
        this.scene.background = new THREE.Color(0x0c1824);
        this.scene.fog = new THREE.FogExp2(0x0c1824, 0.01);
        
        // Criar câmera
        this.camera = new THREE.PerspectiveCamera(
            75, 
            window.innerWidth / window.innerHeight, 
            0.1, 
            1000
        );
        this.camera.position.set(0, 10, 20);
        
        // Criar renderer
        this.renderer = new THREE.WebGLRenderer({ 
            canvas: this.elements.canvas, 
            antialias: true,
            alpha: true
        });
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        this.renderer.shadowMap.enabled = true;
        this.renderer.shadowMap.type = THREE.PCFSoftShadowMap;
        
        // Controles de órbita
        this.controls = new THREE.OrbitControls(this.camera, this.renderer.domElement);
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.05;
        this.controls.screenSpacePanning = false;
        this.controls.minDistance = 5;
        this.controls.maxDistance = 100;
        this.controls.maxPolarAngle = Math.PI / 2;
        
        // Relógio para animações
        this.clock = new THREE.Clock();
        
        // Adicionar elementos à cena
        this.addLights();
        this.addEnvironment();
        
        // Adicionar drone de exemplo
        this.addMockDrone();
    }
    
    /**
     * Adiciona iluminação à cena
     */
    addLights() {
        // Luz ambiente
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
        this.scene.add(ambientLight);
        
        // Luz direcional (sol)
        const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
        directionalLight.position.set(50, 200, 100);
        directionalLight.castShadow = true;
        directionalLight.shadow.mapSize.width = 1024;
        directionalLight.shadow.mapSize.height = 1024;
        directionalLight.shadow.camera.near = 10;
        directionalLight.shadow.camera.far = 400;
        directionalLight.shadow.camera.left = -100;
        directionalLight.shadow.camera.right = 100;
        directionalLight.shadow.camera.top = 100;
        directionalLight.shadow.camera.bottom = -100;
        this.scene.add(directionalLight);
        
        // Luz hemisférica
        const hemiLight = new THREE.HemisphereLight(0x3498db, 0x34495e, 0.3);
        this.scene.add(hemiLight);
    }
    
    /**
     * Adiciona elementos do ambiente à cena
     */
    addEnvironment() {
        // Grade para referência
        const gridHelper = new THREE.GridHelper(200, 50, 0x444444, 0x222222);
        this.scene.add(gridHelper);
        
        // Terreno
        this.terrain = this.createTerrain();
        this.scene.add(this.terrain);
        
        // Skybox
        this.addSkybox();
    }
    
    /**
     * Cria o terreno
     */
    createTerrain() {
        const geometry = new THREE.PlaneGeometry(200, 200, 64, 64);
        
        // Modificar vértices para criar elevações
        const vertices = geometry.attributes.position.array;
        for (let i = 0; i < vertices.length; i += 3) {
            const x = vertices[i];
            const z = vertices[i + 2];
            
            // Criar elevações usando funções de ruído (simplificado)
            vertices[i + 1] = Math.sin(x * 0.05) * Math.cos(z * 0.05) * 5 + 
                             Math.sin(x * 0.02) * Math.cos(z * 0.03) * 2;
        }
        
        geometry.computeVertexNormals();
        
        // Material com textura de grade
        const material = new THREE.MeshPhongMaterial({
            color: 0x3a7e8c,
            wireframe: true,
            transparent: true,
            opacity: 0.6,
            side: THREE.DoubleSide
        });
        
        const terrain = new THREE.Mesh(geometry, material);
        terrain.rotation.x = -Math.PI / 2;
        terrain.position.y = -2;
        terrain.receiveShadow = true;
        
        return terrain;
    }
    
    /**
     * Adiciona um skybox à cena
     */
    addSkybox() {
        const geometry = new THREE.SphereGeometry(500, 60, 40);
        geometry.scale(-1, 1, 1);
        
        const material = new THREE.MeshBasicMaterial({
            color: new THREE.Color(0x0c1824),
            side: THREE.BackSide
        });
        
        const skybox = new THREE.Mesh(geometry, material);
        this.scene.add(skybox);
    }
    
    /**
     * Adiciona um drone de exemplo para desenvolvimento
     */
    addMockDrone() {
        const mockDroneData = {
            id: 'drone-001',
            name: 'EVO Phantom Pro',
            model: 'Phantom Pro',
            status: 'online',
            battery: 78,
            signal: 92,
            altitude: 120,
            speed: 35,
            heading: 45,
            position: { lat: -23.550500, lng: -46.633300, alt: 120 },
            sensors: {
                temperature: 24,
                humidity: 65,
                pressure: 1013
            },
            camera: {
                resolution: '4K',
                fps: 30,
                stabilization: true
            },
            flightTime: 942 // segundos
        };
        
        this.processDroneData(mockDroneData);
        this.selectDrone('drone-001');
    }
    
    /**
     * Conecta ao servidor para receber dados em tempo real
     */
    connectToServer() {
        // Tentar conexão WebSocket
        try {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/api/ws/drones`;
            
            this.websocket = new WebSocket(wsUrl);
            
            this.websocket.onopen = () => {
                console.log('Conexão WebSocket estabelecida');
                this.isConnected = true;
                this.reconnectAttempts = 0;
                this.updateConnectionStatus(true);
            };
            
            this.websocket.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    if (Array.isArray(data)) {
                        data.forEach(drone => this.processDroneData(drone));
                    } else {
                        this.processDroneData(data);
                    }
                } catch (error) {
                    console.error('Erro ao processar mensagem:', error);
                }
            };
            
            this.websocket.onerror = (error) => {
                console.error('Erro na conexão WebSocket:', error);
                this.fallbackToPolling();
            };
            
            this.websocket.onclose = () => {
                console.log('Conexão WebSocket fechada');
                this.isConnected = false;
                this.updateConnectionStatus(false);
                
                // Tentar reconectar
                if (this.reconnectAttempts < this.maxReconnectAttempts) {
                    this.reconnectAttempts++;
                    setTimeout(() => {
                        console.log(`Tentativa de reconexão ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
                        this.connectToServer();
                    }, this.config.reconnectInterval);
                } else {
                    console.log('Número máximo de tentativas de reconexão atingido');
                    this.fallbackToPolling();
                }
            };
        } catch (error) {
            console.error('Erro ao inicializar WebSocket:', error);
            this.fallbackToPolling();
        }
    }
    
    /**
     * Fallback para polling via REST API
     */
    fallbackToPolling() {
        console.log('Usando polling como fallback');
        
        // Limpar intervalo existente, se houver
        if (this.pollingInterval) {
            clearInterval(this.pollingInterval);
        }
        
        // Iniciar polling
        this.pollingInterval = setInterval(() => {
            fetch(this.config.apiEndpoint)
                .then(response => response.json())
                .then(data => {
                    if (Array.isArray(data)) {
                        data.forEach(drone => this.processDroneData(drone));
                    } else {
                        this.processDroneData(data);
                    }
                    this.updateConnectionStatus(true);
                })
                .catch(error => {
                    console.error('Erro ao buscar dados:', error);
                    this.updateConnectionStatus(false);
                });
        }, this.config.updateInterval);
    }
    
    /**
     * Processa os dados recebidos do drone
     */
    processDroneData(data) {
        if (!data || !data.id) return;
        
        const droneId = data.id;
        const isNewDrone = !this.drones.has(droneId);
        
        if (isNewDrone) {
            // Criar novo drone
            const drone = {
                id: droneId,
                name: data.name || `Drone ${droneId}`,
                model: data.model || 'Desconhecido',
                status: data.status || 'offline',
                battery: data.battery || 0,
                signal: data.signal || 0,
                altitude: data.altitude || 0,
                speed: data.speed || 0,
                heading: data.heading || 0,
                position: data.position || { lat: 0, lng: 0, alt: 0 },
                sensors: data.sensors || {},
                camera: data.camera || {},
                flightTime: data.flightTime || 0,
                history: [],
                object3D: this.createDroneObject(data)
            };
            
            // Adicionar à cena
            if (drone.object3D) {
                this.scene.add(drone.object3D);
            }
            
            // Adicionar ao mapa de drones
            this.drones.set(droneId, drone);
            
            // Adicionar ao histórico
            this.addToHistory(drone);
        } else {
            // Atualizar drone existente
            const drone = this.drones.get(droneId);
            
            // Atualizar propriedades
            drone.status = data.status || drone.status;
            drone.battery = data.battery || drone.battery;
            drone.signal = data.signal || drone.signal;
            drone.altitude = data.altitude || drone.altitude;
            drone.speed = data.speed || drone.speed;
            drone.heading = data.heading || drone.heading;
            drone.flightTime = data.flightTime || drone.flightTime;
            
            if (data.position) {
                drone.position = data.position;
            }
            
            if (data.sensors) {
                drone.sensors = { ...drone.sensors, ...data.sensors };
            }
            
            if (data.camera) {
                drone.camera = { ...drone.camera, ...data.camera };
            }
            
            // Atualizar objeto 3D
            this.updateDroneObject(drone);
            
            // Adicionar ao histórico
            this.addToHistory(drone);
            
            // Atualizar UI se este for o drone selecionado
            if (this.selectedDroneId === droneId) {
                this.updateDroneDetails(drone);
            }
        }
        
        // Atualizar timestamp da última atualização
        this.lastUpdateTime = Date.now();
        this.updateLastUpdateTime();
    }
    
    /**
     * Cria um objeto 3D para representar o drone
     */
    createDroneObject(data) {
        const droneGroup = new THREE.Group();
        droneGroup.name = `drone-${data.id}`;
        
        // Corpo do drone
        const bodyGeometry = new THREE.BoxGeometry(2, 0.5, 2);
        const bodyMaterial = new THREE.MeshPhongMaterial({ color: 0x333333 });
        const body = new THREE.Mesh(bodyGeometry, bodyMaterial);
        body.castShadow = true;
        body.receiveShadow = true;
        droneGroup.add(body);
        
        // Braços do drone
        const armGeometry = new THREE.CylinderGeometry(0.1, 0.1, 3, 8);
        const armMaterial = new THREE.MeshPhongMaterial({ color: 0x666666 });
        
        const arm1 = new THREE.Mesh(armGeometry, armMaterial);
        arm1.position.set(1.5, 0, 1.5);
        arm1.rotation.z = Math.PI / 2;
        arm1.castShadow = true;
        droneGroup.add(arm1);
        
        const arm2 = new THREE.Mesh(armGeometry, armMaterial);
        arm2.position.set(-1.5, 0, 1.5);
        arm2.rotation.z = Math.PI / 2;
        arm2.castShadow = true;
        droneGroup.add(arm2);
        
        const arm3 = new THREE.Mesh(armGeometry, armMaterial);
        arm3.position.set(1.5, 0, -1.5);
        arm3.rotation.z = Math.PI / 2;
        arm3.castShadow = true;
        droneGroup.add(arm3);
        
        const arm4 = new THREE.Mesh(armGeometry, armMaterial);
        arm4.position.set(-1.5, 0, -1.5);
        arm4.rotation.z = Math.PI / 2;
        arm4.castShadow = true;
        droneGroup.add(arm4);
        
        // Hélices
        const propellerGeometry = new THREE.BoxGeometry(2, 0.1, 0.2);
        const propellerMaterial = new THREE.MeshPhongMaterial({ color: 0x3498db });
        
        const propeller1 = new THREE.Mesh(propellerGeometry, propellerMaterial);
        propeller1.position.set(1.5, 0.3, 1.5);
        propeller1.castShadow = true;
        droneGroup.add(propeller1);
        
        const propeller2 = new THREE.Mesh(propellerGeometry, propellerMaterial);
        propeller2.position.set(-1.5, 0.3, 1.5);
        propeller2.castShadow = true;
        droneGroup.add(propeller2);
        
        const propeller3 = new THREE.Mesh(propellerGeometry, propellerMaterial);
        propeller3.position.set(1.5, 0.3, -1.5);
        propeller3.castShadow = true;
        droneGroup.add(propeller3);
        
        const propeller4 = new THREE.Mesh(propellerGeometry, propellerMaterial);
        propeller4.position.set(-1.5, 0.3, -1.5);
        propeller4.castShadow = true;
        droneGroup.add(propeller4);
        
        // Câmera do drone
        const cameraGeometry = new THREE.CylinderGeometry(0.3, 0.3, 0.5, 16);
        const cameraMaterial = new THREE.MeshPhongMaterial({ color: 0x222222 });
        const droneCamera = new THREE.Mesh(cameraGeometry, cameraMaterial);
        droneCamera.position.set(0, -0.5, 0);
        droneCamera.rotation.x = Math.PI / 2;
        droneCamera.castShadow = true;
        droneGroup.add(droneCamera);
        
        // Lente da câmera
        const lensGeometry = new THREE.CircleGeometry(0.2, 16);
        const lensMaterial = new THREE.MeshPhongMaterial({ color: 0x88ccff });
        const lens = new THREE.Mesh(lensGeometry, lensMaterial);
        lens.position.set(0, -0.5, -0.26);
        lens.rotation.x = Math.PI / 2;
        droneGroup.add(lens);
        
        // Luzes do drone
        const lightGeometry = new THREE.SphereGeometry(0.1, 16, 16);
        
        const frontLight = new THREE.Mesh(
            lightGeometry,
            new THREE.MeshPhongMaterial({ color: 0x00ff00, emissive: 0x00ff00 })
        );
        frontLight.position.set(0, 0, 1.1);
        droneGroup.add(frontLight);
        
        const backLight = new THREE.Mesh(
            lightGeometry,
            new THREE.MeshPhongMaterial({ color: 0xff0000, emissive: 0xff0000 })
        );
        backLight.position.set(0, 0, -1.1);
        droneGroup.add(backLight);
        
        // Adicionar luz pontual para efeito visual
        const pointLight = new THREE.PointLight(0x3498db, 1, 10);
        pointLight.position.set(0, -1, 0);
        droneGroup.add(pointLight);
        
        // Posicionar drone
        const position = data.position || { lat: 0, lng: 0, alt: 0 };
        droneGroup.position.set(
            (Math.random() - 0.5) * 20, // Posição X aleatória para desenvolvimento
            position.alt / 10 || 5,    // Altura baseada na altitude
            (Math.random() - 0.5) * 20  // Posição Z aleatória para desenvolvimento
        );
        
        // Rotação baseada no heading
        droneGroup.rotation.y = (data.heading || 0) * Math.PI / 180;
        
        // Armazenar referências às hélices para animação
        droneGroup.userData = {
            propellers: [propeller1, propeller2, propeller3, propeller4],
            pointLight: pointLight,
            originalY: droneGroup.position.y
        };
        
        return droneGroup;
    }
    
    /**
     * Atualiza o objeto 3D do drone com base nos novos dados
     */
    updateDroneObject(drone) {
        if (!drone.object3D) return;
        
        // Atualizar posição
        if (drone.position) {
            // Em um cenário real, converteríamos lat/lng para coordenadas X/Z
            // Aqui estamos apenas atualizando a altura
            drone.object3D.position.y = drone.position.alt / 10 || 5;
            drone.object3D.userData.originalY = drone.object3D.position.y;
        }
        
        // Atualizar rotação baseada no heading
        if (drone.heading !== undefined) {
            drone.object3D.rotation.y = drone.heading * Math.PI / 180;
        }
        
        // Atualizar cor da luz com base no status
        if (drone.status) {
            const pointLight = drone.object3D.userData.pointLight;
            if (pointLight) {
                switch (drone.status) {
                    case 'online':
                        pointLight.color.set(0x2ecc71);
                        break;
                    case 'warning':
                        pointLight.color.set(0xf39c12);
                        break;
                    case 'critical':
                        pointLight.color.set(0xe74c3c);
                        break;
                    case 'offline':
                        pointLight.color.set(0x7f8c8d);
                        break;
                    default:
                        pointLight.color.set(0x3498db);
                }
            }
        }
    }
    
    /**
     * Anima os objetos 3D dos drones
     */
    animateDrones(deltaTime) {
        this.drones.forEach(drone => {
            if (!drone.object3D || !drone.object3D.userData) return;
            
            // Animar hélices
            const propellers = drone.object3D.userData.propellers;
            if (propellers) {
                propellers.forEach(propeller => {
                    propeller.rotation.y += (drone.speed / 10 + 0.3) * deltaTime * 5;
                });
            }
            
            // Animar flutuação
            if (drone.object3D.userData.originalY !== undefined) {
                drone.object3D.position.y = drone.object3D.userData.originalY + 
                    Math.sin(Date.now() * 0.001) * 0.2;
            }
        });
    }
    
    /**
     * Seleciona um drone para exibir detalhes
     */
    selectDrone(droneId) {
        if (!this.drones.has(droneId)) return;
        
        this.selectedDroneId = droneId;
        const drone = this.drones.get(droneId);
        
        // Atualizar detalhes na UI
        this.updateDroneDetails(drone);
        
        // Focar câmera no drone
        if (drone.object3D && this.cameraMode === 'follow') {
            this.focusCameraOnDrone(drone);
        }
    }
    
    /**
     * Foca a câmera em um drone específico
     */
    focusCameraOnDrone(drone) {
        if (!drone.object3D) return;
        
        // Animação suave com GSAP
        if (window.gsap) {
            gsap.to(this.controls.target, {
                x: drone.object3D.position.x,
                y: drone.object3D.position.y,
                z: drone.object3D.position.z,
                duration: 1,
                ease: 'power2.out'
            });
            
            const cameraPosition = {
                x: this.camera.position.x,
                y: this.camera.position.y,
                z: this.camera.position.z
            };
            
            const targetPosition = {
                x: drone.object3D.position.x,
                y: drone.object3D.position.y + 5,
                z: drone.object3D.position.z + 10
            };
            
            gsap.to(cameraPosition, {
                x: targetPosition.x,
                y: targetPosition.y,
                z: targetPosition.z,
                duration: 1,
                ease: 'power2.out',
                onUpdate: () => {
                    this.camera.position.set(cameraPosition.x, cameraPosition.y, cameraPosition.z);
                }
            });
        } else {
            // Fallback sem animação
            this.controls.target.copy(drone.object3D.position);
            this.camera.position.set(
                drone.object3D.position.x,
                drone.object3D.position.y + 5,
                drone.object3D.position.z + 10
            );
        }
    }
    
    /**
     * Atualiza os detalhes do drone na interface
     */
    updateDroneDetails(drone) {
        if (!drone) return;
        
        // Atualizar valores na UI
        if (this.elements.altitudeBar) {
            this.elements.altitudeBar.style.height = `${Math.min(100, drone.altitude / 2)}%`;
        }
        
        if (this.elements.altitudeValue) {
            this.elements.altitudeValue.textContent = `${drone.altitude}m`;
        }
        
        if (this.elements.speedValue) {
            this.elements.speedValue.textContent = drone.speed;
        }
        
        if (this.elements.batteryValue) {
            this.elements.batteryValue.textContent = `${drone.battery}%`;
        }
        
        if (this.elements.batteryTime) {
            const remainingMinutes = Math.floor((drone.flightTime / 60) * (drone.battery / 100));
            this.elements.batteryTime.textContent = `${remainingMinutes}min restantes`;
        }
        
        if (this.elements.compassNeedle) {
            this.elements.compassNeedle.style.transform = `translateX(-50%) rotate(${drone.heading}deg)`;
        }
        
        if (this.elements.coordinates && drone.position) {
            this.elements.coordinates.textContent = `${drone.position.lat.toFixed(6)}, ${drone.position.lng.toFixed(6)}`;
        }
        
        if (this.elements.locationInfo && drone.position) {
            // Em um cenário real, faríamos geocodificação reversa
            this.elements.locationInfo.textContent = 'São Paulo, SP - Brasil';
        }
        
        // Atualizar outros elementos da UI conforme necessário
        // ...
    }
    
    /**
     * Adiciona um ponto ao histórico do drone
     */
    addToHistory(drone) {
        if (!drone) return;
        
        const historyPoint = {
            timestamp: Date.now(),
            position: { ...drone.position },
            altitude: drone.altitude,
            speed: drone.speed,
            battery: drone.battery,
            heading: drone.heading
        };
        
        drone.history.push(historyPoint);
        
        // Limitar tamanho do histórico
        if (drone.history.length > this.config.maxHistoryPoints) {
            drone.history.shift();
        }
    }
    
    /**
     * Define o modo de câmera
     */
    setCameraMode(buttonElement) {
        // Remover classe ativa de todos os botões
        this.elements.cameraModeButtons.forEach(btn => {
            btn.classList.remove('active');
        });
        
        // Adicionar classe ativa ao botão clicado
        buttonElement.classList.add('active');
        
        // Obter modo de câmera do texto do botão
        const mode = buttonElement.textContent.toLowerCase();
        
        switch (mode) {
            case 'seguir':
                this.cameraMode = 'follow';
                this.controls.enabled = true;
                
                // Se houver um drone selecionado, focar nele
                if (this.selectedDroneId && this.drones.has(this.selectedDroneId)) {
                    this.focusCameraOnDrone(this.drones.get(this.selectedDroneId));
                }
                break;
                
            case 'livre':
                this.cameraMode = 'free';
                this.controls.enabled = true;
                break;
                
            case 'primeira pessoa':
                this.cameraMode = 'firstPerson';
                this.controls.enabled = false;
                
                // Se houver um drone selecionado, posicionar câmera nele
                if (this.selectedDroneId && this.drones.has(this.selectedDroneId)) {
                    const drone = this.drones.get(this.selectedDroneId);
                    if (drone.object3D) {
                        this.camera.position.copy(drone.object3D.position);
                        this.camera.position.y += 0.5;
                        
                        // Calcular direção com base no heading
                        const angle = drone.heading * Math.PI / 180;
                        const lookAtPoint = new THREE.Vector3(
                            drone.object3D.position.x + Math.sin(angle) * 10,
                            drone.object3D.position.y,
                            drone.object3D.position.z + Math.cos(angle) * 10
                        );
                        
                        this.camera.lookAt(lookAtPoint);
                    }
                }
                break;
        }
    }
    
    /**
     * Manipula os controles do drone
     */
    handleDroneControl(buttonElement) {
        // Remover classe ativa de todos os botões
        this.elements.controlButtons.forEach(btn => {
            btn.classList.remove('active');
        });
        
        // Adicionar classe ativa ao botão clicado
        buttonElement.classList.add('active');
        
        // Obter ação do ícone do botão
        const icon = buttonElement.querySelector('i');
        if (!icon) return;
        
        const action = icon.className;
        
        // Implementar ações do drone
        if (this.selectedDroneId && this.drones.has(this.selectedDroneId)) {
            const drone = this.drones.get(this.selectedDroneId);
            
            if (action.includes('arrow-up-circle')) {
                console.log('Comando: Decolar');
                // Implementar lógica de decolagem
            } else if (action.includes('arrow-down-circle')) {
                console.log('Comando: Pousar');
                // Implementar lógica de pouso
            } else if (action.includes('record-circle')) {
                console.log('Comando: Seguir');
                // Implementar lógica de seguimento
            } else if (action.includes('arrow-repeat')) {
                console.log('Comando: Orbitar');
                // Implementar lógica de órbita
            } else if (action.includes('house')) {
                console.log('Comando: Retornar');
                // Implementar lógica de retorno
            }
        }
    }
    
    /**
     * Toggle da sidebar
     */
    toggleSidebar() {
        if (!this.elements.sidebar) return;
        
        this.elements.sidebar.classList.toggle('collapsed');
        
        // Alternar ícone
        const icon = this.elements.sidebarToggle.querySelector('i');
        if (icon) {
            if (this.elements.sidebar.classList.contains('collapsed')) {
                icon.classList.remove('bi-chevron-left');
                icon.classList.add('bi-chevron-right');
            } else {
                icon.classList.remove('bi-chevron-right');
                icon.classList.add('bi-chevron-left');
            }
        }
    }
    
    /**
     * Inicializa elementos da UI
     */
    initUI() {
        this.updateConnectionStatus(false);
        this.updateLastUpdateTime();
    }
    
    /**
     * Atualiza o status da conexão na UI
     */
    updateConnectionStatus(isConnected) {
        // Implementar atualização do status na UI
        this.isConnected = isConnected;
    }
    
    /**
     * Atualiza o timestamp da última atualização na UI
     */
    updateLastUpdateTime() {
        if (!this.lastUpdateTime) return;
        
        const now = Date.now();
        const diff = now - this.lastUpdateTime;
        
        let timeText = '';
        
        if (diff < 1000) {
            timeText = 'agora mesmo';
        } else if (diff < 60000) {
            timeText = `${Math.floor(diff / 1000)} segundos atrás`;
        } else if (diff < 3600000) {
            timeText = `${Math.floor(diff / 60000)} minutos atrás`;
        } else {
            timeText = `${Math.floor(diff / 3600000)} horas atrás`;
        }
        
        // Atualizar elemento na UI
        // ...
    }
    
    /**
     * Manipula o redimensionamento da janela
     */
    handleResize() {
        if (!this.camera || !this.renderer) return;
        
        this.camera.aspect = window.innerWidth / window.innerHeight;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(window.innerWidth, window.innerHeight);
    }
    
    /**
     * Loop de animação principal
     */
    animate() {
        requestAnimationFrame(this.animate.bind(this));
        
        const delta = this.clock.getDelta();
        
        // Atualizar controles
        if (this.controls && this.controls.enabled) {
            this.controls.update();
        }
        
        // Animar drones
        this.animateDrones(delta);
        
        // Atualizar câmera em primeira pessoa
        if (this.cameraMode === 'firstPerson' && this.selectedDroneId) {
            const drone = this.drones.get(this.selectedDroneId);
            if (drone && drone.object3D) {
                this.camera.position.copy(drone.object3D.position);
                this.camera.position.y += 0.5;
                
                // Calcular direção com base no heading
                const angle = drone.heading * Math.PI / 180;
                const lookAtPoint = new THREE.Vector3(
                    drone.object3D.position.x + Math.sin(angle) * 10,
                    drone.object3D.position.y,
                    drone.object3D.position.z + Math.cos(angle) * 10
                );
                
                this.camera.lookAt(lookAtPoint);
            }
        }
        
        // Renderizar cena
        if (this.renderer && this.scene && this.camera) {
            this.renderer.render(this.scene, this.camera);
        }
    }
}

// Inicializar visualização 3D quando o DOM estiver carregado
document.addEventListener('DOMContentLoaded', function() {
    const viewer = new Drone3DViewer({
        canvasId: 'scene-canvas',
        apiEndpoint: '/api/drones/realtime',
        updateInterval: 1000
    });
});