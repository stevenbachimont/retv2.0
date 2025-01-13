<script lang="ts">
    import { onMount } from 'svelte';
    
    interface CarbonFactor {
        id: string;
        name: string;
        factor: number;
    }

    let factors: CarbonFactor[] = [];
    let loading = true;
    let error: string | null = null;
    let selectedFactor: string | null = null;
    let quantity = 1;

    const API_BASE_URL = 'http://localhost:8080';

    onMount(async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/factors`);
            factors = await response.json();
            loading = false;
        } catch (err) {
            error = 'Erreur lors du chargement des facteurs';
            loading = false;
        }
    });

    $: emissions = calculateEmissions(selectedFactor, quantity);

    function calculateEmissions(selectedFactor: string | null, quantity: number): string {
        if (!selectedFactor) return '0';
        const factor = factors.find(f => f.id === selectedFactor);
        return factor ? (factor.factor * quantity).toFixed(2) : '0';
    }
</script>

<div class="p-4 bg-white rounded-lg shadow">
    <h2 class="text-xl font-bold mb-4">Calculateur d'Émissions CO2</h2>
    
    {#if loading}
        <div>Chargement...</div>
    {:else if error}
        <div class="text-red-500">{error}</div>
    {:else}
        <div class="mb-4">
            <label class="block mb-2">
                Sélectionnez une catégorie :
                <select 
                    class="border p-2 ml-2 rounded"
                    bind:value={selectedFactor}
                >
                    <option value={null}>Choisir...</option>
                    {#each factors as factor}
                        <option value={factor.id}>
                            {factor.name}
                        </option>
                    {/each}
                </select>
            </label>
        </div>

        <div class="mb-4">
            <label class="block mb-2">
                Quantité :
                <input 
                    type="number" 
                    class="border p-2 ml-2 rounded"
                    bind:value={quantity}
                    min="0"
                />
                {#if selectedFactor === 'viande_rouge'}
                    <span class="ml-2">kg</span>
                {:else if selectedFactor && ['voiture', 'train', 'avion'].includes(selectedFactor)}
                    <span class="ml-2">km</span>
                {/if}
            </label>
        </div>

        {#if selectedFactor}
            <div class="mt-4 p-4 bg-gray-100 rounded">
                <p class="text-lg">
                    Émissions CO2 : <strong>{emissions} kg CO2e</strong>
                </p>
            </div>
        {/if}
    {/if}
</div> 