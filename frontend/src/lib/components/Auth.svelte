<script lang="ts">
    import { user } from '../stores';
    
    let formData = {
        username: '',
        email: '',
        password: ''
    };
    let isLogin = true;
    let error = '';

    async function handleSubmit() {
        error = '';
        
        const endpoint = isLogin ? '/api/login' : '/api/register';
        try {
            const response = await fetch(`http://localhost:8080${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: formData.username,
                    password: formData.password,
                    ...(isLogin ? {} : { email: formData.email })
                }),
            });
            
            const responseData = await response.json();
            
            if (response.ok) {
                user.set(responseData.user);
                localStorage.setItem('token', responseData.token);
                window.location.href = '/explanations';
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
    <div class="auth-card">
        <h1>{isLogin ? 'Connexion' : 'Inscription'}</h1>
        
        {#if error}
            <div class="error">{error}</div>
        {/if}

        <form on:submit|preventDefault={handleSubmit}>
            <div class="form-group">
                <label for="username">Pseudonyme</label>
                <input 
                    type="text" 
                    id="username" 
                    bind:value={formData.username} 
                    required
                />
            </div>

            {#if !isLogin}
                <div class="form-group">
                    <label for="email">Email</label>
                    <input 
                        type="email" 
                        id="email" 
                        bind:value={formData.email} 
                        required
                    />
                </div>
            {/if}

            <div class="form-group">
                <label for="password">Mot de passe</label>
                <input 
                    type="password" 
                    id="password" 
                    bind:value={formData.password} 
                    required
                />
            </div>

            <button type="submit">
                {isLogin ? 'Se connecter' : 'S\'inscrire'}
            </button>
        </form>

        <p class="toggle-mode">
            {isLogin ? 'Pas encore de compte ?' : 'Déjà un compte ?'}
            <button class="link-button" on:click={() => isLogin = !isLogin}>
                {isLogin ? 'S\'inscrire' : 'Se connecter'}
            </button>
        </p>
    </div>
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

    .link-button {
        background: none;
        border: none;
        color: #3498db;
        cursor: pointer;
        padding: 0;
        font: inherit;
        text-decoration: underline;
    }

    .link-button:hover {
        color: #2980b9;
    }
</style> 