<script lang="ts">
	import { cn } from 'cn';
	import { highlightCode, render_hast as render_hast } from './markdown';
	import { onMount } from 'svelte';
    import MarkdownNode from "./MarkdownNode.svelte"

	interface Props {
		source: string;
		class?: string;
	}

	const { source, class: _class }: Props = $props();
	let tree = $state<any>(undefined);

	onMount(()=>{
	    tree = render_hast(source)
		// setTimeout(highlightCode)
	})
</script>

<div class={cn('md-content h-fit max-w-full overflow-hidden', _class)}>
    <MarkdownNode node={tree}/>
</div>
