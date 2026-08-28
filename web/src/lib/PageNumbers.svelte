<script lang="ts">
	import type { api } from "./api";

   	interface Props {
        page_data: api.Paginated<any>
        url: string;
        /**
         * Query parameter to use. Defaults to "page". Useful to change if need multiple PageNumbers in a single page.
         */
        query_key?: string;
        class?: string
	}

	const { page_data, url, class: className, query_key}: Props = $props();
	const page_key = $derived(query_key ?? "page");
	

</script>

<nav>
    <ul class="flex gap-2 items-center justify-center">
        <li class:hidden={page_data.page === 1}><a href="{url}?{page_key}={page_data.page-1}">Prev</a></li>
        {#each Array(page_data.total_pages) as _, page_i}
            <li class:font-bold={page_i+1 === page_data.page}><a href="{url}?{page_key}={page_i+1}">{page_i+1}</a></li>
        {/each}
        <li class:hidden={page_data.page === page_data.total_pages}><a href="{url}?{page_key}={page_data.page+1}">Next</a></li>
    </ul>
</nav>