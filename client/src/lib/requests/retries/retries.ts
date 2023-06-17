import {UNAUTHENTICATED} from "../errors/generic";
import {token} from "../../stash";

export const  AUTHENTICATED_RETRY = (failureCount, error) => {
    if (error === UNAUTHENTICATED) {
        token.set('')
        setTimeout(() => window.location.replace('/login'), 250)
        return false
    }
    return failureCount < 5;
}