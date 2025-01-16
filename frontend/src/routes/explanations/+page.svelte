<script>
    import { user } from '../../lib/stores';
    import Navbar from '../../lib/components/Navbar.svelte';
    
    function handleLogout() {
        localStorage.removeItem('token');
        user.set(null);
        window.location.href = '/';
    }

    const steps = [
        {
            title: "1. Saisissez vos données mensuelles",
            description: "Choisissez une catégorie (Transport, Logement, etc.) et remplissez les informations demandées pour le mois en cours.",
            icon: "📝"
        },
        {
            title: "2. Visualisez votre impact",
            description: "Observez les barres de progression colorées qui représentent votre empreinte carbone par catégorie.",
            icon: "📊"
        },
        {
            title: "3. Suivez votre évolution",
            description: "Comparez vos résultats mois après mois et année après année pour voir vos progrès.",
            icon: "📈"
        },
        {
            title: "4. Identifiez les axes d'amélioration",
            description: "Repérez les catégories où votre impact est le plus important pour agir prioritairement.",
            icon: "🎯"
        }
    ];

    const tips = [
        {
            category: "🚗 Transport",
            tips: ["Privilégiez le train à l'avion", "Optez pour le covoiturage", "Utilisez les transports en commun"]
        },
        {
            category: "🏠 Logement",
            tips: ["Réduisez votre consommation d'électricité", "Isolez votre logement", "Utilisez des appareils économes en énergie"]
        },
        {
            category: "🍽️ Alimentation",
            tips: ["Limitez la consommation de viande rouge", "Achetez local et de saison", "Privilégiez le vrac"]
        },
        {
            category: "👕 Vêtements",
            tips: ["Achetez en seconde main", "Privilégiez les vêtements durables", "Choisissez des marques locales"]
        }
    ];
</script>

<div class="container">
    <Navbar {handleLogout} />
    
    <div class="content">
        <h1>Comment utiliser le calculateur carbone ?</h1>
        
        <div class="steps-container">
            {#each steps as step}
                <div class="step-card">
                    <div class="step-icon">{step.icon}</div>
                    <h3>{step.title}</h3>
                    <p>{step.description}</p>
                </div>
            {/each}
        </div>

        <h2>Conseils pour réduire votre empreinte</h2>
        <div class="tips-container">
            {#each tips as category}
                <div class="category-card">
                    <h3>{category.category}</h3>
                    <ul>
                        {#each category.tips as tip}
                            <li>{tip}</li>
                        {/each}
                    </ul>
                </div>
            {/each}
        </div>

        <div class="action-card">
            <h2>Prêt à commencer ?</h2>
            <p>Calculez votre empreinte carbone dès maintenant !</p>
            <a href="/calculator" class="cta-button">Accéder au calculateur</a>
        </div>
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
        margin-bottom: 3rem;
    }

    h2 {
        text-align: center;
        color: #2c3e50;
        margin: 3rem 0 2rem;
    }

    .steps-container {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 2rem;
        margin-bottom: 3rem;
    }

    .step-card {
        background: white;
        padding: 2rem;
        border-radius: 1rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        text-align: center;
        transition: transform 0.3s ease;
    }

    .step-card:hover {
        transform: translateY(-5px);
    }

    .step-icon {
        font-size: 2.5rem;
        margin-bottom: 1rem;
    }

    .tips-container {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 2rem;
        margin-bottom: 3rem;
    }

    .category-card {
        background: white;
        padding: 2rem;
        border-radius: 1rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .category-card h3 {
        color: #2c3e50;
        margin-bottom: 1rem;
    }

    .category-card ul {
        list-style-type: none;
        padding: 0;
    }

    .category-card li {
        margin: 0.5rem 0;
        padding-left: 1.5rem;
        position: relative;
    }

    .category-card li::before {
        content: "→";
        position: absolute;
        left: 0;
        color: #3498db;
    }

    .action-card {
        background: #3498db;
        color: white;
        padding: 3rem;
        border-radius: 1rem;
        text-align: center;
        margin-top: 3rem;
    }

    .cta-button {
        display: inline-block;
        background: white;
        color: #3498db;
        padding: 1rem 2rem;
        border-radius: 2rem;
        text-decoration: none;
        font-weight: bold;
        margin-top: 1rem;
        transition: transform 0.3s ease;
    }

    .cta-button:hover {
        transform: scale(1.05);
    }
</style> 