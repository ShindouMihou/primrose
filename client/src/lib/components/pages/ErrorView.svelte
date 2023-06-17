<script lang="ts">
    import {ChevronLeft} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";

    export let err: Error;
</script>

<div class="relative">
    <div class="fixed transition ease-in-out duration-400 w-screen z-30 hidden md:block" id="header">
        <a href="/" class="flex flex-row items-center px-6 py-4 hover:opacity-80 duration-400 transition ease-in-out">
            <Icon src={ChevronLeft} size="24"/>
            <p class="playfair uppercase font-bold text-lg">{import.meta.env.VITE_APP_NAME}</p>
        </a>
    </div>
    <div class="w-full h-screen bg-gray-400 relative overflow-hidden group-hover:backdrop-blur-xl" id="hero">
        <div class="absolute h-screen w-full bg-cover bg-center bg-red-500 bg-no-repeat"></div>
        <div class="relative drop-shadow shadow-white backdrop-blur group-hover:backdrop-blur-2xl transition ease-in-out duration-300 bg-black bg-opacity-30 h-full overflow-hidden">
            <div class="relative md:my-64">
                <h2 class="text-[16rem] leading-none md:text-[32rem] font-bold break-all text-justify select-none">Error</h2>
            </div>
        </div>
    </div>
    <div class="w-screen bg-black text-white px-12 2xl:px-24 py-8">
        <h2 class="p-1 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">An error occurred</h2>
        <p class="text-lg font-light playfair break-words p-1 mt-0.5 transition group-hover:bg-white group-hover:bg-opacity-5 ease-in-out duration-300 w-fit">
            {err.message}
        </p>
    </div>
</div>