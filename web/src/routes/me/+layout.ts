// app/+layout.ts
import { nullish } from '#lib/utils.js';
import { redirect } from '@sveltejs/kit';


export const load = async (event) => {
    const data = await event.parent()
    const user = data.user
    if (nullish(user)) {
        throw redirect(302, '/')
    }

    return { ...data, user }
};
