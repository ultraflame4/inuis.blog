
import { api } from '$lib/api';
import { nullish } from '$lib/utils';
import type { PageLoad } from './$types';


export const load: PageLoad = async ({ url, ...e }) => {
    const data = await e.parent();


    const page = url.searchParams.get("page") ?? "1";
    const search = url.searchParams.get("q") ?? "";
    const filter = url.searchParams.get("filter") ?? undefined;

    const user_posts = await api.listPosts(
        e.fetch,
        page,
        filter,
        search
    );


    return { ...data, user_posts };
}