<script lang="ts">
    import PrimroseView from "$lib/components/pages/PrimroseView.svelte";
    import Alert from "$lib/components/alerts/Alert.svelte";
    import {
        CrumpledPaper, ExclamationTriangle, Image, LightningBolt, Link1, Rocket, Text
    } from "@steeze-ui/radix-icons";
    import Footer from "$lib/components/Footer.svelte";
    import {createMutation, createQuery} from "@tanstack/svelte-query";
    import {token} from "$lib/stash";
    import {UNAUTHENTICATED} from "$lib/requests/errors/generic";
    import {fetchSelf} from "$lib/requests/user";
    import ErrorView from "$lib/components/pages/ErrorView.svelte";
    import {deletePost, fetchPostForEditor, savePost} from "../../requests/posts";
    import {onMount} from "svelte";
    import type {Post} from "../../types/post";
    import EditorInformationAlert from "$lib/components/pages/editor/EditorInformationAlert.svelte";
    import HideableSection from "$lib/components/pages/editor/HideableSection.svelte";
    import ControlButton from "$lib/components/pages/editor/buttons/ControlButton.svelte";
    import EditorTextField from "$lib/components/pages/editor/inputs/EditorTextField.svelte";
    import EditorCheckbox from "$lib/components/pages/editor/inputs/EditorCheckbox.svelte";
    import SavingAlert from "$lib/components/pages/editor/alerts/SavingAlert.svelte";
    import SimpleErrorAlert from "$lib/components/alerts/SimpleErrorAlert.svelte";

    export let key: string | null;

    const self = createQuery(['self'], () => fetchSelf($token))
    onMount(() => {
        function logout() {
            token.set('')
            window.location.replace('/')
        }
        $: if ($self.error === UNAUTHENTICATED || ($self.isSuccess && $self.data == null)) {
            logout()
            return
        }

        setTimeout(() => {
            //@ts-ignore
            window.lenis.destroy();
            console.log('[site] destroyed lenis')
        }, 500)

        enableTabsOnEditor();
        enableCtrlS();
    })

    let listenersAttached = false;

    function enableCtrlS() {
        document.addEventListener("keydown", (event) => {
            if (event.ctrlKey && event.key === "s") {
                event.preventDefault();
                save();
            }
        });
    }

    function enableTabsOnEditor() {
        if (listenersAttached) return;

        setTimeout(() => {
            listenersAttached = true;
            document
                .querySelector("textarea")!
                .addEventListener("keydown", function (event) {
                    if (event.key === "Tab") {
                        event.preventDefault();
                        const start = this.selectionStart;
                        const end = this.selectionEnd;

                        this.value =
                            this.value.substring(0, start) +
                            "\t" +
                            this.value.substring(end);
                        this.selectionStart = this.selectionEnd = start + 1;
                    }
                });
        }, 1000);
    }

    let lock: boolean = false;
    let title = ""
    let content = ""
    let slug = ""
    let image = ""
    let published = false

    const post = createQuery(['post', key], () => fetchPostForEditor(key, $token))

    const mutator = createMutation(['post', key], (newPost: Post) => savePost($token, newPost), {
        onSuccess: data => {
            if (key === null || key === "") {
                key = data.id
            }
            $post.refetch()
        },
        onError: () => lock = false
    })

    post.subscribe((post) => {
        if (!post.isLoading) {
            if (post.isSuccess) {
                paste(post.data)
            }
            lock = false
        }
    })
    function transfer(): Post {
        let copy: Post = structuredClone($post.data)
        copy.title = title
        copy.slug = slug
        copy.image = image
        copy.content = content
        copy.published = published
        return copy
    }

    function paste(post: Post) {
        title = post.title
        slug = post.slug
        image = post.image ?? ""
        content = post.content
        published = post.published
        key = post.id

        let url = new URL(window.location.toString());
        url.pathname = "/editor/" + key
        window.history.replaceState({}, "", url)
    }

    function save() {
        $mutator.mutate(transfer());
    }

    let showDeleteWarning = false
    $: deleteWarningClass = showDeleteWarning ? "block" : "hidden"

    async function del() {
        if (key != null && key !== '') {
            let response = await deletePost($token, key)
            if (response?.ack ?? false) {
                setTimeout(() => window.location.replace("/creator"), 250)
            }
        }
    }

    function open() {
        if ($post.data.slug != null && $post.data.slug !== '') {
            document.getElementById('_open_new_tab')?.click()
        }
    }

    let showContentEditor = false
    let showControls = false
    let showSettings = false

    $: hasKey = key !== null && key !== ''
</script>

{#if $post.error != null || $self.error != null}
    <ErrorView err={$post.error == null ? $self.error : $post.error}/>
{:else}
    <div>
        {#if $post.data?.slug != null && $post.data.slug !== ""}
            <a href="/read/{$post.data.slug}" target="_blank" id="_open_new_tab">{""}</a>
        {/if}
        <PrimroseView title="Editor" headline={import.meta.env.VITE_APP_NAME}>
            <EditorInformationAlert/>
            {#if $mutator.isLoading}
                <SavingAlert/>
            {:else}
                {#if $mutator.isError}
                    <SimpleErrorAlert error={$mutator.error}/>
                {/if}
            {/if}
            <HideableSection bind:show={showSettings} title="Settings">
                <div class="flex flex-wrap gap-2 mt-4">
                    <div class="w-full bg-gray-400 h-96">
                        <img src="{image}" class="w-full h-full bg-gray-400 object-cover" alt="Hero">
                    </div>
                    <EditorTextField icon={Text} placeholder="Title" bind:value={title} bind:lock/>
                    <EditorTextField icon={Image} placeholder="Image" bind:value={image} bind:lock/>
                    <EditorTextField icon={Link1} placeholder="Slug" bind:value={slug} bind:lock/>
                    <EditorCheckbox bind:value={published}>Published</EditorCheckbox>
                </div>
            </HideableSection>
            <HideableSection bind:show={showControls} title="Controls">
                {#if showDeleteWarning}
                    <div class="py-4">
                        <Alert icon="{ExclamationTriangle}"
                               title="You are about to delete this post!"
                               description="Are you sure that you want to delete this post? This will be an irreversible action!">
                            <div class="flex flex-row gap-4 pt-4" id="delete_warning">
                                <button class="text-gray-500 font-bold text-lg" on:click={del}>Continue</button>
                                <button class="text-blue-500 font-bold text-lg" on:click={() => showDeleteWarning = !showDeleteWarning}>Cancel</button>
                            </div>
                        </Alert>
                        {(() => { setTimeout(() => document.getElementById('delete_warning')?.scroll({ behavior: 'smooth' }), 500); return "";})()}
                    </div>
                {/if}
                <div class="flex flex-wrap gap-2 mt-4">
                    <ControlButton icon="{Rocket}"
                                   on:click={save}>Save</ControlButton>
                    <ControlButton icon="{LightningBolt}"
                                   on:click={open}
                                   enabled={hasKey}>Open</ControlButton>
                    <ControlButton icon="{CrumpledPaper}"
                                   on:click={() => showDeleteWarning = !showDeleteWarning}
                                   background={'bg-red-500'}
                                   enabled={hasKey}>Delete</ControlButton>
                </div>
            </HideableSection>
            <HideableSection bind:show={showContentEditor} title="Content">
                <div class="flex flex-col gap-8 py-4 pl-1 border border-[#7a7a7a] border-opacity-20 mt-2">
                        <textarea
                                class="text-neutral-50 bg-transparent p-4 outline-none text-base placeholder:text-neutral-600 min-h-screen resize-none selection:text-black selection:bg-white"
                                bind:value={content}
                                placeholder="The beautiful moon is charming always."
                                id="text"
                                disabled={lock}/>
                </div>
            </HideableSection>
        </PrimroseView>
        <Footer/>
    </div>
{/if}