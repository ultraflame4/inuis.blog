<script lang="ts">
    import { page } from "$app/state";
	import OidcLoginButton from "$lib/OidcLoginButton.svelte";

    const queries = page.url.searchParams
    const reason = queries.get("reason")

    const reason_descriptions: Record<string, string> = {
        "add_identity": `You have connected your ${queries.get("provider") ?? 'error'} identity to your account!`,
        "session_elevated": `You have elevated your session!`,
    }

    const desc_lines = reason_descriptions[reason]?.split("\n") ?? ["WHAT DID YOU DO??"]
</script>
<main class="p-4 space-y-4">
    <h2 class="classic text-center">
        Action Success
    </h2>
    <h3 class="text-center">
        <code class="ring px-2 py-1">{reason ?? '_null_'}</code>
    </h3>
    <div class="text-center font-mono space-y-4 my-4">
        {#each desc_lines as l }
            <p class="mb-4">{ l }</p>
        {/each}
    </div>

    <p class="text-center">
        You can now close this page.
    </p>
    <nav class="text-center">
        <a href="/me">Back to Home</a>
    </nav>

</main>
