<script lang="ts">
	import SvelteMarkdown from '@humanspeak/svelte-markdown';
	import { codeToHtml, bundledLanguages } from 'shiki/dist/bundle-web.mjs';
	

	
	interface Props {
		source: string;
	}

	const { source }: Props = $props();


	function highlightCode() {
	    const langs = new Set(Object.keys(bundledLanguages));
        Array.from(document.querySelectorAll('.shiki-highlightme')).map(async (el) => {
            
                let lang = el.getAttribute('data-lang') ?? '';
           	    console.log("shiki", lang);

                if (!langs.has(lang)) {
                    lang = 'text';
                }
                
                const text = el.querySelector('code')!.textContent;
                el.classList.remove('shiki-highlightme');
                el.classList.add('shiki-highlighted');
                try{
                    const html = await codeToHtml(text, { lang, theme: 'github-dark' });
                    console.log("html", html);
                    el.outerHTML = html;
                } catch (e) {
                    console.log("shiki", e);
                }
      
            });
        }
	
	$effect( () => {
	    console.log("highlightCode", source);
		highlightCode()
	});
</script>



<div class="cm-md-content h-fit max-w-full overflow-hidden">

	<SvelteMarkdown {source}>
		{#snippet heading({ depth, children, text, slug })}
			<a href="#__{slug(text)}">
				{#if depth === 1}
					<h1 class="cm-md-heading" id="__{slug(text)}">{@render children?.()}</h1>
				{:else if depth === 2}
					<h2 class="cm-md-heading" id="__{slug(text)}">{@render children?.()}</h2>
				{:else if depth === 3}
					<h3 class="cm-md-heading" id="__{slug(text)}">{@render children?.()}</h3>
				{:else if depth === 4}
					<h4 class="cm-md-heading" id="__{slug(text)}">{@render children?.()}</h4>
				{:else if depth === 5}
					<h5 class="cm-md-heading" id="__{slug(text)}">{@render children?.()}</h5>
				{:else if depth === 6}
					<h6 class="cm-md-heading" id="__{slug(text)}">{@render children?.()}</h6>
				{/if}
			</a>
		{/snippet}
		{#snippet link({ href, title, children })}
			{#if href?.startsWith('#')}
				<a href={`#__${href.slice(1)}`} {title}>
					{@render children?.()}
				</a>
			{:else}
				<a {href} {title} target="_blank" rel="noopener noreferrer"> {@render children?.()} </a>
			{/if}
		{/snippet}
		{#snippet code({ lang, text })}
			<pre class="shiki-highlightme" data-lang={lang}><code>{text}</code></pre>
		{/snippet}
	</SvelteMarkdown>
</div>
