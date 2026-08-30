<script lang="ts">
	interface Props {
		provider_id: string;
		intent?: string;
		prefix?: string;
		disabled?: boolean;
		class?: string;
		inviteCode?: string;
		requireInviteCode?: boolean;
	}
	const { provider_id, intent, prefix, disabled, class: className, requireInviteCode, inviteCode}: Props = $props();

	const provider = $derived(provider_id.toLowerCase());

	const invite_valid = $derived(requireInviteCode ? !!inviteCode : true);

	function getHref(): string {
		if (disabled) return '';
		if (!invite_valid) return '';
    	let href = `/api/oidc/${provider}/login`;
    	const searchParams = new URLSearchParams();
    	if (intent) {
    		searchParams.append('intent', intent);
    	}
        if (requireInviteCode && intent === 'signup' && inviteCode) {
    		searchParams.append('invite_code', inviteCode);
    	}
		return `${href}?${searchParams.toString()}`;
	}
</script>

<a
	href={getHref()}
	class="block w-48 rounded-lg p-2 text-center font-semibold border border-ink data-disabled:border-divider data-disabled:text-divider {className}"
	data-disabled={disabled || (!invite_valid) || undefined}
>
	{ prefix ?? "Sign in with"}
	<span class="capitalize">{provider}</span>
</a>
