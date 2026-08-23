<script lang="ts">
	import { goto } from '$app/navigation';


    const {data} = $props()

</script>
<main class="classic">
    <p>How is life, {data.user.Username}?</p>
    <a href="/me/editor">Create Post</a>
</main>