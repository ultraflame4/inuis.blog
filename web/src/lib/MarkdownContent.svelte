<script lang="ts">
	import { cn } from 'cn';
	import { highlightCode, render_md } from './markdown';
	import { onMount } from 'svelte';

	interface Props {
		source: string;
		class?: string;
	}

	const { source, class: _class }: Props = $props();
	let rendered = $state<any>(undefined);

	onMount(()=>{
	    rendered = render_md(source)
		setTimeout(highlightCode)
	})
</script>

<div class={cn('md-content h-fit max-w-full overflow-hidden', _class)}>
    {@html rendered}
</div>
