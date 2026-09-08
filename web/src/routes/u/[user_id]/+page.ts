import { api } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async (event) => {
    const parent = await event.parent()

    const author_hid = event.params["user_id"]
    const author = await api.fetchUser(event.fetch, author_hid)

    const page = event.url.searchParams.get("page") ?? "1"
    const author_posts = await api.listUserPosts(event.fetch, author_hid, page)

    return {
        ...parent,
        author,
        author_posts,
    }
}