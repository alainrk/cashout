// Package web implements the web server functionalities for the "static" dashboard.
package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	basePath = "/web"
)

func Router(s *Server) http.Handler {
	mux := http.NewServeMux()

	// Serve static files with directory listing disabled and null byte protection
	staticFS := noDirFileSystem{http.Dir("./web/static")}
	mux.Handle(basePath+"/static/", http.StripPrefix(basePath+"/static/", http.FileServer(staticFS)))

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
	mux.HandleFunc(basePath+"/api/transactions/delete", s.requireAuth(s.handleAPIDeleteTransaction))
	mux.HandleFunc(basePath+"/api/categories", s.requireAuth(s.handleAPICategories))
	mux.HandleFunc(basePath+"/api/stats", s.requireAuth(s.handleAPIStats))

	// WebAuthn/Passkey management (protected)
	mux.HandleFunc(basePath+"/api/passkey/begin-register", s.requireAuth(s.handlePasskeyBeginRegister))
	mux.HandleFunc(basePath+"/api/passkey/finish-register", s.requireAuth(s.handlePasskeyFinishRegister))
	mux.HandleFunc(basePath+"/api/passkey/list", s.requireAuth(s.handlePasskeyList))
	mux.HandleFunc(basePath+"/api/passkey/delete", s.requireAuth(s.handlePasskeyDelete))

	return s.securityHeadersMiddleware(s.loggingMiddleware(mux))
}

// noDirFileSystem wraps http.FileSystem to disable directory listing
// and protect against null byte injection
type noDirFileSystem struct {
	fs http.FileSystem
}

func (nfs noDirFileSystem) Open(path string) (http.File, error) {
	// Protect against null byte injection
	if strings.Contains(path, "\x00") {
		return nil, os.ErrNotExist
	}

	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	// Deny directory listing
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			_ = f.Close()
			return nil, os.ErrNotExist
		}
	}

	return f, nil
}
