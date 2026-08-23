import { Decoration, EditorView, ViewPlugin, ViewUpdate } from "@codemirror/view";
import { syntaxTree } from "@codemirror/language";
import type { Range } from "@codemirror/state";
import { markdownLanguage } from "@codemirror/lang-markdown";

export function headingPlugin() {

    return ViewPlugin.fromClass(class {
        decorations;

        constructor(view: EditorView) {
            this.decorations = this.build(view);
        }

        update(update: ViewUpdate) {
            if (update.docChanged || update.viewportChanged) {
                this.decorations = this.build(update.view);
            }
        }

        build(view: EditorView) {
            const ranges: Range<Decoration>[] = [];

            syntaxTree(view.state).iterate({
                enter: node => {
                    if (node.name.startsWith("ATXHeading")) {
                        const level = Number(node.name.slice(-1));
                        ranges.push(
                            Decoration.mark({
                                tagName: `h${Math.min(6, level)}`,
                                class: `cm-md-heading`,
                            }).range(node.from, node.to)
                        );
                    }
                }
            });

            return Decoration.set(ranges);
        }
    }, {
        decorations: v => v.decorations
    });
}