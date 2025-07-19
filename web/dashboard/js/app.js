// API endpoints
const API_BASE = 'http://localhost:8080/api/v1';
const EXPERIMENTS_ENDPOINT = `${API_BASE}/experiments`;
const TARGETS_ENDPOINT = `${API_BASE}/targets`;

// DOM elements
const dashboardLink = document.getElementById('dashboard-link');
const experimentsLink = document.getElementById('experiments-link');
const targetsLink = document.getElementById('targets-link');
const resultsLink = document.getElementById('results-link');
const settingsLink = document.getElementById('settings-link');

const dashboardView = document.getElementById('dashboard-view');
const experimentsView = document.getElementById('experiments-view');
const targetsView = document.getElementById('targets-view');
const resultsView = document.getElementById('results-view');
const settingsView = document.getElementById('settings-view');

const newExperimentBtn = document.getElementById('new-experiment-btn');
const newExperimentBtn2 = document.getElementById('new-experiment-btn-2');
const newTargetBtn = document.getElementById('new-target-btn');

const experimentTypeSelect = document.getElementById('experiment-type');
const externalParams = document.getElementById('external-params');
const k8sParams = document.getElementById('k8s-params');

const createExperimentBtn = document.getElementById('create-experiment-btn');
const createTargetBtn = document.getElementById('create-target-btn');

// Bootstrap modals
const newExperimentModal = new bootstrap.Modal(document.getElementById('new-experiment-modal'));
const newTargetModal = new bootstrap.Modal(document.getElementById('new-target-modal'));

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    // Load initial data
    loadDashboardData();
    
    // Set up event listeners
    setupEventListeners();
    
    // Set up settings event listeners
    setupSettingsEventListeners();
});

// Set up event listeners
function setupEventListeners() {
    // Navigation
    dashboardLink.addEventListener('click', showDashboard);
    experimentsLink.addEventListener('click', showExperiments);
    targetsLink.addEventListener('click', showTargets);
    resultsLink.addEventListener('click', showResults);
    settingsLink.addEventListener('click', showSettings);
    
    // New experiment buttons
    newExperimentBtn.addEventListener('click', () => newExperimentModal.show());
    newExperimentBtn2.addEventListener('click', () => newExperimentModal.show());
    
    // New target button
    newTargetBtn.addEventListener('click', () => newTargetModal.show());
    
    // Experiment type change
    experimentTypeSelect.addEventListener('change', handleExperimentTypeChange);
    
    // Create buttons
    createExperimentBtn.addEventListener('click', createExperiment);
    createTargetBtn.addEventListener('click', createTarget);
    
    // Range slider
    const percentageSlider = document.getElementById('k8s-percentage');
    const percentageValue = document.getElementById('percentage-value');
    percentageSlider.addEventListener('input', () => {
        percentageValue.textContent = `${percentageSlider.value}%`;
    });
}

// Set up settings event listeners
function setupSettingsEventListeners() {
    // Save general settings
    const saveGeneralSettingsBtn = document.getElementById('save-general-settings-btn');
    if (saveGeneralSettingsBtn) {
        saveGeneralSettingsBtn.addEventListener('click', () => {
            const apiUrl = document.getElementById('api-url').value;
            const refreshInterval = document.getElementById('refresh-interval').value;
            
            // In a real application, this would save the settings to local storage or API
            alert('General settings saved successfully');
        });
    }
    
    // Save safety settings
    const saveSafetySettingsBtn = document.getElementById('save-safety-settings-btn');
    if (saveSafetySettingsBtn) {
        saveSafetySettingsBtn.addEventListener('click', () => {
            const protectedNamespaces = document.getElementById('protected-namespaces').value;
            const maxConcurrentExperiments = document.getElementById('max-concurrent-experiments').value;
            const autoTerminate = document.getElementById('auto-terminate').checked;
            
            // In a real application, this would save the settings to local storage or API
            alert('Safety settings saved successfully');
        });
    }
    
    // Export results
    const exportResultsBtn = document.getElementById('export-results-btn');
    if (exportResultsBtn) {
        exportResultsBtn.addEventListener('click', () => {
            // In a real application, this would export results to CSV or JSON
            alert('Results exported successfully');
        });
    }
}

// Show dashboard view
function showDashboard(e) {
    if (e) e.preventDefault();
    setActiveLink(dashboardLink);
    dashboardView.style.display = 'block';
    experimentsView.style.display = 'none';
    targetsView.style.display = 'none';
    resultsView.style.display = 'none';
    settingsView.style.display = 'none';
    loadDashboardData();
}

// Show experiments view
function showExperiments(e) {
    if (e) e.preventDefault();
    setActiveLink(experimentsLink);
    dashboardView.style.display = 'none';
    experimentsView.style.display = 'block';
    targetsView.style.display = 'none';
    resultsView.style.display = 'none';
    settingsView.style.display = 'none';
    loadAllExperiments();
}

// Show targets view
function showTargets(e) {
    if (e) e.preventDefault();
    setActiveLink(targetsLink);
    dashboardView.style.display = 'none';
    experimentsView.style.display = 'none';
    targetsView.style.display = 'block';
    resultsView.style.display = 'none';
    settingsView.style.display = 'none';
    loadTargets();
}

// Show results view
function showResults(e) {
    if (e) e.preventDefault();
    setActiveLink(resultsLink);
    dashboardView.style.display = 'none';
    experimentsView.style.display = 'none';
    targetsView.style.display = 'none';
    resultsView.style.display = 'block';
    settingsView.style.display = 'none';
    // TODO: Implement loadResults() function
}

// Show settings view
function showSettings(e) {
    if (e) e.preventDefault();
    setActiveLink(settingsLink);
    dashboardView.style.display = 'none';
    experimentsView.style.display = 'none';
    targetsView.style.display = 'none';
    resultsView.style.display = 'none';
    settingsView.style.display = 'block';
}

// Set active navigation link
function setActiveLink(link) {
    const links = [dashboardLink, experimentsLink, targetsLink, resultsLink, settingsLink];
    links.forEach(l => l.classList.remove('active'));
    link.classList.add('active');
}

// Handle experiment type change
function handleExperimentTypeChange() {
    const type = experimentTypeSelect.value;
    if (type === 'external-target') {
        externalParams.style.display = 'block';
        k8sParams.style.display = 'none';
    } else {
        externalParams.style.display = 'none';
        k8sParams.style.display = 'block';
    }
}

// Load dashboard data
async function loadDashboardData() {
    try {
        // Show loading indicators
        document.getElementById('loading-experiments').style.display = 'block';
        
        // Fetch experiments
        const experiments = await fetchExperiments();
        
        // Update summary cards
        updateSummaryCards(experiments);
        
        // Update recent experiments table
        updateRecentExperimentsTable(experiments);
        
        // Hide loading indicators
        document.getElementById('loading-experiments').style.display = 'none';
    } catch (error) {
        console.error('Error loading dashboard data:', error);
        // Hide loading indicators
        document.getElementById('loading-experiments').style.display = 'none';
    }
}

// Load all experiments
async function loadAllExperiments() {
    try {
        // Show loading indicator
        document.getElementById('loading-all-experiments').style.display = 'block';
        
        // Fetch experiments
        const experiments = await fetchExperiments();
        
        // Update experiments table
        updateAllExperimentsTable(experiments);
        
        // Hide loading indicator
        document.getElementById('loading-all-experiments').style.display = 'none';
    } catch (error) {
        console.error('Error loading experiments:', error);
        // Hide loading indicator
        document.getElementById('loading-all-experiments').style.display = 'none';
    }
}

// Load targets
async function loadTargets() {
    try {
        // Show loading indicator
        document.getElementById('loading-targets').style.display = 'block';
        
        // Fetch targets
        const targets = await fetchTargets();
        
        // Update targets table
        updateTargetsTable(targets);
        
        // Update target select in new experiment form
        updateTargetSelect(targets);
        
        // Hide loading indicator
        document.getElementById('loading-targets').style.display = 'none';
    } catch (error) {
        console.error('Error loading targets:', error);
        // Hide loading indicator
        document.getElementById('loading-targets').style.display = 'none';
    }
}

// Fetch experiments from API
async function fetchExperiments() {
    const response = await fetch(EXPERIMENTS_ENDPOINT);
    if (!response.ok) {
        throw new Error(`Failed to fetch experiments: ${response.statusText}`);
    }
    return await response.json();
}

// Fetch targets from API
async function fetchTargets() {
    const response = await fetch(TARGETS_ENDPOINT);
    if (!response.ok) {
        throw new Error(`Failed to fetch targets: ${response.statusText}`);
    }
    return await response.json();
}

// Update summary cards
function updateSummaryCards(experiments) {
    document.getElementById('total-experiments').textContent = experiments.length;
    
    const successful = experiments.filter(exp => exp.status === 'completed').length;
    document.getElementById('successful-experiments').textContent = successful;
    
    const failed = experiments.filter(exp => exp.status === 'failed').length;
    document.getElementById('failed-experiments').textContent = failed;
    
    const active = experiments.filter(exp => exp.status === 'running').length;
    document.getElementById('active-experiments').textContent = active;
}

// Update recent experiments table
function updateRecentExperimentsTable(experiments) {
    const tableBody = document.getElementById('recent-experiments-table');
    tableBody.innerHTML = '';
    
    // Sort experiments by creation date (newest first)
    const sortedExperiments = [...experiments].sort((a, b) => {
        return new Date(b.created_at) - new Date(a.created_at);
    });
    
    // Take only the 5 most recent experiments
    const recentExperiments = sortedExperiments.slice(0, 5);
    
    if (recentExperiments.length === 0) {
        const row = document.createElement('tr');
        row.innerHTML = '<td colspan="6" class="text-center">No experiments found</td>';
        tableBody.appendChild(row);
        return;
    }
    
    recentExperiments.forEach(experiment => {
        const row = document.createElement('tr');
        
        // Format the status badge
        let statusBadgeClass = 'bg-secondary';
        if (experiment.status === 'running') statusBadgeClass = 'bg-primary';
        if (experiment.status === 'completed') statusBadgeClass = 'bg-success';
        if (experiment.status === 'failed') statusBadgeClass = 'bg-danger';
        
        row.innerHTML = `
            <td>${experiment.name}</td>
            <td>${experiment.type}</td>
            <td>${experiment.target}</td>
            <td><span class="badge ${statusBadgeClass} status-badge">${experiment.status}</span></td>
            <td>${experiment.duration}s</td>
            <td>
                <button class="btn btn-sm btn-outline-primary btn-action view-experiment" data-id="${experiment.id}">
                    <i class="bi bi-eye"></i>
                </button>
                ${experiment.status === 'pending' ? `
                <button class="btn btn-sm btn-outline-success btn-action execute-experiment" data-id="${experiment.id}">
                    <i class="bi bi-play"></i>
                </button>
                ` : ''}
                ${experiment.status === 'running' ? `
                <button class="btn btn-sm btn-outline-danger btn-action stop-experiment" data-id="${experiment.id}">
                    <i class="bi bi-stop"></i>
                </button>
                ` : ''}
            </td>
        `;
        
        tableBody.appendChild(row);
    });
    
    // Add event listeners to buttons
    addExperimentButtonListeners();
}

// Update all experiments table
function updateAllExperimentsTable(experiments) {
    const tableBody = document.getElementById('all-experiments-table');
    tableBody.innerHTML = '';
    
    // Sort experiments by creation date (newest first)
    const sortedExperiments = [...experiments].sort((a, b) => {
        return new Date(b.created_at) - new Date(a.created_at);
    });
    
    if (sortedExperiments.length === 0) {
        const row = document.createElement('tr');
        row.innerHTML = '<td colspan="7" class="text-center">No experiments found</td>';
        tableBody.appendChild(row);
        return;
    }
    
    sortedExperiments.forEach(experiment => {
        const row = document.createElement('tr');
        
        // Format the status badge
        let statusBadgeClass = 'bg-secondary';
        if (experiment.status === 'running') statusBadgeClass = 'bg-primary';
        if (experiment.status === 'completed') statusBadgeClass = 'bg-success';
        if (experiment.status === 'failed') statusBadgeClass = 'bg-danger';
        
        // Format the date
        const createdDate = new Date(experiment.created_at);
        const formattedDate = createdDate.toLocaleString();
        
        row.innerHTML = `
            <td>${experiment.name}</td>
            <td>${experiment.type}</td>
            <td>${experiment.target}</td>
            <td><span class="badge ${statusBadgeClass} status-badge">${experiment.status}</span></td>
            <td>${formattedDate}</td>
            <td>${experiment.duration}s</td>
            <td>
                <button class="btn btn-sm btn-outline-primary btn-action view-experiment" data-id="${experiment.id}">
                    <i class="bi bi-eye"></i>
                </button>
                ${experiment.status === 'pending' ? `
                <button class="btn btn-sm btn-outline-success btn-action execute-experiment" data-id="${experiment.id}">
                    <i class="bi bi-play"></i>
                </button>
                ` : ''}
                ${experiment.status === 'running' ? `
                <button class="btn btn-sm btn-outline-danger btn-action stop-experiment" data-id="${experiment.id}">
                    <i class="bi bi-stop"></i>
                </button>
                ` : ''}
                <button class="btn btn-sm btn-outline-danger btn-action delete-experiment" data-id="${experiment.id}">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        `;
        
        tableBody.appendChild(row);
    });
    
    // Add event listeners to buttons
    addExperimentButtonListeners();
}

// Update targets table
function updateTargetsTable(targets) {
    const tableBody = document.getElementById('targets-table');
    tableBody.innerHTML = '';
    
    // Sort targets by creation date (newest first)
    const sortedTargets = [...targets].sort((a, b) => {
        return new Date(b.created_at) - new Date(a.created_at);
    });
    
    if (sortedTargets.length === 0) {
        const row = document.createElement('tr');
        row.innerHTML = '<td colspan="6" class="text-center">No targets found</td>';
        tableBody.appendChild(row);
        return;
    }
    
    sortedTargets.forEach(target => {
        const row = document.createElement('tr');
        
        // Format the date
        const createdDate = new Date(target.created_at);
        const formattedDate = createdDate.toLocaleString();
        
        row.innerHTML = `
            <td>${target.name}</td>
            <td>${target.type}</td>
            <td>${target.namespace}</td>
            <td>${target.selector}</td>
            <td>${formattedDate}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary btn-action view-target" data-id="${target.id}">
                    <i class="bi bi-eye"></i>
                </button>
                <button class="btn btn-sm btn-outline-danger btn-action delete-target" data-id="${target.id}">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        `;
        
        tableBody.appendChild(row);
    });
    
    // Add event listeners to buttons
    addTargetButtonListeners();
}

// Update target select in new experiment form
function updateTargetSelect(targets) {
    const targetSelect = document.getElementById('experiment-target');
    
    // Clear existing options except the first one
    while (targetSelect.options.length > 1) {
        targetSelect.remove(1);
    }
    
    // Add targets as options
    targets.forEach(target => {
        const option = document.createElement('option');
        option.value = target.selector;
        option.textContent = `${target.name} (${target.type})`;
        option.dataset.type = target.type;
        targetSelect.appendChild(option);
    });
}

// Add event listeners to experiment buttons
function addExperimentButtonListeners() {
    // View experiment buttons
    document.querySelectorAll('.view-experiment').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.dataset.id;
            viewExperiment(id);
        });
    });
    
    // Execute experiment buttons
    document.querySelectorAll('.execute-experiment').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.dataset.id;
            executeExperiment(id);
        });
    });
    
    // Stop experiment buttons
    document.querySelectorAll('.stop-experiment').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.dataset.id;
            stopExperiment(id);
        });
    });
    
    // Delete experiment buttons
    document.querySelectorAll('.delete-experiment').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.dataset.id;
            deleteExperiment(id);
        });
    });
}

// Add event listeners to target buttons
function addTargetButtonListeners() {
    // View target buttons
    document.querySelectorAll('.view-target').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.dataset.id;
            viewTarget(id);
        });
    });
    
    // Delete target buttons
    document.querySelectorAll('.delete-target').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.dataset.id;
            deleteTarget(id);
        });
    });
}

// View experiment details
function viewExperiment(id) {
    // In a real application, this would show a modal with experiment details
    alert(`View experiment ${id}`);
}

// Execute experiment
async function executeExperiment(id) {
    try {
        const response = await fetch(`${EXPERIMENTS_ENDPOINT}/${id}/execute`, {
            method: 'POST'
        });
        
        if (!response.ok) {
            throw new Error(`Failed to execute experiment: ${response.statusText}`);
        }
        
        // Reload data
        loadDashboardData();
        loadAllExperiments();
        
        alert(`Experiment ${id} started successfully`);
    } catch (error) {
        console.error('Error executing experiment:', error);
        alert(`Error executing experiment: ${error.message}`);
    }
}

// Stop experiment
async function stopExperiment(id) {
    // In a real application, this would call an API endpoint to stop the experiment
    alert(`Stop experiment ${id}`);
}

// Delete experiment
async function deleteExperiment(id) {
    if (!confirm('Are you sure you want to delete this experiment?')) {
        return;
    }
    
    try {
        const response = await fetch(`${EXPERIMENTS_ENDPOINT}/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error(`Failed to delete experiment: ${response.statusText}`);
        }
        
        // Reload data
        loadDashboardData();
        loadAllExperiments();
        
        alert('Experiment deleted successfully');
    } catch (error) {
        console.error('Error deleting experiment:', error);
        alert(`Error deleting experiment: ${error.message}`);
    }
}

// View target details
function viewTarget(id) {
    // In a real application, this would show a modal with target details
    alert(`View target ${id}`);
}

// Delete target
async function deleteTarget(id) {
    if (!confirm('Are you sure you want to delete this target?')) {
        return;
    }
    
    try {
        const response = await fetch(`${TARGETS_ENDPOINT}/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error(`Failed to delete target: ${response.statusText}`);
        }
        
        // Reload data
        loadTargets();
        
        alert('Target deleted successfully');
    } catch (error) {
        console.error('Error deleting target:', error);
        alert(`Error deleting target: ${error.message}`);
    }
}

// Create experiment
async function createExperiment() {
    const name = document.getElementById('experiment-name').value;
    const description = document.getElementById('experiment-description').value;
    const type = document.getElementById('experiment-type').value;
    const target = document.getElementById('experiment-target').value;
    const duration = parseInt(document.getElementById('experiment-duration').value);
    
    if (!name || !type || !target || !duration) {
        alert('Please fill in all required fields');
        return;
    }
    
    let parameters = {};
    
    if (type === 'external-target') {
        parameters = {
            target_type: 'external',
            endpoint: document.getElementById('external-endpoint').value,
            type: 'service-failure'
        };
        
        const authToken = document.getElementById('external-auth-token').value;
        if (authToken) {
            parameters.auth_token = authToken;
        }
        
        const cleanupEndpoint = document.getElementById('external-cleanup').value;
        if (cleanupEndpoint) {
            parameters.cleanup_endpoint = cleanupEndpoint;
        }
    } else {
        parameters = {
            namespace: document.getElementById('k8s-namespace').value,
            percentage: document.getElementById('k8s-percentage').value
        };
    }
    
    const experimentData = {
        name,
        description,
        type,
        target,
        parameters,
        duration
    };
    
    try {
        const response = await fetch(EXPERIMENTS_ENDPOINT, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(experimentData)
        });
        
        if (!response.ok) {
            throw new Error(`Failed to create experiment: ${response.statusText}`);
        }
        
        // Close modal
        newExperimentModal.hide();
        
        // Reset form
        document.getElementById('new-experiment-form').reset();
        
        // Reload data
        loadDashboardData();
        loadAllExperiments();
        
        alert('Experiment created successfully');
    } catch (error) {
        console.error('Error creating experiment:', error);
        alert(`Error creating experiment: ${error.message}`);
    }
}

// Create target
async function createTarget() {
    const name = document.getElementById('target-name').value;
    const description = document.getElementById('target-description').value;
    const type = document.getElementById('target-type').value;
    const namespace = document.getElementById('target-namespace').value;
    const selector = document.getElementById('target-selector').value;
    
    if (!name || !type || !namespace || !selector) {
        alert('Please fill in all required fields');
        return;
    }
    
    const targetData = {
        name,
        description,
        type,
        namespace,
        selector
    };
    
    try {
        const response = await fetch(TARGETS_ENDPOINT, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(targetData)
        });
        
        if (!response.ok) {
            throw new Error(`Failed to create target: ${response.statusText}`);
        }
        
        // Close modal
        newTargetModal.hide();
        
        // Reset form
        document.getElementById('new-target-form').reset();
        
        // Reload data
        loadTargets();
        
        alert('Target created successfully');
    } catch (error) {
        console.error('Error creating target:', error);
        alert(`Error creating target: ${error.message}`);
    }
}