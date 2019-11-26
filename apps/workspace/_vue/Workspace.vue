<template>
  <div>
    This is the workspace page!

    <div v-for="item in dataItems">
      {{ item }}
    </div>

    <br>

    !!!


  </div>
</template>


<script>

    // -----------------------------------------------------------------------
    export default {
        name: 'Workspace',
        created () {
            this.getDataItems();
        },
        data: () => {
            return {
                dataItems: [],
            };
        },
        methods: {
            getDataItems() {
                fetch('http://localhost:3003/workspace/v1/list', {
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
                    console.log(this.dataItems);
                });
            }
        },

    };
</script>
