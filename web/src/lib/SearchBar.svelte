<script lang="ts">
	import { page } from '$app/state';
	import { updateSearchParams } from '$lib/utils';

	let searchQuery = $derived(page.url.searchParams.get('q') ?? '');
	const update_query = updateSearchParams;
</script>

<form
	class="my-0 flex h-9 grow gap-2 sm:grow-0"
	onsubmit={(e) => {
		console.log('test', e);
		update_query('q', searchQuery);
		e.preventDefault();
	}}
>
	<button
		class="h-full btn align-middle"
		class:hidden={searchQuery.length == 0}
		type="button"
		onclick={(e) => {
			update_query('q', null);
			e.preventDefault();
		}}
	>
		Clear
	</button>
	<input
		class="my-0 h-full grow btn px-1 outline-none focus:border-accent"
		bind:value={searchQuery}
		type="search"
	/>
	<button type="submit" class="h-full btn align-middle">Search</button>
</form>
