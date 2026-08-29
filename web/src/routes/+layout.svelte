<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	let { data, children } = $props();
	$inspect(data)
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header class="mx-4 flex items-center gap-2" style="height: var(--navbar-height);">
    <h1 class="text-2xl font-bold"><a href="/">{ data.app_meta.brand_name }</a></h1>

    
    {#if !!data.user}
        <a class="classic m-0" href="/me">/me</a>
        <div class="ml-auto"></div>
        <a class="classic m-0" href="/me/settings">{data.user.Username}</a>
    {:else}
        <div class="ml-auto"></div>
        <a class="classic m-0" href="/signup">Signup</a>
        <a class="classic m-0" href="/login">Login</a>
    {/if}
</header>
{@render children()}
