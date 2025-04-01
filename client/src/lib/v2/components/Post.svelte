<script lang="ts">
    import type {Post} from "$lib/types/post";
    import {Bookmark} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";
    import removeMarkdown from "remove-markdown";
    import PostPadding from "$lib/v2/components/PostPadding.svelte";

    export let post: Post;
</script>

<a href="/read/{post.slug}" class="article-container">
    <PostPadding>
        <div class="flex flex-row justify-between items-center gap-4 text-xs">
            <div class="flex flex-row items-center gap-2">
                <Icon src={Bookmark} size="18"/>
                <p>
                    READ
                </p>
            </div>
        </div>
    </PostPadding>
    <img
            alt="{post.title}'s Hero Image"
            src={post.image}
            class="article-image"
    />
    <PostPadding>
        <div class="flex flex-col gap-2">
            <h2 class={"text-2xl font-medium"}>
                {post.title}
            </h2>
            <p class="text-sm">
                {removeMarkdown(post.content).slice(0, 200)}...
                <span class="relative font-bold group-hover:animate-pulse">Read more</span>
            </p>
        </div>
    </PostPadding>
</a>
