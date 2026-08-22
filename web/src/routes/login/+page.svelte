<script lang="ts">
	import { onMount } from 'svelte';

	let providers: string[] = $state([]);

	async function loadProviders() {
		const res = await fetch('/api/oidc/providers');

		if (!res.ok) {
			throw new Error('Failed to load providers');
		}

		providers = await res.json();
	}

	onMount(() => {
		loadProviders();
	});
</script>

<main>
	<h2>Login</h2>

	<nav>
		<ul>
			{#each providers as provider}
				<li>
					<a href="/api/oidc/{provider}/login">Sign in with {provider}</a>
				</li>
			{/each}
		</ul>
	</nav>
</main>
