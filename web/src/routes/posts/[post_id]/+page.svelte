<script lang="ts">
	import MarkdownContent from "$lib/MarkdownContent.svelte";
	import { humanizeDate } from "$lib/utils.js";

	const { data } = $props()
    const body = $derived(data.post_data.Content)
    const author = $derived(data.post_data.User)

</script>


<main class="p-4">
    <header class="flex flex-wrap gap-x-4 ">
        <p><span class="font-semibold">Author</span>: <a href="/u/{author.hid}">@{author.Username}</a></p>
        {#if author.hid === data.user?.hid}
            <a href="/me/editor?hid={data.post_data.hid}" class="underline">/edit</a>
        {/if}
        <div class="mr-auto"></div>
        <p><span class="font-semibold">Updated</span>: {humanizeDate(data.post_data.UpdatedAt)}</p>
    </header>
    <MarkdownContent source={body} />
    <hr/>
    <p class="font-thin opacity-50 text-xl uppercase text-center">End of page</p>
</main>
