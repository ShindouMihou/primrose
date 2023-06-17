<div class="w-screen py-8 border-t-2 border-t-[#797979] border-opacity-30">
    <div class="backdrop-blur flex flex-col p-4 px-12">
        <div class="flex flex-col text-center">
            <h1 class="text-6xl uppercase font-bold leading-none playfair">{import.meta.env.VITE_APP_NAME}</h1>
            <h3 class="text-sm font-light leading-none">{import.meta.env.VITE_SEO_DESCRIPTION}</h3>
            {#if (import.meta.env.VITE_SHOW_MADE_WITH_PRIMROSE ?? 'true') === 'true'}
                <a href="https://github.com/ShindouMihou/primrose" class="hover:scale-105 transition ease-in-out duration-700 text-xs font-light leading-none py-2">made with primrose</a>
            {/if}
        </div>
    </div>
</div>