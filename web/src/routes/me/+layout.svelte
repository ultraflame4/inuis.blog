<script lang="ts">
	import { resolve } from '$app/paths';


    const { children} = $props()


</script>
<aside class="w-48 border-r border-r-divider min-h-page-content float-left mr-8">
    <nav class="flex flex-col">

        <a href={resolve("/me")}>Home</a>
        <a href={resolve("/me/editor")}>New Post</a>
    </nav>
</aside>
{@render children()}
