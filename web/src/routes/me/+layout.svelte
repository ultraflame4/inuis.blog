<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';

	const { children } = $props();

	type SidebarRoute = { path: string; label: string };

	
	
	const links_top: SidebarRoute[] = [
		{ path: '/me', label: 'Your Posts' },
		{ path: '/me/editor', label: 'Editor' }
	];
	const links_bottom: SidebarRoute[] = [{ path: '/me/settings', label: 'Settings' }];

</script>

<aside
	class="fixed top-0 left-0 h-dvh w-48 border-r border-r-divider bg-paper pt-(--navbar-height)"
>
	{#snippet sidebar_routes(routes: SidebarRoute[])}
		{#each routes as route}
			<a
				class="border border-transparent px-3 py-2
			 hover:border-divider
			 data-active:border-divider data-active:bg-paper-muted"
				href="{resolve('')}{route.path}"
				onclick={(e) => {
					e.preventDefault();
					goto(route.path, {
						replaceState: true
					});
				}}
				data-active={page.route.id === route.path || undefined}>{route.label}</a
			>
		{/each}
	{/snippet}
	<nav class="flex h-full w-full flex-col gap-2 px-2 py-4">
		{@render sidebar_routes(links_top)}
		<div class="my-auto"></div>
		{@render sidebar_routes(links_bottom)}
	</nav>
</aside>
<div class="ml-48 min-h-page-content relative">
	{@render children()}
</div>
