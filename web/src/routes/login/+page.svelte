<script lang="ts">
	import { api } from '$lib/api';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';
	import { onMount } from 'svelte';

	const { data } = $props();
	const providers = $derived(data['providers'])

</script>

<div class="w-full min-h-page-content flex items-center justify-center">
	<main class="mx-auto border border-ink p-8 flex flex-col gap-4">
		<h2 class="classic my-0">Log In</h2>
		<span>Hello there!</span>
		<p>Sign in with your preferred provider below.</p>
		<nav class="w-full mx-auto my-2">
			<ul class="flex flex-col gap-2">
				{#each providers as provider}
					<li class="block">
						<OidcLoginButton provider_id={provider} />
					</li>
				{/each}
			</ul>
		</nav>
	</main>
</div>
