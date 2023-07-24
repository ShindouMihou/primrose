<script lang="ts">
    import {ChevronRight, Cross1, EnvelopeClosed, LockClosed} from "@steeze-ui/radix-icons";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {fetchSelf, signIn} from "../../lib/requests/user";
    import {token} from "../../lib/stash";
    import {createQuery} from "@tanstack/svelte-query";
    import {UNAUTHENTICATED} from "../../lib/requests/errors/generic";
    import PrimroseView from "$lib/components/pages/PrimroseView.svelte";
    import {page} from "$app/stores";

    const self = createQuery(['self', $token], () => fetchSelf($token))
    $: if ($self.error === UNAUTHENTICATED) {
        token.set('')
    }

    $: if ($self.data != null) {
        window.location.replace('/')
    }

    let email: string = ''
    let password: string = ''

    let error: string | null = null
    let loading = false
    $: loadingClass = loading ? 'animate-pulse animated duration-700' : ''

    $: if (email ||  password) {
        error = null
    }

    let callback: string = $page.url.searchParams.get('callback')

    async function login() {
        if (loading) {
            return
        }

        if (callback != null) {
            callback = callback.replace( /^[a-zA-Z]{3,5}\:\/{2}[a-zA-Z0-9_.:-]+\//, '');
        }

        if (!callback.startsWith('/')) {
            callback = "/" + callback
        }

        loading = true
        return signIn(email, password)
            .then(value => {
                token.set(value)
                setTimeout(() => window.location.replace(callback != null ? callback : '/'), 512)
            })
            .catch(reason => { error = reason.message; loading = false })
    }
</script>

<PrimroseView title="Login" headline={import.meta.env.VITE_APP_NAME}>
    {#if error != null}
        <div class="border border-[#797979] border-opacity-30 p-4 w-fit rounded bg-gray-300 bg-opacity-5 my-2 w-full md:w-full md:max-w-sm">
            <div class="flex flex-row items-center gap-2 pb-2 text-red-500">
                <Icon src="{Cross1}" size="18"/>
                <h3  class="font-bold">An error occurred</h3>
            </div>
            <p class="text-sm font-light max-w-2xl">{error}</p>
        </div>
    {/if}
    <div class="px-2">
        <div class="flex flex-row gap-4 items-center pt-4 border-b pb-2 md:max-w-sm {loadingClass}">
            <Icon src={EnvelopeClosed} size="16" class="text-white"/>
            <input disabled={loading} type="email" placeholder="Email" class="outline-none bg-transparent" bind:value={email}/>
        </div>
        <div class="flex flex-row gap-4 items-center pt-4 border-b pb-2 md:max-w-sm {loadingClass}">
            <Icon src={LockClosed} size="16" class="text-white"/>
            <input disabled={loading} type="password" placeholder="Password" class="outline-none bg-transparent" bind:value={password}/>
        </div>
        <div class="flex flex-row gap-1 items-center">
            <button disabled={loading} on:click={login} class="{loadingClass} p-2 pr-4 my-4 flex flex-row gap-2 items-center hover:bg-white hover:text-black animated duration-700">
                <Icon src={ChevronRight} size="16"/>
                <p>Login</p>
            </button>
            <a href="/signup{callback !== null ? '?callback='+encodeURIComponent(callback) : ''}" class="hover:opacity-60 animated duration-700">
                <p class="text-sm font-light">or <span class="text-blue-300">create an account here</span>.</p>
            </a>
        </div>
    </div>
</PrimroseView>