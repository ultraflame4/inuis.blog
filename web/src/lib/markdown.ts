// Modified markdown Renderer taken from ultraflame4.github.io/v7 (clyh.me/v7) source code

import { visit } from 'unist-util-visit';
import * as unified from 'unified';
import type * as mdast from 'mdast';
import * as _ from 'lodash-es';
import type * as hast from 'hast'
import remarkParse from 'remark-parse';
import remarkRehype from 'remark-rehype';
import rehypeStringify from 'rehype-stringify';
import remarkDirective from 'remark-directive';
import remarkFrontmatter from 'remark-frontmatter';
import { matter } from 'vfile-matter';




export const remarkIconDirective: unified.Plugin<any, mdast.Root> = () => {
    return (tree) => {
        visit(tree, "textDirective", (node) => {
            if (node.name === 'icon') {
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


const alertRegex = /^\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]/i;
export const remarkAlerts: unified.Plugin<any, mdast.Root> = () => {
    return (tree) => {
        visit(tree, "blockquote", (node) => {
            // Check the first line for alerts.
            const firstParagraph = node.children.find(x => x.type == "paragraph");
            if (!firstParagraph) return
            const firstNode = firstParagraph?.children[0]
            const match = firstNode?.type == 'text' && firstNode.value.match(alertRegex);
            if (!match) return

            let [full, inner] = match
            let variant = inner.toLowerCase()
            let stripped = _.trim(firstNode.value.slice(full.length), "\n")
            firstNode.value = stripped

            // Trim the paragraph of newlines and whitespaces
            firstParagraph.children = _.dropWhile(firstParagraph.children, x => {
                if (x.type == "break") return true
                if (x.type == "text" && x.value.length == 0) return true
                return false
            })

            node.children = [
                {
                    type: "html",
                    value: `<p class="cm-md-alert alert-${variant}">${variant}</p>`
                },
                ...node.children
            ]

        })
    }
}


export const remarkHeading: unified.Plugin<any, mdast.Root> = () => {

    return (tree) => {
        visit(tree, "heading", (node) => {
            const text_content = node.children.reduce((a, b) => a + (b.type == "text" ? b.value : ''), "")
            const slugified = _.kebabCase(text_content)
            node.data = {

                hProperties: {
                    id: "md-" + slugified,
                }
            }

        })
    }
}

export const remarkAnchorHeading: unified.Plugin<any, mdast.Root> = () => {

    return (tree) => {
        visit(tree, "link", (node) => {
            if (!node.url.startsWith("#")) return
            node.url = "#md-" + node.url.slice(1)
        })
    }
}

const headingRegex = /^h[0-9]/i;
// Inserts custom attribute into the html elements
export const rehypeInsertAttrs: unified.Plugin<any, hast.Root> = () => {

    return (tree) => {
        visit(tree, (node) => {
            if (node.type !== "element") return
            if (node.tagName.match(headingRegex)) {
                let class_ = node.properties["class"] ?? ''
                node.properties["class"] = class_ + "cm-md-heading hover:underline"
                console.log(node)
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
        .use(remarkHeading)
        .use(remarkAnchorHeading)
        .use(remarkAlerts)
        .use(remarkDirective)
        .use(remarkIconDirective)
        .use(remarkRehype, { allowDangerousHtml: true })
        .use(rehypeInsertAttrs)
        .use(rehypeStringify, { allowDangerousHtml: true })
    return processor.processSync(content);
}
