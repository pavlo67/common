<template>
  <div class="workspace">
    This is a workspace page

    <div v-for="item in workspaceItems">  <!--  v-bind:key="item.path" :to="item.path" -->
      {{ JSON.stringify(item) }}<br>
    </div>


  </div>
</template>


<script>

    // -----------------------------------------------------------------------
    export default {
        name: 'Workspace',
        created () {
            this.getFlowItems();
        },
        data: () => {
            return {
                workspaceItems: [],
            };
        },
        methods: {
            getFlowItems() {
                fetch('http://localhost:3333/workspace/v1/list', {
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
                    this.workspaceItems = data;
                    console.log(this.workspaceItems);
                });
            }
        },

    };
</script>
