<template>
    <div id="storage_item" class="small">
        <div class="title"><b>Новий запис</b></div>
        <DataItemView v-bind:dataItem="dataItem"/>
    </div>
</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        title: () => 'новий запис',
        data: () => {
            return {
                dataItem: {},
            };
        },
        methods: {
        },
    }
</script>

<style lang="scss">
</style>
