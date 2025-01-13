<script lang="ts">
    import { onMount } from 'svelte';
    
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
        Vetements: 0,
        Services_communs: 1500
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
                            <label class="form-label">
                                Consommation électrique annuelle (kWh) :
                                <input type="number" bind:value={userInputs.electricityKwh} class="form-input" />
                            </label>
                            <label class="form-label">
                                Consommation de gaz annuelle (kWh) :
                                <input type="number" bind:value={userInputs.gasKwh} class="form-input" />
                            </label>
                            <label class="form-label">
                                Type de logement :
                                <select bind:value={userInputs.housingType} class="form-input">
                                    <option value="apartment">Appartement</option>
                                    <option value="house">Maison</option>
                                </select>
                            </label>
                            <label class="form-label">
                                Nombre d'appareils électroménagers :
                                <input type="number" bind:value={userInputs.applianceCount} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                Nombre d'appareils électroniques :
                                <input type="number" bind:value={userInputs.electronicCount} class="form-input" min="0" />
                            </label>
                        {/if}

                        {#if selectedCategory === 'Alimentation'}
                            <label class="form-label">
                                Consommation annuelle de viande rouge (kg) :
                                <input type="number" bind:value={userInputs.redMeatKg} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                Consommation annuelle de viande blanche (kg) :
                                <input type="number" bind:value={userInputs.whiteMeatKg} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                Consommation annuelle de porc (kg) :
                                <input type="number" bind:value={userInputs.porkKg} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                Achats en vrac :
                                <select bind:value={userInputs.bulkPurchase} class="form-input">
                                    <option value="none">Jamais</option>
                                    <option value="partial">Parfois</option>
                                    <option value="total">Toujours</option>
                                </select>
                            </label>
                        {/if}

                        {#if selectedCategory === 'Vetements'}
                            <label class="form-label">
                                Nombre de grands vêtements achetés par an :
                                <input type="number" bind:value={userInputs.largeItems} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                Nombre de petits vêtements achetés par an :
                                <input type="number" bind:value={userInputs.smallItems} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                Origine principale des vêtements :
                                <select bind:value={userInputs.origin} class="form-input">
                                    <option value="france">France</option>
                                    <option value="autre">Autre pays</option>
      </select>
                            </label>
                        {/if}
    </div>

                    <button class="calculate-button" on:click={calculateEmissions}>
                        Calculer la catégorie
                    </button>
                {/if}

                <div class="results-section">
                    <h3 class="title">Résumé des émissions</h3>
                    {#each Object.entries(categoryEmissions) as [category, emissions]}
                        <div class="result-card" style="--progress: {Math.min((emissions / totalGlobalEmissions) * 100, 100)}%">
                            <span>
                                {category.replace('_', ' ')}
                                {#if category === 'Services_communs'}
                                    <span class="info-text">(Services publics, infrastructures, etc.)</span>
                                {/if}
                            </span>
                            <div>
                                <span>{emissions.toFixed(2)} kg CO2e</span>
                                {#if emissions > 0 && category !== 'Services_communs'}
                                    <button class="reset-button" on:click={() => resetCategory(category as keyof CarbonData)}>
                                        ✕
                                    </button>
                                {/if}
                            </div>
                        </div>
                    {/each}

                    <div class="total-card">
                        <div class="total-card-content">
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

    .calculator-container {
        min-height: 100vh;
        background: linear-gradient(135deg, #f0fff4 0%, #e6fffa 100%);
        padding: 2rem 1rem;
    }

    .info-text {
        font-size: 0.8em;
        color: #666;
        font-style: italic;
    }

    .result-card {
        position: relative;
        overflow: hidden;
        z-index: 1;
    }

    .result-card::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        width: var(--progress);
        background: linear-gradient(to right, rgba(5, 150, 105, 0.1), rgba(5, 150, 105, 0.2));
        z-index: -1;
        transition: width 0.3s ease-in-out;
    }
</style> 