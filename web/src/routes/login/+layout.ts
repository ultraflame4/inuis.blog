// app/+layout.ts
import { redirect } from '@sveltejs/kit';


export const load = async (event) => {
    const data = await event.parent()

    if (data.user) {
        throw redirect(302, '/me')
    }

    return data
};