<script lang="ts">
	import { page } from '$app/state';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';

	const queries = page.url.searchParams;
	const fail_reason = queries.get('reason') ?? 'unknown';
	const _allowed_providers = queries.get('allowed')?.split(',') ?? [];
	const allowed_providers = _allowed_providers.filter((x) => x.length > 0);

	const reason_descriptions: Record<string, () => string> = {
		bad_identity: () => `
        Found Account (email: ${queries.get('email')})
        Identity you are using is not connected to the your account!
        `,
		bad_subject: () => `
        Found Account (email: ${queries.get('email')})
        Identity you are using does not match the expected subject!
        Expected: ${queries.get('expected')}
        Received: ${queries.get('received')}
        `
	};

	const desc_lines = reason_descriptions[fail_reason]().split('\n');
</script>

<main class="space-y-4 p-4">
	<h2 class="classic text-center">Login failed</h2>
	<h3 class="text-center">
		<code class="px-2 py-1 ring">{fail_reason}</code>
	</h3>
	<div class="my-4 space-y-4 text-center font-mono">
		{#each desc_lines as l}
			<p class="mb-4">{l}</p>
		{/each}
	</div>
	{#if allowed_providers.length > 0}
		<div class="flex flex-col items-center">
		    <p>Please login with these providers first!</p>
			{#each allowed_providers as provider}
				<OidcLoginButton provider_id={provider} />
			{/each}
		</div>
	{:else}
		<div class="classic text-center font-mono">
    		<h3>[Game Over]</h3>
    		<p>No identity providers detected.</p>
    		<p>Please contact the administrator to recover access.</p>
		</div>
	{/if}
</main>
