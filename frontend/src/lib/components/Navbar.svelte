<script lang="ts">
    import { user } from '../stores';

    export let handleLogout: () => void;
    export let currentPath: string = '/calculator';
    let isMenuOpen = false;

    const routes = [
        { path: '/calculator', label: 'Calculateur', icon: '🌡️' },
        { path: '/facilitator', label: 'Facilitateur', icon: '👥' },
        { path: '/explanations', label: 'Explications', icon: '📚' }
    ];

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }
</script>

<nav class="navbar">
    <div class="nav-brand">
        <span class="logo">🌱</span>
        <span class="brand-name">R&V</span>
    </div>
    
    <button class="menu-toggle" on:click={toggleMenu}>
        {#if isMenuOpen}✕{:else}☰{/if}
    </button>
    
    <div class="nav-menu" class:open={isMenuOpen}>
        {#each routes as route}
            <a 
                href={route.path} 
                class="nav-link {currentPath === route.path ? 'active' : ''}"
                on:click={() => isMenuOpen = false}
            >
                <span class="nav-icon">{route.icon}</span>
                {route.label}
            </a>
        {/each}

        {#if $user}
            <div class="user-menu">
                <a href="/account" class="user-name" on:click={() => isMenuOpen = false}>
                    👤 {$user.username || 'Utilisateur'}
                </a>
                <button class="logout-button" on:click={() => {
                    isMenuOpen = false;
                    handleLogout();
                }}>
                    Déconnexion
                </button>
            </div>
        {/if}
    </div>
</nav>

<style>
    @media (max-width: 768px) {
        .navbar {
            flex-wrap: wrap;
        }

        .nav-menu {
            display: flex;
            width: 100%;
            flex-direction: column;
            padding: 0;
            gap: 1rem;
            max-height: 0;
            overflow: hidden;
            transition: max-height 0.3s ease-out;
        }

        .nav-menu.open {
            max-height: 500px;
            padding: 1rem 0;
            transition: all 0.3s ease-out;
        }

        .user-menu {
            padding: 0;
            border-left: none;
            border-top: 2px solid hsl(162, 30%, 90%);
            width: 100%;
            justify-content: center;
        }

        .nav-menu.open .user-menu {
            padding: 1rem 0 0 0;
        }
    }

    .menu-toggle {
        display: none;
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
        color: hsl(162, 10%, 30%);
    }

    @media (max-width: 768px) {
        .menu-toggle {
            display: block;
        }
    }

    .navbar {
        background: white;
        padding: 1rem 2rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        margin-bottom: 2rem;
        border-radius: 0.5rem;
    }

    .nav-brand {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .logo {
        font-size: 1.5rem;
    }

    .brand-name {
        font-size: 1.25rem;
        font-weight: 600;
        color: hsl(162, 85%, 32%);
    }

    .nav-menu {
        display: flex;
        gap: 1.5rem;
        align-items: center;
    }

    .nav-link {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        text-decoration: none;
        color: hsl(162, 10%, 30%);
        font-weight: 500;
        border-radius: 0.25rem;
        transition: all 0.2s ease;
    }

    .nav-link:hover {
        background: hsl(162, 85%, 95%);
        color: hsl(162, 85%, 32%);
    }

    .nav-link.active {
        background: hsl(162, 85%, 95%);
        color: hsl(162, 85%, 32%);
        font-weight: 600;
    }

    .nav-icon {
        font-size: 1.2rem;
    }

    .user-menu {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding-left: 1rem;
        border-left: 2px solid hsl(162, 30%, 90%);
    }

    .user-name {
        color: hsl(162, 10%, 30%);
        font-weight: 600;
        font-size: 1rem;
        text-decoration: none;
        cursor: pointer;
        transition: color 0.2s ease;
    }

    .user-name:hover {
        color: hsl(162, 85%, 32%);
    }

    .logout-button {
        padding: 0.5rem 1rem;
        background: #dc3545;
        color: white;
        border: none;
        border-radius: 0.25rem;
        cursor: pointer;
        transition: background-color 0.3s;
    }

    .logout-button:hover {
        background: #c82333;
    }
</style> 