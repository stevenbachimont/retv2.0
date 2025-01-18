<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '../stores';
    import { writable } from 'svelte/store';
    import { slide } from 'svelte/transition';
    
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
            shortCircuit: {
                none: number;
                partial: number;
                majority: number;
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
        Consommation: {
            ecommerce: {
                amazon: number;
                leboncoin: number;
                artisanat: number;
            };
            commerce: {
                brocante: number;
                localShops: number;
            };
        };
    }

    let carbonData: CarbonData | null = null;
    const selectedCategoryStore = writable<keyof CarbonData | null>(null);
    let userInputs: Record<string, number | string> = {};
    let categoryEmissions: Record<string, number> = {
        Transports: 0,
        Logement_electromenagers: 0,
        Alimentation: 0,
        Vetements: 0,
        Numerique: 0,
        Consommation: 0,
        Services_communs: 1500
    };

    $: totalGlobalEmissions = Object.values(categoryEmissions).reduce((sum, val) => sum + val, 0);
    $: colorIntensity = Math.max(0, Math.min(1, 1 - (totalGlobalEmissions / 10000)));

    // Stocker les inputs sauvegardés par catégorie
    let savedInputsByCategory: Record<string, any> = {};

    // Effet pour mettre à jour les inputs quand la catégorie change
    $: if ($selectedCategoryStore && savedInputsByCategory[$selectedCategoryStore]) {
        userInputs = { ...savedInputsByCategory[$selectedCategoryStore] };
    }

    // Ajout du mois sélectionné
    const currentMonth = new Date().toISOString().slice(0, 7); // Format: "2024-01"
    let selectedMonth = currentMonth;
    
    // Grouper les résultats par mois
    let resultsByMonth: Record<string, Record<string, number>> = {};
    let inputsByMonth: Record<string, Record<string, any>> = {};

    // Ajouter une structure pour stocker les totaux mensuels
    let monthlyTotals: Record<string, number> = {};

    // Calculer le total pour un mois donné
    function calculateMonthlyTotal(month: string) {
        if (!resultsByMonth[month]) return 0;
        return Object.values(resultsByMonth[month]).reduce((sum, val) => sum + val, 0) + (1500/12); // +125 pour Services_communs (1500/12)
    }

    // Mettre à jour les totaux quand les résultats changent
    $: if (resultsByMonth) {
        Object.keys(resultsByMonth).forEach(month => {
            monthlyTotals[month] = calculateMonthlyTotal(month);
        });
    }

    // Charger les résultats avec gestion des mois
    async function loadUserResults() {
        try {
            const response = await fetch('http://localhost:8080/api/results', {
                headers: {
                    'Authorization': localStorage.getItem('token') || '',
                }
            });
            if (response.ok) {
                const results = await response.json();
                console.log("Résultats reçus:", results);
                if (Array.isArray(results)) {
                    // Réinitialiser les structures
                    resultsByMonth = {};
                    inputsByMonth = {};
                    
                    results.forEach((result: any) => {
                        const month = new Date(result.month).toISOString().slice(0, 7);
                        console.log("Traitement du mois:", month, "pour la catégorie:", result.category);
                        if (!resultsByMonth[month]) {
                            resultsByMonth[month] = {};
                            inputsByMonth[month] = {};
                        }
                        resultsByMonth[month][result.category] = result.value;
                        if (result.inputs) {
                            inputsByMonth[month][result.category] = result.inputs;
                        }
                    });
                    
                    console.log("ResultsByMonth:", resultsByMonth);
                    console.log("MonthlyTotals:", monthlyTotals);
                    // Mettre à jour les émissions pour le mois sélectionné
                    updateEmissionsForMonth(selectedMonth);
                }
            }
        } catch (error) {
            console.error('Erreur lors du chargement des résultats:', error);
        }
    }

    // Mettre à jour les émissions quand le mois change
    function updateEmissionsForMonth(month: string) {
        // Réinitialiser d'abord toutes les émissions
        Object.keys(categoryEmissions).forEach(category => {
            if (category !== 'Services_communs') {
                categoryEmissions[category] = 0;
            }
        });
        userInputs = {};

        if (resultsByMonth[month]) {
            Object.entries(resultsByMonth[month]).forEach(([category, value]) => {
                categoryEmissions[category] = value;
            });
            
            // Si une catégorie est sélectionnée, charger ses inputs
            if ($selectedCategoryStore && inputsByMonth[month]?.[$selectedCategoryStore]) {
                userInputs = { ...inputsByMonth[month][$selectedCategoryStore] };
            }
        }
    }

    // Surveiller les changements de mois sélectionné
    $: selectedMonth && updateEmissionsForMonth(selectedMonth);

    async function saveResult(category: string, value: number, inputs: any) {
        try {
            // Sauvegarder pour le mois sélectionné
            await fetch('http://localhost:8080/api/results', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': localStorage.getItem('token') || '',
                },
                body: JSON.stringify({
                    category,
                    value,
                    inputs,
                    month: selectedMonth,
                }),
            });

            // Mettre à jour les structures locales
            if (!resultsByMonth[selectedMonth]) {
                resultsByMonth[selectedMonth] = {};
                inputsByMonth[selectedMonth] = {};
            }
            resultsByMonth[selectedMonth][category] = value;
            inputsByMonth[selectedMonth][category] = inputs;
            
            // Recalculer le total mensuel
            monthlyTotals[selectedMonth] = calculateMonthlyTotal(selectedMonth);

        } catch (error) {
            console.error('Erreur lors de la sauvegarde:', error);
        }
    }

    // Ajouter un état pour suivre quel mois a ses détails affichés
    let expandedMonth: string | null = null;

    // Fonction pour gérer l'affichage/masquage des détails
    function toggleMonthDetails(month: string) {
        expandedMonth = expandedMonth === month ? null : month;
        // Mettre à jour le mois sélectionné
        selectedMonth = month;
    }

    onMount(async () => {
        if (!$user) {
            window.location.href = '/';
            return;
        }

        try {
            await Promise.all([
                loadUserResults(),
                (async () => {
            const response = await fetch('http://localhost:8080/api/factors');
            carbonData = await response.json();
                })()
            ]);
        } catch (error) {
            console.error('Erreur lors du chargement des données:', error);
        }
    });

    async function calculateEmissions() {
        if (!$selectedCategoryStore) return;

        try {
            const response = await fetch('http://localhost:8080/api/calculate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': localStorage.getItem('token') || '',
                },
                body: JSON.stringify({
                    category: $selectedCategoryStore,
                    userInputs: userInputs,
                }),
            });

            const data = await response.json();
            categoryEmissions[$selectedCategoryStore] = data.result;

            if ($user) {
                await saveResult($selectedCategoryStore, data.result, userInputs);
            }
        } catch (error) {
            console.error('Erreur lors du calcul:', error);
        }
    }

    function resetCategory(category: keyof CarbonData) {
        categoryEmissions[category] = 0;
        if (category === $selectedCategoryStore) {
            userInputs = {};
        }
  }

    function handleLogout() {
        localStorage.removeItem('token');
        user.set(null);
        window.location.href = '/';
    }

    // Fonction pour obtenir une couleur par catégorie
    const colors = {
        'Transports': '#FF5252',
        'Logement_electromenagers': '#4CAF50',
        'Alimentation': '#FF9800',
        'Vetements': '#03A9F4',
        'Numerique': '#7C4DFF',
        'Consommation': '#FFC107',
        'Services_communs': '#9E9E9E'
    } as const;

    function getCategoryColor(category: keyof typeof colors | string): string {
        return colors[category as keyof typeof colors] || '#9E9E9E';
    }

    // Fonction pour calculer le total annuel
    function calculateYearlyTotal(year: string) {
        const monthCount = Object.keys(resultsByMonth)
            .filter(month => month.startsWith(year))
            .length;
        
        return Object.entries(monthlyTotals)
            .filter(([month]) => month.startsWith(year))
            .reduce((sum, [_, total]) => sum + total, 0);
    }

    // Fonction pour obtenir les totaux par catégorie pour une année
    function getYearlyTotalsByCategory(year: string) {
        const totals: Record<string, number> = {
            'Services_communs': 1500 // 1500 kg CO2e par an
        };
        
        Object.entries(resultsByMonth)
            .filter(([month]) => month.startsWith(year))
            .forEach(([_, monthData]) => {
                Object.entries(monthData).forEach(([category, value]) => {
                    totals[category] = (totals[category] || 0) + value;
                });
            });
        return totals;
    }

    // Fonction pour obtenir toutes les années disponibles
    function getAvailableYears(): string[] {
        return [...new Set(
            Object.keys(monthlyTotals)
                .map(month => month.slice(0, 4))
        )].sort().reverse();
    }

    // Fonction pour obtenir les mois d'une année spécifique
    function getMonthsForYear(year: string): string[] {
        return Object.keys(monthlyTotals)
            .filter(month => month.startsWith(year))
            .sort((a, b) => b.localeCompare(a));
  }
</script>

{#if $user}
    <div class="calculator-container" style="--color-intensity: {colorIntensity}">
        <h1>Calculateur de CO2</h1>
    <div class="calculator-card">
        <h2 class="title">
                Saisissez vos données mensuelles :
        </h2>

        {#if carbonData}
            <div class="form-section">
                        <label class="form-label">
                            Mois :
                            <input 
                                type="month" 
                                bind:value={selectedMonth}
                                max={currentMonth}
                                class="form-input"
                            />
                        </label>

                <label class="form-label">
                    Catégorie :
      <select 
                        class="category-select"
                            bind:value={$selectedCategoryStore}
                    >
                        <option value={null}>Choisir une catégorie</option>
                        <option value="Transports">🚗 Transports</option>
                        <option value="Logement_electromenagers">🏠 Logement et Électroménagers</option>
                        <option value="Alimentation">🍽️ Alimentation</option>
                        <option value="Vetements">👕 Vêtements</option>
                            <option value="Numerique">💻 Numérique</option>
                            <option value="Consommation">🛍️ Consommation</option>
                    </select>
                </label>

                        {#if $selectedCategoryStore}
                    <div class="input-group">
                                {#if $selectedCategoryStore === 'Transports'}
                            <label class="form-label">
                                    Type de vol :
                                    <select bind:value={userInputs.flightType} class="form-input">
                                        <option value="domestic">Vol intérieur</option>
                                        <option value="international">Vol international</option>
                                    </select>
                                </label>
                            <label class="form-label">
                                    Kilomètres en avion :
                                <input type="number" bind:value={userInputs.flightKm} class="form-input" />
                            </label>
                            <label class="form-label">
                                <label class="form-label">
                                    Kilomètres en train :
                                <input type="number" bind:value={userInputs.trainKm} class="form-input" />
                            </label>
                                Type de voiture :
                                <select bind:value={userInputs.carType} class="form-input">
                                    <option value="small">Petite</option>
                                    <option value="medium">Moyenne</option>
                                    <option value="big">Grande</option>
                                </select>
                            </label>
                            <label class="form-label">
                                    Kilomètres en voiture :
                                <input type="number" bind:value={userInputs.carKm} class="form-input" />
                            </label>
                            <label class="form-label">
                                Nombre d'occupants en voiture :
                                <input type="number" bind:value={userInputs.carOccupants} class="form-input" min="1" max="9" placeholder="1" />
                            </label>
                        {/if}

                                {#if $selectedCategoryStore === 'Logement_electromenagers'}
                            <label class="form-label">
                                Nombre d'occupants dans le logement :
                                <input type="number" bind:value={userInputs.homeOccupants} class="form-input" min="1" placeholder="1" />
                            </label>
                            <label class="form-label">
                                Surface du logement (m²) :
                                <input type="number" bind:value={userInputs.homeSize} class="form-input" min="1" placeholder="50" />
                            </label>
                            <label class="form-label">
                                    Consommation électrique mensuelle (kWh) :
                                <input type="number" bind:value={userInputs.electricityKwh} class="form-input" />
                            </label>
                            <label class="form-label">
                                    Consommation de gaz mensuelle (kWh) :
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

                                {#if $selectedCategoryStore === 'Alimentation'}
                            <label class="form-label">
                                    Consommation de viande rouge (kg) :
                                <input type="number" bind:value={userInputs.redMeatKg} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                    Consommation de viande blanche (kg) :
                                <input type="number" bind:value={userInputs.whiteMeatKg} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                    Consommation de porc (kg) :
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
                                    <label class="form-label">
                                        Achats en circuits courts :
                                        <select bind:value={userInputs.shortCircuit} class="form-input">
                                            <option value="none">Rarement ou jamais</option>
                                            <option value="partial">Parfois (environ 50%)</option>
                                            <option value="majority">Majoritairement (>80%)</option>
                                </select>
                            </label>
                        {/if}

                                {#if $selectedCategoryStore === 'Vetements'}
                            <label class="form-label">
                                    Nombre de grands vêtements achetés :
                                <input type="number" bind:value={userInputs.largeItems} class="form-input" min="0" />
                            </label>
                            <label class="form-label">
                                    Nombre de petits vêtements achetés :
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

                                {#if $selectedCategoryStore === 'Numerique'}
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
                                        Achat de smartphone :
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
                                                <option value="old">Gardé depuis 3+ ans</option>
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
                                        Les valeurs seront automatiquement multipliées par 30 pour obtenir l'impact mensuel
                                    </p>
                                    <p class="info-text">
                                        Les PC sont comptés dans la section logement et électroménagers (électronique)
                                    </p>
                                {/if}

                                {#if $selectedCategoryStore === 'Consommation'}
                                    <label class="form-label">
                                        Nombre de commandes Amazon :
                                        <input 
                                            type="number" 
                                            bind:value={userInputs.amazonOrders} 
                                            class="form-input"
                                            min="0"
                                            placeholder="0"
                                        />
                                    </label>
                                    <label class="form-label">
                                        Nombre d'achats Le Bon Coin :
                                        <input 
                                            type="number" 
                                            bind:value={userInputs.leboncoinOrders} 
                                            class="form-input"
                                            min="0"
                                            placeholder="0"
                                        />
                                    </label>
                                    <label class="form-label">
                                        Nombre d'achats artisanaux :
                                        <input 
                                            type="number" 
                                            bind:value={userInputs.artisanatOrders} 
                                            class="form-input"
                                            min="0"
                                            placeholder="0"
                                        />
                                    </label>
                                    <label class="form-label">
                                        Nombre d'achats en brocante :
                                        <input 
                                            type="number" 
                                            bind:value={userInputs.brocanteItems} 
                                            class="form-input"
                                            min="0"
                                            placeholder="0"
                                        />
                                    </label>
                                    <label class="form-label">
                                        Nombre d'achats chez les commerçants locaux :
                                        <input 
                                            type="number" 
                                            bind:value={userInputs.localShopOrders} 
                                            class="form-input"
                                            min="0"
                                            placeholder="0"
                                        />
                                    </label>
                                    <p class="info-text">
                                        Les achats alimentaires sont à compter dans la section Alimentation
                                    </p>
                                {/if}
                                
    </div>

                    <button class="calculate-button" on:click={calculateEmissions}>
                        Calculer la catégorie
                    </button>
                {/if}

                <div class="monthly-summary">
                    {#if getAvailableYears().length > 0}
                        <div class="yearly-summary">
                            {#each getAvailableYears() as year}
                                <div class="year-section">
                                    <h3 class="year-title">
                                        Année {year}
                                        <span class="yearly-total">
                                            {calculateYearlyTotal(year).toFixed(2)} kg CO2e
                                        </span>
                                    </h3>
                                    <div class="year-progress-container">
                                        {#each Object.entries(getYearlyTotalsByCategory(year))
                                            .sort(([catA], [catB]) => catA === 'Services_communs' ? 1 : catB === 'Services_communs' ? -1 : 0) 
                                            as [category, value]}
                                            <div 
                                                class="progress-bar" 
                                                style="--width: {(value / calculateYearlyTotal(year) * 100)}%; --color: {getCategoryColor(category)};"
                                                title="{category.replace('_', ' ')}: {value.toFixed(2)} kg CO2e"
                                            ></div>
                                        {/each}
                                    </div>

                                    <div class="months-container">
                                        {#each getMonthsForYear(year) as month}
                                            <div class="monthly-card">
                                                <div class="month-header">
                                                    <div class="month-info">
                                                        <span class="month-label">
                                                            {new Date(month).toLocaleDateString('fr-FR', { month: 'long', year: 'numeric' })}
                                                            {#if month === selectedMonth}
                                                                <span class="active-month-indicator">Mois sélectionné</span>
                                {/if}
                            </span>
                                                        <span class="month-value">
                                                            {monthlyTotals[month].toFixed(2)} kg CO2e
                                                        </span>
                                                    </div>
                                                    <button 
                                                        class="view-month-button"
                                                        on:click={() => toggleMonthDetails(month)}
                                                    >
                                                        {expandedMonth === month ? 'Masquer détails' : 'Voir détails'}
                                    </button>
                                                </div>
                                                
                                                <div class="progress-container">
                                                    {#each Object.entries(resultsByMonth[month] || {})
                                                        .sort(([catA], [catB]) => catA === 'Services_communs' ? 1 : catB === 'Services_communs' ? -1 : 0)
                                                        as [category, value]}
                                                        <div 
                                                            class="progress-bar" 
                                                            style="--width: {(value / monthlyTotals[month] * 100)}%; --color: {getCategoryColor(category)};"
                                                            title="{category.replace('_', ' ')}: {value.toFixed(2)} kg CO2e"
                                                        ></div>
                                                    {/each}
                                                    <div 
                                                        class="progress-bar"
                                                        style="--width: {((1500/12) / monthlyTotals[month] * 100)}%; --color: {getCategoryColor('Services_communs')};"
                                                        title="Services communs: 125 kg CO2e"
                                                    ></div>
                                                </div>
                                                
                                                {#if expandedMonth === month}
                                                    <div class="month-details" transition:slide={{ duration: 300 }}>
                                                        <div class="emissions-breakdown">
                                                            {#each Object.entries(resultsByMonth[month] || {}) as [category, value]}
                                                                <div class="category-row">
                                                                    <div class="category-label">
                                                                        <div class="color-indicator" style="background: {getCategoryColor(category)}"></div>
                                                                        <span>{category.replace('_', ' ')}</span>
                                                                    </div>
                                                                    <span>{value.toFixed(2)} kg CO2e</span>
                                                                </div>
                                                            {/each}
                                                            <div class="category-row">
                                                                <div class="category-label">
                                                                    <div class="color-indicator" style="background: {getCategoryColor('Services_communs')}"></div>
                                                                    <span>Services communs</span>
                                                                </div>
                                                                <span>125.00 kg CO2e</span>
                                                            </div>
                                                        </div>
                                                    </div>
                                {/if}
                                            </div>
                                        {/each}
                            </div>
                        </div>
                    {/each}
                        </div>
                    {/if}
                </div>
            </div>
        {:else}
            <div class="loading-spinner"></div>
        {/if}
    </div>
</div> 

<style>
  
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }

    h1 {
        color: #2c3e50;
        margin-bottom: 2rem;
        text-align: center;
    }   

    h2 {
        color: #2c3e50;
        margin-bottom: 2rem;
        text-align: center;
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
            background: hsl(0, 0%, calc(100% - (10% * (1 - var(--color-intensity)))));
            border-radius: 1rem;
            padding: 2rem;
            box-shadow: 0 4px 6px rgba(0, 0, 0, calc(0.15 + (0.25 * (1 - var(--color-intensity)))));
            max-width: 800px;
            margin: 0 auto;
            color: hsl(162, 10%, calc(15% + (25% * var(--color-intensity))));
        }

        .title {
            color: hsl(162, calc(85% * var(--color-intensity)), 20%);
            margin-bottom: 2rem;
            font-weight: 600;
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

        .form-label {
            color: hsl(162, 10%, 20%);
            font-weight: 500;
        }

        .form-input, .category-select {
            border: 2px solid hsl(162, calc(30% * var(--color-intensity)), 75%);
            background: white;
            color: hsl(162, 10%, 15%);
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
            background: hsl(0, 0%, calc(100% - (5% * (1 - var(--color-intensity)))));
            padding: 1rem;
            margin: 0.5rem 0;
            border-radius: 0.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            box-shadow: 0 2px 4px rgba(0, 0, 0, calc(0.1 + (0.2 * (1 - var(--color-intensity)))));
            position: relative;
            overflow: hidden;
            z-index: 1;
            color: hsl(162, 10%, 20%);
            font-weight: 500;
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
            color: hsl(162, 10%, 35%);
            font-weight: 500;
            font-style: italic;
        }

        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 2rem;
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

        .form-section {
            display: flex;
            flex-direction: column;
            gap: 1rem;
            max-width: 600px;
            margin: 0 auto;
        }

        .input-group {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1rem;
            margin: 0.5rem 0;
        }

        .form-label {
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
            font-size: 0.9rem;
            color: hsl(162, 10%, 20%);
        }

        .form-input, .category-select {
            width: 100%;
            padding: 0.5rem;
            border: 1px solid hsl(162, calc(30% * var(--color-intensity)), 75%);
            border-radius: 0.25rem;
            background: white;
            color: hsl(162, 10%, 15%);
            font-size: 0.9rem;
        }

        .calculate-button {
            width: auto;
            margin: 1rem auto;
            padding: 0.5rem 2rem;
            font-size: 1rem;
        }

        .category-select {
            max-width: 300px;
            margin: 0 auto;
            margin-left: 0;
        }

        input[type="month"] {
            padding: 0.5rem;
            border: 1px solid hsl(162, calc(30% * var(--color-intensity)), 75%);
            border-radius: 0.25rem;
            background: white;
            color: hsl(162, 10%, 15%);
            font-size: 0.9rem;
            width: 200px;
        }

        .monthly-summary {
            margin-top: 2rem;
            padding-top: 2rem;
            border-top: 2px solid hsl(162, calc(30% * var(--color-intensity)), 90%);
        }

        .monthly-totals {
            display: grid;
            gap: 1rem;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        }

        .monthly-card {
            background: white;
            border-radius: 0.5rem;
            padding: 1rem;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            border: 1px solid hsl(162, calc(30% * var(--color-intensity)), 85%);
            transition: all 0.3s ease;
        }

        .month-info {
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
        }

        .month-label {
            font-weight: 500;
            color: hsl(162, 10%, 30%);
        }

        .month-value {
            font-size: 1.1rem;
            font-weight: 600;
            color: hsl(162, calc(85% * var(--color-intensity)), 32%);
        }

        .current-month {
            border: 2px solid hsl(162, calc(85% * var(--color-intensity)), 32%);
            background: hsl(162, calc(85% * var(--color-intensity)), 95%);
        }

        .month-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .view-month-button {
            padding: 0.25rem 0.75rem;
            background: hsl(162, calc(85% * var(--color-intensity)), 32%);
            color: white;
            border: none;
            border-radius: 0.25rem;
            cursor: pointer;
            font-size: 0.8rem;
            transition: all 0.2s ease;
        }

        .view-month-button:hover {
            background: hsl(162, calc(85% * var(--color-intensity)), 28%);
        }

        .month-details {
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid hsl(162, calc(30% * var(--color-intensity)), 90%);
        }

        .emissions-breakdown {
            margin-bottom: 1rem;
        }

        .category-row {
            display: flex;
            justify-content: space-between;
            padding: 0.25rem 0;
            color: hsl(162, 10%, 30%);
        }

        .category-label {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .color-indicator {
            width: 12px;
            height: 12px;
            border-radius: 2px;
            flex-shrink: 0;
        }

        .month-total {
            display: flex;
            justify-content: space-between;
            padding-top: 0.5rem;
            border-top: 1px solid hsl(162, calc(30% * var(--color-intensity)), 90%);
            font-weight: 500;
            color: hsl(162, calc(85% * var(--color-intensity)), 32%);
    }

    .progress-container {
        margin-top: 1rem;
        height: 8px;
        background: #f0f0f0;
        border-radius: 4px;
        overflow: hidden;
        display: flex;
    }

    .progress-bar {
        height: 100%;
        width: var(--width);
        background: var(--color);
        transition: width 0.3s ease;
    }

    .active-month-indicator {
        font-size: 0.8rem;
        color: hsl(162, calc(85% * var(--color-intensity)), 32%);
        font-weight: normal;
        margin-left: 0.5rem;
    }

    .yearly-summary {
        margin-bottom: 2rem;
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }

    .year-section {
        padding: 1rem;
        background: white;
        border-radius: 0.5rem;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        margin-bottom: 2rem;
    }

    .year-title {
        display: flex;
        justify-content: space-between;
        align-items: center;
        color: hsl(162, calc(85% * var(--color-intensity)), 32%);
        font-size: 1.5rem;
        font-weight: 600;
        margin-bottom: 1rem;
    }

    .year-progress-container {
        height: 12px;
        background: #f0f0f0;
        border-radius: 6px;
        overflow: hidden;
        display: flex;
        margin-top: 1rem;
    }

    .yearly-total {
        font-size: 1.2rem;
        color: hsl(162, calc(85% * var(--color-intensity)), 28%);
        background: hsl(162, calc(85% * var(--color-intensity)), 95%);
        padding: 0.5rem 1rem;
        border-radius: 0.5rem;
    }

    .months-container {
        margin-top: 1.5rem;
        display: grid;
        gap: 1rem;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    }
</style> 
{/if} 