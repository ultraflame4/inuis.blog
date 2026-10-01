<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { base } from '$app/paths';
	let { data, children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header class="flex h-(--navbar-height) items-center gap-2 border-b border-b-divider px-4 sticky top-0 z-20 bg-paper">
	<h1 class="text-xl text-ink"><a href="/">[{data.app_meta.brand_name}]</a></h1>
	{#if !data.user}
		<div class="ml-auto"></div>
		{#if data.app_meta.allow_signup}
			<a class="m-0 btn" href="{base}/signup">Signup</a>
		{/if}
		<a class="m-0 btn" href="{base}/login">Login</a>
	{:else}
		<a class="m-0" href="{base}/me">/me</a>
	{/if}
</header>
{@render children()}

<footer
	class="flex h-(--footer-height) items-center justify-center border-t border-divider text-center text-sm text-ink-muted z-20 bottom-0 sticky bg-paper"
>
	<p>
		<span class="font-serif">&copy;&nbsp;</span>{data.app_meta.copyright_name}
		{new Date().getFullYear()} | Powered with 💖 by
		<a class="font-mono underline" href="https://github.com/ultraflame4/inuis.blog">I.N.U.I.S</a>
	</p>
</footer>
