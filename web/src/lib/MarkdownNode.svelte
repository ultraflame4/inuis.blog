<script lang="ts">
    import type { Snippet } from 'svelte';
    import type * as hast from 'hast';
	import CopyButton from './CopyButton.svelte';

    type Node =  hast.Nodes

    let {node}: {node?: Node;} = $props();

    // Cast to any here because I can't be arsed to manually do type checks.
    // Its either they have children or don't.
    const children = () => (node as any)?.children ?? [];

    function asProps(properties: hast.Element['properties']) {
        return {
            id: properties.id,
            class: Array.isArray(properties.className)
                ? properties.className.join(' ')
                : properties.className,
    
            title: properties.title,
            style: properties.style,
    
            href: properties.href,
            target: properties.target,
            rel: properties.rel,
    
            src: properties.src,
            alt: properties.alt,
            width: properties.width,
            height: properties.height,
    
            disabled: properties.disabled,
            checked: properties.checked,
            selected: properties.selected,
            readonly: properties.readOnly,
    
            name: properties.name,
            value: properties.value,
            placeholder: properties.placeholder,
    
            role: properties.role,
    
            ...Object.fromEntries(
                Object.entries(properties).filter(([key]) =>
                    key.startsWith('data-') ||
                    key.startsWith('aria-')
                )
            )
        }
    }
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

{#snippet _self(node: Node)}
    <svelte:self node={node} />
{/snippet}

{#if node === undefined}
    <p></p>
{:else if node.type === 'text'}
    {node.value}

{:else if node.type === 'root'}
    {#each children() as child}
        {@render _self(child)}
    {/each}

{:else if node.type === 'element'}
    <!-- Intercept rendering calls here -->
    {#if node.properties?.className?.includes("md-heading")}
        <!-- <svelte:element this={node.tagName!} {...node.properties}> -->
        
        <h1 class="md-heading" id={node.properties.id}>
            <a class="md-heading-ref" href="#{node.properties.id}">#</a>
            {#each children() as child}
                {@render _self(child)}
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

			<svelte:element this={"pre"} class="shiki-highlightme" data-lang={node.properties['data-lang']} >
                {#each children() as child}
                    {@render _self(child)}
                {/each}
            </svelte:element>
		</div>
		
    {:else}
        
        <svelte:element this={node.tagName!} {...asProps(node.properties)}>
            {#each children() as child}
                {@render _self(child)}
            {/each}
        </svelte:element>
    {/if}

{/if}