<script lang="ts">
	import PageNumbers from '$lib/PageNumbers.svelte';
	import { page } from '$app/state';
	import { humanizeDate, nullish } from '$lib/utils';
	import { goto } from '$app/navigation';

	const { data } = $props();

	const author = $derived(data.author);
	const author_posts = $derived(data.author_posts);


	let searchQuery = $state(page.url.searchParams.get('q') ?? '');
	async function update_query(key: string, value: string | null | undefined) {
		// await invalidate('data:posts')
		const searchParams = new URLSearchParams(page.url.searchParams);
		if (!nullish(value)) {
			searchParams.set(key, value);
		}
		await goto(`?${searchParams.toString()}`, { replaceState: true });
	}

	$inspect(author_posts)
</script>

<main class="min-h-page-content p-4">
	<div class="mb-6">
		<h1 class="mb-1 text-2xl font-bold">{author.Username}</h1>
		<p class="my-0 text-sm text-ink-muted">
			Joined {humanizeDate(author.CreatedAt)}
		</p>
	</div>

	<div class="classic flex flex-wrap items-baseline gap-x-4">
		<h2>Posts</h2>
		<form
			class="my-0 ml-auto flex h-6 grow gap-2 sm:grow-0"
			onsubmit={(e) => {
				update_query('q', searchQuery);
				e.preventDefault();
			}}
		>
			<input class="my-0 h-6 grow" bind:value={searchQuery} />
			<button type="submit" class="my-0">Search</button>
		</form>
	</div>

	{#if author_posts.posts.length === 0}
		<p class="text-ink-muted">No posts yet.</p>
	{:else}
		<ul class="flex flex-wrap items-center justify-center gap-4 sm:justify-start">
			{#each author_posts.posts as post, index}
				<li
					class="grid h-38 w-full gap-2 overflow-hidden rounded-2xl border px-4 py-3 sm:h-36 sm:w-96"
					style="grid-template-columns: auto 1fr; grid-template-rows: auto 1fr auto;"
				>
					<span class="col-1 row-span-2 mr-1 w-4 text-right font-bold"
						>{index + 1 + (author_posts.page - 1) * author_posts.per_page}.</span
					>

					<a href="/posts/{post.hid}" class="font-semibold underline">{post.Title}</a>
					<p class="my-0 text-sm text-ink-light">{post.Summary.slice(0, 96)}...</p>
					<div class="col-2 row-3 m-0 flex justify-between text-xs text-ink-muted">
						<span>
							Updated: {humanizeDate(post.UpdatedAt)}
						</span>
					</div>
				</li>
			{/each}
		</ul>

		<br />
		<PageNumbers page_data={author_posts} url="/u/{author.hid}" />
	{/if}
</main>
