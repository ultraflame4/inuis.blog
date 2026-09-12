<script lang="ts">
	import { page } from '$app/state';
	import { updateSearchParams } from '$lib/utils';

	let searchQuery = $derived(page.url.searchParams.get('q') ?? '');
	const update_query = updateSearchParams;
</script>

<form
	class="my-0 flex h-7 grow gap-2 sm:grow-0"
	onsubmit={(e) => {
		update_query('q', searchQuery);
		e.preventDefault();
	}}
>
	<button
		class="my-0 cursor-pointer no-underline hover:underline"
		class:hidden={searchQuery.length == 0}
		onclick={(e) => {
			update_query('q', null);
			e.preventDefault();
		}}
	>
		Clear
	</button>
	<input class="my-0 h-full grow" bind:value={searchQuery} type="search" />
	<button type="submit" class="my-0 h-full">Search</button>
</form>
