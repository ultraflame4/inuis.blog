<script lang="ts">
	import { api } from '$lib/api';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';
	import { onMount } from 'svelte';

	const { data } = $props();
	const providers = $derived(data['providers']);
	let inviteCode = $state(data.invite_code);
</script>

<div class="flex min-h-page-content w-full items-center justify-center">
	<main class="mx-auto flex flex-col gap-4 border border-ink p-8">
		<h2 class=" my-0">Sign Up</h2>
		<span>Hello there!</span>
		<p>Sign up with your preferred provider below.</p>

		<ul class="flex flex-col gap-2">
			{#each providers as provider}
				<li class="block">
					<OidcLoginButton
						provider_id={provider}
						intent="signup"
						prefix="Sign up with"
						{inviteCode}
						requireInviteCode={data.app_meta.require_invite_code}
					/>
				</li>
			{/each}
		</ul>

		{#if data.app_meta.require_invite_code || inviteCode}
			<div class="flex flex-col gap-1">
				<label for="inp-invite_code">This server has enabled invite only.</label>
				<input
					bind:value={inviteCode}
					type="text"
					id="inp-invite_code"
					name="invite_code"
					class="h-10 border border-ink text-center text-lg tracking-wider placeholder:text-ink-muted"
					placeholder="INVITE CODE"
				/>
			</div>
		{/if}
	</main>
</div>
