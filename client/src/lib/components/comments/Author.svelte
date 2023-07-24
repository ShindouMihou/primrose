<script lang="ts">
    import type {Comment} from "../../types/comment";
    import TimeAgo from 'javascript-time-ago'
    import en from 'javascript-time-ago/locale/en'

    TimeAgo.setDefaultLocale(en)
    TimeAgo.addLocale(en)

    export let comment: Comment
    const timeago = new TimeAgo(['en-US'])

    const avatar = createAvatar()

    function createAvatar() {
        let name = comment.author instanceof String ? comment.author : comment.author.name
        return "https://source.boringavatars.com/beam/120/" + encodeURI(name) + "?colors=264653,2a9d8f,e9c46a,f4a261,e76f51"
    }
</script>
<div class="flex flex-row gap-2 items-center">
    <div class="w-10 mr-2 rounded-full">
        <img src={avatar} alt="{comment.author.name}'s avatar"/>
    </div>
    <div class="flex flex-col">
        <h2 class="text-lg font-bold">{comment.author.name}</h2>
        <p class="text-sm text-gray-200 font-light">{timeago.format(new Date(comment.created_at))}</p>
    </div>
</div>