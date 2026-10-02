<script lang="ts">
	import { resolve } from '$app/paths';
	import MarkdownContent from '$lib/MarkdownContent.svelte';
	import { humanizeDate } from '$lib/utils.js';

	const { data } = $props();
	const body = $derived(data.post_data.Content);
	const author = $derived(data.post_data.User);
</script>

<main class="mt-4 px-4">
	<header class="mx-auto flex max-w-396 flex-wrap space-y-4 gap-x-4">
		{#if !data.post_data.AllowPublic}
			<div
				class="w-full border border-danger py-2 text-center text-danger"
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
	<MarkdownContent source={body} class="mx-auto max-w-396 border px-8" />
	<div class="my-8"></div>
</main>
