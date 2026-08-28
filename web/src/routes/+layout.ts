
import type { LoadEvent } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import type { Nullish } from '$lib/utils';
import { api } from '$lib/api';


async function fetchUser(event: LoadEvent) {
    const response = await event.fetch('/api/me');
    if (!response.ok) {
        return null;
    }
    return await response.json();
}

export const load: LayoutLoad = async (event) => {

    const { user } = await api.fetchSession(event.fetch)

    return {
        user
    }
};