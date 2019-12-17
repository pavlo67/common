<template>
    <div id="data">
        <b>Мій каталог</b>

        <div v-for="item in dataItems">
            <span v-html="announce(item)" class="data_announce"></span><br>&nbsp;
        </div>


    </div>
</template>


<script>
    import b from '../../components.js/basis';

    export default {
        name: 'Workspace',
        created () {
            this.getTags();
        },
        data: () => {
            return {
                tags: [],
            };
        },
        methods: {
            announce(j) {
                if (!(typeof j === "object")) return j;

                let text =
                    "[" +  b.dateStr(j.CreatedAt) + "]" +
                    " &nbsp; " + j.Title +
                    "&nbsp;" + "<span class=\"control\">[" +  "докладно" + "][" +  "ред." + "]</span>" +
                    "<br>" + j.Summary;

                return text;
            },

            getTags() {
                fetch('http://localhost:3003/storage/v1/list', {
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
                    }).then(data => {
                    this.dataItems = data;
                    console.log(this.tags);
                });
            }
        },
    }
</script>

<style lang="scss">
    .data_announce {
        color: brown;
        font-size: small;
    }
    .control {
        color: blue;
        font-size: xx-small;
    }
    #data {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
