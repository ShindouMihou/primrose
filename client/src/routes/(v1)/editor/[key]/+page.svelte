<script>
    import PrimroseEditor from "$lib/components/pages/PrimroseEditor.svelte";
    import {page} from "$app/stores";
    let key = $page.params.key;
</script>

<PrimroseEditor key={key}/>