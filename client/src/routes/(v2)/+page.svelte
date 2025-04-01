<script lang="ts">
    import {createQuery} from "@tanstack/svelte-query";
    import {fetchList} from "$lib/requests/posts";
    import Header from "$lib/v2/components/Header.svelte";
    import PostSkeleton from "$lib/v2/components/PostSkeleton.svelte";
    import Post from "$lib/v2/components/Post.svelte";
    import ErrorPost from "$lib/v2/components/ErrorPost.svelte";

    const latest = createQuery(['latest'], fetchList, {
        staleTime: 10 * 24 * 60 * 1000
    })

    $: headline = $latest.data?.data?.at(0)
</script>
<svelte:head>
    <title>{import.meta.env.VITE_SEO_TITLE}</title>
    <meta name="title" content="{import.meta.env.VITE_SEO_TITLE}"/>
    <meta name="image" content="{import.meta.env.VITE_SEO_IMAGE}"/>
    <meta name="og:image" content="{import.meta.env.VITE_SEO_IMAGE}"/>
    <meta name="description" content="{import.meta.env.VITE_SEO_DESCRIPTION}"/>
    <meta name="og:type" content="website"/>
</svelte:head>

<div class="flex flex-col" lenis>
    <div class="">
        <Header/>
        <main class={"py-8 pb-16 flex flex-col gap-8 w-full"}>
            <div class={"w-full flex flex-col gap-4"}>
                {#if !$latest.isLoading}
                    {#if $latest.error == null && $latest.data != null}
                        {#each $latest.data.data as post}
                            <Post post={post}/>
                        {/each}
                    {:else}
                        <ErrorPost error={$latest.error}/>
                    {/if}
                {:else}
                    <PostSkeleton/>
                {/if}
            </div>
        </main>
    </div>
</div>
