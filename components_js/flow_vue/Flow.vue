<template>
    <div id="flow">
        <b>Новини!</b>

        <div v-for="sourcePack in sourcePacks" class="small">

            &nbsp;
            <div>
                {{ sourcePack.url }} &nbsp [{{ dateStr(sourcePack.Time) }}]

                <span v-for="item in sourcePack.flowItems">
                    <br><span :id=announceId(item)  class="control" v-on:click="importData">[імпорт]</span> &nbsp;
                    <span v-html="announce(item)" @mouseover="showSummary" @mouseleave="hideSummary" class="flow_announce" :id=announceId(item,true)></span>&nbsp;
                </span>

            </div>
        </div>

    </div>
</template>


<script>
    import b       from '../basis';
    import sh      from '../show_hide/show_hide';
    import { cfg } from './init';

    let showHide       = sh.NewShowHide("_summary");
    let flowItemPrefix = "flow_item_";

    export default {
        title: 'новини',
        created () {
            this.getFlowItems();
        },
        data: () => {
            return {
                sourcePacks: [],
            };
        },
        methods: {
            dateStr: b.dateStr,

            announceId(j, prefixed) {
               return (prefixed ? flowItemPrefix : "") + j.ID;
            },

            announce(j) {
                if (typeof j !== "object" || !j) return j;

                let href = ' &nbsp; <a href="' + j.URL + '" target="_blank">>>></a>';
                let text = j.Title + href +
                    "<div class=\"flow_summary\" id=\"" + flowItemPrefix + j.ID + "_summary\">" + j.Summary + "</div>";

                return text;
            },

            showSummary(ev) {
                showHide.showContent(ev);
            },

            hideSummary(ev) {
                showHide.hideContent(ev);
            },

            importData(ev) {
                fetch(cfg.readEp + "/" + ev.target.id, {
                    method: 'GET', // *GET, POST, PUT, DELETE, etc.
                    headers: {
                        'content-type': 'application/json',
                        'authorization': cfg.user && cfg.user.Creds.jwt,
                    },
                    mode: 'cors', // no-cors, cors, *same-origin

                }).then(response => {
                    return response.json();

                }).then(flowItem => {
                    this.$router.push({ name: 'StorageItemImport',  params: { dataItem: flowItem } })

                });
            },
            
            getFlowItems() {
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

                }).then(flow => {
                    this.sourcePacks = [];

                    let source = "";
                    let sourceTime = "";
                    let sourcePack = {flowItems: []};

                    for (let item of flow) {
                        if (!item.Origin) {
                            item.Origin = {};
                        }

                        if (item.Origin.Source != source || item.Origin.Time != sourceTime) {
                            if (sourcePack.flowItems.length > 0) {
                                this.sourcePacks.push(sourcePack);
                            }

                            sourcePack = {url: (source = item.Origin.Source), Time: (sourceTime = item.Origin.Time), flowItems: []};
                        }

                        sourcePack.flowItems.push(item);
                    }

                    if (sourcePack.flowItems.length > 0) {
                        this.sourcePacks.push(sourcePack);
                    }

                });
            }
        },
    };
</script>

<style lang="scss">
    .flow_announce {
        color: brown;
    }
    .flow_summary {
        color: black;
        background-color: #97c9be;
        padding: 10px;
        position: absolute;
        display: none;
    }
</style>
