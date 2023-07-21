import { fetchPost } from '$lib/requests/posts';
import type { PageServerLoad } from './$types';
import isbot from 'isbot'

export const load = (async ({ params, request }) => {
    const useragent = request.headers.get('User-Agent')
    if (isbot(useragent))  {
        const post = await fetchPost(params.key, null)
        return {
            preload: { post }
        }
    }
    return { post: null };
}) satisfies PageServerLoad;