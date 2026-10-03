
import { api } from '#lib/api.js';
import { nullish } from '#lib/utils.js';
import type { PageLoad } from './$types';

export const load: PageLoad = async (e) => {
    const data = await e.parent();
    const post_hid = e.url.searchParams.get('hid');

    let content = '# Untitled New Post\nThe beginnings of a new post!'
    let content_offset = 0
    let published = false
    if (!nullish(post_hid)) {

        const post = await api.fetchPost(e.fetch, post_hid)
        content = post.Content;
        content_offset = post.ContentOffset;
        published = post.AllowPublic;

    }

    return { ...data, post_hid, content, published, content_offset };
}
