<script lang="ts">
	interface Props {
		provider_id: string;
		intent?: string;
		prefix?: string;
	}
	const { provider_id, intent, prefix}: Props = $props();
	const provider = provider_id.toLowerCase();


	function getHref(): string {
    	let href = `/api/oidc/${provider}/login`;
    	const searchParams = new URLSearchParams();
    	if (intent) {
    		searchParams.append('intent', intent);
    	}
		return `${href}?${searchParams.toString()}`;
	}
</script>

<a
	href={getHref()}
	class="block w-48 rounded-lg p-2 text-center font-semibold ring"
>
	{ prefix ?? "Sign in with"}
	<span class="capitalize">{provider}</span>
</a>
