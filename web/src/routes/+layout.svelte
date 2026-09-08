<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	let { data, children } = $props();
	$inspect(data);
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header class="mx-4 flex h-(--navbar-height) items-center gap-2">
	<h1 class="text-xl text-ink"><a href="/">[{data.app_meta.brand_name}]</a></h1>

	{#if !!data.user}
		<a class="classic m-0" href="/me">/me</a>
		<div class="ml-auto"></div>
		<a class="classic m-0 text-ink" href="/me/settings">Settings</a>
	{:else}
		<div class="ml-auto"></div>
		{#if data.app_meta.allow_signup}
			<a class="classic m-0" href="/signup">Signup</a>
		{/if}
		<a class="classic m-0" href="/login">Login</a>
	{/if}
</header>
{@render children()}

<footer class="flex h-(--footer-height) items-center justify-center text-center text-sm">
	<p>
		<span class="font-serif">&copy;&nbsp;</span>{data.app_meta.copyright_name}
		{new Date().getFullYear()} | Powered with 💖 by
		<a class="font-mono underline" href="https://github.com/ultraflame4/inuis.blog">I.N.U.I.S</a>
	</p>
</footer>
