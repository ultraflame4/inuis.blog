<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	let { data, children } = $props();
	$inspect(data)
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header class="mx-4 flex items-center gap-2 h-(--navbar-height)">
    <h1 class="text-2xl font-bold"><a href="/">{ data.app_meta.brand_name }</a></h1>


    {#if !!data.user}
        <a class="classic m-0" href="/me">/me</a>
        <div class="ml-auto"></div>
        <a class="classic text-ink m-0" href="/me/settings">Account Settings</a>
    {:else}
        <div class="ml-auto"></div>
        <a class="classic m-0" href="/signup">Signup</a>
        <a class="classic m-0" href="/login">Login</a>
    {/if}
</header>
{@render children()}

<footer class="text-center text-sm h-(--footer-height) flex items-center justify-center">
    <p>&copy; { data.app_meta.brand_name } { new Date().getFullYear() } | Powered with 💖 by <a class="font-mono underline" href="https://github.com/ultraflame4/bytesblog">bytesblog</a> </p>
</footer>
