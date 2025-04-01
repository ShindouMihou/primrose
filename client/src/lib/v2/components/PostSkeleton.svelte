<script>
    import {Bookmark} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";
    import PostPadding from "$lib/v2/components/PostPadding.svelte";
</script>

<div class="article-container animate-pulse">
    <PostPadding>
        <div class="flex flex-row justify-between items-center text-xs">
            <div class="flex flex-row items-center gap-2">
                <Icon src={Bookmark} size="18"/>
                <p>
                    READ
                </p>
            </div>
        </div>
    </PostPadding>
    <div
            class="article-image bg-gray-400 h-80"
    />
    <PostPadding>
        <div class="flex flex-col gap-2">
            <div class="w-80 bg-gray-400 p-2 rounded-md my-2 animate-pulse"></div>
            <div class="bg-gray-500 p-1 rounded-md animate-pulse"></div>
            <div class="bg-gray-500 p-1 rounded-md animate-pulse"></div>
            <div class="w-[95%] bg-gray-500 p-1 rounded-md animate-pulse"></div>
        </div>
    </PostPadding>
</div>
