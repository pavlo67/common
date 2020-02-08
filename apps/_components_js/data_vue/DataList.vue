<template>
    <div>

        <table align="right" class="table_right">
            <tr><td>
                <ActionSteps
                        v-bind:activationText="'експортувати'"
                        v-bind:activationAction="exportData"
                        v-bind:preparationText="'дані готуються...'"
                        v-bind:actionButtonText="'зберегти собі JSON файл'"
                />
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
    import {createdAt} from '../date';
    import DataEdit from "./DataEdit";

    let cfg = {};

    export default {
        name: 'DataList',
        props: ['dataList'],
        methods: {
            createdAt: createdAt,
            object: b.object,
            href: e.href,

            prepare(dataItems, cfgCommon) {
                cfg = cfgCommon;
                return dataItems;
            },

            exportData(cb) {
                if (typeof cb === "function") {
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
                        console.log('EXPORTED DATA ARE READY TO BE SAVED: ', data);
                        cb(JSON.stringify(data));

                    });
                } else {
                    console.error('THERE IS NO CALLBACK TO RETURN EXPORTED DATA!');

                }
            },
        }
    }
</script>

<style lang="scss">
</style>
