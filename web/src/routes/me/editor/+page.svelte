<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { highlightSpecialChars } from '@codemirror/view';
	import { markdown } from '@codemirror/lang-markdown';
	import { headingPlugin } from '$lib/editor';

	let showPreview = $state(false);
	let unsaved_changes = $state(true);
	let title = $state('');

	const warn_unsaved = (e: BeforeUnloadEvent) => {
		if (unsaved_changes) {
			e.preventDefault();
			return 'You have unsaved changes!';
		}
	};
	onMount(() => {
		function updateTitle(view: EditorView) {
			let t = view.state.doc.lineAt(0).text;
			t = t.startsWith('#') ? t.slice(1) : t;
			title = t.trim();
		}

		let view = new EditorView({
			parent: document.getElementById('editor') ?? undefined,
			doc: '# Example post title\nSome text content',
			extensions: [
				basicSetup,
				highlightSpecialChars(),
				markdown({}),
				headingPlugin(),
				EditorView.updateListener.of((update) => {
					if (update.docChanged) {
						updateTitle(view);
					}
				})
			]
		});

		updateTitle(view);

		return () => {
			view.destroy();
		};
	});

	let saving_post = $state<Promise<void> | null>(null);

	async function save_post() {
		async function _save_post() {
			
		}
		if (saving_post == null) {
			saving_post = _save_post();
		}
		await saving_post;
	}
</script>

<svelte:window on:beforeunload={warn_unsaved} />
<main class="page-content-height grid gap-2 p-4" style="grid-template-rows: auto 1fr auto;">
	<div class="flex items-center gap-1">
		<h2 class="font-bold">Title: {title}</h2>
		<button
			class="classic ml-auto"
			class:hidden={!showPreview}
			onclick={() => (showPreview = false)}
		>
			View Editor
		</button>
		<button class="classic ml-auto" class:hidden={showPreview} onclick={() => (showPreview = true)}>
			View Preview
		</button>
		<button class="classic" disabled={saving_post != null} onclick={save_post}>Save</button>
		<button class="classic">Publish</button>
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

	<div></div>
</main>
