// Modified markdown Renderer taken from ultraflame4.github.io/v7 (clyh.me/v7) source code

import { visit } from 'unist-util-visit';
import * as unified from 'unified';
import type * as mdast from 'mdast';
// Needed for typesscript support
import type * as _ from 'remark-directive'
import type * as hast from 'hast'
import remarkParse from 'remark-parse';
import remarkRehype from 'remark-rehype';
import rehypeStringify from 'rehype-stringify';
import remarkDirective from 'remark-directive';
import remarkFrontmatter from 'remark-frontmatter';
import { matter } from 'vfile-matter';




export const remarkIconDirective: unified.Plugin<any, mdast.Root> = () => {
    return (tree) => {
        visit(tree, (node) => {
            if (
                node.type === 'textDirective' &&
                node.name === 'icon'
            ) {
                // const icon_name = node.children[0].data
                const get_text = (a?: mdast.PhrasingContent) => {
                    if (a === undefined) return ''
                    if (a.type === "text") return a.value
                    if (a.type === "textDirective") return ":" + a.name
                }
                const value = node.children.map(get_text).reduce((a, b) => (a ?? '') + (b ?? ''), '')
                if (!value && value !== "false") return
                node.children = [
                    {
                        type: 'html',
                        value: `<span class="icon remark-icon-directive"><image src="https://api.iconify.design/${value}.svg"></span>`
                    },
                ]
            }
        })
    }
}


// Inserts custom attribute into the html elements
export const rehypeInsertAttrs: unified.Plugin<any, hast.Root> = () => {
    return (tree) => {
        visit(tree, (node) => {
            if (node.type !== "element") return
            if (node.tagName.startsWith("h")) {
                let class_ = node.properties["class"] ?? ''
                node.properties["class"] = class_ + "cm-md-heading "
            }

        })
    }
}


export const remarkExtractFrontmatter: unified.Plugin<void[], mdast.Root> = () => {
    return (_tree, file) => {
        matter(file, { strip: true })
    }
}

export function render_md(content: string) {
    const processor = unified.unified()
        .use(remarkParse)
        .use(remarkFrontmatter, ['yaml', 'toml'])
        .use(remarkExtractFrontmatter)
        .use(remarkDirective)
        .use(remarkIconDirective)
        .use(remarkRehype, { allowDangerousHtml: true })
        .use(rehypeInsertAttrs)
        .use(rehypeStringify, { allowDangerousHtml: true })
    return processor.processSync(content);
}
