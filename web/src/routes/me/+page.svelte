<script lang="ts">
	import { goto } from '$app/navigation';
	import PageNumbers from '$lib/PageNumbers.svelte';
	import { humanizeDate } from '$lib/utils';


    const {data} = $props()

</script>
<main class="classic">
    <p class="text-lg">Hello <span class="italic">{data.user.Username}</span>, writing today?</p>
    <nav class="ring px-4 flex gap-3">
        <a href="/me/editor">New Post</a>
        <div class="mr-auto my-0"></div>
        <a href="/me/change-username">Change Username</a>
        <a href="/me/settings/login">Login Settings</a>
    </nav>

    <h2>Your Posts</h2>
    <ul class="space-y-2">
        {#each  data.user_posts.items as post, index}
            <li class="flex gap-2 ">
                <span class="w-8 m-0 text-right font-bold">{index + 1 + (data.user_posts.page - 1) * data.user_posts.per_page}.</span>
                <div class="my-0">
                    <a href="/me/editor?post_id={post.ID}">{post.Title}</a>
                    <p class="my-0 text-sm">{post.Summary}...</p>
                    <div class="m-0">
                        <span class="text-xs text-gray-500 m-0">
                            Created {humanizeDate(post.CreatedAt)}
                            &nbsp;|&nbsp;
                            Updated {humanizeDate(post.UpdatedAt)}
                        </span>

                    </div>
                </div>
            </li>
        {/each}
    </ul>
    <PageNumbers page_data={data.user_posts} url="/me" />
</main>
