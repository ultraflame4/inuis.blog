<script lang="ts">
	import SvelteMarkdown, { defaultRenderers } from '@humanspeak/svelte-markdown';
	import {
		markedFootnote,
		FootnoteRef,
		FootnoteSection
	} from '@humanspeak/svelte-markdown/extensions/footnote';
	import { codeToHtml, bundledLanguages } from 'shiki/dist/bundle-web.mjs';
	import CopyButton from './CopyButton.svelte';
	import type { Nullish } from './utils';
	import { cn } from 'cn';
	import markedAlert from 'marked-alert';

	import type { Snippet } from 'svelte';
	import { createDirectives, presetDirectiveConfigs } from 'marked-directive';
	import iconDirective from './icon-directive';

	interface Props {
		source: string;
		class?: string;
	}

	const { source, class: _class }: Props = $props();

	function highlightCode() {
		const langs = new Set(Object.keys(bundledLanguages));
		Array.from(document.querySelectorAll('.shiki-highlightme')).map(async (el) => {
			let lang = el.getAttribute('data-lang') ?? '';
			console.log('shiki', lang);

			if (!langs.has(lang)) {
				lang = 'text';
			}

			const text = el.querySelector('code')!.textContent;
			el.classList.remove('shiki-highlightme');
			el.classList.add('shiki-highlighted');
			try {
				const html = await codeToHtml(text, { lang, theme: 'github-dark' });
				// console.log('html', html);
				el.outerHTML = html;
			} catch (e) {
				console.log('shiki', e);
			}
		});
	}

	function stripAlertPrefix(text: Nullish<string>) {
		if (!text) return null;
		return text.replace(/\[!.*?\]\s*/, '');
	}

	$effect(() => {
		console.log('highlightCode', source);
		highlightCode();
		document
			.querySelectorAll('blockquote:has(> .cm-md-alert) > p:not(.cm-md-alert)')
			.forEach((el) => {
				el.innerHTML = stripAlertPrefix(el.innerHTML) ?? '';
			});
	});

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
	type GenericRendererProp<T extends Record<string, any> = {}> = {
		children: Snippet;
		meta: any;
	} & T;

	const extensions = [
		markedAlert(),
		// markedFootnote(),
		createDirectives([...presetDirectiveConfigs, iconDirective])
	];
	const renderers = {
		...defaultRenderers,
		footnoteRef: FootnoteRef,
		footnoteSection: FootnoteSection
	};
</script>

<div class={cn('cm-md-content h-fit max-w-full overflow-hidden', _class)}>
	<SvelteMarkdown {source} {extensions} {renderers}>
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
			<div class="md-escape group relative">
				<div
					class="md-escape absolute top-1 right-1 text-xs opacity-0 drop-shadow-2xl group-hover:opacity-100"
				>
					<button
						class="md-escape relative my-0 btn border-paper btn-tiny font-mono"
						onclick={selectCode}
					>
						Select
					</button>
					<CopyButton {text} />
				</div>

				<pre class="shiki-highlightme" data-lang={lang}><code>{text}</code></pre>
			</div>
		{/snippet}
		{#snippet image({ text, href, title })}
			<img src={href} alt={text} {title} />
		{/snippet}
		{#snippet html_img({ attributes })}
			<img
				{...attributes}
				style:height="{attributes?.['height']}px;"
				style:width="{attributes?.['width']}px;"
			/>
		{/snippet}
		{#snippet blockquote({ text, raw, children })}
			<blockquote>
				{@render children?.()}
			</blockquote>
		{/snippet}
		{#snippet alert({ meta, children }: GenericRendererProp)}
			{@const variant: string = meta['variant']}
			{@const title: string = meta['title']}
			<blockquote>
				<p class="cm-md-alert alert-{variant}">{title}</p>

				{@render children?.()}
			</blockquote>
		{/snippet}
	</SvelteMarkdown>
</div>
