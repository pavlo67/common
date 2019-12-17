<template>
    <div class="tags_index">
        <div v-for="tag in tags">
            <li>{{ tag.Label }} [<a href="">{{ tag.Immediate }} / {{ tag.Full }}</a>]</li>
        </div>
    </div>
</template>


<script>
    // import b from '../../components.js/basis';
    // import e from '../../components.js/elements';

    export default {
        name: 'TagsIndex',
        props: ['tags'],
    }
</script>

<style lang="scss">
    .tags_index {
        font-size: small;
    }
</style>
