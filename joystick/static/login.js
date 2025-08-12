document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorDiv = document.getElementById('error');
    
    // Clear previous error
    errorDiv.classList.add('hidden');
    errorDiv.textContent = '';
    
    try {
        const response = await fetch('/api/v1/users/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });
        
        if (response.ok) {
            // Login successful, store token and redirect to dashboard
            const data = await response.json();
            localStorage.setItem('authToken', data.data.token);
            window.location.href = '/dashboard';
        } else {
            // Login failed, show error
            const errorData = await response.json();
            errorDiv.textContent = errorData.message || 'Login failed. Please try again.';
            errorDiv.classList.remove('hidden');
        }
    } catch (error) {
        errorDiv.textContent = 'Network error. Please try again.';
        errorDiv.classList.remove('hidden');
    }
});