<script lang="ts">
	import { api } from '$lib/api';
	import OidcLoginButton from '$lib/OidcLoginButton.svelte';
	import { onMount } from 'svelte';

	const { data } = $props();
	const providers = $derived(data['providers'])
	let inviteCode = $state("")
</script>

<div class="w-full min-h-page-content flex items-center justify-center">
	<main class="mx-auto border border-ink p-8 flex flex-col gap-4">
		<h2 class="classic my-0">Sign Up</h2>
		<span>Hello there!</span>
		<p>Sign up with your preferred provider below.</p>

		<ul class="flex flex-col gap-2">
			{#each providers as provider}
				<li class="block">
					<OidcLoginButton provider_id={provider} intent="signup" prefix="Sign up with" inviteCode={inviteCode} requireInviteCode={data.app_meta.requires_invite_code} />
				</li>
			{/each}
		</ul>

		{#if data.app_meta.requires_invite_code}
			<div class="flex flex-col gap-1">
				<label for="inp-invite_code">This server has enabled invite only.</label>
				<input bind:value={inviteCode} type="text" id="inp-invite_code" name="invite_code" class="border border-ink h-10 text-center text-lg placeholder:text-ink-muted  tracking-wider" placeholder="INVITE CODE" />
			</div>
		{/if}

	</main>
</div>
