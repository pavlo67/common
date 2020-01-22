<template>
    <div class="tags_index">
        <div v-for="tag in tagsIndex">
            <li>{{ tag.Label }} [<router-link :to="{ name: 'StorageTagged', params: {tag: tag.Label} }">{{ tag.Immediate }} / {{ tag.Full }}</router-link>]</li>
        </div>
    </div>
</template>


<script>
    export default {
        name: 'TagsIndex',
        props: ['tagsIndex'],
    }
</script>

<style lang="scss">
</style>
