<script lang="ts">
    import {Icon} from "@steeze-ui/svelte-icon";
    import type {IconProps} from "@steeze-ui/svelte-icon/Icon.svelte";

    export let lock: boolean;
    export let value: string;
    export let placeholder: string;
    export let icon: IconProps;
</script>

<div class="flex flex-row gap-4 items-center pt-4 border-b pb-2 w-full">
    <Icon src={icon} size="16" class="text-white"/>
    <input disabled={lock} type="text" placeholder="{placeholder}" class="outline-none bg-transparent w-full" bind:value={value}/>
</div>