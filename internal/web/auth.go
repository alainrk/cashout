package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
)

// handleHome redirects to dashboard if authenticated, otherwise to login
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	session, _ := s.getSession(r)
	if session != nil && session.IsValid() {
		http.Redirect(w, r, basePath+"/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, basePath+"/login", http.StatusSeeOther)
}

// handleLogin shows the login page
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	session, _ := s.getSession(r)
	if session != nil && session.IsValid() {
		http.Redirect(w, r, basePath+"/dashboard", http.StatusSeeOther)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Cashout - Login</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .login-container {
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            width: 100%;
            max-width: 400px;
        }
        h1 {
            margin: 0 0 2rem 0;
            text-align: center;
            color: #333;
        }
        .form-group {
            margin-bottom: 1.5rem;
        }
        label {
            display: block;
            margin-bottom: 0.5rem;
            color: #555;
            font-weight: 500;
        }
        input {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
            box-sizing: border-box;
        }
        input:focus {
            outline: none;
            border-color: #0088cc;
        }
        button {
            width: 100%;
            padding: 0.75rem;
            background: #0088cc;
            color: white;
            border: none;
            border-radius: 4px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
        }
        button:hover {
            background: #006ba1;
        }
        button:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        .message {
            margin-top: 1rem;
            padding: 0.75rem;
            border-radius: 4px;
            text-align: center;
        }
        .error {
            background: #fee;
            color: #c33;
            border: 1px solid #fcc;
        }
        .success {
            background: #efe;
            color: #3c3;
            border: 1px solid #cfc;
        }
        .info {
            background: #e6f2ff;
            color: #0066cc;
            border: 1px solid #b3d9ff;
        }
        #verifySection {
            display: none;
        }
        .telegram-hint {
            font-size: 0.875rem;
            color: #666;
            margin-top: 0.5rem;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>Cashout Login</h1>
        
        <div id="loginSection">
            <form id="loginForm">
                <div class="form-group">
                    <label for="username">Telegram Username</label>
                    <input type="text" id="username" name="username" placeholder="@username" required>
                    <div class="telegram-hint">Enter your Telegram username (without @)</div>
                </div>
                <button type="submit" id="submitBtn">Send Login Code</button>
            </form>
        </div>

        <div id="verifySection">
            <form id="verifyForm">
                <div class="form-group">
                    <label for="code">Verification Code</label>
                    <input type="text" id="code" name="code" placeholder="Enter 6-digit code" maxlength="6" required>
                    <div class="telegram-hint">Check your Telegram for the verification code</div>
                </div>
                <button type="submit" id="verifyBtn">Verify</button>
            </form>
        </div>

        <div id="message"></div>
    </div>

    <script>
        const basePath = '/web';
        const loginForm = document.getElementById('loginForm');
        const verifyForm = document.getElementById('verifyForm');
        const loginSection = document.getElementById('loginSection');
        const verifySection = document.getElementById('verifySection');
        const messageDiv = document.getElementById('message');

        function showMessage(text, type) {
            messageDiv.className = 'message ' + type;
            messageDiv.textContent = text;
        }

        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('username').value.replace('@', '');
            const submitBtn = document.getElementById('submitBtn');
            
            submitBtn.disabled = true;
            submitBtn.textContent = 'Sending...';
            
            try {
                const response = await fetch(basePath+'/auth/request', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({username})
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    loginSection.style.display = 'none';
                    verifySection.style.display = 'block';
                    showMessage('Verification code sent to your Telegram!', 'success');
                } else {
                    showMessage(data.error || 'Failed to send code', 'error');
                }
            } catch (error) {
                showMessage('Network error. Please try again.', 'error');
            } finally {
                submitBtn.disabled = false;
                submitBtn.textContent = 'Send Login Code';
            }
        });

        verifyForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const code = document.getElementById('code').value;
            const verifyBtn = document.getElementById('verifyBtn');
            
            verifyBtn.disabled = true;
            verifyBtn.textContent = 'Verifying...';
            
            try {
                const response = await fetch(basePath+'/auth/verify', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({code})
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    showMessage('Login successful! Redirecting...', 'success');
                    setTimeout(() => {
                        window.location.href = basePath+'/dashboard';
                    }, 1000);
                } else {
                    showMessage(data.error || 'Invalid code', 'error');
                }
            } catch (error) {
                showMessage('Network error. Please try again.', 'error');
            } finally {
                verifyBtn.disabled = false;
                verifyBtn.textContent = 'Verify';
            }
        });
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, tmpl)
	if err != nil {
		s.logger.Errorf("Failed to send login page: %v", err)
	}
}

// handleAuthRequest handles the initial auth request
func (s *Server) handleAuthRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Clean username
	username := strings.TrimSpace(strings.TrimPrefix(req.Username, "@"))
	if username == "" {
		s.sendJSONError(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get user by username
	user, exists, err := s.repositories.Users.GetByUsername(username)
	if err != nil || !exists {
		s.sendJSONError(w, "Invalid username or credentials", http.StatusNotFound)
		return
	}

	// Create auth token
	authToken, err := s.repositories.Auth.CreateAuthToken(user.TgID)
	if err != nil {
		s.logger.Errorf("Failed to create auth token: %v", err)
		s.sendJSONError(w, "Failed to create auth token", http.StatusInternalServerError)
		return
	}

	// Send code via Telegram
	message := fmt.Sprintf("🔐 Your Cashout login code is: <b>%s</b>\n\nThis code will expire in 5 minutes.", authToken.Token)
	_, err = s.bot.SendMessage(user.TgID, message, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		s.logger.Errorf("Failed to send auth code: %v", err)
		s.sendJSONError(w, "Failed to send code. Please make sure the bot is not blocked.", http.StatusInternalServerError)
		return
	}

	s.sendJSONSuccess(w, map[string]interface{}{
		"message": "Code sent successfully",
	})
}

// handleAuthVerify verifies the auth code
func (s *Server) handleAuthVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	code := strings.TrimSpace(strings.ToUpper(req.Code))
	if code == "" {
		s.sendJSONError(w, "Code is required", http.StatusBadRequest)
		return
	}

	// Verify auth token
	user, err := s.repositories.Auth.VerifyAuthToken(code)
	if err != nil {
		s.sendJSONError(w, "Invalid or expired code", http.StatusUnauthorized)
		return
	}

	// Create web session
	session, err := s.repositories.Auth.CreateWebSession(user.TgID)
	if err != nil {
		s.logger.Errorf("Failed to create session: %v", err)
		s.sendJSONError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

	s.sendJSONSuccess(w, map[string]interface{}{
		"message":  "Login successful",
		"redirect": basePath + "/dashboard",
	})
}

// handleLogout handles user logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Delete session from database
		err = errors.Join(err, s.repositories.Auth.DeleteWebSession(cookie.Value))
		if err != nil {
			s.logger.Errorf("Failed to delete session: %v", err)
		}
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, basePath+"/login", http.StatusSeeOther)
}

// Helper functions for JSON responses
func (s *Server) sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		s.logger.Errorf("Failed to send error response: %v", err)
	}
}

func (s *Server) sendJSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.logger.Errorf("Failed to send success response: %v", err)
	}
}
