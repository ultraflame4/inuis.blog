import { api } from "$lib/api";
import type { PageLoad, PageLoadEvent } from "./$types";


export const load: PageLoad = async (event: PageLoadEvent) => {
    const parent = await event.parent();
    const connected_providers = await api.fetchConnectedProviders(event.fetch);

    return {
        ...parent,
        connected_providers,
    }
};