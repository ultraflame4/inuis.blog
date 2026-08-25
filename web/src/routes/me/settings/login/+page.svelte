
<script lang="ts">
	import OidcLoginButton from "$lib/OidcLoginButton.svelte";
	import { onMount } from "svelte";

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
</script>
<main class="px-4">

    <h2 class="text-xl font-bold">Account Options</h2>
    <ul class="classic">
        <li><a href="/me/settings/update-username">Change Usernames</a></li>

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
                {:else}
                    <OidcLoginButton provider_id={p} intent={"add_identity"} prefix="Connect to" />
                {/if}
            </li>
        {/each}
    </ul>
    <p class="italic my-4">
        To disconnect an identity provider, please go to identity provider's settings page,
        and remove TeenyBlog from the list of connected applications.

        TeenyBlog may still show a disconnected identity provider as connected.
        <br/>
        TeenyBlog does not support this feature due to complexity and security concerns.
    </p>
</main>
