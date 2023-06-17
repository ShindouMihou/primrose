<script lang="ts">
    import {colors} from "../../utils/colors";
    import {ChevronLeft} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {onMount} from "svelte";

    export let headline, title: string;
    export let color: string;

    if (color == null) {
        color = colors.randomBackgroundColor()
    }

    onMount(() => {
        window.addEventListener('scroll', function () {
            const target = document.getElementById("hero")
            if (target == null) {
                return
            }
            if (window.scrollY > (target.offsetTop + target.offsetHeight)) {
                document.getElementById("header")?.classList.add("backdrop-blur")
                document.getElementById("header")?.classList.remove("hidden", "md:block")
            } else {
                document.getElementById("header")?.classList.remove("backdrop-blur")
                document.getElementById("header")?.classList.add("hidden", "md:block")
            }
        })
    })
</script>
<div class="relative">
    <div class="fixed transition ease-in-out duration-400 w-screen z-30 hidden md:block" id="header">
        <div on:click={() => { document.referrer ? window.location = document.referrer : history.back() }} class="flex hover:cursor-pointer flex-row items-center px-6 py-4 hover:opacity-80 duration-400 transition ease-in-out">
            <Icon src={ChevronLeft} size="24"/>
            <p class="playfair uppercase font-bold text-lg">{import.meta.env.VITE_APP_NAME}</p>
        </div>
    </div>
    <div class="w-full h-screen bg-gray-400 relative overflow-hidden group-hover:backdrop-blur-xl" id="hero">
        <div class="absolute h-screen w-full bg-cover bg-center {color} bg-no-repeat"></div>
        <div class="relative drop-shadow shadow-white backdrop-blur group-hover:backdrop-blur-2xl transition ease-in-out duration-300 bg-black bg-opacity-30 h-full overflow-hidden">
            <div class="relative md:my-64">
                <h2 class="leading-none text-[32rem] font-bold break-all text-justify select-none">{headline}</h2>
            </div>
        </div>
    </div>
    <div class="w-screen bg-black text-white px-12 2xl:px-24 py-8">
        <h2 class="p-1 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">{title}</h2>
        <slot/>
    </div>
</div>