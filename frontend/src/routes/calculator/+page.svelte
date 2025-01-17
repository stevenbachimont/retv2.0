<script>
    import CompletCalculator from '../../lib/components/CompletCalculator.svelte';
    import SimpleCalculator from '../../lib/components/SimpleCalculator.svelte';
    import { user } from '../../lib/stores';

    let isSimpleMode = false;
</script>

<div class="container">
    <div class="mode-switch">
        <button 
            class="switch-button" 
            class:active={!isSimpleMode} 
            on:click={() => isSimpleMode = false}
        >
            Calculateur Complet
        </button>
        <button 
            class="switch-button" 
            class:active={isSimpleMode} 
            on:click={() => isSimpleMode = true}
        >
            Calculateur Simplifié
        </button>
    </div>

    {#if isSimpleMode}
        <SimpleCalculator />
    {:else}
        <CompletCalculator />
    {/if}
</div>

<style>
    .container {
        padding: 1rem;
    }

    .mode-switch {
        display: flex;
        justify-content: center;
        gap: 1rem;
        margin-bottom: 2rem;
        padding: 0.5rem;
        background: white;
        border-radius: 0.5rem;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        max-width: 600px;
        margin: 0 auto 2rem;
    }

    .switch-button {
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 0.25rem;
        cursor: pointer;
        background: #f0f0f0;
        color: #666;
        transition: all 0.3s ease;
        flex: 1;
        font-size: 1rem;
    }

    .switch-button.active {
        background: #3498db;
        color: white;
        box-shadow: 0 2px 4px rgba(52, 152, 219, 0.2);
    }

    .switch-button:hover:not(.active) {
        background: #e0e0e0;
    }
</style> 