package web

import (
	"cashout/internal/ai"
	"cashout/internal/client"
	"cashout/internal/model"
	"cashout/internal/repository"
	"net/http"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/sirupsen/logrus"
)

type Repositories struct {
	Users        repository.Users
	Transactions repository.Transactions
	Auth         repository.Auth
}

type Server struct {
	logger       *logrus.Logger
	repositories Repositories
	bot          *gotgbot.Bot
	llm          ai.LLM
}

func NewServer(logger *logrus.Logger, repos Repositories, bot *gotgbot.Bot, llm ai.LLM) *Server {
	return &Server{
		logger:       logger,
		repositories: repos,
		bot:          bot,
		llm:          llm,
	}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Auth routes
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/auth/request", s.handleAuthRequest)
	mux.HandleFunc("/auth/verify", s.handleAuthVerify)
	mux.HandleFunc("/logout", s.handleLogout)

	// Dashboard routes (protected)
	mux.HandleFunc("/dashboard", s.requireAuth(s.handleDashboard))
	mux.HandleFunc("/api/transactions", s.requireAuth(s.handleAPITransactions))
	mux.HandleFunc("/api/stats", s.requireAuth(s.handleAPIStats))

	return s.loggingMiddleware(mux)
}

// Middleware to log requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware to require authentication
func (s *Server) requireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.getSession(r)
		if err != nil || session == nil || !session.IsValid() {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user to context
		ctx := r.Context()
		ctx = client.SetUserInContext(ctx, session.User)
		handler(w, r.WithContext(ctx))
	}
}

// Helper to get session from cookie
func (s *Server) getSession(r *http.Request) (*model.WebSession, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	session, err := s.repositories.Auth.GetWebSession(cookie.Value)
	if err != nil {
		return nil, err
	}

	// Load user data
	user, err := s.repositories.Users.GetByTgID(session.TgID)
	if err != nil {
		return nil, err
	}
	session.User = &user

	return session, nil
}
