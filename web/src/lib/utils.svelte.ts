

import { onMount } from 'svelte';

export function useMediaQuery(query: string) {
    let matches = $state(false);

    onMount(() => {
        const media = window.matchMedia(query);

        const update = () => {
            matches = media.matches;
        };

        update();
        media.addEventListener('change', update);

        return () => {
            media.removeEventListener('change', update);
        };
    });

    return {
        get matches() {
            return matches;
        }
    };
}