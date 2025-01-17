<script lang="ts">
    import { onMount } from 'svelte';
    import CompletCalculator from '../lib/components/CompletCalculator.svelte';
    import Auth from '../lib/components/Auth.svelte';
    import { user } from '../lib/stores';

    let loading = true;
    console.log("Initial state:", { loading, user: $user });

    onMount(async () => {
        console.log("onMount started");
        const token = localStorage.getItem('token');
        console.log("Token:", token);
        
        if (!token) {
            user.set(null);
            loading = false;
            console.log("No token found");
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/api/verify', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            console.log("Verify response:", response);
            
            if (response.ok) {
                const data = await response.json();
                user.set(data);
                console.log("User set:", data);
            } else {
                localStorage.removeItem('token');
                user.set(null);
                console.log("Invalid token");
            }
        } catch (error) {
            console.error('Erreur de vérification:', error);
            user.set(null);
        }
        loading = false;
        console.log("Final state:", { loading, user: $user });
    });
</script>

{#if loading}
    <div class="loading">Chargement...</div>
{:else}
    {#if $user}
        <CarbonCalculator />
    {:else}
        <Auth />
    {/if}
{/if}

<style>
    .loading {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        font-size: 1.5rem;
        color: #333;
    }
</style> 