import {BACKEND_ADDRESS} from "../constants";
import {CONTENT_NOT_ACCESSIBLE, UNAUTHENTICATED} from "./errors/generic";
import type {User} from "../types/user";

export const fetchSelf = async (token: string): Promise<User | null> => {
    if (token === '') {
        return Promise.resolve(null)
    }
    const response = await fetch(BACKEND_ADDRESS + "/users/@me", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    if (response.status === 403) {
        throw UNAUTHENTICATED
    }
    if (!response.ok) {
        throw CONTENT_NOT_ACCESSIBLE
    }
    return response.json()
}

export const signIn = async (email: string, password: string): Promise<string> => {
    if (email === '' || password === '') {
        throw Error('Email or password cannot be empty.')
    }
    const response = await fetch(BACKEND_ADDRESS + "/users/", {
        method: 'POST',
        body: JSON.stringify({
            email,
            password
        })
    })
    if (!response.ok) {
        if (response.status >= 400 &&  response.status < 500) {
            const body = await response.json()
            if (body.error === 'Invalid payload.') {
                throw Error('Invalid email or password.')
            }
            throw Error(body.error)
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
    const body = await response.json()
    return body.token
}

export const signUp = async (username: string, email: string, password: string) => {
    if (username === '' || email === '' || password === '') {
        throw Error('Username, Email or password cannot be empty.')
    }
    const response = await fetch(BACKEND_ADDRESS + "/users/", {
        method: 'PUT',
        body: JSON.stringify({
            username,
            email,
            password
        })
    })
    if (!response.ok) {
        if (response.status >= 400 &&  response.status < 500) {
            const body = await response.json()
            if (body.error === 'Invalid payload.') {
                throw Error('Invalid email or password.')
            }
            throw Error(body.error)
        }
        throw CONTENT_NOT_ACCESSIBLE
    }
}