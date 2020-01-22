<template>
    <div id="storage_item" class="small">
        <div v-if="dataItem">
            <b>Мій каталог: {{ dataItem.Title }}</b>
            <br>&nbsp;
            <DataItemView v-bind:dataItem="dataItem"/>
        </div>
        <div v-else>
            <b>Мій каталог:</b> відсутній запис для показу
        </div>


    </div>
</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        created () {
            this.getDataItem();
        },
        data: () => {
            return {
                dataItem: {},
            };
        },
        methods: {
            getDataItem() {
                fetch(cfg.readEp + "/" + encodeURIComponent(this.$route.params.id), {
                    method: 'GET', // *GET, POST, PUT, DELETE, etc.
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    mode: 'cors', // no-cors, cors, *same-origin

                    // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
                    // credentials: 'same-origin', // include, *same-origin, omit
                    // redirect: 'follow', // manual, *follow, error
                    // referrer: 'no-referrer', // no-referrer, *client
                }).then(response => {
                    return response.json();
                }).then(data => {
                    this.dataItem = data;
                    console.log(66666666, this.dataItem);
                });
            }
        },
    }
</script>

<style lang="scss">
</style>
