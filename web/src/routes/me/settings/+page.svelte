<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { app } from '#lib';
	import { api } from '#lib/api.js';
	import OidcLoginButton from '#lib/OidcLoginButton.svelte';
	import { humanizeDate } from '#lib/utils.js';
	import { onMount } from 'svelte';

	const { data } = $props();

	const added_providers: string[] = $derived(data['connected_providers'] ?? []);
	const providers: string[] = $derived(data['providers']);

	async function disconnectProvider(p: string) {
		if (!confirm('Are you sure you want to disconnect this provider?')) return;
		if (!(await app.elevateSession())) {
			alert('Failed to verify user.');
			return;
		}
		try {
			await api.disconnectOidcProvider(fetch, p);
			invalidateAll();
		} catch (e) {
			alert('Action failed: ' + e);
		}
	}

	let updating_username = $state<Promise<void> | null>(null);
	async function updateUsername() {
		if (updating_username == null) {
			const inp = document.getElementById('inp-change-username')! as HTMLInputElement;
			if (inp.value == data.user.Username) return;
			if (inp.value == '') return;
			updating_username = api.updateUsername(fetch, inp.value).then((x) => {
				updating_username = null;
				inp.value = '';
				invalidateAll();
			});
		}
		await updating_username;
	}
</script>

<main class="m-4 space-y-4 flow-root">
	<h2 class="text-xl font-bold">Account Details</h2>
	<p class="">Username: <code>{data.user.Username}</code></p>
	<p class="">
		Created at:
		<code>{data.user.CreatedAt}</code>
		<small>({humanizeDate(data.user.CreatedAt)} )</small>
	</p>
	<h2 class="text-xl font-bold">Account Options</h2>
	<ul class="">
		<li>
			<label for="inp-change-username">Update Username:</label>
			<input id="inp-change-username" type="text" placeholder={data.user.Username} />
			<button onclick={updateUsername} disabled={updating_username != null}>Save</button>
		</li>
		<li class="text-danger"><a href="/api/me/logout">Logout</a></li>
	</ul>
	<hr />
	<h2 class="text-xl font-bold">Connect Identity Providers</h2>
	<p class="">Connect with more identity providers to add more ways to login to your account.</p>
	<ul class="space-y-2">
		{#each providers as p}
			<li class="flex h-full items-center gap-4">
				{#if added_providers.includes(p)}
					<OidcLoginButton
						class="data-disabled:text-success"
						provider_id={p}
						intent={'add_identity'}
						prefix="Connected"
						disabled
					/>
					<button
						onclick={() => disconnectProvider(p)}
						class="btn btn-danger h-fit"
					>
						Disconnect
					</button>
				{:else}
					<OidcLoginButton provider_id={p} intent={'add_identity'} prefix="Connect to" />
				{/if}
			</li>
		{/each}
	</ul>
	<p class="my-4 italic">
		Note: These are OIDC compliant identity providers. TeenyBlog only uses them to authenticate the
		user. TeenyBlog will only fetch the user's identity from the identity provider, it will not do
		any other form of integration. (i.e cross-posting etc.)
	</p>


	<h2 class="text-xl font-bold mb-4">Advanced Options</h2>
	<ul class="space-y-4">
		<li class="btn btn-danger w-fit"><a href="/api/me/logout/all">Logout Everywhere</a></li>
		<li class="btn btn-danger w-fit"><a href="/todo">Delete Account</a></li>
	</ul>
</main>
