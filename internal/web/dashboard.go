package web

import (
	"cashout/internal/client"
	"cashout/internal/model"
	"html/template"
	"net/http"
	"time"
)

const (
	monthLayout = "2006-01"
)

// handleDashboard shows the main dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, basePath+"/login", http.StatusSeeOther)
		return
	}

	// Parse month from query, default to current month
	monthStr := r.URL.Query().Get("month")
	currentMonth, err := time.Parse(monthLayout, monthStr)
	if err != nil || currentMonth.After(time.Now()) {
		currentMonth = time.Now()
	}

	// Calculate previous and next months
	prevMonth := currentMonth.AddDate(0, -1, 0)
	nextMonth := currentMonth.AddDate(0, 1, 0)

	// Disable next month button if it's the future
	now := time.Now()
	isCurrentMonth := currentMonth.Format(monthLayout) == now.Format(monthLayout)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Cashout Dashboard - {{.User.Name}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        * {
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
        }
        .header {
            background: white;
            border-bottom: 1px solid #e0e0e0;
            padding: 1rem 0;
            position: sticky;
            top: 0;
            z-index: 100;
        }
        .header-content {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 1rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .logo {
            font-size: 1.5rem;
            font-weight: bold;
            color: #333;
        }
        .user-info {
            display: flex;
            align-items: center;
            gap: 1rem;
        }
        .logout-btn {
            padding: 0.5rem 1rem;
            background: #dc3545;
            color: white;
            border: none;
            border-radius: 4px;
            text-decoration: none;
            font-size: 0.9rem;
            transition: background 0.2s;
        }
        .logout-btn:hover {
            background: #c82333;
        }
        .container {
            max-width: 1200px;
            margin: 2rem auto;
            padding: 0 1rem;
        }
		.month-navigation {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 2rem;
		}
		.month-navigation a {
			padding: 0.5rem 1rem;
			background: #007bff;
			color: white;
			text-decoration: none;
			border-radius: 4px;
			transition: background 0.2s;
		}
		.month-navigation a:hover {
			background: #0056b3;
		}
		.month-navigation a.disabled {
			background: #6c757d;
			pointer-events: none;
		}
		.month-navigation h2 {
			margin: 0;
			font-size: 1.5rem;
			font-weight: 600;
		}
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
        }
        .stat-card {
            background: white;
            padding: 1.5rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .stat-label {
            color: #666;
            font-size: 0.875rem;
            margin-bottom: 0.5rem;
        }
        .stat-value {
            font-size: 2rem;
            font-weight: bold;
            color: #333;
        }
        .stat-change {
            font-size: 0.875rem;
            margin-top: 0.5rem;
        }
        .positive {
            color: #28a745;
        }
        .negative {
            color: #dc3545;
        }
        .section {
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 2rem;
        }
        .section-title {
            font-size: 1.25rem;
            font-weight: 600;
            margin-bottom: 1.5rem;
            color: #333;
        }
				.section-header {
					display: flex;
					justify-content: space-between;
					align-items: center;
					margin-bottom: 1.5rem;
				}
				.view-toggle button {
					padding: 0.5rem 1rem;
					border: 1px solid #007bff;
					background: white;
					color: #007bff;
					cursor: pointer;
					transition: all 0.2s;
				}
				.view-toggle button.active {
					background: #007bff;
					color: white;
				}
				.view-toggle button:first-child {
					border-top-left-radius: 4px;
					border-bottom-left-radius: 4px;
				}
         .view-toggle button:last-child {
            border-top-right-radius: 4px;
            border-bottom-right-radius: 4px;
            margin-left: -1px;
        }
        .transactions-table {
            width: 100%;
            border-collapse: collapse;
        }
        .transactions-table th {
            text-align: left;
            padding: 0.75rem;
            border-bottom: 2px solid #e0e0e0;
            color: #666;
            font-weight: 500;
        }
        .transactions-table td {
            padding: 0.75rem;
            border-bottom: 1px solid #f0f0f0;
        }
        .transactions-table tr:hover {
            background: #f8f9fa;
        }
        .cluster {
            margin-bottom: 1.5rem;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            overflow: hidden;
        }
        .cluster-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1rem;
            background: #f8f9fa;
            border-bottom: 1px solid #e0e0e0;
            cursor: pointer;
            user-select: none;
        }
        .cluster-header:hover {
            background: #e9ecef;
        }
        .cluster-title {
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        .cluster-icon {
            display: inline-block;
            width: 0;
            height: 0;
            border-left: 5px solid transparent;
            border-right: 5px solid transparent;
            border-top: 6px solid #666;
            transition: transform 0.2s ease;
        }
        .cluster-icon.expanded {
            transform: rotate(180deg);
        }
        .cluster-total {
            font-weight: 500;
        }
        .cluster-content {
            display: none;
        }
        .cluster-content.expanded {
            display: block;
        }
        .cluster-transactions {
            padding: 0;
        }
        .amount {
            font-weight: 500;
        }
        .income {
            color: #28a745;
        }
        .expense {
            color: #dc3545;
        }
        .loading {
            text-align: center;
            padding: 2rem;
            color: #666;
        }
        .error {
            background: #fee;
            color: #c33;
            padding: 1rem;
            border-radius: 4px;
            margin-bottom: 1rem;
        }
        @media (max-width: 768px) {
            .header-content {
                flex-direction: column;
                gap: 1rem;
            }
            .stat-value {
                font-size: 1.5rem;
            }
			.month-navigation {
				flex-direction: column;
				gap: 1rem;
			}
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="header-content">
            <div class="logo">Cashout</div>
            <div class="user-info">
                <span>Welcome, <strong>{{.User.Name}}</strong></span>
                <a href="/web/logout" class="logout-btn">Logout</a>
            </div>
        </div>
    </div>

    <div class="container">
		<div class="month-navigation">
			<a href="/web/dashboard?month={{.PrevMonth}}">Previous</a>
			<h2>{{.CurrentMonthTitle}}</h2>
			<a href="/web/dashboard?month={{.NextMonth}}" {{if .IsCurrentMonth}}class="disabled"{{end}}>Next</a>
		</div>

		<input type="hidden" id="currentMonth" value="{{.CurrentMonth}}">

        <div class="stats-grid" id="statsGrid">
            <div class="loading">Loading statistics...</div>
        </div>

        <div class="section">
			<div class="section-header">
				<h2 class="section-title">Transactions</h2>
				<div class="view-toggle">
					<button id="listViewBtn" class="active">List</button>
					<button id="clusteredViewBtn">Clustered</button>
				</div>
			</div>
            <div id="transactionsContainer">
                <div class="loading">Loading transactions...</div>
            </div>
        </div>
    </div>

    <script>
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
                statsGrid.innerHTML = ` + "`" + `
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
                ` + "`" + `;
            } catch (error) {
                document.getElementById('statsGrid').innerHTML =
                    '<div class="error">Failed to load statistics: ' + error.message + '</div>';
            }
        }

        let transactionsData = [];
        let currentView = localStorage.getItem('dashboardView') || 'list';

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

            const tableRows = transactionsData.map(tx =>` + "`" + ` 
                <tr>
                    <td>${formatDate(tx.date)}</td>
                    <td>${tx.category}</td>
                    <td>${tx.description || '-'}</td>
                    <td class="amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
                </tr>
            ` + "`" + `).join('');

            container.innerHTML =` + "`" + ` 
                <table class="transactions-table">
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Category</th>
                            <th>Description</th>
                            <th>Amount</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${tableRows}
                    </tbody>
                </table>
            ` + "`" + `;
        }

        // Render clustered view
        function renderClusteredView() {
            const container = document.getElementById('transactionsContainer');

            if (transactionsData.length === 0) {
                container.innerHTML = '<p>No transactions for this month.</p>';
                return;
            }

            const clusteredData = transactionsData.reduce((acc, tx) => {
                const key = ` + "`" + `${tx.type}-${tx.category}` + "`" + `;
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

            container.innerHTML = sortedClusters.map((cluster, index) =>` + "`" + `
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
                                ${cluster.transactions.map(tx => ` + "`" + `
                                    <tr>
                                        <td>${formatDate(tx.date)}</td>
                                        <td>${tx.category}</td>
                                        <td>${tx.description || '-'}</td>
                                        <td class="amount ${tx.type.toLowerCase()}">${tx.type.toLowerCase() === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
                                    </tr>
                                ` + "`" + `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            ` + "`" + `).join('');

            // Add click handlers to cluster headers
            document.querySelectorAll('.cluster-header').forEach(header => {
                header.addEventListener('click', (e) => {
                    const clusterIndex = e.currentTarget.getAttribute('data-cluster-index');
                    const content = document.getElementById(` + "`" + `cluster-${clusterIndex}` + "`" + `);
                    const icon = document.getElementById(` + "`" + `icon-${clusterIndex}` + "`" + `);
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
    </script>
</body>
</html>
`

	t, err := template.New("dashboard").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User              *model.User
		CurrentMonthTitle string
		CurrentMonth      string
		PrevMonth         string
		NextMonth         string
		IsCurrentMonth    bool
	}{
		User:              user,
		CurrentMonthTitle: currentMonth.Format("January 2006"),
		CurrentMonth:      currentMonth.Format(monthLayout),
		PrevMonth:         prevMonth.Format(monthLayout),
		NextMonth:         nextMonth.Format(monthLayout),
		IsCurrentMonth:    isCurrentMonth,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

// handleAPIStats returns user statistics for a given month
func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse month from query, default to current month
	monthStr := r.URL.Query().Get("month")
	currentMonth, err := time.Parse(monthLayout, monthStr)
	if err != nil {
		currentMonth = time.Now()
	}

	// Get transactions for the month
	startDate := time.Date(currentMonth.Year(), currentMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
	transactions, err := s.repositories.Transactions.GetUserTransactionsByDateRange(user.TgID, startDate, endDate)
	if err != nil {
		s.sendJSONError(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Calculate statistics
	var totalIncome, totalExpenses float64

	for _, tx := range transactions {
		if tx.Type == model.TypeIncome {
			totalIncome += tx.Amount
		} else {
			totalExpenses += tx.Amount
		}
	}

	balance := totalIncome - totalExpenses

	stats := map[string]interface{}{
		"balance":           balance,
		"totalIncome":       totalIncome,
		"totalExpenses":     totalExpenses,
		"totalTransactions": len(transactions),
	}

	s.sendJSONSuccess(w, stats)
}

// handleAPITransactions returns user transactions for a given month
func (s *Server) handleAPITransactions(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse month from query, default to current month
	monthStr := r.URL.Query().Get("month")
	currentMonth, err := time.Parse(monthLayout, monthStr)
	if err != nil {
		currentMonth = time.Now()
	}

	// Get transactions for the month
	startDate := time.Date(currentMonth.Year(), currentMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
	transactions, err := s.repositories.Transactions.GetUserTransactionsByDateRange(user.TgID, startDate, endDate)
	if err != nil {
		s.sendJSONError(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	type TransactionResponse struct {
		ID          int64     `json:"id"`
		Date        time.Time `json:"date"`
		Category    string    `json:"category"`
		Description string    `json:"description"`
		Amount      float64   `json:"amount"`
		Type        string    `json:"type"`
	}

	transactionResponses := make([]TransactionResponse, len(transactions))
	for i, tx := range transactions {
		transactionResponses[i] = TransactionResponse{
			ID:          tx.ID,
			Date:        tx.Date,
			Category:    string(tx.Category),
			Description: tx.Description,
			Amount:      tx.Amount,
			Type:        string(tx.Type),
		}
	}

	response := map[string]interface{}{
		"transactions": transactionResponses,
		"count":        len(transactionResponses),
	}

	s.sendJSONSuccess(w, response)
}
