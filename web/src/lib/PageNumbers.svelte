<script lang="ts">
	import type { api } from './api';

	interface Props {
		page_data: api.PageMeta;
		url: string;
		/**
		 * Query parameter to use. Defaults to "page". Useful to change if need multiple PageNumbers in a single page.
		 */
		query_key?: string;
		class?: string;
	}

	const { page_data, url, class: className, query_key }: Props = $props();
	const page_key = $derived(query_key ?? 'page');
</script>

<nav>
	{#snippet link(href: string, label: string)}
		<a class="text-ink no-underline" {href}>{label}</a>
	{/snippet}
	<ul class="flex items-center justify-center gap-2">
		<li class:hidden={page_data.page === 1}>
			<a class="text-ink no-underline hover:underline" href="{url}?{page_key}={page_data.page - 1}">
				Prev
			</a>
		</li>
		{#each Array(page_data.total_pages) as _, page_i}
			<li class:font-bold={page_i + 1 === page_data.page}>
				<a class="text-ink no-underline hover:underline" href="{url}?{page_key}={page_i + 1}">
					{page_i + 1}
				</a>
			</li>
		{/each}
		<li class:hidden={page_data.page >= page_data.total_pages}>
			<a class="text-ink no-underline hover:underline" href="{url}?{page_key}={page_data.page + 1}">
				Next
			</a>
		</li>
	</ul>
</nav>
