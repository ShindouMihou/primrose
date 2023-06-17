<script lang="ts">
    import {Icon} from "@steeze-ui/svelte-icon";
    import type {IconProps} from "@steeze-ui/svelte-icon/Icon.svelte";
    export let icon: IconProps;
    export let title: string;
    export let description: string;
</script>

<div class="border border-[#797979] border-opacity-30 p-4 w-fit rounded bg-gray-300 bg-opacity-5 w-full">
    <div class="flex flex-row items-center gap-2 pb-4 text-red-500">
        <Icon src="{icon}" size="18"/>
        <h3  class="font-bold">{title}</h3>
    </div>
    {#if description != null}
        <p class="text-sm font-light">{description}</p>
    {/if}
    <slot/>
</div>