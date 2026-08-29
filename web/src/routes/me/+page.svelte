<script lang="ts">
	import { goto } from '$app/navigation';
	import PageNumbers from '$lib/PageNumbers.svelte';
	import { humanizeDate } from '$lib/utils';


    const {data} = $props()

</script>
<main class="p-4">
    <p class="text-lg my-4 border p-4">Hello <span class="italic">{data.user.Username}</span>, writing today?</p>

    <div class="flex items-baseline classic gap-4 ">
        <h2 >Your Posts</h2>
        <div class="w-0 border bg-black h-4 "></div>
        <a href="/me/editor" class="self-center">/New Post</a>
    </div>
    <ul class="flex flex-wrap items-center justify-center sm:justify-start  gap-4">
        {#each  data.user_posts.items as post, index}
            <li class="grid gap-2 border rounded-2xl w-full h-38 sm:w-96 sm:h-36 px-4 py-3 overflow-hidden" style="grid-template-columns: auto 1fr; grid-template-rows: auto 1fr auto;">
                <span class="w-4 col-1 row-span-2 text-right font-bold mr-1">{index + 1 + (data.user_posts.page - 1) * data.user_posts.per_page}.</span>
            
                <a href="/me/editor?post_id={post.ID}" class="underline font-semibold">{post.Title}</a>
                <p class="my-0 text-sm">{post.Summary.slice(0, 96)}...</p>
                <div class="m-0 col-2 row-3">
                    <span class="text-xs text-gray-500 m-0">
                        <!-- Created {humanizeDate(post.CreatedAt)} -->
                        Updated: {humanizeDate(post.UpdatedAt)}
                    </span>

                </div>
            </li>
        {/each}

    </ul>
    <br/>
    <PageNumbers page_data={data.user_posts} url="/me" />
</main>
