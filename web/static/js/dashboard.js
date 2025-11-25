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
    const tableRows = sortedData.map(tx =>`
        <tr>
            <td>${formatDate(tx.date)}</td>
            <td>${tx.category}</td>
            <td>${tx.description || '-'}</td>
            <td class="amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
        </tr>
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
                </tr>
            </thead>
            <tbody>
                ${tableRows}
            </tbody>
        </table>
    `;

    // Add click handlers to headers
    document.querySelectorAll('.transactions-table th[data-column]').forEach(th => {
        th.addEventListener('click', () => {
            handleSort(th.getAttribute('data-column'));
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

    container.innerHTML = sortedClusters.map((cluster, index) =>`
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
                        </tr>
                    </thead>
                    <tbody>
                        ${cluster.transactions.map(tx => `
                            <tr>
                                <td>${formatDate(tx.date)}</td>
                                <td>${tx.category}</td>
                                <td>${tx.description || '-'}</td>
                                <td class="amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
        </div>
    `).join('');

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

// Load data on page load
const currentMonth = document.getElementById('currentMonth').value;
loadStats(currentMonth);
loadTransactions(currentMonth);
