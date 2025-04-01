<script lang="ts">
    import {ExclamationTriangle} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";
    import PostPadding from "$lib/v2/components/PostPadding.svelte";
    export let error: string;
</script>

<div class="article-container bg-red-500 bg-opacity-20">
    <PostPadding>
        <div class="flex flex-row justify-between items-center pb-2 text-xs">
            <div class="flex flex-row items-center gap-2">
                <Icon src={ExclamationTriangle} size="18"/>
                <p>
                    ERROR
                </p>
            </div>
            <p class="text-gray-200">/failed-to-load-post</p>
        </div>
    </PostPadding>
    <img
            alt="Error Cat"
            src="https://media1.tenor.com/m/KLAHPW3-BaQAAAAC/loading-error.gif"
            class="article-image bg-red-400 h-96"
    />
    <PostPadding>
        <div class="flex flex-col gap-2 py-2">
            <h2 class={"text-2xl font-medium"}>
                An error occurred when we tried to retrieve the posts.
            </h2>
            <p class="text-xs">Error: {error}</p>
        </div>
    </PostPadding>
</div>
