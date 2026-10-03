<script lang="ts">
	import PageNumbers from '#lib/PageNumbers.svelte';
	import { page } from '$app/state';
	import { humanizeDate, updateSearchParams } from '#lib/utils.js';
	import SearchBar from '#lib/SearchBar.svelte';

	const { data } = $props();

	const author = $derived(data.author);
	const author_posts = $derived(data.author_posts);

	let searchQuery = $derived(page.url.searchParams.get('q') ?? '');
</script>

<main class="min-h-page-content p-4">
	<div class="mb-6">
		<h1 class="mb-1 text-2xl font-bold">{author.Username}</h1>
		<p class="my-0 text-sm text-ink-muted">
			Joined {humanizeDate(author.CreatedAt)}
		</p>
	</div>

	<div class=" flex flex-wrap items-baseline gap-x-4">
		<h2>Posts</h2>
		<div class="ml-auto"></div>
		<SearchBar />
	</div>

	{#if author_posts.posts.length === 0}
		{#if searchQuery !== ''}
			<span class="text-ink-muted">No post matching query: "{searchQuery}"</span>
		{:else}
			<p class="text-ink-muted">No posts yet.</p>
		{/if}
	{:else}
		<ul class="flex flex-wrap items-center justify-center gap-4 sm:justify-start">
			{#each author_posts.posts as post, index}
				<li class="relative flow-root w-full space-y-4 px-4 py-3 pl-8 sm:h-36">
					<span class="absolute left-0 w-4 text-right font-bold"
						>{index + 1 + (author_posts.page - 1) * author_posts.per_page}.</span
					>

					<a href="/posts/{post.hid}" class="font-semibold underline"><h2>{post.Title}</h2></a>

					<p class="text-sm text-ink-light">{post.Summary.slice(0, 106)}...</p>
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
