<script lang="ts">
    import {createInfiniteQuery, createQuery} from "@tanstack/svelte-query";
    import {deleteComment, editComment, fetchList, saveComment} from "$lib/requests/comments";
    import type {Comment} from "$lib/types/comment";
    import type {Arrayed} from "$lib/types/generic";
    import CommentElement from "$lib/components/CommentElement.svelte";
    import {token} from "$lib/stash.js";
    import SimpleErrorAlert from "$lib/components/alerts/SimpleErrorAlert.svelte";
    import {fetchSelf} from "$lib/requests/user";
    import {AUTHENTICATED_RETRY} from "$lib/requests/retries/retries";
    import {UNAUTHENTICATED} from "$lib/requests/errors/generic";
    import InfiniteLoading from 'svelte-infinite-loading';
    import Spinner from 'svelte-infinite-loading/src/Spinner.svelte'
    import CommentTextArea from "$lib/v2/components/comments/CommentTextArea.svelte";
    import DeleteModal from "$lib/components/sections/comments/DeleteModal.svelte";

    export let id: string
    export let slug: string

    const self = createQuery(['self'], () => fetchSelf($token), {
        retry: AUTHENTICATED_RETRY,
        refetchInterval: false,
        refetchOnWindowFocus: false,
        refetchOnMount: false,
        refetchOnReconnect: false
    })

    const query = createInfiniteQuery({
        queryKey: ['comments', id],
        queryFn: async ({ pageParam = null }) => await fetchList(id, pageParam),
        getNextPageParam: (lastPage: Arrayed<Comment>) => {
            if (lastPage.data.length > 0) {
                return lastPage.data[lastPage.data.length - 1].id
            }
            return undefined
        }
    })

    let selectedComment: Comment | null = null

    let lock = false
    let error: Error | null = null

    let callback = encodeURIComponent('/read/' + slug + '#comments')

    function send(event: CustomEvent<string>) {
        if (lock) {
            return
        }
        if (event.detail.length === 0) {
            error = new Error("You cannot post an empty comment.")
            return;
        }
        lock = true
        saveComment($token, id, {content: event.detail}, selectedComment?.id)
            .then(() => $query.refetch())
            .then(() => { lock = false; selectedComment = null; })
            .catch((err) => {
                if (error === UNAUTHENTICATED) {
                    token.set('')
                    window.location.replace('/login?callback=' + callback)
                    return
                }
                error = err
                lock = false
            })
    }

    function select(event: CustomEvent<Comment>) {
        if ($self.data == null) {
            return
        }
        selectedComment = event.detail
        window.lenis.scrollTo(document.querySelector('#comments'))
    }

    let deleting: Comment | null = null
    $: showDeletingModal = deleting != null

    function del(event: CustomEvent<Comment>) {
        if ($self.data == null) {
            return
        }
        if (lock) {
            return;
        }
        deleting = event.detail
    }

    function actuallyDelete() {
        if (deleting == null) {
            return
        }
        lock = true
        deleteComment($token, deleting.id)
            .then(() => $query.refetch())
            .then(() => { lock = false; deleting = null; })
            .catch((err) => {
                if (error === UNAUTHENTICATED) {
                    token.set('')
                    window.location.replace('/login?callback=' + callback)
                    return
                }
                error = err
                lock = false
            })
    }

    function edit(event: CustomEvent<{content: string, comment: Comment}>) {
        if (lock) {
            return
        }
        lock = true
        editComment($token, event.detail.comment.id, event.detail.content)
            .then(() => $query.refetch())
            .then(() => { lock = false; })
            .catch((err) => {
                if (error === UNAUTHENTICATED) {
                    token.set('')
                    window.location.replace('/login?callback=' + callback)
                    return
                }
                error = err
                lock = false
            })
    }
</script>

<div class="pb-8 flex flex-col gap-4 border-t border-t-gray-500 pt-2" id="comments">
    <div class="rounded">
        <div class="py-2">
            <h2 class="text-lg font-medium">Comments</h2>
        </div>
        {#if error != null}
            <SimpleErrorAlert error={error}></SimpleErrorAlert>
        {/if}
        {#if $query.isLoading || $self.isLoading}
            <Spinner/>
        {/if}


        {#if $query.isSuccess && $self.isSuccess}
            <div class="p-2">
                <CommentTextArea bind:selectedComment bind:callback self={$self.data} on:send={send}/>
            </div>
            <div id="comments__section">
                {#each $query.data.pages as { data }}
                    {#each data as comment}
                        <div class="my-0.5 p-2">
                            <CommentElement data={comment}
                                            on:reply={select}
                                            on:edit={edit}
                                            on:delete={del}
                                            authenticated={$self.data != null}
                                            showAuthorControls={$self.data != null && comment.author.id === $self.data.id}/>
                        </div>
                    {/each}
                {/each}
                <InfiniteLoading on:infinite={({ detail: { loaded, complete }}) => {
                    $query.fetchNextPage().then((result) => {
                        if (result.data.pages[result.data.pages.length - 1].data.length === 0) {
                            complete()
                        } else {
                            setTimeout(loaded, 500)
                        }
                    })
                }}>
                    <div slot="noResults"></div>
                    <div slot="noMore"></div>
                    <div slot="error">
                        <SimpleErrorAlert error={new Error("Something went wrong.")}></SimpleErrorAlert>
                    </div>
                </InfiniteLoading>
            </div>
            <DeleteModal bind:show={showDeletingModal} on:delete={actuallyDelete} on:cancel={() => deleting = null}/>
        {/if}
    </div>
</div>
