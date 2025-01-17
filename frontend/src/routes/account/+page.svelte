<script lang="ts">
    import { user } from '../../lib/stores';
    import { onMount } from 'svelte';

    let currentPassword = '';
    let newPassword = '';
    let confirmPassword = '';
    let newUsername = $user?.username || '';
    let newEmail = $user?.email || '';
    let message = { text: '', type: '' };

    async function updateProfile() {
        const token = localStorage.getItem('token');
        if (!token) return;

        try {
            const response = await fetch('http://localhost:8080/api/user/profile', {
                method: 'PUT',
                headers: {
                    'Authorization': token,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    username: newUsername,
                    email: newEmail
                })
            });

            if (response.ok) {
                const updatedUser = await response.json();
                user.set(updatedUser);
                message = { text: 'Profil mis à jour avec succès', type: 'success' };
            } else {
                message = { text: 'Erreur lors de la mise à jour du profil', type: 'error' };
            }
        } catch (error) {
            message = { text: 'Erreur de connexion', type: 'error' };
        }
    }

    async function updatePassword() {
        if (newPassword !== confirmPassword) {
            message = { text: 'Les mots de passe ne correspondent pas', type: 'error' };
            return;
        }

        const token = localStorage.getItem('token');
        if (!token) return;

        try {
            const response = await fetch('http://localhost:8080/api/user/password', {
                method: 'PUT',
                headers: {
                    'Authorization': token,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    currentPassword,
                    newPassword
                })
            });

            if (response.ok) {
                message = { text: 'Mot de passe mis à jour avec succès', type: 'success' };
                currentPassword = '';
                newPassword = '';
                confirmPassword = '';
            } else {
                message = { text: 'Erreur lors de la mise à jour du mot de passe', type: 'error' };
            }
        } catch (error) {
            message = { text: 'Erreur de connexion', type: 'error' };
        }
    }
</script>

<div class="account-container">
    <h1>Paramètres du compte</h1>

    {#if message.text}
        <div class="message" class:error={message.type === 'error'} class:success={message.type === 'success'}>
            {message.text}
        </div>
    {/if}

    <div class="settings-card">
        <h2>Informations du profil</h2>
        <div class="form-group">
            <label for="username">Pseudonyme</label>
            <input 
                type="text" 
                id="username" 
                bind:value={newUsername}
                placeholder="Nouveau pseudonyme"
            />
        </div>

        <div class="form-group">
            <label for="email">Email</label>
            <input 
                type="email" 
                id="email" 
                bind:value={newEmail}
                placeholder="Nouvel email"
            />
        </div>

        <button class="update-button" on:click={updateProfile}>
            Mettre à jour le profil
        </button>
    </div>

    <div class="settings-card">
        <h2>Changer le mot de passe</h2>
        <div class="form-group">
            <label for="current-password">Mot de passe actuel</label>
            <input 
                type="password" 
                id="current-password" 
                bind:value={currentPassword}
                placeholder="Mot de passe actuel"
            />
        </div>

        <div class="form-group">
            <label for="new-password">Nouveau mot de passe</label>
            <input 
                type="password" 
                id="new-password" 
                bind:value={newPassword}
                placeholder="Nouveau mot de passe"
            />
        </div>

        <div class="form-group">
            <label for="confirm-password">Confirmer le mot de passe</label>
            <input 
                type="password" 
                id="confirm-password" 
                bind:value={confirmPassword}
                placeholder="Confirmer le mot de passe"
            />
        </div>

        <button class="update-button" on:click={updatePassword}>
            Mettre à jour le mot de passe
        </button>
    </div>
</div>

<style>
    .account-container {
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem;
    }

    h1 {
        color: #2c3e50;
        text-align: center;
        margin-bottom: 2rem;
    }

    h2 {
        color: #2c3e50;
        margin-bottom: 1.5rem;
    }

    .settings-card {
        background: white;
        padding: 2rem;
        border-radius: 1rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        margin-bottom: 2rem;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
        color: #2c3e50;
        font-weight: 500;
    }

    input {
        width: 100%;
        padding: 0.75rem;
        border: 2px solid #3498db;
        border-radius: 0.5rem;
        font-size: 1rem;
        transition: border-color 0.3s ease;
    }

    input:focus {
        outline: none;
        border-color: #2980b9;
    }

    .update-button {
        background: #3498db;
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 0.5rem;
        font-size: 1rem;
        cursor: pointer;
        transition: background 0.3s ease;
        width: 100%;
    }

    .update-button:hover {
        background: #2980b9;
    }

    .message {
        padding: 1rem;
        border-radius: 0.5rem;
        margin-bottom: 1rem;
        text-align: center;
    }

    .error {
        background: #fee2e2;
        color: #dc2626;
    }

    .success {
        background: #dcfce7;
        color: #16a34a;
    }
</style> 