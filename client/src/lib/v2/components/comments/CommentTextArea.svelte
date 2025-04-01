<script lang="ts">
    import CommentElement from "$lib/components/CommentElement.svelte";
    import type {Comment} from "$lib/types/comment";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {ChevronDown, ChevronUp, Cross2, PaperPlane} from "@steeze-ui/radix-icons";
    import autosize from 'autosize'
    import {createEventDispatcher, onMount} from "svelte";
    import type {User} from "$lib/types/user";

    const dispatcher = createEventDispatcher()

    onMount(() => {
        autosize(document.querySelector('#text-bar')!!)
    })

    export let selectedComment: Comment | null
    export let callback: string
    export let self: User | null


    let contents = "";
    let avatar: string = createAvatar()

    function createAvatar() {
        if (self == null) {
            return "https://api.dicebear.com/9.x/bottts-neutral/svg?seed=primrose"
        }

        let name = self.name
        return "https://api.dicebear.com/9.x/bottts-neutral/svg?seed=" + encodeURI(name)
    }

    let showSelectedCommentFull = false
    function onClickShowSelectedCommentFull() {
        showSelectedCommentFull = !showSelectedCommentFull
    }

    function clear() {
        contents = ""
        autosize.update(document.querySelector('#text-bar')!!)
    }

    function send() {
        dispatcher('send', contents);
        clear();
    }
</script>
{#if selectedComment != null}
    <div class="{!showSelectedCommentFull ? 'max-h-32' : ''} relative overflow-hidden">
        <button on:click={() => { selectedComment = null }}
                class="absolute -top h-12 left-0 pr-4 w-full text-center text-red-500 mx-auto flex justify-end hover:cursor-pointer">
            <Icon src={Cross2} class="w-4 resize-none"></Icon>
        </button>
        <CommentElement data={selectedComment} showControls={false}/>
        {#if !showSelectedCommentFull}
            <button on:click={onClickShowSelectedCommentFull}
                    class="absolute -bottom-1 h-12 left-0 w-full text-center mx-auto from-transparent to-black bg-gradient-to-b flex justify-center hover:cursor-pointer">
                <Icon src={ChevronDown} class="w-4 resize-none"></Icon>
            </button>
        {:else}
            <button on:click={onClickShowSelectedCommentFull}
                    class="absolute -bottom-1 h-12 left-0 w-full text-center mx-auto flex justify-center hover:cursor-pointer">
                <Icon src={ChevronUp} class="w-4 resize-none"></Icon>
            </button>
        {/if}
    </div>
{/if}
{#if self != null}
    <div class="flex flex-col gap-4 md:gap-0 md:flex-row md:justify-between w-full pl-1 px-5 h-fit md:items-center flex-shrink-0 border p-1.5 rounded border-gray-700">
        <div class="ml-2 mt-2 lg:mt-0">
            <img src={avatar} class="rounded-full object-cover w-9" alt="Yourself"/>
        </div>
        <div class="ml-2 w-full max-h-96  rounded-full flex flex-row gap-4 justify-between py-1 pl-0 px-3">
                    <textarea id="text-bar"
                              rows="1"
                              placeholder={selectedComment ? "Reply to " + selectedComment.author.name : "Add a comment."}
                              class="ml-2 bg-transparent w-full max-h-96 outline-none resize-none"
                              bind:value={contents}
                    ></textarea>
            <button class="text-white hover:opacity-60 duration-300" on:click={send} data-tippy-content="Send Comment">
                <Icon src={PaperPlane} class="h-5 w-5 flex-shrink-0"></Icon>
            </button>
        </div>
    </div>
{:else}
    <a href="/login?callback={callback}" class="flex flex-col gap-4 md:gap-0 md:flex-row md:justify-between w-full pl-1 px-5 h-fit md:items-center flex-shrink-0 border p-1.5 rounded border-gray-700">
        <div class="ml-2 mt-2 lg:mt-0">
            <img src={avatar} class="rounded-full object-cover w-9" alt="Yourself"/>
        </div>
        <div class="ml-2 w-full max-h-96  rounded-full flex flex-row gap-4 justify-between py-1 pl-0 px-3">
                    <div id="text-bar" class="ml-2 bg-transparent w-full max-h-96 outline-none resize-none">Login to comment.</div>
        </div>
    </a>
{/if}
