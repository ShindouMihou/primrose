import {BACKEND_ADDRESS} from "../constants";
import {CONTENT_NOT_ACCESSIBLE, NOT_FOUND, UNAUTHENTICATED} from "./errors/generic";
import type {Post} from "../types/post";
import type {Arrayed} from "../types/generic";

export const fetchList = async (): Promise<Arrayed<Post>> => {
    const response = await fetch(BACKEND_ADDRESS + "/posts/list?limit=100")
    if (!response.ok) {
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const fetchListWithUnpublished = async (token: string): Promise<Arrayed<Post>> => {
    const response = await fetch(BACKEND_ADDRESS + "/posts/list?published=false&limit=100", {
        headers: {
            Authorization: "Bearer " + token
        }
    })
    if (!response.ok) {
        if (response.status === 401) {
            throw UNAUTHENTICATED
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const fetchPost = async (key: string, token: string | null = null): Promise<Post> => {
    const response = await fetch(BACKEND_ADDRESS  + "/posts/view?key=" + encodeURI(key) + "&isSlug=true&findClosest=true", {
        headers: token == null ? undefined : {
            Authorization: "Bearer " + token
        }
    })
    if (response.status === 404) {
        throw NOT_FOUND
    }
    if (!response.ok) {
        if (response.status === 401) {
            throw UNAUTHENTICATED
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const fetchPostWithAuthentication = async (key: string, token: string | null): Promise<Post> => {
    const response = await fetch(BACKEND_ADDRESS  + "/posts/view?key=" + encodeURI(key) + "&isSlug=false", {
        headers: {
            Authorization: "Bearer " + token
        }
    })
    if (response.status === 404) {
        throw NOT_FOUND
    }
    if (!response.ok) {
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const fetchPostForEditor = async (key: string | null, token: string): Promise<Post> =>  {
    if (key == "" || key == null) {
        return Promise.resolve({id:"",title: "", content: "", image: null, slug: "", created_at: new Date(), updated_at: new Date(), published: false})
    }
    return fetchPostWithAuthentication(key, token)
}

export const savePost = async (token: string, post: Post): Promise<Post> => {
    const response = await fetch(BACKEND_ADDRESS  + "/posts/save", {
        headers: {
            Authorization: "Bearer " + token
        },
        method: 'PUT',
        body: JSON.stringify(post)
    })
    if (!response.ok) {
        if (response.status === 401) {
            throw UNAUTHENTICATED
        }
        if (response.status === 400) {
            const error = (await response.json()).error
            if (error === "Invalid payload.") {
                throw Error("Cannot save the post with the given settings. It is likely the image, or slug is set to a bad value.")
            }
            throw Error(error)
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const deletePost = async (token: string, key: string): Promise<any> => {
    const response = await fetch(BACKEND_ADDRESS  + "/posts/delete?key="+encodeURI(key), {
        headers: {
            Authorization: "Bearer " + token
        },
        method: 'DELETE'
    })
    if (!response.ok) {
        if (response.status === 401) {
            throw UNAUTHENTICATED
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}