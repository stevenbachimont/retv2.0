<script lang="ts">
  import CarbonCalculator from './lib/components/CarbonCalculator.svelte';
  import CarbonSummary from './lib/components/CarbonSummary.svelte';
</script>

<main class="container mx-auto p-4">
  <h1 class="text-2xl font-bold mb-4">Calculateur d'Empreinte Carbone</h1>
  
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <CarbonCalculator />
    <CarbonSummary />
  </div>
</main> 