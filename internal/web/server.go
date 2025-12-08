package web

import (
	"net/http"
	"sync"
	"time"

	"cashout/internal/ai"
	"cashout/internal/client"
	"cashout/internal/email"
	"cashout/internal/model"
	"cashout/internal/repository"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Repositories struct {
	Users        repository.Users
	Transactions repository.Transactions
	Auth         repository.Auth
	WebAuthn     *repository.WebAuthn
}

type Server struct {
	logger         *logrus.Logger
	repositories   Repositories
	bot            *gotgbot.Bot
	llm            ai.LLM
	loginLimiter   map[string]*rate.Limiter
	loginLimiterMu sync.Mutex
	emailService   *email.EmailService
}

func NewServer(logger *logrus.Logger, repos Repositories, bot *gotgbot.Bot, llm ai.LLM, emailService *email.EmailService) *Server {
	return &Server{
		logger:         logger,
		repositories:   repos,
		bot:            bot,
		llm:            llm,
		loginLimiter:   make(map[string]*rate.Limiter),
		loginLimiterMu: sync.Mutex{},
		emailService:   emailService,
	}
}

func (s *Server) Router() http.Handler {
	return Router(s)
}

func (s *Server) getLimiter(key string) *rate.Limiter {
	s.loginLimiterMu.Lock()
	defer s.loginLimiterMu.Unlock()

	limiter, exists := s.loginLimiter[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute), 5) // 5 requests per minute
		s.loginLimiter[key] = limiter
	}

	return limiter
}

// Middleware to log requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// securityHeadersMiddleware adds security headers to all responses
func (s *Server) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent XSS and restrict resource loading
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Control referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Disable unnecessary browser features
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Legacy XSS protection (for older browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}

func (s *Server) rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter := s.getLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			s.sendJSONError(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Middleware to require authentication
func (s *Server) requireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.getSession(r)
		if err != nil || session == nil || !session.IsValid() {
			http.Redirect(w, r, basePath+"/login", http.StatusSeeOther)
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
