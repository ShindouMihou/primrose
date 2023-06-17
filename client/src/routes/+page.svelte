<script lang="ts">
    import {createQuery} from "@tanstack/svelte-query";
    import {fetchList} from "../lib/requests/posts";
    import Headline from "$lib/components/Headline.svelte";
    import HeadlineSkeleton from "$lib/components/HeadlineSkeleton.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import ErrorView from "$lib/components/pages/ErrorView.svelte";

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
        <div>
            {#if !$latest.isLoading}
                {#if $latest.error == null}
                    {#if headline != null}
                        <Headline post={headline}/>
                    {/if}
                {:else}
                    <ErrorView err={$latest.error}/>
                {/if}
            {:else}
                <HeadlineSkeleton/>
            {/if}
        </div>
        <div class="gap-4">
            {#if !$latest.isLoading && $latest.error  == null}
                {#each $latest.data.data.slice(1) as post}
                    <Headline post={post}/>
                {/each}
            {/if}
        </div>
    </div>
    <Footer/>
</div>