import { api } from "$lib/api";
import type { PageLoad, PageLoadEvent } from "./$types";


export const load: PageLoad = async (ev: PageLoadEvent) => {
    const parent = await ev.parent()

    const fetch_post = parent.user ? api.fetchPost : api.fetchPublicPost
    const post_data = await fetch_post(ev.fetch, ev.params.post_id)


    return {
        ...parent,
        post_data,
    }
};