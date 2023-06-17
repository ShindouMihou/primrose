<script>
    import PrimroseView from "$lib/components/pages/PrimroseView.svelte";
    import Alert from "$lib/components/alerts/Alert.svelte";
    import {
        Component1,
        ExitFullScreen,
        Person,
        QuestionMarkCircled
    } from "@steeze-ui/radix-icons";
    import Footer from "$lib/components/Footer.svelte";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {createQuery} from "@tanstack/svelte-query";
    import {fetchListWithUnpublished} from "$lib/requests/posts";
    import {token} from "$lib/stash";
    import {Dates} from "$lib/utils/dates";
    import {UNAUTHENTICATED} from "$lib/requests/errors/generic";
    import {fetchSelf} from "$lib/requests/user";
    import ErrorView from "$lib/components/pages/ErrorView.svelte";
    import {onMount} from "svelte";

    const self = createQuery(['self'], () => fetchSelf($token))
    onMount(() => {
        function logout() {
            token.set('')
            setTimeout(() => window.location.replace('/login'), 250)
        }
        $: if ($self.error === UNAUTHENTICATED || ($self.isSuccess && $self.data == null)) {
            logout()
            return
        }
        $: if ($posts.error === UNAUTHENTICATED) {
            logout()
        }
    })

    $: posts = createQuery(['list', 'unpublished'], () => fetchListWithUnpublished($token), {
        enabled: ($self.isSuccess && $self.data != null)
    })
</script>

{#if $posts.error != null || $self.error != null}
    <ErrorView err={$posts.error == null ? $self.error : $posts.error}/>
{:else}
    <div>
        <PrimroseView title="Creator's Page" headline={import.meta.env.VITE_APP_NAME}>
            <div class="py-4">
                <Alert icon={QuestionMarkCircled}
                       title="What to do here?"
                       description="The Creator's Page is the place to create, edit and delete posts that you've made. Primrose is designed with the same editor
                        as its predecessor and supports many of the same features. You can start by creating a post by clicking the Create Post button."
                />
            </div>
            <div class="flex flex-col md:flex-row-reverse gap-4 justify-between">
                <div class="block border border-[#797979] border-opacity-30 bg-gray-300 bg-opacity-5 p-4">
                    <h2 class="p-1 mb-4 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">
                        Creator
                    </h2>
                    <div class="block mb-2 hover:opacity-80 transition duration-300 ease-in-out border border-[#797979] border-opacity-30 p-6 pr-12 w-full rounded bg-gray-300 bg-opacity-5">
                        <Icon src="{Person}" size="32"/>
                        <h3 class="font-bold pt-4 text-2xl">
                            {#if $self.isLoading}
                                Loading
                            {:else}
                                {$self.data.name}
                            {/if}
                        </h3>
                    </div>
                    <div class="flex flex-wrap gap-2">
                        <a href="/editor" class="block hover:opacity-80 transition duration-300 ease-in-out border border-[#797979] border-opacity-30 p-6 pr-12 w-full rounded bg-gray-300 bg-opacity-5">
                            <Icon src="{Component1}" size="32"/>
                            <h3 class="font-bold pt-4 text-2xl">Create Post</h3>
                        </a>
                        <a href="/logout" class="block hover:opacity-80 transition duration-300 ease-in-out border border-[#797979] border-opacity-30 p-6 pr-12 w-full rounded bg-red-400 bg-opacity-5">
                            <Icon src="{ExitFullScreen}" size="32"/>
                            <h3 class="font-bold pt-4 text-2xl">Logout</h3>
                        </a>
                    </div>
                </div>
                <div class="py-4 md:hidden"></div>
                <div class="block border border-[#797979] border-opacity-30 bg-gray-300 bg-opacity-5 p-4 w-full">
                    <h2 class="p-1 group-hover:bg-white group-hover:bg-opacity-5 w-fit transition ease-in-out duration-300 text-2xl font-bold playfair break-words">
                        Posts
                    </h2>
                    <div class="flex flex-col gap-8 py-4 md:pl-4">
                        {#if !$posts.isLoading}
                            {#each $posts.data.data as post, i}
                                <a href={"/editor/" + post.id } class="hover:scale-105 transition ease-in-out duration-700 flex flex-row gap-4 items-center w-full">
                                    <div class="w-7 h-7 text-center text-blue-200 playfair">
                                        <p>{i + 1}</p>
                                    </div>
                                    <div class="flex flex-row-reverse md:flex-row w-full gap-4">
                                        <img src="{post.image}" class="w-16 h-16 bg-gray-200 object-cover"/>
                                        <div class="max-w-lg w-full">
                                            <h2 class="font-bold text-3xl max-w-lg playfair">{post.title}</h2>
                                            <p class="font-light text-zinc-200">{post.published ? "Published" : "Draft"} at {Dates.toDateString(new Date(post.created_at))}</p>
                                        </div>
                                    </div>
                                </a>
                            {/each}
                        {:else}
                            <div class="flex flex-row gap-4 pl-4">
                                <div class="w-7 h-7 text-center text-blue-200 playfair">
                                    <p>1</p>
                                </div>
                                <div class="flex flex-row gap-4">
                                    <div class="max-w-sm w-full">
                                        <h2 class="font-bold text-3xl max-w-sm playfair">Loading</h2>
                                        <p class="font-light text-zinc-200">Draft at {Dates.toDateString(new Date())}</p>
                                    </div>
                                    <div class="w-16 h-16 bg-gray-200 object-cover"></div>
                                </div>
                            </div>
                        {/if}
                    </div>
                </div>
            </div>
        </PrimroseView>
        <Footer/>
    </div>
{/if}