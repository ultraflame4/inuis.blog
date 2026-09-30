<script lang="ts">
	interface Props {
		text: string;
	}

	const { text }: Props = $props();

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
	class="btn md-escape my-0 border-paper font-mono data-copied:text-success"
	data-copied={isCopied || undefined}
	onclick={() => copyCode(text)}
>
	{#if isCopied}
		Copied!
	{:else}
		Copy
	{/if}
</button>
