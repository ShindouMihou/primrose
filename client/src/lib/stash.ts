import {type Writable, writable} from "svelte/store";
import {browser} from "$app/environment";
import type {User} from "./types/user";

export const token: Writable<string> = writable(stored('primrose.acc', ''))

if (browser) {
    token.subscribe((value)  => localStorage.setItem('primrose.acc', value))
}

function stored(key: string, def: string) {
    if (browser) {
        return (localStorage.getItem(key) ?? def)
    }

    return def;
}