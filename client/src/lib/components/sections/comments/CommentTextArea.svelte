<script lang="ts">
    import CommentElement from "$lib/components/CommentElement.svelte";
    import type {Comment} from "../../../types/comment";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {Avatar, ChevronDown, ChevronUp, Cross2, TriangleRight} from "@steeze-ui/radix-icons";
    import {colors} from "$lib/utils/colors.js";
    import autosize from 'autosize'
    import {createEventDispatcher, onMount} from "svelte";
    import type {CreateQueryResult} from "@tanstack/svelte-query";
    import type {User} from "../../../types/user";

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
            return "https://source.boringavatars.com/beam/120/primrose?colors=264653,2a9d8f,e9c46a,f4a261,e76f51"
        }

        let name = self.name
        return "https://source.boringavatars.com/beam/120/" + encodeURI(name) + "?colors=264653,2a9d8f,e9c46a,f4a261,e76f51"
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
    <div class="flex flex-row justify-between bg-neutral-900 bg-opacity-80 w-full p-1 py-2 pl-1 px-5 h-fit items-center flex-shrink-0">
        <div class="w-10 ml-2 rounded-full object-cover border border-orange-400">
            <img src={avatar} alt="Yourself"/>
        </div>
        <textarea id="text-bar"
                  rows="1"
                  placeholder={"Nothing is going to change my love for you."}
                  class="ml-2 playfair bg-transparent w-full max-h-96 p-1 outline-none resize-none"
                  bind:value={contents}
        ></textarea>
        <button class="p-2 bg-black text-white flex flex-row gap-1 items-center hover:opacity-60 duration-300" on:click={send} data-tippy-content="Send Comment">
            <Icon src={TriangleRight} class="h-6 w-6 flex-shrink-0"></Icon>
            <p>Send</p>
        </button>
    </div>
{:else}
    <div class="flex flex-row {colors.randomBackgroundColor()} overflow-hidden bg-opacity-30 w-full p-1 py-2 pl-1 px-5">
        <a href="/login?callback={callback}" class="ml-2 playfair flex flex-row items-center gap-2 bg-transparent w-full max-h-96 p-1 outline-none resize-none">
            <Icon src={Avatar} class="h-6 w-6 flex-shrink-0"></Icon>
            <p class="text-lg font-bold playfair">Login with {import.meta.env.VITE_APP_NAME}</p>
        </a>
    </div>
{/if}