<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '../stores';
    import { writable } from 'svelte/store';

    interface SimpleInputs {
        transport_km: number;
        logement_m2: number;
        viande_portions: number;
        achats_montant: number;
    }

    let inputs: SimpleInputs = {
        transport_km: 0,
        logement_m2: 0,
        viande_portions: 0,
        achats_montant: 0
    };

    let totalEmissions = 0;
    const selectedMonth = new Date().toISOString().slice(0, 7);

    // Facteurs d'émission simplifiés
    const factors = {
        transport_km: 0.2,        // 200g CO2/km (moyenne)
        logement_m2: 30,         // 30kg CO2/m2/an
        viande_portions: 2,      // 2kg CO2/portion
        achats_montant: 0.5      // 500g CO2/€
    };

    function calculateEmissions() {
        totalEmissions = 
            inputs.transport_km * factors.transport_km +
            (inputs.logement_m2 * factors.logement_m2) / 12 +
            inputs.viande_portions * factors.viande_portions +
            inputs.achats_montant * factors.achats_montant;
    }

    async function saveResults() {
        const token = localStorage.getItem('token');
        if (!token) return;

        try {
            await fetch('http://localhost:8080/api/results', {
                method: 'POST',
                headers: {
                    'Authorization': token,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    month: selectedMonth,
                    category: 'Total_simplifie',
                    value: totalEmissions,
                    inputs: inputs
                })
            });
        } catch (error) {
            console.error('Erreur:', error);
        }
    }

    $: inputs, calculateEmissions();
</script>

<div class="calculator-container">
    <h1>Calculateur Simplifié</h1>
    
    <div class="calculator-card">
        <div class="input-group">
            <label>
                Transport (km/mois)
                <input 
                    type="number" 
                    bind:value={inputs.transport_km}
                    min="0"
                    placeholder="Kilomètres parcourus"
                />
            </label>
        </div>

        <div class="input-group">
            <label>
                Logement (m²)
                <input 
                    type="number" 
                    bind:value={inputs.logement_m2}
                    min="0"
                    placeholder="Surface habitable"
                />
            </label>
        </div>

        <div class="input-group">
            <label>
                Portions de viande/mois
                <input 
                    type="number" 
                    bind:value={inputs.viande_portions}
                    min="0"
                    placeholder="Nombre de portions"
                />
            </label>
        </div>

        <div class="input-group">
            <label>
                Achats du mois (€)
                <input 
                    type="number" 
                    bind:value={inputs.achats_montant}
                    min="0"
                    placeholder="Montant des achats"
                />
            </label>
        </div>

        <div class="result-section">
            <h2>Votre empreinte carbone mensuelle</h2>
            <div class="total-emissions">
                {Math.round(totalEmissions)} kg CO2e
            </div>
            <button class="save-button" on:click={saveResults}>
                Sauvegarder
            </button>
        </div>
    </div>
</div>

<style>
    .calculator-container {
        padding: 2rem;
        max-width: 800px;
        margin: 0 auto;
    }

    h1 {
        color: #2c3e50;
        text-align: center;
        margin-bottom: 2rem;
    }

    .calculator-card {
        background: white;
        padding: 2rem;
        border-radius: 1rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .input-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
        color: #2c3e50;
        font-weight: 500;
    }

    input {
        width: 100%;
        padding: 0.5rem;
        border: 2px solid #3498db;
        border-radius: 0.5rem;
        font-size: 1rem;
    }

    .result-section {
        margin-top: 2rem;
        text-align: center;
        padding: 1.5rem;
        background: #f8f9fa;
        border-radius: 0.5rem;
    }

    .total-emissions {
        font-size: 2rem;
        color: #27ae60;
        font-weight: bold;
        margin: 1rem 0;
    }

    .save-button {
        background: #3498db;
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 0.5rem;
        font-size: 1rem;
        cursor: pointer;
        transition: background 0.3s ease;
    }

    .save-button:hover {
        background: #2980b9;
    }
</style> 