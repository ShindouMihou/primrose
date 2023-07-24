<script lang="ts">
    import {Icon} from "@steeze-ui/svelte-icon";
    import {createEventDispatcher} from "svelte";
    import type {IconSource} from "@steeze-ui/radix-icons/types";

    const dispatcher = createEventDispatcher()
    export let icon: IconSource
</script>

<button class="text-sm text-white font-light hover:opacity-60 transition duration-300 hover:cursor-pointer" on:click={() => dispatcher('click')}>
    <Icon src={icon} class="w-5"></Icon>
</button>