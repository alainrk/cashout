// WebAuthn utility functions for passkey authentication

/**
 * Base64URL encoding/decoding utilities
 */
const base64url = {
    encode: (buffer) => {
        // Convert buffer to binary string without using spread operator (avoids stack overflow)
        const bytes = new Uint8Array(buffer);
        let binary = '';
        for (let i = 0; i < bytes.length; i++) {
            binary += String.fromCharCode(bytes[i]);
        }
        const base64 = btoa(binary);
        return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
    },

    decode: (base64url) => {
        const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
        const binary = atob(base64);
        return Uint8Array.from(binary, c => c.charCodeAt(0));
    }
};

/**
 * Convert server options to WebAuthn API format
 */
function preparePublicKeyOptions(options) {
    // Convert challenge
    if (options.challenge) {
        options.challenge = base64url.decode(options.challenge);
    }

    // Convert user.id for registration
    if (options.user && options.user.id) {
        options.user.id = base64url.decode(options.user.id);
    }

    // Convert allowCredentials for authentication
    if (options.allowCredentials) {
        options.allowCredentials = options.allowCredentials.map(cred => ({
            ...cred,
            id: base64url.decode(cred.id)
        }));
    }

    // Convert excludeCredentials for registration
    if (options.excludeCredentials) {
        options.excludeCredentials = options.excludeCredentials.map(cred => ({
            ...cred,
            id: base64url.decode(cred.id)
        }));
    }

    return options;
}

/**
 * Convert WebAuthn credential to server format
 */
function prepareCredentialForServer(credential) {
    const response = credential.response;

    const result = {
        id: credential.id,
        rawId: base64url.encode(credential.rawId),
        type: credential.type,
        response: {}
    };

    // Common response fields
    if (response.clientDataJSON) {
        result.response.clientDataJSON = base64url.encode(response.clientDataJSON);
    }

    // Registration-specific fields
    if (response.attestationObject) {
        result.response.attestationObject = base64url.encode(response.attestationObject);
    }

    // Authentication-specific fields
    if (response.authenticatorData) {
        result.response.authenticatorData = base64url.encode(response.authenticatorData);
    }
    if (response.signature) {
        result.response.signature = base64url.encode(response.signature);
    }
    if (response.userHandle) {
        result.response.userHandle = base64url.encode(response.userHandle);
    }

    return result;
}

/**
 * Check if WebAuthn is supported
 */
function isWebAuthnSupported() {
    return window.PublicKeyCredential !== undefined &&
           navigator.credentials !== undefined;
}

/**
 * Check if platform authenticator is available
 */
async function isPlatformAuthenticatorAvailable() {
    if (!isWebAuthnSupported()) {
        return false;
    }

    try {
        return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
    } catch (err) {
        console.error('Error checking platform authenticator:', err);
        return false;
    }
}

/**
 * Authenticate with passkey
 */
async function authenticateWithPasskey(email) {
    try {
        // Request authentication options from server
        const beginResponse = await fetch('/web/auth/passkey/begin-login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({email})
        });

        if (!beginResponse.ok) {
            const error = await beginResponse.json();
            throw new Error(error.error || 'Failed to start authentication');
        }

        const {options} = await beginResponse.json();

        // Prepare options for WebAuthn API
        const publicKeyOptions = preparePublicKeyOptions(options.publicKey);

        // Get credential
        const credential = await navigator.credentials.get({
            publicKey: publicKeyOptions
        });

        if (!credential) {
            throw new Error('Failed to get credential');
        }

        // Convert credential to server format
        const assertionData = prepareCredentialForServer(credential);

        // Send assertion to server
        const finishResponse = await fetch('/web/auth/passkey/finish-login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(assertionData)
        });

        if (!finishResponse.ok) {
            const error = await finishResponse.json();
            throw new Error(error.error || 'Authentication failed');
        }

        return await finishResponse.json();

    } catch (err) {
        console.error('Passkey authentication error:', err);
        throw err;
    }
}

/**
 * Check if user has passkey registered
 */
async function checkPasskeyAvailable(email) {
    try {
        const response = await fetch('/web/auth/passkey/check', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({email})
        });

        if (!response.ok) {
            return false;
        }

        const {hasPasskey} = await response.json();
        return hasPasskey;

    } catch (err) {
        console.error('Error checking passkey:', err);
        return false;
    }
}

/**
 * Create credential from server options (for manual registration flow)
 */
async function createCredential(credentialOptions) {
    try {
        // Prepare options for WebAuthn API
        const publicKeyOptions = preparePublicKeyOptions(credentialOptions.publicKey);

        // Create credential
        const credential = await navigator.credentials.create({
            publicKey: publicKeyOptions
        });

        if (!credential) {
            throw new Error('Failed to create credential');
        }

        // Convert to server format
        return prepareCredentialForServer(credential);

    } catch (err) {
        console.error('Credential creation error:', err);
        throw err;
    }
}

// Export functions
window.WebAuthnClient = {
    isSupported: isWebAuthnSupported,
    isPlatformAuthenticatorAvailable,
    authenticate: authenticateWithPasskey,
    checkAvailable: checkPasskeyAvailable,
    createCredential
};
