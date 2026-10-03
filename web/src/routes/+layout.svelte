<script lang="ts">
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import { resolve } from '$app/paths';
	let { data, children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<header
	class="relative z-20 flex h-(--navbar-height) items-center gap-2 border-b border-b-divider bg-paper px-4"
>
	<h1 class="text-xl text-ink"><a href="/">[{data.app_meta.brand_name}]</a></h1>
	{#if !data.user}
		<div class="ml-auto"></div>
		{#if data.app_meta.allow_signup}
			<a class="m-0 btn" href={resolve('/signup')}>Signup</a>
		{/if}
		<a class="m-0 btn" href={resolve('/login')}>Login</a>
	{:else}
		<a class="m-0" href={resolve('/me')}>/me</a>
	{/if}
</header>
<div class="flow-root min-h-page-content">
	{@render children()}
</div>

<footer
	class="z-20 flex h-(--footer-height) items-center justify-center border-t border-divider bg-paper text-center text-sm text-ink-muted"
>
	<p>
		<span class="font-serif">©</span>
		{data.app_meta.copyright_name}
		{new Date().getFullYear()} | Powered with 💖 by

		<a class="font-mono underline" href="https://github.com/ultraflame4/inuis.blog">I.N.U.I.S</a>
	</p>
</footer>
