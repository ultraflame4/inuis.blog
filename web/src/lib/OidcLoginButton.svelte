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
	class="border-dividerink group flex w-64 overflow-hidden rounded-xl border
	p-4
	    font-semibold text-nowrap text-ellipsis text-ink-muted
		no-underline hover:border-accent
		 data-disabled:cursor-default data-disabled:border-divider data-disabled:text-divider {className}"
	data-disabled={disabled || !invite_valid || undefined}
>
	<span class="mr-auto text-ink-light">{prefix ?? 'Sign in with'}</span>

	<span class="mx-auto text-ink-muted  group-hover:text-accent">
		[<span class="font-bold text-ink capitalize group-hover:underline">{provider}</span>]
	</span>
</a>
