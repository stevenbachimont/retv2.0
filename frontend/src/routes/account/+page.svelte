<script>
    import Navbar from '../../lib/components/Navbar.svelte';
    import { user } from '../../lib/stores';
    
    function handleLogout() {
        localStorage.removeItem('token');
        user.set(null);
        window.location.href = '/';
    }
</script>

<div class="container">
    <Navbar {handleLogout} />
    
    <div class="content">
        <h1>Mon Compte</h1>
        {#if $user}
            <p>Email : {$user.email}</p>
            <!-- Ajoutez ici d'autres informations du compte -->
        {/if}
    </div>
</div>

<style>
    .container {
        padding: 2rem;
    }
    
    .content {
        background: white;
        padding: 2rem;
        border-radius: 0.5rem;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
</style> 