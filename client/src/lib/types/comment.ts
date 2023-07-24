import type {User} from "./user";

export type Comment = {
    id: string,
    content: string,
    post: string,
    author: string | User,
    reply_to: string | null,
    created_at: Date,
    updated_at: Date,
    resolved_parent: Comment | null,
}