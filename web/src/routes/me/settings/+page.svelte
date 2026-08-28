
<script lang="ts">
	import { app } from "$lib";
	import { api } from "$lib/api";
	import OidcLoginButton from "$lib/OidcLoginButton.svelte";
	import { humanizeDate } from "$lib/utils";
	import { onMount } from "svelte";

	
	const { data } = $props()
   	let providers: string[] = $state([]);
   	let added_providers: string[] = $state([]);

	async function loadProviders() {
		const res = await fetch('/api/oidc/providers');
		const res2 = await fetch('/api/me/oidc-providers');

		if (!res.ok) {
			throw new Error('Failed to load providers');
		}

		providers = await res.json();
		const res2_ = await res2.json();
		added_providers = res2_["providers"];
	}

	onMount(() => {
		loadProviders();
	});


	async function disconnectProvider(p: string) {
		if (!confirm("Are you sure you want to disconnect this provider?")) return;
		if (!await app.elevateSession()) {
			alert("Failed to verify user.");
			return;
		}
		await api.disconnectOidcProvider(fetch, p);
		location.reload()
	}
	
</script>
<main class="px-4">

    <h2 class="text-xl font-bold">Account Details</h2>
    <p class="classic">Username: <code>{data.user.Username}</code></p>
    <p class="classic">
        Created at: 
        <code>{data.user.CreatedAt}</code>
        <small>({humanizeDate(data.user.CreatedAt)} )</small>
    </p>
    <h2 class="text-xl font-bold">Account Options</h2>
    <ul class="classic">
        <li><a href="/me/settings/update-username">Change Usernames</a></li>
        <li class="text-danger"><a href="/api/me/logout">Logout</a></li>
    </ul>
    <hr/>
    <h2 class="text-xl font-bold">Connect Identity Providers</h2>
    <p class="classic">Connect with more identity providers to add more ways to login to your account.</p>
    <ul class="space-y-2">
        {#each providers as p}
            <li class="flex items-center gap-4">
                {#if added_providers.includes(p)}
                    <OidcLoginButton
                        class="text-green-700"
                        provider_id={p} intent={"add_identity"} prefix="Connected" disabled />
                    <button onclick={()=>disconnectProvider(p)} class="text-danger cursor-pointer underline">Disconnect</button>
                {:else}
                    <OidcLoginButton provider_id={p} intent={"add_identity"} prefix="Connect to" />
                {/if}
            </li>
        {/each}
    </ul>
    <p class="italic my-4">
        Note: These are OIDC compliant identity providers. TeenyBlog only uses them to authenticate the user.
        TeenyBlog will only fetch the user's identity from the identity provider, it will not do any other form of integration. (i.e cross-posting etc.)
    </p>

    <hr/>

    <h2 class="text-xl font-bold">Advanced Options</h2>
    <ul class="classic">
        <li class="text-danger"><a href="/api/me/logout">Logout Everywhere</a></li>
        <li class="text-danger"><a href="/todo">Delete Account</a></li>
    </ul>
</main>
