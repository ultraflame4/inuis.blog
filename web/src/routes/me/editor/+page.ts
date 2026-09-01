
import { nullish } from '$lib/utils';
import type { PageLoad } from './$types';

export const load: PageLoad = async (e) => {
    const data = await e.parent();
    const post_hid = e.url.searchParams.get('hid');

    let content = '# Example post title\nSome text content'
    let published = false
    if (!nullish(post_hid)) {

        const resp = await e.fetch(`/api/me/posts/${post_hid}`);
        if (resp.ok) {
            const data = await resp.json();
            content = data.post.Content;
            published = data.post.AllowPublic;
        }
    }

    return { ...data, post_hid, content, published };
}
