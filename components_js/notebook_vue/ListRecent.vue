<template>
    <div id="storage_tagged" class="small">
        <div class="title"><b>Нещодавні записи</b></div>

        <DataList v-bind:dataList="dataList"/>
    </div>

</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        title: () => 'нещодавні записи',
        name: 'ListRecent',
        created () {
            this.getDataList();
        },
        data: () => {
            return {
                dataList: [],
            };
        },
        methods: {
            getDataList() {
                console.log(cfg.taggedEp + "?tag=" + encodeURIComponent(this.$route.params.tag));

                fetch(cfg.taggedEp + "?key=storage", {
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
                    this.dataList = data;
                    // console.log(this.tags);
                });
            }
        },
    }
</script>

<style lang="scss">
</style>
