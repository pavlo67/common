<template>
    <div id="item_import">
        <b>Імпорт запису в приватний каталог</b>

        <br>&nbsp;

        <div v-if="itemToImport">
            <DataItemView v-bind:itemToImport="itemToImport"/>
            <button v-on:click="saveToStorage()">Save to the storage</button>
        </div>
        <div v-else>
            немає запису для імпорту...
        </div>

    </div>
</template>


<script>

    import { cfg } from './init';

    export default {
        name: 'DataItemImport',

        beforeMount() {
            this.itemToImport = this.$route.params.dataItem;
        },

        methods: {
            saveToStorage() {
                delete this.itemToImport.ID;
                delete this.itemToImport.Status;

                // TODO!!! backup it
                delete this.itemToImport.ExportID;

                fetch(cfg.saveEp, {
                    method: 'POST',
                    headers: {
                        // 'Content-Type': 'application/json;charset=utf-8'
                        'Content-Type': 'application/json',
                    },
                    mode: 'cors', // no-cors, cors, *same-origin
                    body: JSON.stringify(this.itemToImport),
                }).then(response => {
                    console.log(99999999, response);
                });
            }
        },

    };
</script>

<style lang="scss">
    #item_import {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
