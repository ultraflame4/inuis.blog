import type { PageLoadEvent, PageLoad } from "./$types";

export const load: PageLoad = async (ev: PageLoadEvent) => {
    const parent_data = await ev.parent()
    const data = { ...parent_data }
    const invite_code = ev.url.searchParams.get('invite')
    if (invite_code) {
        data.invite_code = invite_code
    }
    return data
};