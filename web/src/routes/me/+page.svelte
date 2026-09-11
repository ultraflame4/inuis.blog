<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { page } from '$app/state';
	import PageNumbers from '$lib/PageNumbers.svelte';
	import { humanizeDate, nullish } from '$lib/utils';

	const { data } = $props();

	let searchQuery = $state(page.url.searchParams.get('q') ?? '');

	async function update_query(key: string, value: string | null | undefined) {
		// await invalidate('data:posts')
		const searchParams = new URLSearchParams(page.url.searchParams);
		if (!nullish(value)) {
			searchParams.set(key, value);
		}
		await goto(`?${searchParams.toString()}`, { replaceState: true });
	}
</script>

<main class="min-h-page-content p-4">
	<p class="my-4 text-lg">
		Welcome back: <code>{data.user.Username}</code>
	</p>

	<div class="classic flex flex-wrap items-baseline gap-x-4">
		<h2>Your Posts</h2>

		<select
			id="posts-filter"
			class="relative bottom-0.5"
			onchange={(e) => update_query('filter', e.currentTarget.value)}
			value={page.url.searchParams.get('filter') ?? 'all'}
		>
			<option value="all">All</option>
			<option value="published">Published</option>
			<option value="private">Private</option>
		</select>
		<div class="mr-auto h-4 w-0 border bg-black sm:mr-0"></div>

		<a href="/me/editor" class="self-center">/New Post</a>
		<form
			class="my-0 ml-auto flex h-7 grow gap-2 sm:grow-0"
			onsubmit={(e) => {
				update_query('q', searchQuery);
				e.preventDefault();
			}}
		>
			<input class="h-full  my-0 grow" bind:value={searchQuery} />
			<button type="submit" class="h-full my-0">Search</button>
		</form>
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

				<a href="/posts/{post.hid}" class="font-semibold underline">{post.Title}</a>
				<p class="my-0 text-sm text-ink-light">{post.Summary.slice(0, 96)}...</p>
				<div class="col-2 row-3 m-0 flex justify-between text-xs text-ink-muted">
					<span>
						[{post.AllowPublic ? 'Public' : 'Private'}]
					</span>
					<span>
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
