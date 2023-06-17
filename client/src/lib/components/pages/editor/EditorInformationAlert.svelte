<script>
import {QuestionMarkCircled} from "@steeze-ui/radix-icons";
import Alert from "$lib/components/alerts/Alert.svelte";
</script>

<div class="py-4">
    <Alert icon={QuestionMarkCircled}
           title="The Primrose Editor"
           description="The Primrose Editor is the markdown editor for Primrose. To write posts, you have to use Markdown Text Language, if you are
                       unfamiliar with Markdown, you should learn about it first before using the editor. You can view changes to draft in semi-realtime after
                       saving the post since Primrose has constant refresh enabled for drafts. CTRL+S is supported for quick-saving."
    />
</div>