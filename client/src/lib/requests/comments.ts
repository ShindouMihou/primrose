import {BACKEND_ADDRESS} from "../constants";
import {CONTENT_NOT_ACCESSIBLE, NOT_FOUND, UNAUTHENTICATED} from "./errors/generic";
import type {Post} from "../types/post";
import type {Arrayed} from "../types/generic";
import type {Comment} from "../types/comment";

export const fetchList = async (post: string, after: string | null = null): Promise<Arrayed<Comment>> => {
    let query = "?limit=10"
    if (after != null) {
        query += "&after=" + encodeURI(after)
    }
    const response = await fetch(BACKEND_ADDRESS + "/comments/" + encodeURI(post) + query)
    if (!response.ok) {
        throw CONTENT_NOT_ACCESSIBLE
    }
    const comments: Arrayed<Comment> = await response.json()

    let pendingResolves: string[] = []
    for(let comment of comments.data) {
        if (comment.reply_to == null) {
            continue
        }
        let resolved: Comment | null = comments.data.find((cmt) => cmt.id == comment.reply_to) ?? null
        if (resolved != null) {
            comment.resolved_parent = resolved
            continue
        }
        pendingResolves.push(comment.reply_to)
    }
    const resolveResponse = await fetch(BACKEND_ADDRESS + "/comments/many", {
        method: "POST",
        body: JSON.stringify({
            ids: pendingResolves
        })
    })
    if (!resolveResponse.ok) {
        throw CONTENT_NOT_ACCESSIBLE
    }
    const resolves: Arrayed<Comment> = await resolveResponse.json()
    for(let comment of comments.data) {
        if (comment.reply_to == null) {
            continue
        }
        let resolved: Comment | null = resolves.data.find((cmt) => cmt.id == comment.reply_to) ?? null
        if (resolved != null) {
            comment.resolved_parent = resolved
            continue
        }
        console.error("failed to find resolve parent for comment id ", comment.id, " with resolve id ", comment.reply_to)
    }
    return comments
}

export const fetchComment = async (key: string): Promise<Comment> => {
    const response = await fetch(BACKEND_ADDRESS  + "/comments/view/" + encodeURI(key))
    if (response.status === 404) {
        throw NOT_FOUND
    }
    if (!response.ok) {
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const saveComment = async (token: string, post: string, comment: { content: string }, replyTo: string | undefined): Promise<Comment> => {
    post = encodeURI(post)
    const query = replyTo == null ? "" : "?reply_to=" + encodeURI(replyTo)
    const response = await fetch(BACKEND_ADDRESS  + "/comments/"+post+query, {
        headers: {
            Authorization: "Bearer " + token
        },
        method: 'PUT',
        body: JSON.stringify(comment)
    })
    if (!response.ok) {
        if (response.status === 401) {
            throw UNAUTHENTICATED
        }
        if (response.status === 400) {
            const error = (await response.json()).error
            if (error === "Invalid payload.") {
                throw Error("Cannot post the comment. It is likely that one of the parameters is wrong.")
            }
            throw Error(error)
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const editComment = async (token: string, comment: string, content: string): Promise<Comment> => {
    if (content.length === 0) {
        throw Error("Content of the comment cannot be empty.")
    }
    const response = await fetch(BACKEND_ADDRESS  + "/comments/edit/"+encodeURI(comment), {
        headers: {
            Authorization: "Bearer " + token
        },
        method: 'PATCH',
        body: JSON.stringify({
            content: content
        })
    })
    if (!response.ok) {
        if (response.status === 401) {
            throw UNAUTHENTICATED
        }
        if (response.status === 400) {
            const error = (await response.json()).error
            if (error === "Invalid payload.") {
                throw Error("Cannot post the comment. It is likely that one of the parameters is wrong.")
            }
            throw Error(error)
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const deleteComment = async (token: string, comment: string): Promise<any> => {
    comment = encodeURI(comment)
    const response = await fetch(BACKEND_ADDRESS  + "/comments/delete/"+comment, {
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