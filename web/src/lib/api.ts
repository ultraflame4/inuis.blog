
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

    }

    export type FetchLike = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>

    export async function listPosts(f: FetchLike, page: string | number): Promise<Paginated<any>> {
        const resp = await f(`/api/me/posts?page=${page}`);
        if (!resp.ok) {
            console.error("Failed to fetch user posts. \nStatus: " + resp.status, "\nBody:\n" + await resp.text())
            throw new Error('Failed to fetch user posts');
        }
        const user_posts: Paginated<any> = await resp.json();
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
        const is_elevated = data.is_elevated;
        return { user, is_elevated };
    }
}