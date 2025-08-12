// Check if user is authenticated on page load
function checkAuth() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

// Helper function to make authenticated API calls
async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem('authToken');
    const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        ...options.headers
    };
    
    const response = await fetch(url, {
        ...options,
        headers
    });
    
    // If 401, redirect to login
    if (response.status === 401) {
        localStorage.removeItem('authToken');
        window.location.href = '/login';
        return null;
    }
    
    return response;
}

let pollingInterval;

document.addEventListener('DOMContentLoaded', function() {
    // Check auth on page load
    if (!checkAuth()) return;
    
    // Load services data
    loadServices();
    
    // Event listeners
    setupEventListeners();
    
    // Start polling for updates every 5 seconds
    pollingInterval = setInterval(loadServices, 5000);
});

function setupEventListeners() {
    // Logout functionality
    document.getElementById('logoutBtn').addEventListener('click', function() {
        clearInterval(pollingInterval);
        localStorage.removeItem('authToken');
        window.location.href = '/login';
    });
    
    // Create service button
    document.getElementById('createServiceBtn').addEventListener('click', function() {
        document.getElementById('createServiceModal').classList.remove('hidden');
    });
    
    // Modal close buttons
    document.getElementById('closeEnvModal').addEventListener('click', function(e) {
        e.preventDefault();
        closeModal('envVarsModal');
    });
    
    document.getElementById('closeCreateModal').addEventListener('click', function(e) {
        e.preventDefault();
        closeModal('createServiceModal');
    });
    
    // Close modals on outside click
    ['envVarsModal', 'createServiceModal'].forEach(modalId => {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.addEventListener('click', function(e) {
                if (e.target === this) {
                    closeModal(modalId);
                }
            });
        }
    });
    
    // Close modals on Escape key
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            closeModal('envVarsModal');
            closeModal('createServiceModal');
        }
    });
    
    // Add port mapping button
    document.getElementById('addPortBtn').addEventListener('click', addPortMapping);
    
    // Create service form submission
    document.getElementById('createServiceForm').addEventListener('submit', handleCreateService);
}

async function loadServices() {
    try {
        const response = await authenticatedFetch('/api/v1/services');
        if (!response) return; // 401 handled by authenticatedFetch
        
        const data = await response.json();
        displayServices(data.data.services);
    } catch (error) {
        showToast('Failed to load services', 'error');
    }
}

function displayServices(services) {
    const loadingSpinner = document.getElementById('loadingSpinner');
    const servicesTable = document.getElementById('servicesTable');
    const noServices = document.getElementById('noServices');
    const tableBody = document.getElementById('servicesTableBody');
    
    loadingSpinner.classList.add('hidden');
    
    if (!services || services.length === 0) {
        noServices.classList.remove('hidden');
        return;
    }
    
    // Clear existing rows
    tableBody.innerHTML = '';
    
    services.forEach(service => {
        const row = createServiceRow(service);
        tableBody.appendChild(row);
    });
    
    servicesTable.classList.remove('hidden');
}

function createServiceRow(service) {
    const row = document.createElement('tr');
    
    // Format created date
    const createdDate = new Date(service.createdAt).toLocaleDateString();
    
    // Format port mapping with clear labels
    const portMapping = service.portMapping.map(pm => `${pm.hostPort} → ${pm.containerPort}`).join(', ');
    
    // Determine status
    let status = 'unknown';
    let statusText = 'NA';
    if (service.latestDeployment) {
        status = service.latestDeployment.status;
        statusText = status;
    }
    
    row.innerHTML = `
        <td>
            <div class="service-name">
                <span>${service.name}</span>
                <button class="copy-btn" onclick="copyToClipboard('${service.uuid}', 'Service UUID copied!')">📋</button>
            </div>
        </td>
        <td>${createdDate}</td>
        <td>
            <button class="env-vars-btn" onclick="showEnvVars('${service.name}', '${btoa(JSON.stringify(service.envVars))}')">
                view
            </button>
        </td>
        <td>${portMapping}</td>
        <td>
            <span class="status-badge status-${status.toLowerCase()}">${statusText}</span>
        </td>
        <td class="actions-cell">
            <button class="deploy-btn" onclick="deployService('${service.uuid}', '${service.name}')">🚀</button>
            <button class="delete-btn" onclick="deleteService('${service.uuid}', '${service.name}')">🗑️</button>
        </td>
    `;
    
    return row;
}

function copyToClipboard(text, message) {
    navigator.clipboard.writeText(text).then(() => {
        showToast(message);
    }).catch(() => {
        showToast('Failed to copy to clipboard', 'error');
    });
}

function showEnvVars(serviceName, encodedEnvVars) {
    const envVars = JSON.parse(atob(encodedEnvVars));
    const modal = document.getElementById('envVarsModal');
    const content = document.getElementById('envVarsContent');
    
    // Update modal title
    document.querySelector('.modal-header h3').textContent = `Environment Variables - ${serviceName}`;
    
    // Format and display env vars
    const formattedEnvVars = JSON.stringify(envVars, null, 2);
    content.textContent = formattedEnvVars;
    
    modal.classList.remove('hidden');
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('hidden');
    }
}

async function deleteService(serviceUuid, serviceName) {
    if (!confirm(`Are you sure you want to delete service "${serviceName}"?`)) {
        return;
    }
    
    try {
        const response = await authenticatedFetch(`/api/v1/services/${serviceUuid}`, {
            method: 'DELETE'
        });
        
        if (!response) return; // 401 handled by authenticatedFetch
        
        if (response.ok) {
            showToast('Service deleted successfully');
            loadServices(); // Reload the services list
        } else {
            const error = await response.json();
            showToast(error.message || 'Failed to delete service', 'error');
        }
    } catch (error) {
        showToast('Failed to delete service', 'error');
    }
}

async function deployService(serviceUuid, serviceName) {
    try {
        const response = await authenticatedFetch(`/api/v1/services/${serviceUuid}/deployments`, {
            method: 'POST'
        });
        
        if (!response) return; // 401 handled by authenticatedFetch
        
        if (response.ok) {
            showToast(`Deployment initiated for ${serviceName}`);
            loadServices(); // Refresh the services list
        } else {
            const error = await response.json();
            showToast(error.message || 'Failed to deploy service', 'error');
        }
    } catch (error) {
        showToast('Failed to deploy service', 'error');
    }
}

function addPortMapping() {
    const container = document.getElementById('portMappingContainer');
    const newRow = document.createElement('div');
    newRow.className = 'port-mapping-row';
    newRow.innerHTML = `
        <div>
            <div class="port-label">Host Port</div>
            <input type="number" class="host-port" placeholder="8080" required>
        </div>
        <div>
            <div class="port-label">Container Port</div>
            <input type="number" class="container-port" placeholder="8080" required>
        </div>
        <button type="button" class="remove-port-btn" onclick="removePortMapping(this)">×</button>
    `;
    container.appendChild(newRow);
}

function removePortMapping(button) {
    const container = document.getElementById('portMappingContainer');
    if (container.children.length > 1) {
        button.parentElement.remove();
    }
}

async function handleCreateService(e) {
    e.preventDefault();
    
    const form = e.target;
    const formData = new FormData(form);
    
    // Collect port mappings
    const portMappings = [];
    const hostPorts = document.querySelectorAll('.host-port');
    const containerPorts = document.querySelectorAll('.container-port');
    
    for (let i = 0; i < hostPorts.length; i++) {
        if (hostPorts[i].value && containerPorts[i].value) {
            portMappings.push({
                hostPort: parseInt(hostPorts[i].value),
                containerPort: parseInt(containerPorts[i].value)
            });
        }
    }
    
    // Parse environment variables
    let envVars = {};
    try {
        envVars = JSON.parse(formData.get('envVarsInput') || '{}');
    } catch (error) {
        showToast('Invalid JSON in environment variables', 'error');
        return;
    }
    
    const serviceData = {
        name: formData.get('serviceName'),
        repositoryUrl: formData.get('repositoryUrl'),
        branch: formData.get('branch'),
        envVars: envVars,
        portMapping: portMappings
    };
    
    try {
        const response = await authenticatedFetch('/api/v1/services', {
            method: 'POST',
            body: JSON.stringify(serviceData)
        });
        
        if (!response) return; // 401 handled by authenticatedFetch
        
        if (response.ok) {
            showToast('Service created successfully');
            closeModal('createServiceModal');
            form.reset();
            // Reset port mappings to default
            const container = document.getElementById('portMappingContainer');
            container.innerHTML = `
                <div class="port-mapping-row">
                    <div>
                        <div class="port-label">Host Port</div>
                        <input type="number" class="host-port" placeholder="8080" required>
                    </div>
                    <div>
                        <div class="port-label">Container Port</div>
                        <input type="number" class="container-port" placeholder="8080" required>
                    </div>
                    <button type="button" class="remove-port-btn" onclick="removePortMapping(this)">×</button>
                </div>
            `;
            document.getElementById('envVarsInput').value = '{}';
            loadServices(); // Refresh the services list
        } else {
            const error = await response.json();
            showToast(error.message || 'Failed to create service', 'error');
        }
    } catch (error) {
        showToast('Failed to create service', 'error');
    }
}

function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type}`;
    toast.classList.remove('hidden');
    
    setTimeout(() => {
        toast.classList.add('hidden');
    }, 3000);
}