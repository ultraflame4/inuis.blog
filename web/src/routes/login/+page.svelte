<script lang="ts">
	import { api } from '$lib/api';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';
	import { onMount } from 'svelte';

	const { data } = $props();
	const providers = $derived(data['providers'])

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
						<OidcLoginButton provider_id={provider} />
					</li>
				{/each}
			</ul>
		</nav>
	</main>
</div>
