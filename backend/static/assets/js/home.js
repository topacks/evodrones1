// Matrix Rain Effect
document.addEventListener('DOMContentLoaded', function() {
    const canvas = document.getElementById('matrix-canvas');
    if (canvas) {
        const ctx = canvas.getContext('2d');
        
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
        
        const matrix = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789@#$%^&*()*&^%+-/~{[|`]}";
        const matrixArray = matrix.split("");
        
        const fontSize = 10;
        const columns = canvas.width / fontSize;
        
        const drops = [];
        for(let x = 0; x < columns; x++) {
            drops[x] = 1;
        }
        
        function drawMatrix() {
            ctx.fillStyle = 'rgba(0, 0, 0, 0.04)';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            
            ctx.fillStyle = '#00ffff';
            ctx.font = fontSize + 'px monospace';
            
            for(let i = 0; i < drops.length; i++) {
                const text = matrixArray[Math.floor(Math.random() * matrixArray.length)];
                ctx.fillText(text, i * fontSize, drops[i] * fontSize);
                
                if(drops[i] * fontSize > canvas.height && Math.random() > 0.975) {
                    drops[i] = 0;
                }
                drops[i]++;
            }
        }
        
        setInterval(drawMatrix, 35);
    }
    
    // Floating Particles
    function createParticles() {
        const container = document.getElementById('particles-container');
        if (container) {
            for(let i = 0; i < 50; i++) {
                const particle = document.createElement('div');
                particle.className = 'particle';
                particle.style.left = Math.random() * 100 + '%';
                particle.style.animationDelay = Math.random() * 8 + 's';
                particle.style.animationDuration = (Math.random() * 3 + 5) + 's';
                container.appendChild(particle);
            }
        }
    }
    
    createParticles();
    
    // Three.js 3D Background
    const threeContainer = document.getElementById('three-container');
    if (threeContainer && typeof THREE !== 'undefined') {
        const scene = new THREE.Scene();
        const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        const renderer = new THREE.WebGLRenderer({ alpha: true });
        
        renderer.setSize(window.innerWidth, window.innerHeight);
        threeContainer.appendChild(renderer.domElement);
        
        // Create 3D objects
        const geometry = new THREE.BoxGeometry(1, 1, 1);
        const material = new THREE.MeshBasicMaterial({ 
            color: 0x00ffff, 
            wireframe: true,
            transparent: true,
            opacity: 0.3
        });
        
        const cubes = [];
        for(let i = 0; i < 20; i++) {
            const cube = new THREE.Mesh(geometry, material);
            cube.position.x = (Math.random() - 0.5) * 20;
            cube.position.y = (Math.random() - 0.5) * 20;
            cube.position.z = (Math.random() - 0.5) * 20;
            scene.add(cube);
            cubes.push(cube);
        }
        
        camera.position.z = 5;
        
        function animate() {
            requestAnimationFrame(animate);
            
            cubes.forEach(cube => {
                cube.rotation.x += 0.01;
                cube.rotation.y += 0.01;
            });
            
            renderer.render(scene, camera);
        }
        
        animate();
        
        // Handle window resize
        window.addEventListener('resize', function() {
            camera.aspect = window.innerWidth / window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.setSize(window.innerWidth, window.innerHeight);
        });
    }
    
    // Smooth Scrolling
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                window.scrollTo({
                    top: target.offsetTop - 100,
                    behavior: 'smooth'
                });
            }
        });
    });
});