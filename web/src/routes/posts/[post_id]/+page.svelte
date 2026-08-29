<script lang="ts">
	import MarkdownContent from "$lib/MarkdownContent.svelte";
	import { humanizeDate } from "$lib/utils.js";

	const { data } = $props()
    const body = $derived(data.post_data.Content)
    const author = $derived(data.post_data.User)
    
</script>


<main class="p-4">
    <header class="flex flex-wrap gap-4">
        <p class="mr-auto">Author: {author.Username}</p>
        <p>Last updated: {humanizeDate(data.post_data.UpdatedAt)}</p>
    </header>
    <MarkdownContent source={body} />
</main>