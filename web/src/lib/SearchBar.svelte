<script lang="ts">
	import { page } from '$app/state';
	import { updateSearchParams } from '$lib/utils';

	let searchQuery = $derived(page.url.searchParams.get('q') ?? '');
	const update_query = updateSearchParams;
</script>

<form
	class="my-0 flex h-9 grow gap-2 sm:grow-0"
	onsubmit={(e) => {
		update_query('q', searchQuery);
		e.preventDefault();
	}}
>
	<button
		class="btn h-full align-middle"
		class:hidden={searchQuery.length == 0}
		onclick={(e) => {
			update_query('q', null);
			e.preventDefault();
		}}
	>
		Clear
	</button>
	<input class="btn my-0 h-full grow px-1 focus:border-accent outline-none" bind:value={searchQuery} type="search" />
	<button type="submit" class="btn h-full align-middle">Search</button>
</form>
