<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { highlightSpecialChars } from '@codemirror/view';
	import { markdown } from '@codemirror/lang-markdown';
	import { headingPlugin } from '$lib/editor';

	let showPreview = $state(false);
	let title = $state('');
	let view = $state<EditorView | null>(null);
	onMount(() => {
		view = new EditorView({
			parent: document.getElementById('editor') ?? undefined,
			doc: '# Your title here\nSome text content',
			extensions: [basicSetup, highlightSpecialChars(), markdown({}), headingPlugin()]
		});

		return () => {
			view?.destroy();
		};
	});

	$effect(() => {
		if (view) {
			let t = view.state.doc.lineAt(0).text
			title = t.startsWith('#') ? t.slice(1) : t;
			
		}
	});
</script>

<main class="page-content-height grid gap-2 p-4" style="grid-template-rows: auto 1fr auto;">
	<div class="flex gap-1 items-center">
		<h2 class="font-bold">{title}</h2>
		<button class="classic ml-auto" class:hidden={!showPreview} onclick={() => (showPreview = false)}>
		    See Editor
		</button>
		<button class="classic ml-auto" class:hidden={showPreview} onclick={() => (showPreview = true)}>
    		See Preview
		</button>
	</div>
	<div id="editor" class:hidden={showPreview}>
		<style>
			.cm-editor {
				height: 100%;
				border: 1px solid silver;
			}
		</style>
	</div>
	<div id="preview" class="ring" class:hidden={!showPreview}></div>

	<div>
		<button class="classic">Save</button>
	</div>
</main>
