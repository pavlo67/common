<template>
    <div id="flow">
        <b>Новини!</b>

        <div v-for="sourcePack in sourcePacks">

            &nbsp;
            <div>
                {{ sourcePack.url }} &nbsp [{{ dateStr(sourcePack.createdAt) }}]

                <div v-for="item in sourcePack.flowItems">  <!--  v-bind:key="item.path" :to="item.path" -->
                    <span v-html="announce(item)" class="announce"></span><br>&nbsp;
                </div>
          </div>

        </div>

    </div>
</template>


<script>
    import b from '../../components.js/basis';

    let cfg = {};

    function init(backend) {
        cfg.backend = backend;

        // TODO: do it safely!!!
        cfg.listEp = window.location.protocol + "//" + window.location.hostname + backend.host + backend.endpoints.flow.path;
    }

    export {init};

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

            announce(j) {
                if (typeof j !== "object" || !j) return j;

                let href = "";

                if (j.Embedded instanceof Array) {
                    for (let embedded of j.Embedded) {
                        href = ' &nbsp; <a href="' + embedded.URL + '" target="_blank">>>></a>';
                        break;
                    }
                }

                let text =
                    "<span class=\"control\">[" + "імпорт" + "]</span>" +
                    " &nbsp; " + j.Title +
                    href;

                return text;
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
                })
                .then(response => {
                    return response.json();
                }).then(flow => {
                    this.flow = flow;
                    this.sourcePacks = [];

                    let source = "";
                    let createdAt = "";
                    let sourcePack = {flowItems: []};

                    for (let item of flow) {
                      if (item.Source+item.CreatedAt !== source+createdAt) {
                        if (sourcePack.flowItems.length > 0) {
                            this.sourcePacks.push(sourcePack);
                        }

                        source = item.Source;
                        createdAt = item.CreatedAt;
                        sourcePack = {url: item.Source, createdAt: item.CreatedAt, flowItems: []};
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
  .control {
    color: blue;
    font-size: xx-small;
  }
  #flow {
    padding: 0px 10px 10px 10px;
    text-align: left;
  }
</style>
