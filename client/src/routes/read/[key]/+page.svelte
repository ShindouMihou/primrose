<script lang="ts">
    import {page} from "$app/stores";
    import removeMarkdown from 'remove-markdown'
    import {createQuery} from "@tanstack/svelte-query";
    import {fetchPost,} from "../../../lib/requests/posts";
    import {NOT_FOUND} from "../../../lib/requests/errors/generic";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {ChevronLeft, ExclamationTriangle} from "@steeze-ui/radix-icons";
    import HeadlineSkeleton from "$lib/components/HeadlineSkeleton.svelte";
    import {toHTML} from "$lib/renderer/markdown.js";
    import Footer from "$lib/components/Footer.svelte";
    import ErrorView from "$lib/components/pages/ErrorView.svelte";
    import type {Post} from "../../../lib/types/post";
    import Alert from "$lib/components/alerts/Alert.svelte";
    import {token} from "../../../lib/stash";

    const DISALLOW_SELECT_CODEBLOCKS =
        import.meta.env.VITE_DISALLOW_SELECT_CODEBLOCKS == null
            ? false
            : import.meta.env.VITE_DISALLOW_SELECT_CODEBLOCKS === "true";

    let key = $page.params.key
    let post = createQuery(['post', key], () => fetchPost(key, $token), {
        retry: (failureCount, error) => {
            if (error === NOT_FOUND) {
                return false
            }
            return failureCount < 5;
        },
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
                document.getElementById("header")?.classList.remove("hidden", "md:block")
            } else {
                document.getElementById("header")?.classList.remove("backdrop-blur")
                document.getElementById("header")?.classList.add("hidden", "md:block")
            }
        })
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
    {#if !$post.isLoading && $post.error == null}
        <title>{$post.data?.title ?? import.meta.env.VITE_APP_NAME}</title>
        <meta name="title" content={$post.data?.title ?? import.meta.env.VITE_APP_NAME} />
        <meta name="image" content={$post.data?.image} />
        <meta name="og:image" content={$post.data?.image} />
        <meta name="description" content={removeMarkdown($post.data.content)} />
        <meta name="article:published_time" content={$post.data?.created_at} />
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

{#if !$post.isLoading}
    {#if $post.error == null}
        <div class="relative">
            <div class="fixed transition ease-in-out duration-400 w-screen z-30 hidden md:block" id="header">
                <a href="/" class="flex flex-row items-center px-6 py-4 hover:opacity-80 duration-400 transition ease-in-out">
                    <Icon src={ChevronLeft} size="24"/>
                    <p class="playfair uppercase font-bold text-lg">{import.meta.env.VITE_APP_NAME}</p>
                </a>
            </div>
            <div class="w-full h-screen bg-gray-400 relative overflow-hidden group-hover:backdrop-blur-xl" id="hero">
                <div class="absolute h-screen w-full bg-cover bg-center bg-gray-400 bg-no-repeat" style="background-image: url('{$post.data.image}');"></div>
                <div class="relative drop-shadow shadow-white backdrop-blur group-hover:backdrop-blur-2xl transition ease-in-out duration-300 bg-black bg-opacity-30 h-full overflow-hidden">
                    <div class="relative md:my-64">
                        <h2 class="text-[16rem] leading-none md:text-[32rem] font-bold break-all text-justify select-none">{$post.data.title}</h2>
                    </div>
                </div>
            </div>
            <div class="w-screen bg-black text-white px-12 2xl:px-24 py-8">
                {#if !$post.data.published}
                    <div class="py-4">
                        <Alert icon="{ExclamationTriangle}"
                               title="You are reading a draft"
                               description="Primrose will automatically refresh this post once every 2.5 seconds until it is published, if it consumes too much bandwidth,
                                we recommend either editing on WiFi or not opening this page until you are done editing."/>
                    </div>
                {/if}
                <h2 class="p-1 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">{$post.data.title}</h2>
                <div class="mkdown playfair">
                    {@html html($post.data.content)}
                </div>
            </div>
        </div>
        <Footer/>
    {:else}
        <ErrorView err={$post.error}/>
    {/if}
{:else}
    <HeadlineSkeleton/>
{/if}