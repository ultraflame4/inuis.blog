<script lang="ts">
	import { resolve } from '$app/paths';
	import MarkdownContent from '$lib/MarkdownContent.svelte';
	import { humanizeDate } from '$lib/utils.js';

	const { data } = $props();
	const body = $derived(data.post_data.Content);
	const author = $derived(data.post_data.User);
</script>

<main class="px-4 mt-4">
	<header class="flex flex-wrap gap-x-4 space-y-4 max-w-396 mx-auto">
		{#if !data.post_data.AllowPublic}
			<div
				class="py-2 w-full border border-danger text-center text-danger"
				title="This post has not been published can cannot be viewed publicly"
			>
				<span class="tracking-widest">[PRIVATE POST]</span>
				<p>This post has not been published can cannot be viewed publicly</p>
			</div>
		{/if}
		<p>
			<span class="font-semibold">Author</span>: <a href="/u/{author.hid}">@{author.Username}</a>
		</p>
		{#if author.hid === data.user?.hid}
			<a href={resolve(`/me/editor?hid=${data.post_data.hid}`)} class="underline">/edit</a>
		{/if}

		<div class="mr-auto"></div>
		<p><span class="font-semibold">Updated</span>: {humanizeDate(data.post_data.UpdatedAt)}</p>
	</header>
	<MarkdownContent source={body} class="max-w-396 mx-auto border px-8"/>
	<hr />
	<p class="text-center text-xl font-thin uppercase opacity-50">End of page</p>
</main>
