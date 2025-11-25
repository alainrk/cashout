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
