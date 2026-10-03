<script lang="ts">
    import type { Snippet } from 'svelte';
    import type * as hast from 'hast';
	import CopyButton from './CopyButton.svelte';

    type Node =  hast.Nodes

    let {node}: {node?: Node;} = $props();

    // Cast to any here because I can't be arsed to manually do type checks.
    // Its either they have children or don't.
    const children = () => (node as any)?.children ?? [];

   	function selectCode(ev: MouseEvent) {
		const target = ev.target as HTMLElement;
		const pre = target.parentElement!.parentElement!.querySelector('pre');
		const selection = window.getSelection();
		const range = document.createRange();
		if (!pre) throw new Error('no pre');
		if (!selection) throw new Error('no selection');
   
		range.selectNodeContents(pre);
		selection.removeAllRanges();
		selection.addRange(range);
	}
</script>

{#if node === undefined}
    <p></p>
{:else if node.type === 'text'}
    {node.value}

{:else if node.type === 'root'}
    {#each children() as child}
        <svelte:self node={child} />
    {/each}

{:else if node.type === 'element'}
    <!-- Intercept rendering calls here -->
    {#if node.properties?.className?.includes("md-heading")}
        <!-- <svelte:element this={node.tagName!} {...node.properties}> -->
        
        <h1 class="md-heading" id={node.properties.id}>
            <a class="md-heading-ref" href="#{node.properties.id}">#</a>
            {#each children() as child}
                <svelte:self node={child} />
            {/each}
        </h1>
        
    {:else if node.tagName=="pre"}
  		
		<div class="md-escape group relative">
			<div
				class="md-escape absolute top-1.5 right-2 text-xs opacity-0 drop-shadow-2xl group-hover:opacity-100"
			>
				<button
					class="md-escape relative my-0 btn border-divider btn-tiny font-mono"
					onclick={selectCode}
				>
					Select
				</button>
				<CopyButton text={"wip"} />
			</div>

			<svelte:element this={"pre"} data-lang={node.properties['data-lang']} >
                {#each children() as child}
                    <svelte:self node={child} />
                {/each}
            </svelte:element>
		</div>
		
    {:else}
        
        <svelte:element this={node.tagName!} {...node.properties}>
            {#each children() as child}
                <svelte:self node={child} />
            {/each}
        </svelte:element>
    {/if}

{/if}