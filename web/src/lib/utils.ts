
export type Nullish<T> = T | null | undefined
export function nullish<T>(value: Nullish<T>): value is null | undefined {
    return value === null || value === undefined
}

export function jsonFormdata(obj: any) {
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
    return date.toLocaleString();
}

export function waitTimeout(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
}