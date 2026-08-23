<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { highlightSpecialChars } from '@codemirror/view';
	import { markdown } from '@codemirror/lang-markdown';
	import { headingPlugin } from '$lib/editor';
	import { jsonFormdata } from '$lib/utils';
	import { goto } from '$app/navigation';

	const { data } = $props()

	let post_id = $state(data.post_id)
	
	let showPreview = $state(false);
	let unsaved_changes = $state(post_id == null);
	let title = $state('');
	let view = $state<EditorView | null>(null);

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

		view = new EditorView({
			parent: document.getElementById('editor') ?? undefined,
			doc: data.content,
			extensions: [
				basicSetup,
				highlightSpecialChars(),
				markdown({}),
				headingPlugin(),
				EditorView.updateListener.of((update) => {
					if (update.docChanged) {
						unsaved_changes = true;
						updateTitle(update.view);
					}
				})
			]
		});

		updateTitle(view);

		return () => {
			view?.destroy();
		};
	});

	let saving_post = $state<Promise<void> | null>(null);

	async function save_post() {
		async function _save_post() {
		    if (!view) return
			const url = new URL(`/api/me/posts/save`, window.location.origin);
			if (post_id) url.searchParams.set('id', post_id);

			console.log(url.toString());
			const resp = await fetch(url, {
				method: 'POST',
				body: jsonFormdata({
					content: view.state.doc.toString(),
				}),
				
			});

			if (resp.ok) {
				unsaved_changes = false;
				saving_post = null;
				const data = await resp.json();
				post_id = data.post.ID;
				goto(`/me/editor?post_id=${post_id}`);
			}
			
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
