
import type { LoadEvent } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';


async function fetchUser(event: LoadEvent) {
    const response = await event.fetch('/api/me');
    if (!response.ok) {
        return null;
    }
    return await response.json();
}

export const load: LayoutLoad = async (event) => {

    const user_resp = await fetchUser(event);

    return {
        user: user_resp?.["user"] as any | null | undefined
    }
};