const basePath = "/web";
const loginForm = document.getElementById("loginForm");
const verifyForm = document.getElementById("verifyForm");
const loginSection = document.getElementById("loginSection");
const verifySection = document.getElementById("verifySection");
const messageDiv = document.getElementById("message");
const verifyHint = document.getElementById("verifyHint");
const passkeyTabButton = document.getElementById("passkeyTabButton");
const passkeyEmailInput = document.getElementById("passkeyEmail");
const emailInput = document.getElementById("email");
const submitBtn = document.getElementById("submitBtn");

let activeLoginMethod = "telegram"; // Track which login method is active
let passkeySupported = false;

function showMessage(text, type) {
  messageDiv.className = "message " + type;
  messageDiv.textContent = text;
}

// Save credentials to sessionStorage
function saveCredentials(username, email, passkeyEmail) {
  if (username) {
    sessionStorage.setItem("cashout_telegram_username", username);
  }
  if (email) {
    sessionStorage.setItem("cashout_email", email);
  }
  if (passkeyEmail) {
    sessionStorage.setItem("cashout_passkey_email", passkeyEmail);
  }
}

// Load persisted credentials from sessionStorage
function loadPersistedCredentials() {
  const savedUsername = sessionStorage.getItem("cashout_telegram_username");
  const savedEmail = sessionStorage.getItem("cashout_email");
  const savedPasskeyEmail = sessionStorage.getItem("cashout_passkey_email");

  if (savedUsername) {
    document.getElementById("username").value = savedUsername;
  }
  if (savedEmail) {
    emailInput.value = savedEmail;
  }
  if (savedPasskeyEmail) {
    passkeyEmailInput.value = savedPasskeyEmail;
  }
}

// Check WebAuthn support on load
window.addEventListener("DOMContentLoaded", async () => {
  passkeySupported = WebAuthnClient.isSupported();

  if (passkeySupported) {
    const platformAvailable =
      await WebAuthnClient.isPlatformAuthenticatorAvailable();
    if (platformAvailable) {
      passkeyTabButton.style.display = "block";
    }
  }

  // Load persisted credentials
  loadPersistedCredentials();
});

// Update submit button text based on active tab
function updateSubmitButtonText() {
  if (activeLoginMethod === "passkey") {
    submitBtn.textContent = "Sign in with Passkey";
  } else {
    submitBtn.textContent = "Send Login Code";
  }
}

// Tab switching functionality
document.querySelectorAll(".tab-button").forEach((button) => {
  button.addEventListener("click", () => {
    const tabName = button.getAttribute("data-tab");

    // Update active tab button
    document
      .querySelectorAll(".tab-button")
      .forEach((btn) => btn.classList.remove("active"));
    button.classList.add("active");

    // Update active tab content
    document
      .querySelectorAll(".tab-content")
      .forEach((content) => content.classList.remove("active"));
    document.getElementById(tabName + "Tab").classList.add("active");

    // Track active method
    activeLoginMethod = tabName;

    // Update submit button text
    updateSubmitButtonText();

    // Clear any previous messages
    messageDiv.textContent = "";
    messageDiv.className = "message";
  });
});

loginForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  submitBtn.disabled = true;

  try {
    // Handle passkey authentication separately
    if (activeLoginMethod === "passkey") {
      submitBtn.textContent = "Authenticating...";

      const email = passkeyEmailInput.value.trim();
      if (!email) {
        showMessage("Please enter your email address", "error");
        return;
      }

      const result = await WebAuthnClient.authenticate(email);

      // Save passkey email
      saveCredentials(null, null, email);

      showMessage("Login successful! Redirecting...", "success");
      setTimeout(() => {
        window.location.href = result.redirect || basePath + "/dashboard";
      }, 1000);
      return;
    }

    // Handle email/telegram login
    submitBtn.textContent = "Sending...";

    let requestBody = {};

    if (activeLoginMethod === "telegram") {
      const username = document
        .getElementById("username")
        .value.replace("@", "");
      if (!username) {
        showMessage("Please enter your Telegram username", "error");
        return;
      }
      requestBody = { username };
    } else {
      const email = document.getElementById("email").value;
      if (!email) {
        showMessage("Please enter your email address", "error");
        return;
      }
      requestBody = { email };
    }

    const response = await fetch(basePath + "/auth/request", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(requestBody),
    });

    const data = await response.json();

    if (response.ok) {
      // Save credentials on successful code send
      if (activeLoginMethod === "telegram") {
        saveCredentials(requestBody.username, null, null);
      } else {
        saveCredentials(null, requestBody.email, null);
      }

      loginSection.style.display = "none";
      verifySection.style.display = "block";

      // Update verification hint based on login method
      if (activeLoginMethod === "email") {
        verifyHint.textContent = "Check your email for the verification code";
      } else {
        verifyHint.textContent =
          "Check your Telegram for the verification code";
      }

      showMessage(data.message || "Verification code sent!", "success");
    } else {
      showMessage(data.error || "Failed to send code", "error");
    }
  } catch (error) {
    if (activeLoginMethod === "passkey") {
      showMessage(error.message || "Passkey authentication failed", "error");
    } else {
      showMessage("Network error. Please try again.", "error");
    }
  } finally {
    submitBtn.disabled = false;
    updateSubmitButtonText();
  }
});

verifyForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const code = document.getElementById("code").value;
  const verifyBtn = document.getElementById("verifyBtn");

  verifyBtn.disabled = true;
  verifyBtn.textContent = "Verifying...";

  try {
    const response = await fetch(basePath + "/auth/verify", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code }),
    });

    const data = await response.json();

    if (response.ok) {
      showMessage("Login successful! Redirecting...", "success");
      setTimeout(() => {
        window.location.href = basePath + "/dashboard";
      }, 1000);
    } else {
      showMessage(data.error || "Invalid code", "error");
    }
  } catch (error) {
    showMessage("Network error. Please try again.", "error");
  } finally {
    verifyBtn.disabled = false;
    verifyBtn.textContent = "Verify";
  }
});
