import { api } from "$lib/api";
import type { PageLoad, PageLoadEvent } from "./$types";


export const load: PageLoad = async (ev: PageLoadEvent) => {
    const parent = await ev.parent()

    const post_data = await api.fetchPost(ev.fetch, ev.params.post_id)


    return {
        ...parent,
        post_data,
    }
};