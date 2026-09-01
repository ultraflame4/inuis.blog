// place files you want to import through the `$lib` alias in this folder.

import { api } from "./api";
import { waitTimeout } from "./utils";


export namespace app {
    /**
     * Returns true if session was elevated successfully
     */
    export async function elevateSession(): Promise<boolean> {

        const c = await api.fetchSession(fetch)
        if (c.is_elevated) return true;

        const child = window.open("/me/elevate", undefined, "popup=true")
        if (!child) {
            throw new Error("Failed to open elevation popup")
        }
        child.focus()
        console.log("Elevation", child)

        while (!child.closed) {
            await waitTimeout(0)
        }
        console.log("Elevation window closed")
        const s = await api.fetchSession(fetch)
        return s.is_elevated
    }
}
