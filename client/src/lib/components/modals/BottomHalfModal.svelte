<script lang="ts">
    export let title: string
    export let show: boolean
</script>

{#if show}
    <div class="fixed z-30 w-screen min-h-screen md:backdrop-blur animate-in fade-in duration-700 animated bottom-0 left-0 select-none">
        <div class="min-h-screen md:min-h-fit bg-black bg-opacity-100 md:backdrop-blur border-t border-[#797979] border-opacity-30 p-12 pt-18 md:pt-12 md:rounded-t-xl fixed w-screen bottom-0 left-0 animate-in slide-in-from-bottom animated">
            <div class="flex flex-col" id="modal-content">
                <h1 class="font-bold text-[#797979">{title}</h1>
                <slot/>
            </div>
        </div>
    </div>
{/if}