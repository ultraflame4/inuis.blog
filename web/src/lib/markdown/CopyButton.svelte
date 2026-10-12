<script lang="ts">
	import { cn } from 'cn';
	interface Props {
		text: string;
		class?: string;
	}

	const { text, class: class_ }: Props = $props();

	let isCopied = $state(false);
	let timeout = null as number | null;
	async function copyCode(text: string) {
		await navigator.clipboard.writeText(text);
		isCopied = true;
		resetState();
	}

	function resetState() {
		if (timeout) clearTimeout(timeout);
		timeout = setTimeout(() => {
			isCopied = false;
		}, 1000);
	}
</script>

<button
	class={cn('md-escape my-0 btn btn-tiny border-divider font-mono data-copied:text-success', class_)}
	data-copied={isCopied || undefined}
	onclick={() => copyCode(text)}
>
	{#if isCopied}
		Copied!
	{:else}
		Copy
	{/if}
</button>
