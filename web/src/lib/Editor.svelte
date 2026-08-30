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
	import { languages  } from '@codemirror/language-data';

	interface Props {
		view?: EditorView;
		initialContent?: string;
		contentChanged?: (view: EditorView) => void;
		hidden?: boolean;
	}


	let {  view = $bindable(), initialContent, contentChanged, hidden }: Props = $props();


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
						contentChanged?.(update.view);
					}
				}),
				EditorView.lineWrapping,
			    keymap.of([indentWithTab]),
				theme.of(prefersDark.matches ? darkTheme : lightTheme)
			]
		});

		contentChanged?.(view);
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


<div id="editor" class:hidden={hidden} class="w-full overflow-hidden rounded-lg border grid relative group" >
    <p class="absolute top-2 right-2 z-10 bg-paper-subtle border-paper border px-2 py-1 rounded-lg hover:opacity-0 text-sm invisible group-focus-within:visible">
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
