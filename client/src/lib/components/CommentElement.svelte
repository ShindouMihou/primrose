<script lang="ts">
    import type {Comment} from "../types/comment";
    import {toHTML} from "../renderer/markdown";
    import {createEventDispatcher} from "svelte";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {
        ChatBubble, Check, Cross2,
        EyeClosed,
        EyeOpen,
        Pencil1, Share1,
        Trash,
    } from "@steeze-ui/radix-icons";
    import autosize from 'autosize'
    import ControlButton from "$lib/components/comments/ControlButton.svelte";
    import Controls from "$lib/components/comments/Controls.svelte";
    import Author from "$lib/components/comments/Author.svelte";

    export let data: Comment
    export let authenticated: boolean = false
    export let showControls: boolean = true
    export let showAuthorControls: boolean = false

    const dispatcher = createEventDispatcher()

    function html(content: string) {
        content = content.replaceAll('#', '\\#')
        const translated = toHTML(content);

        if (translated.error && !translated.content) {
            return translated.error;
        }

        return translated.content!;
    }

    let showParentComment = false
    function onClickShowParentComment() {
        showParentComment = !showParentComment
    }

    let editing = false
    let edit = data.content

    function onClickEdit() {
        if (!editing) {
            edit = data.content
        }
        editing = !editing
        if (editing) {
            setTimeout(() => autosize(document.querySelector('#edit-bar')!!), 500)
        }
    }

    function shareUrl(): string {
        window.location.hash = "#comment-" + data.id
        return window.location.toString()
    }
</script>

<div class="w-full flex flex-col gap-2 p-4 border rounded border-gray-700 group" id="comment-{data.id}">
    <div class="pb-2 flex flex-col gap-4 md:gap-0 md:flex-row md:justify-between md:items-center">
        <Author bind:comment={data}/>
        <Controls bind:show={showControls}>
            {#if !editing}
                {#if navigator.canShare != null && navigator.canShare({ url: shareUrl() })}
                    <ControlButton icon={Share1} on:click={() => navigator.share({ url: shareUrl() })}/>
                {/if}
                {#if authenticated}
                    <ControlButton icon={ChatBubble} on:click={() => dispatcher('reply', data)}/>
                {/if}
                {#if showAuthorControls}
                    <ControlButton icon={Trash} on:click={() => dispatcher('delete', data)}/>
                    <ControlButton icon={Pencil1} on:click={onClickEdit}/>
                {/if}
            {:else}
                <ControlButton icon={Cross2} on:click={onClickEdit}/>
                <ControlButton icon={Check} on:click={() => {
                        onClickEdit();
                        dispatcher('edit', { content: edit, comment: data })
                    }}/>
            {/if}
        </Controls>
    </div>
    {#if data.resolved_parent != null}
        <div on:click={onClickShowParentComment}
              class="text-gray-300 font-light text-sm hover:opacity-60 duration-300 transition hover:cursor-pointer py-2 flex flex-row gap-2 items-center">
            <Icon src={showParentComment ? EyeClosed : EyeOpen} class="w-4 resize-none"></Icon>
            <p>Quotes comment from {data.resolved_parent.author.name}</p>
        </div>
        <div class="relative overflow-hidden {showParentComment ? '' : 'hidden'}">
            <svelte:self data={data.resolved_parent} showControls={false}/>
        </div>
    {/if}
    <div class="mkdown">
        {#if !editing}
            {@html html(data.content)}
        {:else}
              <textarea id="text-bar"
                        rows="1"
                        placeholder={"Edit your comment."}
                        class="bg-transparent w-full outline-none resize-none"
                        bind:value={edit}
              ></textarea>
        {/if}
    </div>
</div>
