const API_URL = 'http://localhost:8080';

function switchTab(tabName) {

    document.querySelectorAll('.tab').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelector(`button[onclick="switchTab('${tabName}')"]`).classList.add('active');
    

    const tabBorder = document.querySelector('.tab-border');
    tabBorder.style.left = tabName === 'login' ? '0' : '50%';
    

    document.querySelectorAll('.form').forEach(form => {
        form.classList.remove('active');
    });
    document.getElementById(`${tabName}Form`).classList.add('active');
    

    document.getElementById('authMessage').textContent = '';
}

async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();
        if (response.ok) {
            localStorage.setItem('token', data.token);
            window.location.href = 'dashboard.html';
        } else {
            showMessage(data.error || 'Login failed', 'error');
        }
    } catch (error) {
        showMessage('Network error', 'error');
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(`${API_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            showMessage('Registration successful! Please login.', 'success');
            switchTab('login');
        } else {
            const error = await response.json();
            showMessage(error.error || 'Registration failed', 'error');
        }
    } catch (error) {
        showMessage('Network error', 'error');
    }
}

function showMessage(text, type) {
    const msgDiv = document.getElementById('authMessage');
    msgDiv.textContent = text;
    msgDiv.className = `message ${type}`;
}

function logout() {
    localStorage.removeItem('token');
    window.location.href = 'index.html';
}


if (window.location.pathname.endsWith('dashboard.html')) {
    const token = localStorage.getItem('token');
    if (!token) window.location.href = 'index.html';
}