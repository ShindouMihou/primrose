<script lang="ts">
    import {Cross2} from "@steeze-ui/radix-icons";
    import Alert from "$lib/components/alerts/Alert.svelte";
    import {onMount} from "svelte";

    export let error: Error

    onMount(() => {
        setTimeout(() => document.getElementById('error_alert')?.scroll({ behavior: 'smooth' }), 500);
    })
</script>
<div class="pb-4" id="error_alert">
    <Alert icon={Cross2} title="An error occurred" description={error.message}/>
</div>