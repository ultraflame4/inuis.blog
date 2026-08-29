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
