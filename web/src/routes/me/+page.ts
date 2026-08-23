
import { api } from '$lib/api';
import { nullish } from '$lib/utils';
import type { PageLoad } from './$types';

export const load: PageLoad = async (e) => {
    const data = await e.parent();

    const user_posts = await api.listPosts(e.fetch);

    return { ...data, user_posts };
}