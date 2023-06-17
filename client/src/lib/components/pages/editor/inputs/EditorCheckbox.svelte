<script lang="ts">
    import {Check, Cross2} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";

    export let value: boolean;
</script>
<div on:click={() => value = !value} class="hover:cursor-pointer flex flex-row gap-2 items-center pt-4 pb-2 w-full">
    <Icon src={value ? Check : Cross2} size="24" class="text-red-500 bg-black"/>
    <p class="text-lg playfair"><slot/></p>
</div>