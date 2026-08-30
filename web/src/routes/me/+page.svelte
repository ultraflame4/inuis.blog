<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { page } from '$app/state';
	import PageNumbers from '$lib/PageNumbers.svelte';
	import { humanizeDate } from '$lib/utils';

	const { data } = $props();

	async function update_query(key: string, value: string) {
		// await invalidate('data:posts')
		const searchParams = new URLSearchParams(page.url.searchParams);
		searchParams.set(key, value);
		await goto(`?${searchParams.toString()}`, { replaceState: true });
	}
</script>

<main class="p-4">
	<p class="my-4 border p-4 text-lg">
		Hello <span class="italic">{data.user.Username}</span>, writing today?
	</p>

	<div class="classic flex items-baseline gap-4">
		<h2>Your Posts</h2>

		<select
			id="posts-filter"
			onchange={(e) => update_query('filter', e.target.value)}
			value={page.url.searchParams.get('filter') ?? 'all'}
		>
			<option value="all">All</option>
			<option value="published">Published</option>
			<option value="private">Private</option>
		</select>

		<div class="h-4 w-0 border bg-black"></div>
		<a href="/me/editor" class="self-center">/New Post</a>
	</div>
	<ul class="flex flex-wrap items-center justify-center gap-4 sm:justify-start">
		{#each data.user_posts.items as post, index}
			<li
				class="grid h-38 w-full gap-2 overflow-hidden rounded-2xl border px-4 py-3 sm:h-36 sm:w-96"
				style="grid-template-columns: auto 1fr; grid-template-rows: auto 1fr auto;"
			>
				<span class="col-1 row-span-2 mr-1 w-4 text-right font-bold"
					>{index + 1 + (data.user_posts.page - 1) * data.user_posts.per_page}.</span
				>

				<a href="/me/editor?post_id={post.ID}" class="font-semibold underline">{post.Title}</a>
				<p class="my-0 text-sm text-ink-light">{post.Summary.slice(0, 96)}...</p>
				<div class="col-2 row-3 m-0 text-xs text-ink-muted flex justify-between">
					<span>
						[{post.AllowPublic ? 'Public' : 'Private'}]
					</span>
					<span >
						<!-- Created {humanizeDate(post.CreatedAt)} -->
						Updated: {humanizeDate(post.UpdatedAt)}
					</span>
				</div>
			</li>
		{/each}
	</ul>
	<br />
	<PageNumbers page_data={data.user_posts} url="/me" />
</main>
