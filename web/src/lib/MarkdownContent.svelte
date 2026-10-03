<script lang="ts">
	import { cn } from 'cn';
	import { render_md } from './markdown';

	interface Props {
		source: string;
		class?: string;
	}

	const { source, class: _class }: Props = $props();
	const rendered = $derived(render_md(source));
</script>

<div class={cn('cm-md-content h-fit max-w-full overflow-hidden', _class)}>
    {@html rendered}
</div>
