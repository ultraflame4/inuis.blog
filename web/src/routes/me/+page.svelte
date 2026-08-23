<script lang="ts">
	import { goto } from '$app/navigation';


    const {data} = $props()

</script>
<main class="classic">
    <p>How is life, {data.user.Username}?</p>
    <a href="/me/editor">Create Post</a>

    <ul class="space-y-2">
        {#each  data.user_posts.items as post, index}
            <li class="flex gap-2 ">
                <span>{index + 1}</span>
                <div class="my-0">
                    <a href="/me/editor?post_id={post.ID}">{post.Title}</a>
                    <p class="my-1 text-sm">{post.Summary}...</p>
                </div>
            </li>
        {/each}
    </ul>
</main>