<script lang="ts">

    const {data} = $props()
    
</script>
<main class="classic">
    <p>How is life, {data.user.name}?</p>
    <a href="/me/posts/new">Create Post</a>
</main>