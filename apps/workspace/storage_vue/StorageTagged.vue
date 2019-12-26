<template>
    <div id="storage_tagged" class="small">
        <b>Мій каталог: всі записи з міткою '{{ $route.params.tag }}'</b>

        <br>&nbsp;

        <DataList v-bind:dataList="dataList"/>

        <!--<br>&nbsp;-->

        <!--<div v-if="dataItems">-->
            <!--<DataList v-bind:dataItems="dataItems"/>-->
        <!--</div><div v-else>-->
        <!--немає записів для перегляду...-->
    </div>

</template>


<script>
    import b       from '../../../components.js/basis';
    import { cfg } from './init';

    export default {
        name: 'StorageTagged',
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

                fetch(cfg.taggedEp + "?key=storage&tag=" + encodeURIComponent(this.$route.params.tag), {
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
