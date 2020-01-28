<template>
    <div id="storage_item" class="small">
        <div v-if="dataItem">
            <div class="title"><b>Нотатник: {{ dataItem.Title }}</b></div>

            <DataItemView v-bind:dataItem="dataItem"/>
        </div>
        <div v-else>
            <div class="title"><b>Нотатник</b></div>

            Відсутній запис для показу...
        </div>


    </div>
</template>


<script>
    // import b     from '../basis';
    import {cfg}    from './init';
    import DataView from "../data_vue/DataView";

    export default {
        mounted() {
            this.getDataItem(this.$route.params.id);
        },
        beforeRouteUpdate (to, from, next) {
            this.getDataItem(to.params.id);
            next();         // it's necessary!!! else no hook is generated when we returns to the original route
        },
        data: () => {
            return {
                dataItem: {},
            };
        },
        methods: {
            getDataItem(id) {
                fetch(cfg.readEp + "/" + encodeURIComponent(id), {
                    method: 'GET', // *GET, POST, PUT, DELETE, etc.
                    headers: {
                        'content-type': 'application/json;charset=utf-8',
                        'authorization': cfg.common.user && cfg.common.user.Creds && cfg.common.user.Creds.jwt,
                    },
                    mode: 'cors', // no-cors, cors, *same-origin

                    // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
                    // credentials: 'same-origin', // include, *same-origin, omit
                    // redirect: 'follow', // manual, *follow, error
                    // referrer: 'no-referrer', // no-referrer, *client
                }).then(response => {
                    return response.json();
                }).then(data => {
                    this.dataItem = DataView.methods.prepare(data, cfg);
                    console.log("TO SHOW:", this.dataItem);
                });
            }
        },
    }
</script>

<style lang="scss">
</style>
