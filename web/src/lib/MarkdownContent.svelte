<script lang="ts">
    import SvelteMarkdown from '@humanspeak/svelte-markdown'


    interface Props {
        source: string
    }
    const {source}: Props = $props()
</script>


<div class="w-fit h-fit md-content">
    <SvelteMarkdown {source}>
       
    </SvelteMarkdown>
</div>
