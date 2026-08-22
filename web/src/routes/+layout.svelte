<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header class="mx-4 my-2">
    <h1 class="text-2xl font-bold"><a href="/">TeenyBlog</a></h1>
</header>
{@render children()}
