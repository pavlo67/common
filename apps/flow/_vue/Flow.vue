<template>
  <div class="flow">
    This is a flow page

    <div v-for="item in flowItems">  <!--  v-bind:key="item.path" :to="item.path" -->
      {{ JSON.stringify(item) }}<br>
    </div>


  </div>
</template>


<script>

    // -----------------------------------------------------------------------
    export default {
        name: 'Flow',
        created () {
            this.getFlowItems();
        },
        data: () => {
            return {
                flowItems: [],
            };
        },
        methods: {
            getFlowItems() {
                fetch('http://localhost:3333/flow/v1/list', {
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
                    this.flowItems = data;
                    console.log(this.flowItems);
                });
            }
        },

    };
</script>
