<script lang="ts">
    import { page } from "$app/state";
	import OidcLoginButton from "$lib/OidcLoginButton.svelte";

    const queries = page.url.searchParams
    const fail_reason = queries.get("reason") ?? "unknown"
    const allowed_providers = queries.get("allowed")?.split(",")

    const reason_descriptions: Record<string, ()=> string> = {
        "bad_identity": () => `
        Found Account (email: ${queries.get("email")})
        Identity you are using is not connected to the your account!
        Please login with these providers first!
        `
    }

    const desc_lines = reason_descriptions[fail_reason]().split("\n")
</script>
<main class="p-4 space-y-4">
    <h2 class="classic text-center">
        Login failed
    </h2>
    <h3 class="text-center">
        <code class="ring px-2 py-1">{fail_reason}</code>
    </h3>
    <div class="text-center font-mono space-y-4 my-4">
        {#each desc_lines as l }
            <p class="mb-4">{ l }</p>
        {/each}
    </div>
    <div class="flex flex-col items-center">
        {#each allowed_providers as provider}
            <OidcLoginButton provider_id={provider}/>
        {/each}
    </div>
</main>
