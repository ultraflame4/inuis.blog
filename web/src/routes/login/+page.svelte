<script lang="ts">
	import { api } from '$lib/api';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';
	import { onMount } from 'svelte';

	const { data } = $props();
	const providers = $derived(data['providers']);
</script>

<div class="flex min-h-page-content w-full items-center justify-center">
	<main class="mx-auto flex flex-col gap-4 border border-divider p-8">
		<h2 class="text-2xl font-bold my-0">Log In</h2>
		<span>Hello there!</span>
		<p>Sign in with your preferred provider below.</p>
		<nav class="mx-auto my-2 w-full">
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
