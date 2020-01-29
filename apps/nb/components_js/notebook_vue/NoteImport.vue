<template>
    <div id="storage_import">
        <b>Імпорт запису в приватний каталог</b>

        <br>&nbsp;

        <div v-if="itemToImport">
            <DataItemView v-bind:dataItem="itemToImport" v-bind:showTitle="true" />
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
                        'content-type': 'application/json;charset=utf-8',
                        'authorization': cfg.common.user && cfg.common.user.Creds && cfg.common.user.Creds.jwt,
                    },
                    mode: 'cors', // no-cors, cors, *same-origin
                    body: JSON.stringify(this.itemToImport),

                }).then(response => {
                    return response.json();

                }).then(data => {
                    if (data.ID > 0) {
                        cfg.eventBus.$emit('message', "запис з id = " + data.ID + " збережено");
                        this.$router.push({ name: 'StorageItem',  params: { id: data.ID } })

                    } else {
                        console.log(data);
                        cfg.eventBus.$emit('message', "не вдалось зберегти запис: " + data.Error);
                    }


                });
            }
        },

    };
</script>

<style lang="scss">
</style>
