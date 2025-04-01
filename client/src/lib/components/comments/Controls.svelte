<script lang="ts">
    export let show: boolean
</script>

{#if show}
    <div class="p-2 items-center shadow rounded flex flex-row justify-between md:justify-normal gap-4 hidden group-hover:flex">
        <slot/>
    </div>
{/if}
