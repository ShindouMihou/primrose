<script lang="ts">
    import BottomHalfModal from "$lib/components/modals/BottomHalfModal.svelte";
    import {createEventDispatcher} from "svelte";

    const dispatcher = createEventDispatcher()
    export let show: boolean

</script>

<BottomHalfModal title={"Are you sure you want to delete this comment?"} bind:show>
    <p class="text-sm font-medium text-[#7a7a7a]">
        Deleting a comment is an irreversible action, this will also delete any replies that are linked
        to the comment. Do you wish to continue?
    </p>
    <div class={"flex flex-row gap-2 pt-6"}>
        <button on:click={() => dispatcher('cancel')} class={"mr-12 font-medium uppercase hover:scale-105 duration-700 animated hover:cursor-pointer"}>
            Cancel
        </button>
        <button on:click={() => dispatcher('delete')} class={"mr-12 font-medium text-red-500 uppercase hover:scale-105 duration-700 animated hover:cursor-pointer"}>
            Delete
        </button>
    </div>
</BottomHalfModal>