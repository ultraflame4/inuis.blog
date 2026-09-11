<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { highlightSpecialChars, keymap } from '@codemirror/view';
	import { indentWithTab } from '@codemirror/commands';
	import { Compartment } from '@codemirror/state';
	import { markdown } from '@codemirror/lang-markdown';
	import { headingPlugin } from '$lib/editor';
	import { githubDark } from '@ddietr/codemirror-themes/github-dark';
	import { githubLight } from '@ddietr/codemirror-themes/github-light';
	import { useMediaQuery } from '$lib/utils.svelte.js';
	import { languages } from '@codemirror/language-data';
	import { nullish } from './utils';

	interface Props {
		view?: EditorView;
		initialContent?: string;
		update?: (view: EditorView, contentChanged: boolean) => void;
		/**
		 * Called when the  save keybind is pressed. (Ctrl+s)
		 *
		 * If not set, Ctrl+S keybind will revert back to default behaviour (open browser save dialog for most browsers)
		 * @param view
		 */
		handleSave?: (view: EditorView) => void;
		hidden?: boolean;
	}

	let {
		view = $bindable(),
		initialContent,
		update: contentChanged,
		hidden,
		handleSave
	}: Props = $props();

	let showPreview = $state(false);

	const theme = new Compartment();
	const darkTheme = githubDark;
	const lightTheme = githubLight;
	let prefersDark = useMediaQuery('(prefers-color-scheme: dark)');

	onMount(() => {
		view = new EditorView({
			parent: document.getElementById('editor') ?? undefined,
			doc: initialContent,
			extensions: [
				basicSetup,
				highlightSpecialChars(),
				markdown({
					codeLanguages: languages
				}),
				headingPlugin(),
				EditorView.updateListener.of((update) => {
					if (update.docChanged) {
						contentChanged?.(update.view, true);
					}
				}),
				EditorView.lineWrapping,
				keymap.of([
					indentWithTab,
					{
						key: 'Ctrl-s',
						run: (view) => {
							if (nullish(handleSave)) {
								return false;
							}
							handleSave(view);
							return true;
						}
					}
				]),
				theme.of(prefersDark.matches ? darkTheme : lightTheme)
			]
		});

		contentChanged?.(view, false);
		return () => {
			view?.destroy();
		};
	});

	$effect(() => {
		view?.dispatch({
			effects: theme.reconfigure(prefersDark.matches ? darkTheme : lightTheme)
		});
	});
</script>

<div id="editor" class:hidden class="group relative grid w-full overflow-hidden rounded-lg border">
	<p
		class="invisible absolute top-2 right-2 z-10 rounded-lg border border-paper bg-paper-subtle px-2 py-1 text-sm group-focus-within:visible hover:opacity-0"
	>
		Press <kbd>ESC</kbd>+<kbd>Tab</kbd> escape Editor focus.
	</p>
	<style>
		.cm-editor {
			height: 100%;
			width: 100%;
			overflow: auto;
		}
	</style>
</div>
