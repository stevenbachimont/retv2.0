<script>
    import { onMount } from 'svelte';
    import { user } from '../lib/stores';
    import Auth from '../lib/components/Auth.svelte';
    import Navbar from '../lib/components/Navbar.svelte';

    function handleLogout() {
        localStorage.removeItem('token');
        user.set(null);
        window.location.href = '/';
    }

    onMount(async () => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const response = await fetch('http://localhost:8080/api/verify', {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                if (response.ok) {
                    const userData = await response.json();
                    user.set(userData);
                } else {
                    localStorage.removeItem('token');
                    user.set(null);
                }
            } catch (error) {
                console.error('Error:', error);
            }
        }
    });
</script>

{#if $user}
    <Navbar {handleLogout} />
{/if}
<slot /> 