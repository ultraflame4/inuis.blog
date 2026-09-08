import { asFormData } from "./utils";

export namespace api {

    export interface Paginated<T> {
        items: T[];
        page: number;
        per_page: number;
        total_pages: number;
    }

    export interface User {
        ID: string;
        Username: string;
        CreatedAt: string;
        DeletedAt: string;
        UpdatedAt: string;
        // Hid used to fetch user content in the database.
        hid: string

    }

    export interface Post {
        ID: string;
        hid: string;
        CreatedAt: string;
        DeletedAt: string;
        UpdatedAt: string;
        // Title of post
        Title: string
        // Generated summary of post. (First few lines of the post)
        Summary: string
        // Contents of the post in markdown. This also includes the post title.
        Content: string
        // Author id of the post
        UserID: number
        // Author of the post
        User: User
        AllowPublic: boolean

    }

    export interface AppMeta {
        brand_name: string;
        requires_invite_code: boolean;
    }

    export type FetchLike = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>

    export async function fetchAppMeta(f: FetchLike): Promise<AppMeta> {
        const resp = await f('/api/meta');
        if (!resp.ok) {
            console.error("Failed to fetch app meta. \nStatus: " + resp.status, "\nBody:\n" + await resp.text())
            throw new Error('Failed to fetch app meta');
        }
        const app_meta = await resp.json() as AppMeta;
        return app_meta;
    }

    export async function listPosts(f: FetchLike, page: string | number, filter?: string, search?: string): Promise<Paginated<Post>> {

        const params = new URLSearchParams({ page: page.toString(), filter: filter ?? "", q: search ?? "" });

        const resp = await f(`/api/me/posts?${params.toString()}`);
        if (!resp.ok) {
            console.error("Failed to fetch user posts. \nStatus: " + resp.status, "\nBody:\n" + await resp.text())
            throw new Error('Failed to fetch user posts');
        }
        const user_posts: Paginated<any> = await resp.json();
        console.info(user_posts)
        return user_posts;
    }

    export async function deletePost(f: FetchLike, id: string): Promise<void> {
        const resp = await f(`/api/me/posts/${id}`, { method: 'DELETE' });
        if (!resp.ok) {
            console.error("Failed to delete post. \nStatus: " + resp.status, "\nBody:\n" + await resp.text())
            throw new Error('Failed to delete post');
        }

    }

    export async function fetchSession(f: FetchLike): Promise<{ is_elevated: boolean; user: User | null }> {
        const response = await f('/api/me');
        if (!response.ok) {
            return { user: null, is_elevated: false };
        }
        const data = await response.json();
        const user: User = data.user;
        const is_elevated = data.session_elevated;
        return { user, is_elevated };
    }

    export async function fetchProviders(f: FetchLike) {
        const res = await f('/api/oidc/providers');

        if (!res.ok) {
            throw new Error('Failed to load providers');
        }

        return await res.json() as string[];
    }

    export async function disconnectOidcProvider(f: FetchLike, provider: string) {
        const resp = await f(`/api/me/oidc/${provider}/disconnect`);
        if (!resp.ok) {
            console.error("Failed to disconnect provider. \nStatus: " + resp.status, "\nBody:\n" + await resp.text())
            if (resp.status == 409) {
                throw new Error('You need atleast 1 OIDC provider connected!');
            }
            throw new Error('Failed to disconnect provider');
        }
    }

    export async function updateUsername(f: FetchLike, username: string) {
        const resp = await f(`/api/me/username`, {
            method: 'PUT',
            body: asFormData({ username }),
        });
        if (!resp.ok) {
            console.error("Failed to update username. \nStatus: " + resp.status, "\nBody:\n" + await resp.text())
            throw new Error('Failed to update username');
        }
    }

    export async function fetchConnectedProviders(f: FetchLike) {
        const res = await f('/api/me/oidc-providers');

        if (!res.ok) {
            throw new Error('Failed to load connectedproviders');
        }
        const data = await res.json()
        return data['providers'] as string[];
    }


    /**
     * Fetches a public post using it's hid
     * @param f Fetch function
     * @param hid post's hid
     * @returns 
     */
    export async function fetchPublicPost(f: FetchLike, hid: string): Promise<Post> {
        const res = await f(`/api/posts/${hid}`);
        if (!res.ok) {
            throw new Error('Failed to fetch post');
        }
        const data = await res.json();
        return data["post"] as Post;
    }

    /**
     * Fetches a post using it's hid.
     * 
     * This is similar to `fetchPublicPost` except it works for all posts the user has access to.
     * 
     * This includes:
     * - The user's private posts
     * - All public posts
     * 
     * @param f Fetch function
     * @param hid post's hid
     * @returns 
     */
    export async function fetchPost(f: FetchLike, post_hid: string): Promise<Post> {
        const res = await f(`/api/me/posts/${post_hid}`);
        if (!res.ok) {
            throw new Error('Failed to fetch post');
        }
        const data = await res.json();
        return data["post"] as Post;
    }
}
