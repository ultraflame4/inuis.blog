<script lang="ts">
	import { api } from '$lib/api';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';
	import { onMount } from 'svelte';

	let providers: string[] = $state([]);

	onMount(async () => {
		providers = await api.fetchConnectedProviders(fetch);
	});
</script>

<main class="text-center min-h-page-content flow-root">
	<h1 class="classic">Verify User</h1>
	<p>This actions requires elevated permissions.</p>
	<p>We need to verify your identity before we can proceed.</p>
	<p class="my-4">Authorise with a provider below:</p>
	<nav class="mx-auto my-2 w-full">
		<ul class="flex flex-col items-center gap-2">
			{#each providers as provider}
				<li class="block">
					<OidcLoginButton provider_id={provider} intent="elevate" prefix="Verify with"/>
				</li>
			{/each}
		</ul>
	</nav>
</main>
