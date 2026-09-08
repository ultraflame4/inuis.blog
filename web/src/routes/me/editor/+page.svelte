<script lang="ts">
	import { EditorView } from 'codemirror';
	import { asFormData } from '$lib/utils';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api.js';

	import MarkdownContent from '$lib/MarkdownContent.svelte';

	import Editor from '$lib/Editor.svelte';

	const { data } = $props();

	let post_hid = $derived(data.post_hid);
	let published = $derived(data.published);

	let showPreview = $state(false);
	let unsaved_changes = $derived(post_hid == null);
	let title = $state('');
	let view = $state<EditorView | undefined>();
	let preview_content = $state('');

	const warn_unsaved = (e: BeforeUnloadEvent) => {
		if (unsaved_changes) {
			e.preventDefault();
			return 'You have unsaved changes!';
		}
	};

	function updateTitle(view: EditorView) {
		let t = view.state.doc.lineAt(0).text;
		t = t.startsWith('#') ? t.slice(1) : t;
		title = t.trim();
	}

	let saving_post = $state<Promise<void> | null>(null);
	let deleting_post = $state<Promise<void> | null>(null);
	let publishing_post = $state<Promise<void> | null>(null);

	const can_delete = $derived(
		publishing_post == null && post_hid != null && saving_post == null && !published
	);

	async function save_post() {
		async function _save_post() {
			if (!view) return;
			const url = new URL(`/api/me/posts/save`, window.location.origin);
			if (post_hid) url.searchParams.set('hid', post_hid);

			console.log(url.toString());
			const resp = await fetch(url, {
				method: 'POST',
				body: asFormData({
					content: view.state.doc.toString()
				})
			});

			if (resp.ok) {
				unsaved_changes = false;
				saving_post = null;
				const data = await resp.json();
				post_hid = data.post.hid;
				goto(`/me/editor?hid=${post_hid}`, {
					replaceState: true
				});
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
			const url = new URL(`/api/me/posts/${post_hid}/publish`, window.location.origin);
			url.searchParams.set('public', '' + make_public);
			const resp = await fetch(url, {
				method: 'POST'
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
		if (post_hid === null) return;
		if (!confirm('Are you sure you want to delete this post?\nThis action cannot be undone.'))
			return;
		if (deleting_post == null) {
			deleting_post = api.deletePost(fetch, post_hid);
		}
		await deleting_post;
		goto('/me');
	}

	function enablePreview(x: boolean) {
		showPreview = x;
		preview_content = view?.state.doc.toString() ?? '';
	}
</script>

<svelte:window on:beforeunload={warn_unsaved} />
<nav class="classic px-4">
	<a href="/me">&lt&lt; Back</a>
</nav>
<main
	class="grid h-page-content gap-2 overflow-hidden p-4"
	style="grid-template-rows: auto 1fr auto;"
>
	<div class="flex flex-wrap items-center gap-1">
		<h2 class="font-bold">Title: {title}</h2>
		<div class="ml-auto"></div>
		<span class="mx-2 text-sm" class:font-bold={unsaved_changes}>
			{#if unsaved_changes}
				[Unsaved changes]
			{:else}
				[Up to date]
			{/if}
		</span>
		<div class="flex items-center gap-2">
			<button class="classic" class:hidden={!showPreview} onclick={() => enablePreview(false)}>
				View Editor
			</button>
			<button class="classic" class:hidden={showPreview} onclick={() => enablePreview(true)}>
				View Preview
			</button>
			<button
				class="classic"
				disabled={saving_post != null || deleting_post != null}
				onclick={save_post}>Save</button
			>

			{#if !published}
				<button
					class="classic"
					disabled={publishing_post != null || deleting_post != null}
					onclick={() => publish_post(true)}>Publish</button
				>
			{:else}
				<button
					class="classic"
					disabled={publishing_post != null || deleting_post != null}
					onclick={() => publish_post(false)}>Make Private</button
				>
			{/if}
		</div>
	</div>
	<Editor
		bind:view
		hidden={showPreview}
		initialContent={data.content}
		update={(view, changed) => {
			if (changed) unsaved_changes = true;
			updateTitle(view);
		}}
	/>

	<div
		id="preview"
		class="m-0 w-full overflow-auto rounded-lg border px-4"
		class:hidden={!showPreview}
	>
		<MarkdownContent source={preview_content} />
	</div>
	<div class="classic flex flex-wrap">
		<button class="text-danger disabled:text-inherit" disabled={!can_delete} onclick={delete_post}>
			{#if published}
				Cannot delete published posts. Make it private first.
			{:else if post_hid == null}
				-
			{:else}
				Delete Post
			{/if}
		</button>
		<div class="my-0 mr-auto"></div>

		<a href="/posts/{post_hid}" class="my-0">View post page</a>
	</div>
	<br />
</main>

<dialog id="confirmDialog">
	<form method="dialog">
		<p>What is your favorite programming language?</p>
		<input type="text" id="userInput" placeholder="Type here..." />
		<div>
			<button value="cancel">Cancel</button>
			<button id="confirmBtn" value="default">Submit</button>
		</div>
	</form>
</dialog>
