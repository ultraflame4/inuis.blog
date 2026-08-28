// place files you want to import through the `$lib` alias in this folder.


export namespace app {
    export async function elevateSession(): Promise<void> {
        window.open("/me/elevate")
    }
}