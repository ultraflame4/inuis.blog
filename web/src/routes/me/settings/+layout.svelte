<script lang="ts">
	let { children } = $props();
</script>

<nav class="classic mx-4">
	<a href="/me/settings">Back to Settings</a>
</nav>
{@render children()}
