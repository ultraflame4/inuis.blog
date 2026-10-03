<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { base } from '$app/paths';
	import { page } from '$app/state';
	import { api } from '$lib/api.js';
	import PageNumbers from '$lib/PageNumbers.svelte';
	import SearchBar from '$lib/SearchBar.svelte';
	import { humanizeDate, nullish, updateSearchParams } from '$lib/utils';

	const { data } = $props();

	let searchQuery = $derived(page.url.searchParams.get('q') ?? '');

	const update_query = updateSearchParams;
</script>

<main class="min-h-page-content space-y-4 p-4">
	<p class="my-8 text-lg">
		Welcome back: <code>{data.user.Username}</code>
	</p>

	<div class=" flex flex-wrap items-center gap-x-4">
		<h2 class="text-2xl font-bold">Your Posts</h2>

		<select
			id="posts-filter"
			onchange={(e) => update_query('filter', e.currentTarget.value)}
			value={page.url.searchParams.get('filter') ?? 'all'}
		>
			<option value="all">All</option>
			<option value="published">Published</option>
			<option value="private">Private</option>
		</select>
		<div class="mr-auto h-4 w-0 border bg-black sm:mr-0"></div>

		<a href="/me/editor" class="btn">New Post</a>
		<div class="ml-auto"></div>
		<SearchBar />
	</div>
	{#snippet post_entry(post: api.PostSummary, index: number)}
		<li
			class="grid h-42 w-full gap-2 overflow-hidden border px-3 py-2 sm:h-42 sm:w-100"
			style="grid-template-columns: auto 1fr; grid-template-rows: auto 1fr auto;"
		>
			<span class="col-1 row-span-2 mr-1 w-4 text-right"
				>{index + 1 + (data.user_posts.page - 1) * data.user_posts.per_page}.</span
			>

			<a href="{base}/me/editor?hid={post.hid}" class="font-semibold underline" target="_blank"
				>{post.Title}</a
			>
			<p class="my-0 text-sm text-ink-light">{post.Summary.slice(0, 96)}...</p>
			<div class="col-span-full row-3 m-0 flex gap-1 text-xs text-ink-muted">
				<span>
					[{post.AllowPublic ? 'Public' : 'Private'}]
				</span>
				<a class="hover:underline" href="{base}/posts/{post.hid}">
					[View]
				</a>
				<span class="ml-auto">
					<!-- Created {humanizeDate(post.CreatedAt)} -->
					Updated: {humanizeDate(post.UpdatedAt)}
				</span>
			</div>
		</li>
	{/snippet}
	<ul class="flex flex-wrap items-center justify-center gap-4 sm:justify-start">
		{#each data.user_posts.items as post, index}
			{@render post_entry(post, index)}
		{/each}
	</ul>
	<br />
	<PageNumbers page_data={data.user_posts} url="/me" />
</main>
