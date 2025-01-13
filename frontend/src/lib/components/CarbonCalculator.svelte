<script lang="ts">
    import { onMount } from 'svelte';
    import '../styles/carbonCalculator.css';
    
    interface CarbonData {
        Transports: {
            train: number;
            flight: number;
            car: {
                small: number;
                medium: number;
                big: number;
            };
        };
        Logement_electromenagers: {
            electricity: number;
            gas: number;
            apartment: number;
            house: number;
            appliance: number;
            electronic: number;
        };
        Alimentation: {
            redMeat: number;
            whiteMeat: number;
            pork: number;
            bulkFoodPurchase: {
                none: number;
                partial: number;
                total: number;
            };
        };
        Vetements: {
            large: number;
            small: number;
            madein: {
                france: number;
                autre: number;
            };
        };
    }

    let carbonData: CarbonData | null = null;
    let selectedCategory: keyof CarbonData | null = null;
    let userInputs: Record<string, number | string> = {};
    let categoryEmissions: Record<string, number> = {
        Transports: 0,
        Logement_electromenagers: 0,
        Alimentation: 0,
        Vetements: 0
    };

    $: totalGlobalEmissions = Object.values(categoryEmissions).reduce((sum, val) => sum + val, 0);

    onMount(async () => {
        try {
            const response = await fetch('http://localhost:8080/api/factors');
            carbonData = await response.json();
        } catch (error) {
            console.error('Erreur lors du chargement des données:', error);
        }
    });

    function calculateEmissions() {
        if (!carbonData || !selectedCategory) return;

        let categoryTotal = 0;
        const category = carbonData[selectedCategory];

        switch (selectedCategory) {
            case 'Transports':
                if (userInputs.trainKm) categoryTotal += Number(userInputs.trainKm) * category.train;
                if (userInputs.flightKm) categoryTotal += Number(userInputs.flightKm) * category.flight;
                if (userInputs.carKm && userInputs.carType) {
                    const occupants = Number(userInputs.carOccupants) || 1;
                    categoryTotal += (Number(userInputs.carKm) * category.car[userInputs.carType as keyof typeof category.car]) / occupants;
                }
                break;

            case 'Logement_electromenagers':
                const homeOccupants = Number(userInputs.homeOccupants) || 1;
                const homeSize = Number(userInputs.homeSize) || 0;
                if (userInputs.electricityKwh) categoryTotal += (Number(userInputs.electricityKwh) * category.electricity) / homeOccupants;
                if (userInputs.gasKwh) categoryTotal += (Number(userInputs.gasKwh) * category.gas) / homeOccupants;
                if (userInputs.housingType === 'apartment') categoryTotal += (category.apartment * homeSize) / homeOccupants;
                if (userInputs.housingType === 'house') categoryTotal += (category.house * homeSize) / homeOccupants;
                if (userInputs.applianceCount) categoryTotal += (Number(userInputs.applianceCount) * category.appliance) / homeOccupants;
                if (userInputs.electronicCount) categoryTotal += (Number(userInputs.electronicCount) * category.electronic) / homeOccupants;
                break;

            case 'Alimentation':
                if (userInputs.redMeatKg) categoryTotal += userInputs.redMeatKg * category.redMeat;
                if (userInputs.whiteMeatKg) categoryTotal += userInputs.whiteMeatKg * category.whiteMeat;
                if (userInputs.porkKg) categoryTotal += userInputs.porkKg * category.pork;
                break;

            case 'Vetements':
                if (userInputs.largeItems) categoryTotal += userInputs.largeItems * category.large;
                if (userInputs.smallItems) categoryTotal += userInputs.smallItems * category.small;
                if (userInputs.origin) {
                    categoryTotal *= category.madein[userInputs.origin as keyof typeof category.madein];
                }
                break;
        }

        categoryEmissions[selectedCategory] = categoryTotal;
    }

    function resetCategory(category: keyof CarbonData) {
        categoryEmissions[category] = 0;
        if (category === selectedCategory) {
            userInputs = {};
        }
    }
</script>

<div class="calculator-container">
    <div class="calculator-card">
        <h2 class="title">
            Calculateur d'Empreinte Carbone Annuelle
        </h2>

        {#if carbonData}
            <div class="form-section">
                <label class="form-label">
                    Catégorie :
                    <select 
                        class="category-select"
                        bind:value={selectedCategory}
                    >
                        <option value={null}>Choisir une catégorie</option>
                        <option value="Transports">🚗 Transports</option>
                        <option value="Logement_electromenagers">🏠 Logement et Électroménagers</option>
                        <option value="Alimentation">🍽️ Alimentation</option>
                        <option value="Vetements">👕 Vêtements</option>
                    </select>
                </label>

                {#if selectedCategory}
                    <div class="input-group">
                        {#if selectedCategory === 'Transports'}
                            <label class="form-label">
                                Kilomètres en train par an :
                                <input type="number" bind:value={userInputs.trainKm} class="form-input" />
                            </label>
                            <label class="form-label">
                                Kilomètres en avion par an :
                                <input type="number" bind:value={userInputs.flightKm} class="form-input" />
                            </label>
                            <label class="form-label">
                                Type de voiture :
                                <select bind:value={userInputs.carType} class="form-input">
                                    <option value="small">Petite</option>
                                    <option value="medium">Moyenne</option>
                                    <option value="big">Grande</option>
                                </select>
                            </label>
                            <label class="form-label">
                                Kilomètres en voiture par an :
                                <input type="number" bind:value={userInputs.carKm} class="form-input" />
                            </label>
                            <label class="form-label">
                                Nombre d'occupants en voiture :
                                <input type="number" bind:value={userInputs.carOccupants} class="form-input" min="1" max="9" placeholder="1" />
                            </label>
                        {/if}

                        {#if selectedCategory === 'Logement_electromenagers'}
                            <label class="form-label">
                                Nombre d'occupants dans le logement :
                                <input type="number" bind:value={userInputs.homeOccupants} class="form-input" min="1" placeholder="1" />
                            </label>
                            <label class="form-label">
                                Surface du logement (m²) :
                                <input type="number" bind:value={userInputs.homeSize} class="form-input" min="1" placeholder="50" />
                            </label>
                            <!-- ... autres champs du logement ... -->
                        {/if}

                        <!-- ... autres catégories ... -->
                    </div>

                    <button class="calculate-button" on:click={calculateEmissions}>
                        Calculer la catégorie
                    </button>
                {/if}

                <div class="results-section">
                    <h3 class="title">Résumé des émissions</h3>
                    {#each Object.entries(categoryEmissions) as [category, emissions]}
                        <div class="result-card">
                            <span>{category.replace('_', ' ')}</span>
                            <div>
                                <span>{emissions.toFixed(2)} kg CO2e</span>
                                {#if emissions > 0}
                                    <button class="reset-button" on:click={() => resetCategory(category as keyof CarbonData)}>
                                        ✕
                                    </button>
                                {/if}
                            </div>
                        </div>
                    {/each}

                    <div class="total-card">
                        <div class="flex justify-between items-center">
                            <span>Total Global</span>
                            <span>{totalGlobalEmissions.toFixed(2)} kg CO2e</span>
                        </div>
                    </div>
                </div>
            </div>
        {:else}
            <div class="loading-spinner"></div>
        {/if}
    </div>
</div>

<style>
    /* Animations optionnelles */
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }

    .animate-fade-in {
        animation: fadeIn 0.3s ease-out forwards;
    }
</style> 