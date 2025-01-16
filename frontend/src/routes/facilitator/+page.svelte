<script lang="ts">
    import { user } from '../../lib/stores';
    import { onMount } from 'svelte';
    
    function handleLogout() {
        localStorage.removeItem('token');
        user.set(null);
        window.location.href = '/';
    }

    interface Recommendation {
        category: string;
        excess: string;
        value: number;
        threshold: number;
        actions: Array<{
            action: string;
            impact: string;
        }>;
    }

    let results = {};
    let recommendations: Recommendation[] = [];
    let selectedMonth = new Date().toISOString().slice(0, 7);
    let isLoading = true;

    // Seuils d'émission mensuels par catégorie (en kg CO2)
    const thresholds = {
        Transports: 200,
        Logement_electromenagers: 300,
        Alimentation: 150,
        Vetements: 50,
        Numerique: 30,
        Consommation: 100
    };

    // Actions d'amélioration par catégorie
    const improvements = {
        Transports: [
            { action: "Passer au vélo pour les trajets courts", impact: "- 50 kg CO2/mois" },
            { action: "Privilégier le train à l'avion", impact: "- 200 kg CO2/trajet" },
            { action: "Opter pour le covoiturage", impact: "- 30% d'émissions" }
        ],
        Logement_electromenagers: [
            { action: "Installer des LED", impact: "- 80% sur l'éclairage" },
            { action: "Réduire le chauffage de 1°C", impact: "- 7% sur le chauffage" },
            { action: "Isoler les fenêtres", impact: "- 15% sur le chauffage" }
        ],
        Alimentation: [
            { action: "Réduire la viande rouge", impact: "- 30 kg CO2/mois" },
            { action: "Privilégier le local", impact: "- 20% sur l'alimentation" },
            { action: "Éviter le gaspillage", impact: "- 15% sur l'alimentation" }
        ],
        Vetements: [
            { action: "Acheter en seconde main", impact: "- 70% sur les vêtements" },
            { action: "Réparer plutôt que remplacer", impact: "- 40% sur les vêtements" }
        ],
        Numerique: [
            { action: "Garder son smartphone plus longtemps", impact: "- 30 kg CO2/an" },
            { action: "Limiter le streaming", impact: "- 5 kg CO2/mois" }
        ],
        Consommation: [
            { action: "Privilégier l'occasion", impact: "- 50% sur les achats" },
            { action: "Acheter local", impact: "- 30% sur le transport" }
        ]
    };

    async function loadResults() {
        try {
            const response = await fetch('http://localhost:8080/api/results', {
                headers: {
                    'Authorization': localStorage.getItem('token')
                }
            });
            if (response.ok) {
                const data = await response.json();
                results = data.reduce((acc, result) => {
                    const month = new Date(result.month).toISOString().slice(0, 7);
                    if (!acc[month]) acc[month] = {};
                    acc[month][result.category] = result.value;
                    return acc;
                }, {});
                analyzeResults();
            }
        } catch (error) {
            console.error('Erreur:', error);
        }
        isLoading = false;
    }

    function analyzeResults() {
        recommendations = [];
        const monthData = results[selectedMonth] || {};

        Object.entries(monthData).forEach(([category, value]) => {
            if (value > thresholds[category]) {
                const excess = ((value - thresholds[category]) / thresholds[category] * 100).toFixed(0);
                const improvement = improvements[category];
                recommendations.push({
                    category,
                    excess: `+${excess}%`,
                    value: Math.round(value),
                    threshold: thresholds[category],
                    actions: improvement
                });
            }
        });

        recommendations.sort((a, b) => 
            (a.value / a.threshold) - (b.value / b.threshold)
        ).reverse();
    }

    $: selectedMonth && analyzeResults();

    onMount(loadResults);
</script>

<div class="container">
    <div class="content">
        <h1>Facilitateur d'ambitions</h1>
        
        <div class="month-selector">
            <label for="month">Sélectionnez le mois à analyser :</label>
            <input 
                type="month" 
                id="month" 
                bind:value={selectedMonth} 
                max={new Date().toISOString().slice(0, 7)}
            />
        </div>

        {#if isLoading}
            <div class="loading">Chargement des données...</div>
        {:else if recommendations.length === 0}
            <div class="congratulations">
                <h2>🎉 Félicitations !</h2>
                <p>Vos émissions sont en dessous des seuils recommandés pour toutes les catégories.</p>
            </div>
        {:else}
            <div class="recommendations">
                {#each recommendations as rec}
                    <div class="recommendation-card">
                        <div class="card-header">
                            <h3>{rec.category.replace('_', ' ')}</h3>
                            <span class="excess">{rec.excess}</span>
                        </div>
                        
                        <div class="progress-container">
                            <div class="progress-bar" 
                                 style="--progress: {Math.min((rec.value / (rec.value * 1.2)) * 100, 100)}%">
                                <span class="value">{rec.value} kg CO2</span>
                            </div>
                            <div class="threshold-marker" 
                                 style="--position: {(rec.threshold / rec.value * 100)}%">
                                <span class="threshold-label">Seuil : {rec.threshold} kg CO2</span>
                            </div>
                        </div>

                        <div class="actions">
                            <h4>Actions recommandées :</h4>
                            {#each rec.actions as action}
                                <div class="action-item">
                                    <span class="action-text">{action.action}</span>
                                    <span class="impact">{action.impact}</span>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<style>
    .container {
        min-height: 100vh;
        background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
    }

    .content {
        max-width: 1200px;
        margin: 0 auto;
        padding: 2rem;
    }

    h1 {
        text-align: center;
        color: #2c3e50;
        margin-bottom: 2rem;
    }

    .month-selector {
        text-align: center;
        margin-bottom: 2rem;
    }

    .month-selector input {
        padding: 0.5rem;
        border: 2px solid #3498db;
        border-radius: 0.5rem;
        margin-left: 1rem;
    }

    .recommendation-card {
        background: white;
        border-radius: 1rem;
        padding: 2rem;
        margin-bottom: 2rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }

    .excess {
        color: #e74c3c;
        font-weight: bold;
    }

    .progress-container {
        position: relative;
        height: 2rem;
        background: #eee;
        border-radius: 1rem;
        margin-bottom: 2rem;
    }

    .progress-bar {
        position: absolute;
        height: 100%;
        width: 100%;
        background: linear-gradient(
            to right, 
            #3498db 0%, 
            #e74c3c 100%
        );
        border-radius: 1rem;
        transition: width 0.3s ease;
    }

    .value {
        position: absolute;
        right: 1rem;
        top: 50%;
        transform: translateY(-50%);
        color: white;
        font-weight: bold;
        z-index: 2;
    }

    .threshold-marker {
        position: absolute;
        left: var(--position);
        height: 100%;
        width: 4px;
        background: #2c3e50;
        z-index: 1;
    }

    .threshold-label {
        position: absolute;
        top: -1.5rem;
        left: 50%;
        transform: translateX(-50%);
        white-space: nowrap;
        font-size: 0.8rem;
        color: #2c3e50;
    }

    .actions {
        background: #f8f9fa;
        padding: 1.5rem;
        border-radius: 0.5rem;
    }

    .action-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 0;
        border-bottom: 1px solid #dee2e6;
    }

    .action-item:last-child {
        border-bottom: none;
    }

    .impact {
        color: #27ae60;
        font-weight: bold;
    }

    .congratulations {
        text-align: center;
        padding: 3rem;
        background: white;
        border-radius: 1rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .loading {
        text-align: center;
        padding: 2rem;
        color: #666;
    }
</style> 
