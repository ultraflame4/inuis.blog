import {
    createDirectives,
    presetDirectiveConfigs,
    type DirectiveConfig
} from 'marked-directive'
// export const iconDirective: unified.Plugin<any, mdast.Root> = () => {
//     return (tree) => {
//         visit(tree, (node) => {
//             if (
//                 node.type === 'textDirective' &&
//                 node.name === 'icon'
//             ) {
//                 // const icon_name = node.children[0].data
//                 const get_text = (a?: mdast.PhrasingContent) => {
//                     if (a === undefined) return ''
//                     if (a.type === "text") return a.value
//                     if (a.type === "textDirective") return ":" + a.name
//                 }
//                 const value = node.children.map(get_text).reduce((a, b) => (a ?? '') + (b ?? ''), '')
//                 if (!value && value !== "false") return
//                 node.children = [
//                     {
//                         type: 'html',
//                         value: `<span class="icon remark-icon-directive"><image src="https://api.iconify.design/${value}.svg"></span>`
//                     },
//                 ]
//             }
//         })
//     }
// }
// 
export const iconDirective: DirectiveConfig = {
    level: 'inline',
    marker: ":",
    tag: "span",
    renderer(token) {

        if (token.meta.name === 'icon' || token.meta.name === 'i') {

            const value = token.text
            if (!value && value !== "false") return
            return `<span class="icon remark-icon-directive"><image src="https://api.iconify.design/${value}.svg"></span>`

        }
    }
}

export default iconDirective