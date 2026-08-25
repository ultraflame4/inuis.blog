<script lang="ts">
    import { page } from "$app/state";

    const queries = page.url.searchParams
    const fail_reason = queries.get("reason") ?? "unknown"
    const allowed_providers = queries.get("allowed")?.split(",")

    const reason_descriptions: Record<string, string> = {
        "bad_identity": "The identity you are using is not connected to your account!\nPlease login with these providers first!:" +allowed_providers?.map(x => `\n- ${x}`)
    }

</script>
<main class="p-4">
    <h2 class="classic text-center">
        Login failed
    </h2>
    <h3 class="text-center">
        <code class="ring px-2 py-1">{fail_reason}</code>
    </h3>
    <pre class="text-center">
        <p>{ reason_descriptions[fail_reason] }</p>
    </pre>
</main>
