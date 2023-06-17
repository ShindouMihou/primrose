<script lang="ts">
    import {Icon} from "@steeze-ui/svelte-icon";
    import {ChevronRight, Cross1} from "@steeze-ui/radix-icons";
    import Alert from "$lib/components/alerts/Alert.svelte";

    export let error: string;
</script>

<Alert icon={Cross1} title="An error occurred" description={error}>
    <div on:click={() => { window.history.back() }}  class="flex flex-row items-center gap-2 w-fit rounded pt-4 hover:opacity-80 animated hover:cursor-pointer animate-in fade-in duration-1000">
        <Icon src="{ChevronRight}" size="14"/>
        <h2 class="font-light text-sm">Go Back</h2>
    </div>
    <a href="/" class="flex flex-row items-center gap-2 w-fit rounded pt-4 hover:opacity-80 animated hover:cursor-pointer animate-in fade-in duration-1000">
        <Icon src="{ChevronRight}" size="14"/>
        <h2 class="font-light text-sm">Home</h2>
    </a>
</Alert>