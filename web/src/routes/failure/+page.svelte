<script lang="ts">
	import { page } from '$app/state';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';

	const { data } = $props();

	const queries = page.url.searchParams;
	const fail_reason = queries.get('reason') ?? 'unknown';
	const _allowed_providers = queries.get('allowed')?.split(',') ?? [];

	let allowed_providers = _allowed_providers.filter((x) => x.length > 0);

	interface Reason {
		title?: () => string;
		desc: () => string;
		hide_providers?: boolean;
		links?: { label?: string; href: string }[];
		show_cancel?: boolean;
	}

	const reasons: Record<string, Reason> = {
		bad_identity: {
			desc: () => `
            Found Account (email: ${queries.get('email')})
            Identity you are using is not connected to the your account!
            Please login with these existing providers first
            `
		},
		bad_subject: {
			desc: () => `
            Found Account (email: ${queries.get('email')})
            Identity you are using does not match the expected subject!
            Expected: ${queries.get('expected')}
            Received: ${queries.get('received')}
            Please login with these existing providers first
        `
		},
		signup_required: {
			hide_providers: true,
			desc: () => `
            Account not found (email: ${queries.get('email')})
            Please sign up to create an account first.
            `,
			links: [
				{
					label: 'Create Account',
					href: '/signup'
				}
			]
		},
		signup_disabled: {
			hide_providers: true,
			title: () => `Signup Failed!`,
			desc: () => `
			Signup has been disabled on this server,

			Sorry, try again next time, or maybe don't. :shrug:
            `,
			links: [
				{
					label: 'Login',
					href: '/login'
				}
			]
		},
		bad_invite: {
			hide_providers: true,
			title: () => `Bad Invite Code`,
			desc: () => `
			This invite code is not valid, or has expired.
			-
			Who sent you this?
            `,
			links: [
				{
					label: 'Login',
					href: '/login'
				}
			]
		},
		elevate_fail: {
			hide_providers: true,
			show_cancel: true,
			title: () => 'Elevate failed',
			desc: () => `
			Failed to verify user to elevate session.
			Session elevation is required for this action.
			-
			Did you verify using correct identity?
			-
			`,
			links: [
				{
					label: 'Retry',
					href: '/me/elevate'
				}
			]
		},
		add_identity_in_use: {
			hide_providers: true,
			title: () => 'Identity Already in use',
			desc: () => `
			The identity (email: ${queries.get('email')}) for ${queries.get('provider')} is already in use by another account.

			As such, you cannot use connect it to your account!

			`,
			links: [
				{
					label: 'Back to Settings',
					href: '/me/settings'
				}
			]
		}
	};

	const title = reasons[fail_reason]?.title?.() ?? 'Login failed';
	const desc_lines = reasons[fail_reason]?.desc()?.split('\n');
	const hide_providers = reasons[fail_reason]?.hide_providers ?? false;
	const links = reasons[fail_reason]?.links ?? [];
	const show_close = reasons[fail_reason]?.show_cancel ?? false;
</script>

<main class="min-h-page-content space-y-4 p-4 text-center">
	<h2 class="classic text-center capitalize">{title}</h2>
	<h3 class="text-center">
		<code class="border border-divider-light px-2 py-1 text-link">{fail_reason}</code>
	</h3>
	<div class="my-4 space-y-4 text-center font-mono">
		{#each desc_lines as l}
			<p class="mb-4">{l}</p>
		{/each}
	</div>

	<!-- svelte-ignore block_empty -->
	{#if hide_providers}

	{:else if allowed_providers.length > 0}
		<div class="flex flex-col items-center gap-4">
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

	<div class="flex justify-center gap-2">
		{#each links as link}
			<a href={link.href} class="classic text-center font-mono no-underline hover:underline"
				>[{link.label ?? link.href}]</a
			>
		{/each}
	</div>

	{#if show_close}
		<br />
		<button class="px-4 text-danger" onclick={() => window.close()}>Close</button>
	{/if}
</main>
