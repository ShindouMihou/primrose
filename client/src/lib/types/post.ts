export type Post = {
    id: string,
    title: string,
    content: string,
    published: boolean,
    slug: string,
    image: string | null,
    created_at: Date,
    updated_at: Date
}