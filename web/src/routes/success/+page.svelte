<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	const queries = page.url.searchParams;
	const reason = queries.get('reason');

	const reason_descriptions: Record<string, string> = {
		add_identity: `You have connected your ${queries.get('provider') ?? 'error'} identity to your account!`,
		session_elevated: `You have elevated your session!`
	};

	const desc_lines = reason_descriptions[reason ?? '']?.split('\n') ?? ['WHAT DID YOU DO??'];

	const auto_close_reasons = new Set(['session_elevated']);
	const auto_close_delay = 2;
	const should_auto_close = auto_close_reasons.has(reason ?? '');

	let time_counter = $state(auto_close_delay);
	onMount(() => {
		if (should_auto_close) {
			setTimeout(() => {
				window.close();
			}, auto_close_delay * 1000);
			const interval = setInterval(() => {
				time_counter -= 1;
			}, 1000);
			return () => clearInterval(interval);
		}
	});
</script>

<main class="space-y-4 p-4">
	<h2 class=" text-center">Action Success</h2>
	<h3 class="text-center">
		<code class="px-2 py-1 ring">{reason ?? '_null_'}</code>
	</h3>
	<div class="my-4 space-y-4 text-center font-mono">
		{#each desc_lines as l}
			<p class="mb-4">{l}</p>
		{/each}
	</div>

	<p class="text-center">
		{should_auto_close ? `Closing in ${time_counter}s...` : 'You can now close this page.'}
	</p>
	<nav class="text-center">
		{#if !should_auto_close}
			<a href="/me">Back to Home</a>
		{:else}
			<a onclick={window.close}>Close</a>
		{/if}
	</nav>
</main>
