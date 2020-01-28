<template>
    <div>

        <table align="right" class="table_right">
            <tr><td>
                <span class="control">
                    [<a v-bind:href="exportData()" download="data.json">експортувати</a>]
                    <!--                    [<span v-on:click="exportData">експортувати</span>]-->
                </span>
            </td></tr>
        </table>

        <div v-for="dataItem in dataList" class="small">
            <!-- TODO: customize the router target -->
            <span class="data_announce">
                <span class="time">[{{ createdAt(dataItem) }}]</span>
                {{ dataItem.Title }}
                <span class="control smallest">
              &nbsp;     [<span v-on:click="$router.push({ name: 'NoteView',  params: { id: dataItem.ID } })">докладно</span>]
                    &nbsp;
                    [<span v-on:click="$router.push({ name: 'NoteEdit',  params: { id: dataItem.ID } })">ред.</span>]
                </span>
                <span v-if="dataItem.Summary">
                    <br><span v-html="dataItem.Summary" class="smaller"></span>
                </span>
                <span v-if="dataItem.Tags">
                    <br>
                    <span v-for="tag in dataItem.Tags" class="tag control smaller" >
                        [<span v-on:click="$router.push({ name: 'ListTagged',  params: { tag: tag.Label } })">{{ tag.Label }}</span>]
                        &nbsp;
                    </span>
                </span>
            </span>
            <br>&nbsp;
        </div>
    </div>
</template>


<script>
    import b from '../basis';
    import e from '../elements';
    import {createdAt} from './data';
    import DataEdit from "./DataEdit";

    let cfg = {};

    export default {
        name: 'DataList',
        props: ['dataList'],
        methods: {
            createdAt: createdAt,
            object:  b.object,
            href:    e.href,

            prepare(dataItems, cfgCommon) {
                cfg = cfgCommon;
                return dataItems;
            },

            exportData() {
                return "data:text/html,aaa";
            },

            exportDataAsync() {
                console.log(555555555555, cfg.exportEp, cfg.common.user && cfg.common.user.Creds && cfg.common.user.Creds.jwt);

                fetch(cfg.exportEp, {
                    method: 'GET',
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
                    // this.dataItem = DataEdit.methods.prepare(data, cfg);
                    console.log("TO EXPORT:", data);
                });

            }
        },
    }
</script>

<style lang="scss">
</style>
