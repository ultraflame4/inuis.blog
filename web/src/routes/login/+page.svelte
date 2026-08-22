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

<div class="flex w-full">
	<main class="mx-auto ring p-8 flex flex-col fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 gap-4">
		<h2 class="classic my-0">Log In</h2>
		<span>Hello there!</span>
		<p>Sign in with your preferred provider below.</p>
		<nav class="w-full mx-auto my-2">
			<ul class="flex flex-col gap-2">
				{#each providers as provider}
					<li class="block">
						<a
							href="/api/oidc/{provider}/login"
							class="font-semibold ring p-2 block w-full text-center rounded-lg min-w-48"
						>
						    Sign in with 
							<span class="capitalize">{provider}</span>
						</a>
					</li>
				{/each}
			</ul>
		</nav>
	</main>
</div>
