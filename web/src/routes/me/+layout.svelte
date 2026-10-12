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

	let sidebar_open = $state(false);
</script>

<div class="fixed top-0 right-0 z-20 h-(--navbar-height) p-2 md:hidden">
	<button class="h-full btn" onclick={() => (sidebar_open = !sidebar_open)}>Menu</button>
</div>
<aside
	class="fixed top-0 right-full z-10 h-dvh w-48 border-r border-r-divider bg-paper pt-(--navbar-height) md:left-0"
	class:left-0={sidebar_open}
>
	{#snippet sidebar_routes(routes: SidebarRoute[])}
		{#each routes as route}
			<a
				class="border border-transparent px-3 py-2
			 hover:border-divider
			 data-active:border-divider data-active:bg-paper-muted"
				href={resolve(route.path as any /*tell typescript to shutup here */)}
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

<div class="relative min-h-page-content md:ml-48">
	{@render children()}
</div>
