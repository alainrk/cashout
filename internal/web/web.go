package web

import "net/http"

const (
	basePath = "/web"
)

func Router(s *Server) http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle(basePath+"/static/", http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./web/static"))))

	// Auth routes
	mux.HandleFunc(basePath+"/", s.handleHome)
	mux.HandleFunc(basePath+"/login", s.handleLogin)
	mux.HandleFunc(basePath+"/auth/request", s.handleAuthRequest)
	mux.HandleFunc(basePath+"/auth/verify", s.rateLimit(s.handleAuthVerify))
	mux.HandleFunc(basePath+"/logout", s.handleLogout)

	// Dashboard routes (protected)
	mux.HandleFunc(basePath+"/dashboard", s.requireAuth(s.handleDashboard))
	mux.HandleFunc(basePath+"/api/transactions", s.requireAuth(s.handleAPITransactions))
	mux.HandleFunc(basePath+"/api/stats", s.requireAuth(s.handleAPIStats))

	return s.loggingMiddleware(mux)
}
