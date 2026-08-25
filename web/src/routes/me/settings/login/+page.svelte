
<script lang="ts">
	import OidcLoginButton from "$lib/OidcLoginButton.svelte";
	import { onMount } from "svelte";

   	let providers: string[] = $state([]);

	async function loadProviders() {
		const res = await fetch('/api/oidc/providers');

		if (!res.ok) {
			throw new Error('Failed to load providers');
		}

		providers = await res.json();
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
            <li>
                <OidcLoginButton provider_id={p} intent={"add_identity"} prefix="Connect to" />
            </li>
        {/each}
    </ul>
</main>
