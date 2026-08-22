<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	let { data, children } = $props();
	console.log(data)
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header class="mx-4 my-2 flex items-center gap-4 co">
    <h1 class="text-2xl font-bold"><a href="/">TeenyBlog</a></h1>
    <div class="ml-auto"></div>
    {#if !!data.user}
        <span>{data.user.name}</span>
        <a class="classic m-0" href="/api/me/logout">Logout</a>
    {:else}
        <a class="classic m-0" href="/signup">Signup</a>
        <a class="classic m-0" href="/login">Login</a>
    {/if}
</header>
{@render children()}
