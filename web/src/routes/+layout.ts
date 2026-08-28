
import type { LoadEvent } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import type { Nullish } from '$lib/utils';
import { api } from '$lib/api';


export const load: LayoutLoad = async (event) => {

    const { user } = await api.fetchSession(event.fetch)
    const providers = await api.fetchProviders(event.fetch)

    return {
        user,
        providers
    }
};