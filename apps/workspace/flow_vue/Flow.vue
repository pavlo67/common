<template>
    <div id="flow">
        <b>Новини!</b>

        <div v-for="sourcePack in sourcePacks">

            &nbsp;
            <div>
                {{ sourcePack.url }} &nbsp [{{ dateStr(sourcePack.Time) }}]

                <span v-for="item in sourcePack.flowItems">  <!--  v-bind:key="item.path" :to="item.path" -->
                    <br><span class="control" v-on:click="importData">[імпорт]</span> &nbsp;
                    <span v-html="announce(item)" @mouseover="showSummary" @mouseleave="hideSummary" class="announce" :id=announceId(item)></span>&nbsp;
                </span>
            </div>

        </div>

    </div>
</template>


<script>
    import b  from '../../components.js/basis';
    import sh from '../../components.js/show_hide/show_hide';

    let showHide = sh.NewShowHide("_summary");
    let listEp, router;

    export default {
        name: 'Flow',
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

            init(data) {
                router = data.router;
                // TODO: do it safely!!!
                listEp = window.location.protocol + "//" + window.location.hostname + data.backend.host + data.backend.endpoints.flow.path;
            },

            announceId(j) {
               return "flow_item_" + j.ID;
            },

            announce(j) {
                if (typeof j !== "object" || !j) return j;

                // let href = "";
                // if (j.Embedded instanceof Array) {
                //     for (let embedded of j.Embedded) {
                //         href = ' &nbsp; <a href="' + embedded.URL + '" target="_blank">>>></a>';
                //         break;
                //     }
                // }

                let href = ' &nbsp; <a href="' + j.URL + '" target="_blank">>>></a>';

                let text = j.Title + href +
                    "<div class=\"summary\" id=\"flow_item_" + j.ID + "_summary\">" + j.Summary + "</div>";

                return text;
            },

            showSummary(ev) {
                showHide.showContent(ev);
            },

            hideSummary(ev) {
                showHide.hideContent(ev);
            },

            importData(ev) {
                console.log(555555555, ev);
                // self.$router ???
                router.push({ name: 'flow_selector',  params: { id: 777, event: ev } })
            },
            
            getFlowItems() {
                fetch(listEp, {
                    method: 'GET', // *GET, POST, PUT, DELETE, etc.
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    mode: 'cors', // no-cors, cors, *same-origin

                    // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
                    // credentials: 'same-origin', // include, *same-origin, omit
                    // redirect: 'follow', // manual, *follow, error
                    // referrer: 'no-referrer', // no-referrer, *client
                })
                .then(response => {
                    return response.json();
                }).then(flow => {
                    this.flow = flow;
                    this.sourcePacks = [];

                    let source = "";
                    let sourceTime = "";
                    let sourcePack = {flowItems: []};

                    for (let item of flow) {
                      if (item.Source != source || item.Time != sourceTime) {
                        if (sourcePack.flowItems.length > 0) {
                            this.sourcePacks.push(sourcePack);
                        }

                        // console.log(item);

                        sourcePack = {url: (source = item.Source), Time: (sourceTime = item.Time), flowItems: []};
                      }

                      sourcePack.flowItems.push(item);
                    }

                    if (sourcePack.flowItems.length > 0) {
                        this.sourcePacks.push(sourcePack);
                    }

                    // console.log(11111111, this.sourcePacks)
                });
            }
        },
    };
</script>


<style lang="scss">
  .announce {
    color: brown;
    font-size: small;
  }
  .summary {
      color: black;
      background-color: #97c9be;
      padding: 10px;
      position: absolute;
      display: none;
  }
  .control {
    color: blue;
    font-size: xx-small;
  }
  #flow {
    padding: 0px 10px 10px 10px;
    text-align: left;
  }
</style>
