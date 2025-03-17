// Three.js Animation Setup
const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
const renderer = new THREE.WebGLRenderer({ alpha: true });

renderer.setSize(window.innerWidth, window.innerHeight);
document.getElementById('animation-container').appendChild(renderer.domElement);

// Create stars
const starsGeometry = new THREE.BufferGeometry();
const starsMaterial = new THREE.PointsMaterial({
    color: 0xFFFFFF,
    size: 2,
    transparent: true
});

const starsVertices = [];
for (let i = 0; i < 1000; i++) {
    const x = (Math.random() - 0.5) * 2000;
    const y = (Math.random() - 0.5) * 2000;
    const z = -Math.random() * 2000;
    starsVertices.push(x, y, z);
}

starsGeometry.setAttribute('position', new THREE.Float32BufferAttribute(starsVertices, 3));
const stars = new THREE.Points(starsGeometry, starsMaterial);
scene.add(stars);

camera.position.z = 500;

// Animation function
function animate() {
    requestAnimationFrame(animate);

    // Rotate stars
    stars.rotation.y += 0.0005;
    stars.rotation.x += 0.0002;

    // Rotate nebulae
    scene.children.forEach(child => {
        if (child instanceof THREE.Points && child.userData.rotationSpeedX) {
            const data = child.userData;
            child.rotation.x += data.rotationSpeedX;
            child.rotation.y += data.rotationSpeedY;
            child.rotation.z += data.rotationSpeedZ;
        }
    });

    // Twinkle effect
    starsMaterial.opacity = 0.5 + Math.sin(Date.now() * 0.001) * 0.2;

    renderer.render(scene, camera);
}

// Handle window resize
window.addEventListener('resize', () => {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
});

// Create colored nebula clouds
const nebulaCount = 5;
const nebulaColors = [
    0x5035ff, // Purple blue
    0xff5555, // Red
    0x50cfff, // Light blue
    0xff8c50, // Orange
    0x4ca2ff  // Blue
];

// Create nebula particle system
for (let i = 0; i < nebulaCount; i++) {
    const particleCount = 200 + Math.floor(Math.random() * 300);
    const nebulaColor = nebulaColors[i % nebulaColors.length];
    
    const particles = new THREE.BufferGeometry();
    const positions = new Float32Array(particleCount * 3);
    const sizes = new Float32Array(particleCount);
    
    // Fill particles in a cloud-like formation
    for (let j = 0; j < particleCount; j++) {
        const x = (Math.random() - 0.5) * 2000;
        const y = (Math.random() - 0.5) * 2000;
        const z = (Math.random() - 0.5) * 2000;
        positions[j * 3] = x;
        positions[j * 3 + 1] = y;
        positions[j * 3 + 2] = z;
        
        // Vary particle sizes for nebula effect
        sizes[j] = 3 + Math.random() * 5; // Decreased size of particles
    }
    
    particles.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    particles.setAttribute('size', new THREE.BufferAttribute(sizes, 1));
    
    const nebulaMaterial = new THREE.PointsMaterial({
        color: nebulaColor,
        transparent: true,
        opacity: 0.2 + Math.random() * 0.3,
        size: 5, // Decreased base size of particles
        sizeAttenuation: true,
        blending: THREE.AdditiveBlending
    });
    
    const nebula = new THREE.Points(particles, nebulaMaterial);
    nebula.rotation.x = Math.random() * Math.PI;
    nebula.rotation.y = Math.random() * Math.PI;
    
    // Store rotation speed for animation
    nebula.userData = {
        rotationSpeedX: (Math.random() - 0.5) * 0.0002,
        rotationSpeedY: (Math.random() - 0.5) * 0.0002,
        rotationSpeedZ: (Math.random() - 0.5) * 0.0002
    };
    
    scene.add(nebula);
}

// Start animation
animate();