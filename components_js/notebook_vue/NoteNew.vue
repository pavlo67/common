<template>
    <div id="storage_item" class="small">
        <div class="title"><b>Новий запис</b></div>
        <DataEdit v-bind:dataItem="dataItem"/>
    </div>
</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        title: () => 'новий запис',
        data: () => {
            return {
                dataItem: {
                    // Key:      "",
                    // Title:    "type title here",
                    // URL:      "type URL here",
                    // Summary:  "type summary here",
                    // _tags:    "type tags here",
                    // _content: "type content here",
                },
            };
        },
        methods: {
        },
    }
</script>

<style lang="scss">
    .edit_button {
        height: 30px;
    }
</style>
