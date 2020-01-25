<template>
    <div id="storage_tagged" class="small">
        <div class="title"><b>Нотатник: всі записи з міткою '{{ $route.params.tag }}'</b></div>

        <DataList v-bind:dataList="dataList"/>

        <!--<br>&nbsp;-->

        <!--<div v-if="dataItems">-->
            <!--<DataList v-bind:dataItems="dataItems"/>-->
        <!--</div><div v-else>-->
        <!--немає записів для перегляду...-->
    </div>

</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        name: 'StorageTagged',
        mounted() {
            this.getDataList(this.$route.params.tag);
        },

        beforeRouteUpdate (to, from, next) {
            this.getDataList(to.params.tag);
            next();         // it's necessary!!! else no hook is generated when we returns to the original route
        },

        data: () => {
            return {
                dataList: [],
            };
        },
        methods: {
            getDataList(tag) {
                // console.log(cfg.taggedEp + "?tag=" + encodeURIComponent(tag));

                fetch(cfg.taggedEp + "?tag=" + encodeURIComponent(tag), {
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
