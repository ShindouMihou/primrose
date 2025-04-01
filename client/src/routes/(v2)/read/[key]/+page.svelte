<script lang="ts">
    import {page} from "$app/stores";
    import removeMarkdown from 'remove-markdown'
    import {createQuery} from "@tanstack/svelte-query";
    import {fetchPost,} from "$lib/requests/posts";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {Bookmark, ExclamationTriangle} from "@steeze-ui/radix-icons";
    import {toHTML} from "$lib/renderer/markdown.js";
    import type {Post} from "$lib/types/post";
    import Alert from "$lib/components/alerts/Alert.svelte";
    import {token} from "$lib/stash";
    import {AUTHENTICATED_RETRY} from "$lib/requests/retries/retries";
	import type { PageData } from "./$types";
    import Comments from "$lib/v2/components/comments/Comments.svelte";
    import PostPadding from "$lib/v2/components/PostPadding.svelte";
    import Header from "$lib/v2/components/Header.svelte";
    import PostSkeleton from "$lib/v2/components/PostSkeleton.svelte";
    import ErrorPost from "$lib/v2/components/ErrorPost.svelte";

    export let data: PageData;

    const DISALLOW_SELECT_CODEBLOCKS =
        import.meta.env.VITE_DISALLOW_SELECT_CODEBLOCKS == null
            ? false
            : import.meta.env.VITE_DISALLOW_SELECT_CODEBLOCKS === "true";

    let key = $page.params.key
    let post = createQuery(['post', key], () => fetchPost(key, $token), {
        retry: AUTHENTICATED_RETRY,
        refetchInterval: (data: Post) => (!data?.published ?? true) ? 2500 : false
    })

    $: $post.error ? (() => console.log($post.error))() : null;

    $: !$post.isLoading && $post.error == null ? (() => {
        let url = new URL(window.location.toString());
        url.pathname = "/read/" + $post.data.slug
        window.history.replaceState({}, "", url)

        window.addEventListener('scroll', function () {
            const target = document.getElementById("hero")
            if (target == null) {
                return
            }
            if (window.scrollY > (target.offsetTop + target.offsetHeight)) {
                document.getElementById("header")?.classList.add("backdrop-blur")
            } else {
                document.getElementById("header")?.classList.remove("backdrop-blur")
            }
        })

        if (window.location.hash !== "") {
            setTimeout(() => {
                const element = document.querySelector(window.location.hash)
                window.lenis.scrollTo(element)
                console.info('scrolled to', element)
            }, 600)
        }
    })() : null;

    function html(content: string) {
        const translated = toHTML(content);

        if (translated.error && !translated.content) {
            return translated.error;
        }

        return translated.content!;
    }
</script>

<svelte:head>
    {#if data.preload != null}
        {#if data.preload.post != null}
            <title>{data.preload.post?.title ?? import.meta.env.VITE_APP_NAME}</title>
            <meta name="title" content={data.preload.post?.title ?? import.meta.env.VITE_APP_NAME} />
            <meta name="image" content={data.preload.post?.image} />
            <meta name="og:image" content={data.preload.post?.image} />
            <meta name="description" content={removeMarkdown(data.preload.post?.content ?? "")} />
            <meta name="article:published_time" content={data.preload.post?.created_at ?? ""} />
            <meta name="preloaded" content="true"/>
        {/if}
    {:else}
        {#if !$post.isLoading && $post.error == null}
        <title>{$post.data?.title ?? import.meta.env.VITE_APP_NAME}</title>
        <meta name="title" content={$post.data?.title ?? import.meta.env.VITE_APP_NAME} />
        <meta name="image" content={$post.data?.image} />
        <meta name="og:image" content={$post.data?.image} />
        <meta name="description" content={removeMarkdown($post.data?.content ?? "")} />
        <meta name="article:published_time" content={$post.data?.created_at} />
        {/if}
    {/if}
    <meta name="og:type" content="article" />
    <meta name="article:author" content={import.meta.env.VITE_DISPLAY_NAME} />
    <meta name="author" content={import.meta.env.VITE_DISPLAY_NAME} />
    <meta name="twitter:card" content="summary_large_image"/>
    <link
            rel="stylesheet"
            href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.4.0/styles/atom-one-dark.min.css"
    />
    {#if DISALLOW_SELECT_CODEBLOCKS}
        <style>
            .mkdown pre {
                user-select: none;
                -web-kit-user-select: none;
                -ms-user-select: none;
            }
        </style>
    {/if}
</svelte:head>

<Header/>
<main class={"py-8 pb-16 flex flex-col gap-8 w-full"}>
    {#if !$post.isLoading}
        {#if $post.error == null && $post.data != null}
            <div class="article-container hover:bg-background">
                <PostPadding>
                    <div class="flex flex-row justify-between items-center gap-4 text-xs">
                        <div class="flex flex-row items-center gap-2">
                            <Icon src={Bookmark} size="18"/>
                            <p>
                                READ
                            </p>
                        </div>
                        <p>{$post.data.title}</p>
                    </div>
                </PostPadding>
                <img
                        alt="{$post.data.title}'s Hero Image"
                        src={$post.data.image}
                        class="article-image"
                />
                <PostPadding>
                    {#if !$post.data.published}
                        <div class="py-4">
                            <Alert icon="{ExclamationTriangle}"
                                   title="You are reading a draft"
                                   description="Primrose will automatically refresh this post once every 2.5 seconds until it is published, if it consumes too much bandwidth,
                                we recommend either editing on WiFi or not opening this page until you are done editing."/>
                        </div>
                    {/if}
                    <div class="flex flex-col gap-2">
                        <p class="mkdown">
                            {@html html($post.data.content)}
                        </p>
                    </div>
                    <Comments id={$post.data.id} slug={$post.data.slug}/>
                </PostPadding>
            </div>
        {:else}
            <ErrorPost err={$post.error}/>
        {/if}
    {:else}
        <PostSkeleton/>
    {/if}
</main>
