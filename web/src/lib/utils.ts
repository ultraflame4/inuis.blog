import { goto } from '$app/navigation';
import { page } from '$app/state';
import { format, formatDistanceToNow, isToday, isYesterday, parseISO } from 'date-fns';


export type Nullish<T> = T | null | undefined
export function nullish<T>(value: Nullish<T>): value is null | undefined {
    return value === null || value === undefined
}

export function asFormData(obj: any) {
    const formData = new FormData();

    for (const key in obj) {
        if (obj.hasOwnProperty(key)) {
            formData.append(key, obj[key]);
        }
    }

    return formData;
}

export function humanizeDate(date: Date | string): string {

    if (typeof date === 'string') {
        date = new Date(date);
    }
    if (isToday(date)) {
        return formatDistanceToNow(date, { addSuffix: true });
    }

    if (isYesterday(date)) {
        return `Yesterday at ${format(date, 'h:mm a')}`;
    }

    return format(date, 'MMM d, yyyy');
}

export function waitTimeout(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
}


/**
 * Updates the current page's search params (queries) with the key and value
 * @param key Query key
 * @param value Value to set to. If null or undefined, removes it from the search params.
 */
export async function updateSearchParams(key: string, value: string | null | undefined) {
    // await invalidate('data:posts')
    const searchParams = new URLSearchParams(page.url.searchParams);
    if (!nullish(value)) {
        searchParams.set(key, value);
    }
    else {
        searchParams.delete(key)
    }
    await goto(`?${searchParams.toString()}`, { replaceState: true });
}