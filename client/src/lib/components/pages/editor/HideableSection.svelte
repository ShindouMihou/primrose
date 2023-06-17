<script lang="ts">
    import {Icon} from "@steeze-ui/svelte-icon";
    import {EyeClosed, EyeOpen} from "@steeze-ui/radix-icons";

    export let show: boolean
    export let title: string
</script>

<div class="block border border-[#797979] border-opacity-30 bg-gray-300 bg-opacity-5 p-4 mb-4">
    <div class="flex flex-row items-center gap-2 hover:cursor-pointer" on:click={() => show = !show}>
        <div>
            <Icon src={show ? EyeClosed : EyeOpen} size="16"/>
        </div>
        <h2 class="p-1 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">
            {title}
        </h2>
    </div>
    <div class="{show ? 'block' : 'hidden'}">
        <slot/>
    </div>
</div>