<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/state';

	const { children } = $props();

	type SidebarRoute = { path: string; label: string };
	const links_top: SidebarRoute[] = [
		{ path: '/me', label: 'Your Posts' },
		{ path: '/me/editor', label: 'Create Post' }
	];
	const links_bottom: SidebarRoute[] = [{ path: '/me/settings', label: 'Settings' }];
</script>

<aside class="float-left mr-8 min-h-page-content w-48 border-r border-r-divider relative">
	{#snippet sidebar_routes(routes: SidebarRoute[])}
		{#each routes as route}
			<a
				class="border border-transparent p-2
			 hover:border-divider
			 data-active:border-divider data-active:bg-paper-muted"
				href="{base}{route.path}"
				data-active={page.route.id === route.path || undefined}>{route.label}</a
			>
		{/each}
	{/snippet}
	<nav class="flex flex-col h-full gap-2 p-2 absolute w-full">
		{@render sidebar_routes(links_top)}
		<div class="my-auto"></div>
		{@render sidebar_routes(links_bottom)}
	</nav>
</aside>
{@render children()}
