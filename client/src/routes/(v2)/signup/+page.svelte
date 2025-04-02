<script lang="ts">
    import {Cross1, EnvelopeClosed, LockClosed, Person, Symbol} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {fetchSelf, signUp} from "$lib/requests/user";
    import {token} from "$lib/stash";
    import {createQuery} from "@tanstack/svelte-query";
    import {UNAUTHENTICATED} from "$lib/requests/errors/generic";
    import {page} from "$app/stores";
    import Header from "$lib/v2/components/Header.svelte";
    import PostPadding from "$lib/v2/components/PostPadding.svelte";
    import {colors} from "$lib/utils/colors";

    export let color: string;

    if (color == null) {
        color = colors.randomBackgroundColor()
    }

    const self = createQuery(['self', $token], () => fetchSelf($token))
    $: if ($self.error === UNAUTHENTICATED) {
        token.set('')
    }

    $: if ($self.data != null) {
        window.location.replace('/')
    }

    let name: string = ''
    let email: string = ''
    let password: string = ''

    let error: string | null = null
    let loading = false

    let callback: string | null = $page.url.searchParams.get('callback')

    $: loadingClass = loading ? 'animate-pulse animated duration-700' : ''

    $: if (email || name || password) {
        error = null
    }

    async function create() {
        if (loading) {
            return
        }

        if (callback != null) {
            callback = callback.replace( /^[a-zA-Z]{3,5}:\/{2}[a-zA-Z0-9_.:-]+\//, '');

            if (!callback.startsWith('/')) {
                callback = "/" + callback
            }
        }

        loading = true
        return signUp(name, email, password)
            .then(() => setTimeout(() => window.location.replace('/login' + (callback == null ? '' : '?callback='+encodeURI(callback))), 512))
            .catch(reason => { error = reason.message; loading = false })
    }
</script>

<div class="flex flex-col" lenis>
    <div class="max-w-sm m-auto">
        <Header/>
        <main class={"py-8 pb-16 flex flex-col gap-8 w-full"}>
            <div class="article-container hover:bg-zinc-950 w-[20rem]">
                <PostPadding>
                    <div class="flex flex-row justify-between items-center gap-4 text-xs">
                        <div class="flex flex-row items-center gap-2">
                            <Icon src={LockClosed} size="18"/>
                            <p>
                                SIGN UP
                            </p>
                        </div>
                    </div>
                </PostPadding>
                <div
                        class="article-image {color} h-0.5">
                </div>
                <PostPadding>
                    {#if error != null}
                        <div class="border border-[#797979] border-opacity-30 p-4 rounded bg-gray-300 bg-opacity-5 my-2 w-full md:w-full md:max-w-sm">
                            <div class="flex flex-row items-center gap-2 pb-2 text-red-500">
                                <Icon src="{Cross1}" size="18"/>
                                <h3  class="font-bold">An error occurred</h3>
                            </div>
                            <p class="text-sm font-light max-w-2xl">{error}</p>
                        </div>
                    {/if}
                    <div class="px-2">
                        <div class="flex flex-row gap-4 items-center pt-4 border-b pb-2 md:max-w-sm {loadingClass}">
                            <Icon src={Person} size="16" class="text-white"/>
                            <input disabled={loading} type="text" placeholder="Username" class="outline-none bg-transparent w-full" bind:value={name}/>
                        </div>
                        <div class="flex flex-row gap-4 items-center pt-4 border-b pb-2 md:max-w-sm {loadingClass}">
                            <Icon src={EnvelopeClosed} size="16" class="text-white"/>
                            <input disabled={loading} type="email" placeholder="Email" class="outline-none bg-transparent w-full" bind:value={email}/>
                        </div>
                        <div class="flex flex-row gap-4 items-center pt-4 border-b pb-2 md:max-w-sm {loadingClass}">
                            <Icon src={LockClosed} size="16" class="text-white"/>
                            <input disabled={loading} type="password" placeholder="Password" class="outline-none bg-transparent w-full" bind:value={password}/>
                        </div>
                        <div class="flex flex-row gap-1 items-center">
                            <button disabled={loading} on:click={create} class="{loadingClass} p-2 px-4 my-4 items-center bg-white text-black hover:opacity-50 transition ease-in-out duration-500">
                                {#if loading}
                                    <Icon src={Symbol} size="16" class="text-black animate-spin"/>
                                {:else}
                                    SIGN UP
                                {/if}
                            </button>
                        </div>
                        <a href="/login{callback !== null ? '?callback='+encodeURIComponent(callback) : ''}" class="hover:opacity-60 animated duration-700">
                            <p class="text-xs font-light">Login into an existing account</p>
                        </a>
                    </div>
                </PostPadding>
            </div>
        </main>
    </div>
</div>
