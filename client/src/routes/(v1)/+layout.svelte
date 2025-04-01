<script>
    import "../../app.css";
    import {QueryClient, QueryClientProvider} from "@tanstack/svelte-query";
    import {browser} from "$app/environment";
    import {onMount} from "svelte";
    import Lenis from "@studio-freight/lenis";

    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                enabled: browser
            }
        }
    })

    onMount(() => {
        window.lenis = new Lenis()
        function raf(time) {
            lenis.raf(time)
            requestAnimationFrame(raf)
        }

        requestAnimationFrame(raf)
    })
</script>

<svelte:head>
    <meta name="theme-color" content="{ import.meta.env.VITE_SEO_COLOR }"/>
    <link rel="icon" href="{ import.meta.env.VITE_FAVICON }" />
</svelte:head>

<QueryClientProvider client={queryClient}>
    <div class="bg-black text-white flex flex-col min-h-screen inter" lenis>
        <slot />
    </div>
</QueryClientProvider>
