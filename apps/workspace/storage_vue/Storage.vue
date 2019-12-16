<template>
    <div id="data">
        <b>Мій каталог</b>

        <br>&nbsp;

        <div v-if="dataItems">
            <DataList v-bind:dataItems="dataItems"/>
        </div>
        <div v-else>
            немає записів для перегляду...
        </div>

    </div>
</template>


<script>
    import b       from '../../components.js/basis';
    import { cfg } from './init';

    export default {
        title: 'мій каталог',
        created () {
            this.getDataItems();
        },
        data: () => {
             return {
                 dataItems: [],
             };
        },
        methods: {
            getDataItems() {
                fetch(cfg.listEp, {
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
                    this.dataItems = data;
                    console.log(this.dataItems);
                });
            }
        },
    }
</script>

<style lang="scss">
    #data {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
