<script lang="ts">
    import {onMount} from "svelte";
    import {token} from "../../lib/stash";

    onMount(() => {
        token.set('')
        window.location.replace('/')
    })
</script>