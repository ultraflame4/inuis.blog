<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { highlightSpecialChars } from '@codemirror/view';
	import { markdown } from '@codemirror/lang-markdown';
	import { headingPlugin } from '$lib/editor';
	import { jsonFormdata } from '$lib/utils';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api.js';

	const { data } = $props()

	let post_id = $derived(data.post_id)
	let published = $derived(data.published);
	
	let showPreview = $state(false);
	let unsaved_changes = $derived(post_id == null);
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
				}),
				EditorView.lineWrapping
			]
		});

		updateTitle(view);

		return () => {
			view?.destroy();
		};
	});

	let saving_post = $state<Promise<void> | null>(null);
	let deleting_post = $state<Promise<void> | null>(null);
	let publishing_post = $state<Promise<void> | null>(null);

	const can_delete = $derived(publishing_post == null && post_id != null && saving_post == null && !published);
	
	
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

	async function publish_post(make_public: boolean) {
		async function _publish_post() {
			await save_post();
			const url = new URL(`/api/me/posts/${post_id}/publish`, window.location.origin);
			url.searchParams.set('public', ''+make_public);
			const resp = await fetch(url, {
				method: 'POST',
			});
			if (resp.ok) {
			    published = make_public;
				publishing_post = null;
			}
		}
		if (publishing_post == null) {
			publishing_post = _publish_post();
		}
		await publishing_post;
	}

	async function delete_post() {
		if (!can_delete) return;
		if (post_id === null) return;
		if (!confirm('Are you sure you want to delete this post?\nThis action cannot be undone.')) return;
		if (deleting_post == null) {
			deleting_post = api.deletePost(fetch, post_id);
		}
		await deleting_post
		goto("/me")
	}
</script>

<svelte:window on:beforeunload={warn_unsaved} />
<main class="page-content-height grid gap-2 p-4" style="grid-template-rows: auto 1fr auto;">
	<div class="flex items-center gap-1">
		<h2 class="font-bold">Title: {title}</h2>
		<div class="ml-auto"></div>
		<span class="mx-2 text-sm" class:font-bold={unsaved_changes}>
		    {#if unsaved_changes}
				[Unsaved changes]
			{:else}
				[Up to date]
			{/if}
		</span>
		<button
			class="classic"
			class:hidden={!showPreview}
			onclick={() => (showPreview = false)}
		>
			View Editor
		</button>
		<button class="classic" class:hidden={showPreview} onclick={() => (showPreview = true)}>
			View Preview
		</button>
		<button class="classic" disabled={saving_post != null || deleting_post != null} onclick={save_post}>Save</button>

		{#if !published}
			<button class="classic" disabled={publishing_post != null || deleting_post != null} onclick={()=>publish_post(true)}>Publish</button>
		{:else}
			<button class="classic" disabled={publishing_post != null || deleting_post != null} onclick={()=>publish_post(false)}>Make Private</button>
		{/if}
	</div>
	<div id="editor" class:hidden={showPreview} class="w-full overflow-hidden">
		<style>
			.cm-editor {
				height: 100%;
				width: 100%;
				border: 1px solid silver;
			}
		</style>
	</div>
	<div id="preview" class="ring" class:hidden={!showPreview}></div>

	<div>
	
		<button class="classic text-danger disabled:text-inherit"  disabled={!can_delete} onclick={delete_post}>
			
			{#if published}
				Cannot delete published posts. Make it private first.
			{:else if post_id == null}
			    -
			{:else}
				Delete Post
			{/if}
		</button>
	</div>
</main>

<dialog id="confirmDialog">
  <form method="dialog">
    <p>What is your favorite programming language?</p>
    <input type="text" id="userInput" placeholder="Type here...">
    <div>
      <button value="cancel">Cancel</button>
      <button id="confirmBtn" value="default">Submit</button>
    </div>
  </form>
</dialog>
