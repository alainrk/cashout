// Format currency
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'EUR',
        minimumFractionDigits: 2
    }).format(amount);
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
    });
}

// Format date for mobile (shorter format)
function formatDateShort(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
    });
}

// Load statistics
async function loadStats(month) {
    try {
        const response = await fetch('/web/api/stats?month=' + month);
        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Failed to load stats');

        const statsGrid = document.getElementById('statsGrid');
        statsGrid.innerHTML = `
            <div class="stat-card">
                <div class="stat-label">Balance</div>
                <div class="stat-value">${formatCurrency(data.balance)}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Total Income</div>
                <div class="stat-value income">${formatCurrency(data.totalIncome)}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Total Expenses</div>
                <div class="stat-value expense">${formatCurrency(data.totalExpenses)}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Transactions</div>
                <div class="stat-value">${data.totalTransactions}</div>
            </div>
        `;
    } catch (error) {
        document.getElementById('statsGrid').innerHTML =
            '<div class="error">Failed to load statistics: ' + error.message + '</div>';
    }
}

let transactionsData = [];
let currentView = localStorage.getItem('dashboardView') || 'list';
let sortColumn = 'date';
let sortDirection = 'desc';

// Load transactions
async function loadTransactions(month) {
    try {
        const response = await fetch('/web/api/transactions?month=' + month);
        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Failed to load transactions');

        transactionsData = data.transactions;
        renderTransactions();

    } catch (error) {
        document.getElementById('transactionsContainer').innerHTML =
            '<div class="error">Failed to load transactions: ' + error.message + '</div>';
    }
}

// Sort transactions data
function sortTransactions(data, column, direction) {
    return [...data].sort((a, b) => {
        let aVal, bVal;

        switch(column) {
            case 'date':
                aVal = new Date(a.date);
                bVal = new Date(b.date);
                break;
            case 'category':
                aVal = a.category.toLowerCase();
                bVal = b.category.toLowerCase();
                break;
            case 'description':
                aVal = (a.description || '').toLowerCase();
                bVal = (b.description || '').toLowerCase();
                break;
            case 'amount':
                aVal = Math.abs(a.amount);
                bVal = Math.abs(b.amount);
                break;
            default:
                return 0;
        }

        if (aVal < bVal) return direction === 'asc' ? -1 : 1;
        if (aVal > bVal) return direction === 'asc' ? 1 : -1;
        return 0;
    });
}

// Handle column header click
function handleSort(column) {
    if (sortColumn === column) {
        sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
        sortColumn = column;
        sortDirection = 'asc';
    }
    renderTransactions();
}

// Render transactions based on the current view
function renderTransactions() {
    if (currentView === 'list') {
        renderListView();
    } else {
        renderClusteredView();
    }
}

// Render list view
function renderListView() {
    const container = document.getElementById('transactionsContainer');

    if (transactionsData.length === 0) {
        container.innerHTML = '<p>No transactions for this month.</p>';
        return;
    }

    const sortedData = sortTransactions(transactionsData, sortColumn, sortDirection);

    // Render table rows for desktop
    const tableRows = sortedData.map(tx =>`
        <tr>
            <td>${formatDate(tx.date)}</td>
            <td>${tx.category}</td>
            <td>${tx.description || '-'}</td>
            <td class="amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
            <td class="actions">
                <button class="delete-btn" data-id="${tx.id}" data-description="${tx.description || tx.category}" title="Delete transaction">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    </svg>
                </button>
            </td>
        </tr>
    `).join('');

    // Render cards for mobile
    const cards = sortedData.map(tx =>`
        <div class="transaction-card">
            <div class="transaction-card-header">
                <span class="transaction-card-date">${formatDateShort(tx.date)}</span>
                <span class="transaction-card-amount amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</span>
            </div>
            <div class="transaction-card-body">
                <div class="transaction-card-info">
                    <div class="transaction-card-category">${tx.category}</div>
                    <div class="transaction-card-description">${tx.description || '-'}</div>
                </div>
                <button class="delete-btn" data-id="${tx.id}" data-description="${tx.description || tx.category}" title="Delete transaction">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    </svg>
                </button>
            </div>
        </div>
    `).join('');

    const getSortClass = (column) => {
        if (sortColumn === column) {
            return sortDirection === 'asc' ? 'sorted-asc' : 'sorted-desc';
        }
        return 'sortable';
    };

    container.innerHTML =`
        <table class="transactions-table">
            <thead>
                <tr>
                    <th class="${getSortClass('date')}" data-column="date">Date</th>
                    <th class="${getSortClass('category')}" data-column="category">Category</th>
                    <th class="${getSortClass('description')}" data-column="description">Description</th>
                    <th class="${getSortClass('amount')}" data-column="amount">Amount</th>
                    <th class="actions-header">Actions</th>
                </tr>
            </thead>
            <tbody>
                ${tableRows}
            </tbody>
        </table>
        <div class="transaction-cards">
            ${cards}
        </div>
    `;

    // Add click handlers to headers
    document.querySelectorAll('.transactions-table th[data-column]').forEach(th => {
        th.addEventListener('click', () => {
            handleSort(th.getAttribute('data-column'));
        });
    });

    // Add click handlers to delete buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.stopPropagation();
            const id = parseInt(btn.getAttribute('data-id'));
            const description = btn.getAttribute('data-description');
            showDeleteConfirmation(id, description);
        });
    });
}

// Render clustered view
function renderClusteredView() {
    const container = document.getElementById('transactionsContainer');

    if (transactionsData.length === 0) {
        container.innerHTML = '<p>No transactions for this month.</p>';
        return;
    }

    const clusteredData = transactionsData.reduce((acc, tx) => {
        const key = `${tx.type}-${tx.category}`;
        if (!acc[key]) {
            acc[key] = {
                type: tx.type,
                category: tx.category,
                total: 0,
                transactions: []
            };
        }
        acc[key].total += tx.amount;
        acc[key].transactions.push(tx);
        return acc;
    }, {});

    const sortedClusters = Object.values(clusteredData).sort((a, b) => b.total - a.total);

    container.innerHTML = sortedClusters.map((cluster, index) => {
        // Render table rows for desktop
        const tableRows = cluster.transactions.map(tx => `
            <tr>
                <td>${formatDate(tx.date)}</td>
                <td>${tx.category}</td>
                <td>${tx.description || '-'}</td>
                <td class="amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
                <td class="actions">
                    <button class="delete-btn" data-id="${tx.id}" data-description="${tx.description || tx.category}" title="Delete transaction">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="3 6 5 6 21 6"></polyline>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                    </button>
                </td>
            </tr>
        `).join('');

        // Render cards for mobile
        const cards = cluster.transactions.map(tx => `
            <div class="transaction-card">
                <div class="transaction-card-header">
                    <span class="transaction-card-date">${formatDateShort(tx.date)}</span>
                    <span class="transaction-card-amount amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</span>
                </div>
                <div class="transaction-card-body">
                    <div class="transaction-card-info">
                        <div class="transaction-card-category">${tx.category}</div>
                        <div class="transaction-card-description">${tx.description || '-'}</div>
                    </div>
                    <button class="delete-btn" data-id="${tx.id}" data-description="${tx.description || tx.category}" title="Delete transaction">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="3 6 5 6 21 6"></polyline>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                    </button>
                </div>
            </div>
        `).join('');

        return `
            <div class="cluster">
                <div class="cluster-header" data-cluster-index="${index}">
                    <span class="cluster-title">
                        <span class="cluster-icon" id="icon-${index}"></span>
                        ${cluster.category} (${cluster.type})
                    </span>
                    <span class="cluster-total ${cluster.type.toLowerCase()}">${cluster.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(cluster.total))}</span>
                </div>
                <div class="cluster-content" id="cluster-${index}">
                    <table class="transactions-table cluster-transactions">
                        <thead>
                            <tr>
                                <th>Date</th>
                                <th>Category</th>
                                <th>Description</th>
                                <th>Amount</th>
                                <th class="actions-header">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${tableRows}
                        </tbody>
                    </table>
                    <div class="transaction-cards">
                        ${cards}
                    </div>
                </div>
            </div>
        `;
    }).join('');

    // Add click handlers to cluster headers
    document.querySelectorAll('.cluster-header').forEach(header => {
        header.addEventListener('click', (e) => {
            const clusterIndex = e.currentTarget.getAttribute('data-cluster-index');
            const content = document.getElementById(`cluster-${clusterIndex}`);
            const icon = document.getElementById(`icon-${clusterIndex}`);
            content.classList.toggle('expanded');
            icon.classList.toggle('expanded');
        });
    });

    // Add click handlers to delete buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.stopPropagation();
            const id = parseInt(btn.getAttribute('data-id'));
            const description = btn.getAttribute('data-description');
            showDeleteConfirmation(id, description);
        });
    });
}


// Event Listeners for view toggle
document.getElementById('listViewBtn').addEventListener('click', () => {
    currentView = 'list';
    localStorage.setItem('dashboardView', 'list');
    document.getElementById('listViewBtn').classList.add('active');
    document.getElementById('clusteredViewBtn').classList.remove('active');
    renderTransactions();
});

document.getElementById('clusteredViewBtn').addEventListener('click', () => {
    currentView = 'clustered';
    localStorage.setItem('dashboardView', 'clustered');
    document.getElementById('clusteredViewBtn').classList.add('active');
    document.getElementById('listViewBtn').classList.remove('active');
    renderTransactions();
});

// Initialize button states based on saved preference
if (currentView === 'clustered') {
    document.getElementById('clusteredViewBtn').classList.add('active');
    document.getElementById('listViewBtn').classList.remove('active');
}

// Transaction Form Handling

// Set today's date as default
document.getElementById('txDate').valueAsDate = new Date();

// Function to load categories for a given type
async function loadCategories(type) {
    const categorySelect = document.getElementById('txCategory');

    if (!type) {
        categorySelect.innerHTML = '<option value="">Select type first</option>';
        return;
    }

    try {
        categorySelect.innerHTML = '<option value="">Loading...</option>';
        const response = await fetch(`/web/api/categories?type=${type}`);
        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Failed to load categories');

        categorySelect.innerHTML = '';

        data.categories.forEach((category, index) => {
            const option = document.createElement('option');
            option.value = category;
            option.textContent = category;
            if (index === 0) {
                option.selected = true; // Auto-select first category
            }
            categorySelect.appendChild(option);
        });
    } catch (error) {
        categorySelect.innerHTML = '<option value="">Error loading categories</option>';
        showTxMessage('Failed to load categories: ' + error.message, 'error');
    }
}

// Load categories for default Expense type on page load
loadCategories('Expense');

// Load categories when type changes
document.getElementById('txType').addEventListener('change', function() {
    loadCategories(this.value);
});

// Handle form submission
document.getElementById('addTransactionForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const submitBtn = document.getElementById('submitTxBtn');
    const messageDiv = document.getElementById('txMessage');

    submitBtn.disabled = true;
    submitBtn.textContent = 'Adding...';
    messageDiv.textContent = '';
    messageDiv.className = 'message';

    // Normalize amount: replace comma with dot and parse
    const amountValue = document.getElementById('txAmount').value.replace(',', '.');
    const amount = parseFloat(amountValue);

    // Validate amount
    if (isNaN(amount) || amount <= 0) {
        showTxMessage('Please enter a valid amount', 'error');
        submitBtn.disabled = false;
        submitBtn.textContent = 'Add Transaction';
        return;
    }

    const formData = {
        type: document.getElementById('txType').value,
        category: document.getElementById('txCategory').value,
        amount: amount,
        date: document.getElementById('txDate').value,
        description: document.getElementById('txDescription').value
    };

    try {
        const response = await fetch('/web/api/transactions/create', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Failed to create transaction');

        showTxMessage('Transaction added successfully!', 'success');

        // Reset form
        document.getElementById('addTransactionForm').reset();
        document.getElementById('txDate').valueAsDate = new Date();

        // Reset to default Expense type and reload categories
        document.getElementById('txType').value = 'Expense';
        loadCategories('Expense');

        // Reload transactions and stats
        loadStats(currentMonth);
        loadTransactions(currentMonth);

        // Navigate to transactions page after a short delay
        setTimeout(() => {
            showPage('transactions');
        }, 1500);

    } catch (error) {
        showTxMessage('Error: ' + error.message, 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Add Transaction';
    }
});

function showTxMessage(text, type) {
    const messageDiv = document.getElementById('txMessage');
    messageDiv.textContent = text;
    messageDiv.className = 'message ' + type;

    // Auto-hide success messages after 3 seconds
    if (type === 'success') {
        setTimeout(() => {
            messageDiv.textContent = '';
            messageDiv.className = 'message';
        }, 3000);
    }
}

// Page Navigation
let currentPage = localStorage.getItem('currentPage') || 'transactions';

function showPage(pageName) {
    // Hide all pages
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });

    // Remove active class from all tabs
    document.querySelectorAll('.nav-tab').forEach(tab => {
        tab.classList.remove('active');
    });

    // Show selected page
    const pageMap = {
        'transactions': 'transactionsPage',
        'add-transaction': 'addTransactionPage',
        'security': 'securityPage'
    };

    const pageId = pageMap[pageName];
    if (pageId) {
        document.getElementById(pageId).classList.add('active');
        document.querySelector(`[data-page="${pageName}"]`).classList.add('active');
        currentPage = pageName;
        localStorage.setItem('currentPage', pageName);
    }
}

// Add click handlers to navigation tabs
document.querySelectorAll('.nav-tab').forEach(tab => {
    tab.addEventListener('click', () => {
        const pageName = tab.getAttribute('data-page');
        showPage(pageName);
    });
});

// Initialize page on load
showPage(currentPage);

// Load data on page load
const currentMonth = document.getElementById('currentMonth').value;
loadStats(currentMonth);
loadTransactions(currentMonth);

// Delete transaction functionality
function showDeleteConfirmation(id, description) {
    // Create modal overlay
    const modal = document.createElement('div');
    modal.className = 'delete-modal-overlay';
    modal.innerHTML = `
        <div class="delete-modal">
            <div class="delete-modal-header">
                <h3>Delete Transaction</h3>
            </div>
            <div class="delete-modal-body">
                <p>Are you sure you want to delete this transaction?</p>
                <p class="delete-transaction-info">${description}</p>
            </div>
            <div class="delete-modal-footer">
                <button class="cancel-btn" id="cancelDelete">Cancel</button>
                <button class="confirm-delete-btn" id="confirmDelete">Delete</button>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    // Handle cancel
    document.getElementById('cancelDelete').addEventListener('click', () => {
        document.body.removeChild(modal);
    });

    // Handle confirm
    document.getElementById('confirmDelete').addEventListener('click', async () => {
        await deleteTransaction(id);
        document.body.removeChild(modal);
    });

    // Handle click outside modal
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            document.body.removeChild(modal);
        }
    });
}

async function deleteTransaction(id) {
    try {
        const response = await fetch('/web/api/transactions/delete', {
            method: 'DELETE',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ id: id })
        });

        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Failed to delete transaction');

        // Reload transactions and stats
        await loadStats(currentMonth);
        await loadTransactions(currentMonth);

    } catch (error) {
        alert('Error deleting transaction: ' + error.message);
    }
}
