<script lang="ts">
	import SvelteMarkdown from '@humanspeak/svelte-markdown';
	import { codeToHtml, bundledLanguages } from 'shiki/dist/bundle-web.mjs';
	import CopyButton from './CopyButton.svelte';
	import type { Nullish } from './utils';

	interface Props {
		source: string;
	}

	const { source }: Props = $props();

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

	function extractAlert(text: Nullish<string>) {
		if (!text) return null;
		const s = text.toLowerCase();
		if (s.startsWith('[!note]')) return 'note';
		if (s.startsWith('[!tip]')) return 'tip';
		if (s.startsWith('[!important]')) return 'important';
		if (s.startsWith('[!warning]')) return 'warn';
		if (s.startsWith('[!caution]')) return 'caution';
		return null;
	}
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
			<div class="md-escape group relative">
				<div
					class="md-escape absolute top-1 right-1 opacity-0 drop-shadow-2xl group-hover:opacity-100 text-xs"
				>
					<button class="btn btn-tiny md-escape my-0 border-paper font-mono relative " onclick={selectCode}>
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
				style:height="{attributes['height']}px;"
				style:width="{attributes['width']}px;"
			/>
		{/snippet}
		{#snippet blockquote({ text, raw, children })}
			<blockquote>
				{#if extractAlert(text)}
					{@const alert_tip = extractAlert(text)}
					<p class="cm-md-alert alert-{alert_tip}">{alert_tip}</p>
				{/if}
				{@render children?.()}
			</blockquote>
		{/snippet}
	</SvelteMarkdown>
</div>
