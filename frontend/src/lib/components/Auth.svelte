<script lang="ts">
    import { user } from '../stores';
    
    let formData = {
        email: '',
        password: '',
        username: ''
    };
    let isLogin = true;
    let error = '';

    async function handleSubmit() {
        error = '';
        console.log("Submitting:", { ...formData, isLogin });
        
        const endpoint = isLogin ? '/api/login' : '/api/register';
        try {
            const response = await fetch(`http://localhost:8080${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });
            
            const responseData = await response.json();
            
            if (response.ok) {
                user.set(responseData.user);
                localStorage.setItem('token', responseData.token);
                window.location.href = '/facilitator';
            } else {
                error = responseData.error || 'Une erreur est survenue';
                console.error('Erreur:', error);
            }
        } catch (error) {
            console.error('Erreur:', error);
            error = 'Erreur de connexion au serveur';
        }
    }
</script>

<div class="auth-container">
    <h2>{isLogin ? 'Connexion' : 'Inscription'}</h2>
    {#if error}
        <div class="error">{error}</div>
    {/if}
    <form on:submit|preventDefault={handleSubmit}>
        <input
            type="email"
            bind:value={formData.email}
            placeholder="Email"
            autocomplete="email"
            required
        />
        <input
            type="password"
            bind:value={formData.password}
            placeholder="Mot de passe"
            autocomplete="current-password"
            required
        />
        {#if !isLogin}
            <div class="form-group">
                <label for="username">Pseudonyme</label>
                <input
                    type="text"
                    id="username"
                    bind:value={formData.username}
                    required
                    minlength="3"
                    maxlength="50"
                />
            </div>
        {/if}
        <button type="submit">
            {isLogin ? 'Se connecter' : 'S\'inscrire'}
        </button>
    </form>
    <button class="switch-button" on:click={() => isLogin = !isLogin}>
        {isLogin ? 'Créer un compte' : 'Déjà un compte ?'}
    </button>
</div>

<style>
    .auth-container {
        max-width: 400px;
        margin: 2rem auto;
        padding: 2rem;
        border-radius: 0.5rem;
        background: white;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    input {
        padding: 0.5rem;
        border: 1px solid #ccc;
        border-radius: 0.25rem;
    }

    button {
        padding: 0.5rem;
        background: var(--primary-color, #4CAF50);
        color: white;
        border: none;
        border-radius: 0.25rem;
        cursor: pointer;
    }

    .switch-button {
        margin-top: 1rem;
        background: none;
        color: var(--primary-color, #4CAF50);
        text-decoration: underline;
    }

    .error {
        color: red;
        margin-bottom: 1rem;
        text-align: center;
    }
</style> 