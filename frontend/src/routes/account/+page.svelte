<script lang="ts">
    import { user } from '../../lib/stores';
    
</script>

<div class="container">
    <div class="content">
        <h1>Mon Compte</h1>
        {#if $user}
            <div class="user-info">
                <p><strong>Pseudonyme:</strong> {$user.username}</p>
                <p><strong>Email:</strong> {$user.email}</p>
            </div>
        {/if}
    </div>
</div>

<style>
    .container {
        min-height: 100vh;
        background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
    }

    .content {
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem;
        background: white;
        border-radius: 1rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .user-info {
        margin-top: 2rem;
    }

    h1 {
        color: #2c3e50;
        margin-bottom: 2rem;
    }
</style> 