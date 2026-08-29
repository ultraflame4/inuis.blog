<script lang="ts">
	import MarkdownContent from "$lib/MarkdownContent.svelte";
	import { humanizeDate } from "$lib/utils.js";

	const { data } = $props()
    const body = $derived(data.post_data.Content)
    const author = $derived(data.post_data.User)
    
</script>


<main class="p-4">
    <header class="flex flex-wrap gap-4">
        <p><span class="font-semibold">Published by</span>: <span class="font-mono bg-neutral-200 px-1 rounded-md">{author.Username}</span></p>
        {#if author.ID === data.user?.ID}
            <a href="/me/editor?post_id={data.post_data.ID}" class="underline">/edit</a>
        {/if}
        <div class="mr-auto"></div>
        <p><span class="font-semibold">Last updated</span>: {humanizeDate(data.post_data.UpdatedAt)}</p>
    </header>
    <MarkdownContent source={body} />
    <hr/>
    <p class="font-thin opacity-50 text-xl uppercase text-center">End of page</p>
</main>