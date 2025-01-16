<script lang="ts">
      import Navbar from './lib/components/Navbar.svelte';
      import { user } from './lib/stores';
      import Page from './routes/+page.svelte';

      function handleLogout() {
          localStorage.removeItem('token');
          user.set(null);
          window.location.href = '/';
      }
</script>

<main>
  <Navbar {handleLogout} />
  <Page />
</main>

<style>
  main {
    width: 100%;
    min-height: 100vh;
  }
</style> 