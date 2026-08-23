
export type Nullish<T> = T | null | undefined
export function nullish<T>(value: Nullish<T>): value is null | undefined {
    return value === null || value === undefined
}