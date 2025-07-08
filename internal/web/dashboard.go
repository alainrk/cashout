package web

import (
	"cashout/internal/client"
	"cashout/internal/model"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// handleDashboard shows the main dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Cashout Dashboard - {{.User.FirstName}}</title>
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
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="header-content">
            <div class="logo">💰 Cashout</div>
            <div class="user-info">
                <span>Welcome, <strong>{{.User.FirstName}}</strong></span>
                <a href="/logout" class="logout-btn">Logout</a>
            </div>
        </div>
    </div>

    <div class="container">
        <div class="stats-grid" id="statsGrid">
            <div class="loading">Loading statistics...</div>
        </div>

        <div class="section">
            <h2 class="section-title">Recent Transactions</h2>
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
                currency: 'USD',
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
                hour: '2-digit',
                minute: '2-digit'
            });
        }

        // Load statistics
        async function loadStats() {
            try {
                const response = await fetch('/api/stats');
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
                        <div class="stat-change positive">+${formatCurrency(data.monthlyIncome)} this month</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Total Expenses</div>
                        <div class="stat-value expense">${formatCurrency(data.totalExpenses)}</div>
                        <div class="stat-change negative">-${formatCurrency(data.monthlyExpenses)} this month</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Transactions</div>
                        <div class="stat-value">${data.totalTransactions}</div>
                        <div class="stat-change">${data.monthlyTransactions} this month</div>
                    </div>
                ` + "`" + `;
            } catch (error) {
                document.getElementById('statsGrid').innerHTML = 
                    '<div class="error">Failed to load statistics: ' + error.message + '</div>';
            }
        }

        // Load transactions
        async function loadTransactions() {
            try {
                const response = await fetch('/api/transactions?limit=10');
                const data = await response.json();
                
                if (!response.ok) throw new Error(data.error || 'Failed to load transactions');
                
                const container = document.getElementById('transactionsContainer');
                
                if (data.transactions.length === 0) {
                    container.innerHTML = '<p>No transactions yet.</p>';
                    return;
                }
                
                const tableRows = data.transactions.map(tx => ` + "`" + `
                    <tr>
                        <td>${formatDate(tx.created_at)}</td>
                        <td>${tx.category}</td>
                        <td>${tx.description || '-'}</td>
                        <td class="amount ${tx.type}">${tx.type === 'income' ? '+' : '-'}${formatCurrency(Math.abs(tx.amount))}</td>
                    </tr>
                ` + "`" + `).join('');
                
                container.innerHTML = ` + "`" + `
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
            } catch (error) {
                document.getElementById('transactionsContainer').innerHTML = 
                    '<div class="error">Failed to load transactions: ' + error.message + '</div>';
            }
        }

        // Load data on page load
        loadStats();
        loadTransactions();
        
        // Refresh data every 30 seconds
        setInterval(() => {
            loadStats();
            loadTransactions();
        }, 30000);
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
		User *model.User
	}{
		User: user,
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// handleAPIStats returns user statistics
func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all transactions for the user
	transactions, err := s.repositories.Transactions.GetUserTransactions(user.TgID)
	if err != nil {
		s.sendJSONError(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Calculate statistics
	var totalIncome, totalExpenses, monthlyIncome, monthlyExpenses float64
	var monthlyTransactions int

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for _, tx := range transactions {
		if tx.Type == model.TypeIncome {
			totalIncome += tx.Amount
			if tx.CreatedAt.After(startOfMonth) {
				monthlyIncome += tx.Amount
			}
		} else {
			totalExpenses += tx.Amount
			if tx.CreatedAt.After(startOfMonth) {
				monthlyExpenses += tx.Amount
			}
		}

		if tx.CreatedAt.After(startOfMonth) {
			monthlyTransactions++
		}
	}

	balance := totalIncome - totalExpenses

	stats := map[string]interface{}{
		"balance":             balance,
		"totalIncome":         totalIncome,
		"totalExpenses":       totalExpenses,
		"monthlyIncome":       monthlyIncome,
		"monthlyExpenses":     monthlyExpenses,
		"totalTransactions":   len(transactions),
		"monthlyTransactions": monthlyTransactions,
	}

	s.sendJSONSuccess(w, stats)
}

// handleAPITransactions returns user transactions
func (s *Server) handleAPITransactions(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
		if limit > 100 {
			limit = 100
		}
	}

	// Get transactions
	transactions, err := s.repositories.Transactions.GetUserTransactions(user.TgID)
	if err != nil {
		s.sendJSONError(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Sort by date (newest first) and limit
	if len(transactions) > limit {
		transactions = transactions[:limit]
	}

	response := map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	}

	s.sendJSONSuccess(w, response)
}
