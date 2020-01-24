<template>
    <div class="list_tags">
        <div v-for="tag in listTags">
            <li>{{ tag.Label }} [<router-link :to="{ name: 'StorageTagged', params: {tag: tag.Label} }">{{ tag.Immediate }} / {{ tag.Full }}</router-link>]</li>
        </div>
    </div>
</template>


<script>
    export default {
        name: 'ListTags',
        props: ['listTags'],
    }
</script>

<style lang="scss">
</style>
