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
	const {
		provider_id,
		intent,
		prefix,
		disabled,
		class: className,
		requireInviteCode,
		inviteCode
	}: Props = $props();

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
	class="flex w-64 overflow-hidden rounded-sm border-dashed border border-ink p-4
	    text-sm font-semibold text-nowrap text-ellipsis no-underline
		text-ink-muted hover:underline
		 data-disabled:cursor-default data-disabled:border-divider data-disabled:text-divider {className}"
	data-disabled={disabled || !invite_valid || undefined}
>
	<span class="mr-auto text-ink-light">{prefix ?? 'Sign in with'}</span>

	<span class="mx-auto text-ink-muted">
	    [<span class="font-bold text-ink capitalize">{provider}</span>]
	</span>

</a>
