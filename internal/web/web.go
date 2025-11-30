// Package web implements the web server functionalities for the "static" dashboard.
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

	// WebAuthn/Passkey routes (public)
	mux.HandleFunc(basePath+"/auth/passkey/check", s.rateLimit(s.handlePasskeyCheck))
	mux.HandleFunc(basePath+"/auth/passkey/begin-login", s.rateLimit(s.handlePasskeyBeginLogin))
	mux.HandleFunc(basePath+"/auth/passkey/finish-login", s.rateLimit(s.handlePasskeyFinishLogin))

	// Dashboard routes (protected)
	mux.HandleFunc(basePath+"/dashboard", s.requireAuth(s.handleDashboard))
	mux.HandleFunc(basePath+"/api/transactions", s.requireAuth(s.handleAPITransactions))
	mux.HandleFunc(basePath+"/api/transactions/create", s.requireAuth(s.handleAPICreateTransaction))
	mux.HandleFunc(basePath+"/api/categories", s.requireAuth(s.handleAPICategories))
	mux.HandleFunc(basePath+"/api/stats", s.requireAuth(s.handleAPIStats))

	// WebAuthn/Passkey management (protected)
	mux.HandleFunc(basePath+"/api/passkey/begin-register", s.requireAuth(s.handlePasskeyBeginRegister))
	mux.HandleFunc(basePath+"/api/passkey/finish-register", s.requireAuth(s.handlePasskeyFinishRegister))
	mux.HandleFunc(basePath+"/api/passkey/list", s.requireAuth(s.handlePasskeyList))
	mux.HandleFunc(basePath+"/api/passkey/delete", s.requireAuth(s.handlePasskeyDelete))

	return s.loggingMiddleware(mux)
}
