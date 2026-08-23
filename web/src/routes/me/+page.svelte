<script lang="ts">
	import { goto } from '$app/navigation';
	import { humanizeDate } from '$lib/utils';


    const {data} = $props()

</script>
<main class="classic">
    <p>How is life, {data.user.Username}?</p>
    <nav>
        <a href="/me/editor">Create Post</a>
    </nav>

    <h2>Your Posts</h2>
    <ul class="space-y-2">
        {#each  data.user_posts.items as post, index}
            <li class="flex gap-2 ">
                <span class="w-8 m-0 text-right font-bold">{index + 1}.</span>
                <div class="my-0">
                    <a href="/me/editor?post_id={post.ID}">{post.Title}</a>
                    <p class="my-1 text-sm">{post.Summary}...</p>
                    <div >
                        <span class="text-xs text-gray-500">
                            Created {humanizeDate(post.CreatedAt)}
                            &nbsp;|&nbsp;
                            Updated {humanizeDate(post.UpdatedAt)}
                        </span>
                        
                    </div>
                </div>
            </li>
        {/each}
    </ul>
</main>