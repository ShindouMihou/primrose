<script lang="ts">
    import type {Post} from "../types/post";
    import removeMarkdown from "remove-markdown";
    import {colors} from "$lib/utils/colors.js";

    export let post: Post
</script>
<a href="/read/{post.slug}" class="group select-none">
    <div class="w-full h-screen bg-gray-400 relative overflow-hidden group-hover:backdrop-blur-xl">
        <div class="absolute h-screen w-full bg-cover bg-center bg-gray-400 bg-no-repeat" style="background-image: url('{post.image}');"></div>
        <div class="relative drop-shadow shadow-white backdrop-blur group-hover:backdrop-blur-2xl transition ease-in-out duration-300 bg-black bg-opacity-30 h-full overflow-hidden">
            <div class="relative md:my-64">
                <h2 class="text-[16rem] leading-none md:text-[32rem] font-bold break-all text-justify">{post.title}</h2>
            </div>
        </div>
    </div>
    <div class="w-full {colors.randomBackgroundColor()} text-black relative overflow-hidden px-12 2xl:px-24 py-8">
        <h2 class="p-1 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">{post.title}</h2>
        <p class="text-lg font-light playfair break-words p-1 mt-0.5 transition group-hover:bg-white group-hover:bg-opacity-5 ease-in-out duration-300 w-fit">{removeMarkdown(post.content).slice(0, 200)}...
            <span class="relative">
                <span class="relative font-bold group-hover:animate-pulse">Read more</span>
            </span>
        </p>
    </div>
</a>