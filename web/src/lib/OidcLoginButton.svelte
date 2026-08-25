<script lang="ts">
	interface Props {
		provider_id: string;
	}
	const { provider_id }: Props = $props();
	const provider = provider_id.toLowerCase();
</script>

<a
	href="/api/oidc/{provider}/login"
	class="block w-48 rounded-lg p-2 text-center font-semibold ring"
>
	Sign in with
	<span class="capitalize">{provider}</span>
</a>
