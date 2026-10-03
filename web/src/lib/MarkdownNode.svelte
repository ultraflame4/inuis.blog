<script lang="ts">
    import type { Snippet } from 'svelte';
    import type * as hast from 'hast';

    type Node =  hast.Nodes

    let {node}: {node?: Node;} = $props();

    // Cast to any here because I can't be arsed to manually do type checks.
    // Its either they have children or don't.
    const children = () => (node as any)?.children ?? [];
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
        
    {:else}
        <svelte:element this={node.tagName!} {...node.properties}>
            {#each children() as child}
                <svelte:self node={child} />
            {/each}
        </svelte:element>
    {/if}

{/if}