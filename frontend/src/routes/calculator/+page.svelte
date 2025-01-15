<script>
    import CarbonCalculator from '../../lib/components/CarbonCalculator.svelte';
    import Navbar from '../../lib/components/Navbar.svelte';
    import { user } from '../../lib/stores';

    function handleLogout() {
        localStorage.removeItem('token');
        user.set(null);
        window.location.href = '/';
    }
</script>

<div class="container">
    <Navbar handleLogout={handleLogout} currentPath="/calculator" />
    <CarbonCalculator />
</div> 