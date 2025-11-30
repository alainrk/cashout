// Passkey Management for Dashboard
const passkeySection = document.getElementById('passkeySection');
const registerPasskeyBtn = document.getElementById('registerPasskeyBtn');
const passkeyList = document.getElementById('passkeyList');

let passkeys = [];

// Check WebAuthn support and load passkeys on page load
document.addEventListener('DOMContentLoaded', async () => {
    const isSupported = WebAuthnClient.isSupported();

    if (!isSupported) {
        console.log('WebAuthn not supported');
        return;
    }

    // Check if platform authenticator is available
    const platformAvailable = await WebAuthnClient.isPlatformAuthenticatorAvailable();
    if (!platformAvailable) {
        console.log('Platform authenticator not available');
        return;
    }

    // Show the passkey section
    passkeySection.style.display = 'block';

    // Load existing passkeys
    await loadPasskeys();
});

// Load user's passkeys
async function loadPasskeys() {
    try {
        const response = await fetch('/web/api/passkey/list', {
            method: 'GET',
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to load passkeys');
        }

        const data = await response.json();
        passkeys = data.passkeys || [];
        renderPasskeys();
    } catch (error) {
        console.error('Error loading passkeys:', error);
        passkeyList.innerHTML = '<p class="error-text">Failed to load passkeys</p>';
    }
}

// Render passkeys list
function renderPasskeys() {
    if (passkeys.length === 0) {
        passkeyList.innerHTML = '<p class="no-passkeys">No passkeys registered yet. Register one to enable fast, secure login!</p>';
        return;
    }

    const html = passkeys.map(passkey => {
        const createdDate = new Date(passkey.createdAt).toLocaleDateString();
        const lastUsed = passkey.lastUsedAt
            ? new Date(passkey.lastUsedAt).toLocaleDateString()
            : 'Never';

        const deviceName = passkey.name || 'Unnamed Device';
        const credId = escapeHtml(passkey.id);

        return `
            <div class="passkey-item" data-id="${credId}">
                <div class="passkey-info">
                    <div class="passkey-name">🔐 ${escapeHtml(deviceName)}</div>
                    <div class="passkey-meta">
                        <span>Created: ${createdDate}</span>
                        <span>Last used: ${lastUsed}</span>
                    </div>
                </div>
                <button class="delete-passkey-btn" onclick="deletePasskey('${credId}')">
                    Delete
                </button>
            </div>
        `;
    }).join('');

    passkeyList.innerHTML = html;
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Register new passkey
if (registerPasskeyBtn) {
    registerPasskeyBtn.addEventListener('click', async () => {
        registerPasskeyBtn.disabled = true;
        const originalText = registerPasskeyBtn.innerHTML;
        registerPasskeyBtn.innerHTML = '<span>⏳</span> Registering...';

        try {
            // Prompt for device name
            const deviceName = prompt('Enter a name for this device (e.g., "My Laptop", "iPhone"):');
            if (!deviceName || deviceName.trim() === '') {
                throw new Error('Device name is required');
            }

            // Call the registration flow from webauthn.js
            // We need to use the WebAuthnClient.register method but with our API endpoints
            await registerPasskeyFlow(deviceName.trim());

            // Reload passkeys
            await loadPasskeys();

            // Show success message
            showMessage('Passkey registered successfully!', 'success');
        } catch (error) {
            console.error('Error registering passkey:', error);
            showMessage(error.message || 'Failed to register passkey', 'error');
        } finally {
            registerPasskeyBtn.disabled = false;
            registerPasskeyBtn.innerHTML = originalText;
        }
    });
}

// Registration flow
async function registerPasskeyFlow(deviceName) {
    console.log('[Passkey] Starting registration flow for device:', deviceName);

    // Step 1: Begin registration
    console.log('[Passkey] Step 1: Calling begin-register...');
    const beginResponse = await fetch('/web/api/passkey/begin-register', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include'
    });

    if (!beginResponse.ok) {
        const error = await beginResponse.json();
        console.error('[Passkey] Begin registration failed:', error);
        throw new Error(error.error || 'Failed to begin registration');
    }

    const beginData = await beginResponse.json();
    console.log('[Passkey] Begin registration successful, received options');

    // Extract options from the success response
    const credentialOptions = beginData.options;

    // Step 2: Create credential using WebAuthn API
    console.log('[Passkey] Step 2: Creating credential with browser API...');
    const credential = await WebAuthnClient.createCredential(credentialOptions);
    console.log('[Passkey] Credential created:', credential);

    // Step 3: Finish registration
    console.log('[Passkey] Step 3: Calling finish-register...');

    const finishResponse = await fetch('/web/api/passkey/finish-register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-Credential-Name': deviceName
        },
        credentials: 'include',
        body: JSON.stringify(credential)
    });

    if (!finishResponse.ok) {
        const errorText = await finishResponse.text();
        console.error('[Passkey] Finish registration failed:', errorText);
        try {
            const error = JSON.parse(errorText);
            throw new Error(error.error || 'Failed to finish registration');
        } catch (e) {
            throw new Error('Failed to finish registration: ' + errorText);
        }
    }

    const result = await finishResponse.json();
    console.log('[Passkey] Registration completed successfully:', result);
    return result;
}

// Delete passkey
async function deletePasskey(credentialId) {
    if (!confirm('Are you sure you want to delete this passkey? You will no longer be able to use it to sign in.')) {
        return;
    }

    try {
        const response = await fetch('/web/api/passkey/delete', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({ credentialId: credentialId })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to delete passkey');
        }

        // Reload passkeys
        await loadPasskeys();

        showMessage('Passkey deleted successfully', 'success');
    } catch (error) {
        console.error('Error deleting passkey:', error);
        showMessage(error.message || 'Failed to delete passkey', 'error');
    }
}

// Show message (reuse from dashboard.js if available, otherwise create simple version)
function showMessage(text, type) {
    // Try to use the existing message div from the transaction form
    const existingMessageDiv = document.getElementById('txMessage');
    if (existingMessageDiv) {
        existingMessageDiv.className = 'message ' + type;
        existingMessageDiv.textContent = text;
        setTimeout(() => {
            existingMessageDiv.textContent = '';
            existingMessageDiv.className = 'message';
        }, 5000);
    } else {
        // Fallback: show alert
        alert(text);
    }
}
