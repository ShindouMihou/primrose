import {NOT_FOUND, UNAUTHENTICATED} from "../errors/generic";
import {token} from "../../stash";

// @ts-ignore
export const  AUTHENTICATED_RETRY = (failureCount, error) => {
    if (error === UNAUTHENTICATED) {
        token.set('')
        setTimeout(() => window.location.replace('/login'), 250)
        return false
    }
    if (error === NOT_FOUND) {
        return false
    }
    return failureCount < 5;
}