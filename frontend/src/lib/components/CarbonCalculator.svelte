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
        Numerique: {
            googleSearch: number;
            chatGPT: number;
            socialMedia: number;
            smartphone: {
                small: number;
                large: number;
                used: number;
                old: number;
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
        Numerique: 0,
        Services_communs: 1500
    };

    $: totalGlobalEmissions = Object.values(categoryEmissions).reduce((sum, val) => sum + val, 0);
    $: colorIntensity = Math.max(0, Math.min(1, 1 - (totalGlobalEmissions / 10000)));

    onMount(async () => {
        try {
            const response = await fetch('http://localhost:8080/api/factors');
            carbonData = await response.json();
        } catch (error) {
            console.error('Erreur lors du chargement des données:', error);
        }
    });

    async function calculateEmissions() {
        if (!selectedCategory) return;

        try {
            const response = await fetch('http://localhost:8080/api/calculate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    category: selectedCategory,
                    userInputs: userInputs,
                }),
            });

            const data = await response.json();
            categoryEmissions[selectedCategory] = data.result;
        } catch (error) {
            console.error('Erreur lors du calcul:', error);
        }
    }

    function resetCategory(category: keyof CarbonData) {
        categoryEmissions[category] = 0;
        if (category === selectedCategory) {
            userInputs = {};
        }
  }
</script>

<div class="calculator-container" style="--color-intensity: {colorIntensity}">
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
                        <option value="Numerique">💻 Numérique</option>
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

                        {#if selectedCategory === 'Numerique'}
                            <label class="form-label">
                                Nombre de recherches Google par jour :
                                <input 
                                    type="number" 
                                    bind:value={userInputs.googleSearches} 
                                    class="form-input"
                                    min="0"
                                    placeholder="0"
                                />
                            </label>
                            <label class="form-label">
                                Nombre de prompts ChatGPT par jour :
                                <input 
                                    type="number" 
                                    bind:value={userInputs.chatgptPrompts} 
                                    class="form-input"
                                    min="0"
                                    placeholder="0"
                                />
                            </label>
                            
                            <label class="form-label">
                                Achat de smartphone cette année :
                                <select bind:value={userInputs.smartphoneType} class="form-input">
                                    <option value="">Aucun achat</option>
                                    <option value="small">Petit modèle</option>
                                    <option value="large">Grand modèle</option>
                                </select>
                            </label>

                            {#if userInputs.smartphoneType}
                                <label class="form-label">
                                    État du smartphone :
                                    <select bind:value={userInputs.smartphoneState} class="form-input">
                                        <option value="new">Neuf</option>
                                        <option value="used">Occasion</option>
                                        <option value="old">Gardé depuis 5+ ans</option>
                                    </select>
                                </label>
                            {/if}

                            <label class="form-label">
                                Temps passé sur les réseaux sociaux par jour (heures) :
                                <input 
                                    type="number" 
                                    bind:value={userInputs.socialHours} 
                                    class="form-input"
                                    min="0"
                                    max="24"
                                    step="0.5"
                                    placeholder="0"
                                />
                            </label>
                            <p class="info-text">
                                Les valeurs seront automatiquement multipliées par 365 pour obtenir l'impact annuel
                            </p>
                            <p class="info-text">
                                Les PC sont comptés dans la section logement et électroménagers (électronique)
                            </p>
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
        background: linear-gradient(135deg, 
            hsl(142, calc(100% * var(--color-intensity)), calc(97% - (40% * (1 - var(--color-intensity))))) 0%, 
            hsl(170, calc(100% * var(--color-intensity)), calc(95% - (80% * (1 - var(--color-intensity))))) 100%);
        padding: 2rem 1rem;
        color: hsl(162, 10%, calc(20% + (60% * var(--color-intensity))));
    }

    .calculator-card {
        background: hsl(0, 0%, calc(100% - (15% * (1 - var(--color-intensity)))));
        border-radius: 1rem;
        padding: 2rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, calc(0.1 + (0.2 * (1 - var(--color-intensity)))));
        max-width: 800px;
        margin: 0 auto;
    }

    .title {
        color: hsl(162, calc(80% * var(--color-intensity)), 25%);
        margin-bottom: 2rem;
    }

    .calculate-button {
        background: hsl(162, calc(85% * var(--color-intensity)), 32%);
        color: white;
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        transition: background-color 0.3s;
    }

    .calculate-button:hover {
        background: hsl(162, calc(85% * var(--color-intensity)), 28%);
    }

    .form-input, .category-select {
        border: 2px solid hsl(162, calc(30% * var(--color-intensity)), 88%);
        border-radius: 0.375rem;
        padding: 0.5rem;
    }

    .form-input:focus, .category-select:focus {
        border-color: hsl(162, calc(60% * var(--color-intensity)), 45%);
        outline: none;
    }

    .result-card::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        width: var(--progress);
        background: linear-gradient(to right, 
            hsla(162, calc(95% * var(--color-intensity)), 30%, 0.1),
            hsla(162, calc(95% * var(--color-intensity)), 30%, 0.2));
        z-index: -1;
        transition: all 0.3s ease-in-out;
    }

    .result-card {
        background: hsl(0, 0%, calc(100% - (15% * (1 - var(--color-intensity)))));
        padding: 1rem;
        margin: 0.5rem 0;
        border-radius: 0.5rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        box-shadow: 0 2px 4px rgba(0, 0, 0, calc(0.05 + (0.15 * (1 - var(--color-intensity)))));
        position: relative;
        overflow: hidden;
        z-index: 1;
    }

    .total-card {
        background: hsl(162, calc(85% * var(--color-intensity)), 32%);
        color: white;
        padding: 1rem;
        margin-top: 1rem;
        border-radius: 0.5rem;
    }

    .reset-button {
        background: hsl(162, calc(40% * var(--color-intensity)), 90%);
        color: hsl(162, calc(85% * var(--color-intensity)), 32%);
        border: none;
        border-radius: 50%;
        width: 24px;
        height: 24px;
        cursor: pointer;
        margin-left: 0.5rem;
        transition: all 0.3s;
    }

    .reset-button:hover {
        background: hsl(162, calc(85% * var(--color-intensity)), 32%);
        color: white;
    }

    .info-text {
        font-size: 0.8em;
        color: #666;
        font-style: italic;
    }
</style> 