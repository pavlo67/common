<template>
    <div class="tags_index">
        <div v-for="tag in tagsIndex">
            <li>{{ tag.Label }} [<router-link :to="{ name: 'StorageTagged', params: {tag: tag.Label} }">{{ tag.Immediate }} / {{ tag.Full }}</router-link>]</li>
        </div>
    </div>
</template>


<script>
    // import b from '../../components.js/basis';
    // import e from '../../components.js/elements';

    export default {
        name: 'TagsIndex',
        props: ['tagsIndex'],
    }
</script>

<style lang="scss">
    .tags_index {
        font-size: small;
    }
</style>
