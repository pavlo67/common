<template>
    <div id="item_import">
        <b>Імпорт запису в приватний каталог</b>

        <p><div>{{ dataItem.Title }}</div>
        <br><div>{{ dataItem.Summary }}</div>
        <br><div>{{ dataItem }}</div>
    </div>
</template>


<script>
    import b       from '../../components.js/basis';
    import { cfg } from './init';

    export default {
        name: 'data item import',

        beforeMount() {
            console.log(6666, this.$route.params.dataItem);
            this.dataItem = this.$route.params.dataItem;
        },

        // data: () => {
        //     return {
        //         sourcePacks: [],
        //     };
        // },

        methods: {
            dateStr: b.dateStr,

            inputData() {
                console.log(7777, this.$route.params.dataItem);
            },

        },
    };
</script>

<style lang="scss">
    #item_import {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
