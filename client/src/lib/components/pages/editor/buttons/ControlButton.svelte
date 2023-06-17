<script lang="ts">
    import type {IconProps} from "@steeze-ui/svelte-icon/Icon.svelte";
    import {createEventDispatcher} from "svelte";
    import {Icon} from "@steeze-ui/svelte-icon";

    let dispatch = createEventDispatcher()

    export let icon: IconProps;
    export let enabled: boolean;
    export let background: string;
</script>
<div on:click={() => dispatch('click')}
     class="block {(enabled ?? true) ? 'hover:cursor-pointer' : ''}
            hover:opacity-80 transition duration-300 ease-in-out border border-[#797979]
            border-opacity-30 p-6 pr-12 w-full rounded {(background ?? 'bg-gray-300')}
            bg-opacity-5">
    <Icon src="{icon}" size="32"/>
    <h3 class="font-bold pt-4 text-2xl"><slot/></h3>
</div>