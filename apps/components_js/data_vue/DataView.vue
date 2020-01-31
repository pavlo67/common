<template>

    <div id="data_item_view" class="smaller">

        <div v-if="showTitle">
            {{ dataItem.Title }}
            <br>&nbsp;
        </div>

        <table align="right" class="table_right">
            <tr><td>
                Створено: <span class="time">{{ createdAt(dataItem) }}</span><br>
                Ключ запису:&nbsp; {{ dataItem.Key }}

                <p/>
                <span class="control">
                    [<span v-on:click="$router.push({ name: 'NoteEdit',  params: { id: dataItem.ID } })">редаґувати</span>] &nbsp;
                    [<span v-on:click="remove">вилучити</span>]
                </span>

                <p/>Теґи:
                <span v-if="dataItem.Tags instanceof Array">
                    <span v-for="tag in dataItem.Tags" class="tag control" >
                    &nbsp; [<span v-on:click="$router.push({ name: 'ListTagged',  params: { tag: tag.Label } })">{{ tag.Label }}</span>]
                    </span><br>
                </span>

            </td></tr>
        </table>
        <span v-if="dataItem.Summary">
            <b>Короткий зміст, анонс:</b>&nbsp; <span v-html="dataItem.Summary"></span><br>
        </span>
        <span v-if="dataItem.URL">
            <b>URL:</b>&nbsp; <span v-html="href(dataItem.URL)" class="href"></span><br>
        </span>

        <span v-if="dataItem.Data instanceof Object">
            <p v-html="dataItem.Data.Content"></p>

        </span>

    </div>
</template>


<script>
    import e  from '../elements';
    import {createdAt} from './data';

    let cfg = {};

    export default {
        name: 'DataItemView',
        props: ["dataItem", "showTitle"],
        methods: {
            href: e.href,
            createdAt: createdAt,

            prepare(dataItem, cfgCommon) {
                cfg = cfgCommon || {common: {}};

                return dataItem;
            },

            remove() {
                fetch(cfg.removeEp + "/" + this.dataItem.ID, {
                    method: 'DELETE',
                    headers: {
                        'content-type' : 'application/json;charset=utf-8',
                        'authorization': cfg.common.user && cfg.common.user.Creds && cfg.common.user.Creds.jwt,
                    },
                    mode: 'cors', // no-cors, cors, *same-origin

                }).then(response => {
                    return response.json();

                }).then(data => {
                    if (data.id) {
                        cfg.eventBus.$emit('message', "запис з id = " + data.id + " вилучено");
                        this.$router.push({ name: 'ListRecent'})

                    } else {
                        console.log(data);
                        cfg.eventBus.$emit('message', "не вдалось вилучити запис: " + data.Error);
                    }


                });

            },

            // itemId(item, postfixed) {
            //     return "item_to_import_" + item.ID + (postfixed ? itemPostfix : "");
            // },
            //
            // showDetails(ev) {
            //     showHide.showContent(ev);
            // },
            //
            // hideDetails(ev) {
            //     showHide.hideContent(ev);
            // },
            //
            // details(item) {
            //     let itemCopy = {};
            //     for (let k in item) {
            //         if (!["Title", "Summary", "URL"].includes(k)) {
            //             itemCopy[k] = item[k];
            //         }
            //     }
            //
            //     return itemCopy;
            // }
        }

    };
</script>

<style lang="scss">
    .data_summary {
        color: black;
        background-color: #97c9be;
        padding: 10px;
        position: absolute;
        display: none;
    }
</style>
