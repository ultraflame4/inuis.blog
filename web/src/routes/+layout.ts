import type { LayoutLoad } from './$types';
import { api } from '$lib/api';
export const load: LayoutLoad = async (event) => {

    const { user } = await api.fetchSession(event.fetch)
    const providers = await api.fetchProviders(event.fetch)
    const app_meta = await api.fetchAppMeta(event.fetch)

    return {
        user,
        providers,
        app_meta
    }
};
