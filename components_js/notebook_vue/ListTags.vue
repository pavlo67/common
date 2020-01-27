<template>
    <div id="storage_index">
        <div class="title"><b>Теми (сиріч теґи, мітки)</b></div>

        <ListTags listTags="listTags"/>

    </div>
</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        title: () => 'теми нотатника',
        created () {
            this.getTags();
        },
        data: () => {
            return {
                listTags: [],
            };
        },
        methods: {
            getTags() {
                console.log(cfg.tagsEp + "?key=storage")
                // fetch(cfg.tagsEp + "?key=storage", {
                //     method: 'GET', // *GET, POST, PUT, DELETE, etc.
                //     headers: {
                //         'content-type': 'application/json',
                //         'authorization': cfg.user && cfg.user.Creds.jwt,
                //     },
                //     mode: 'cors', // no-cors, cors, *same-origin
                //
                //     // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
                //     // credentials: 'same-origin', // include, *same-origin, omit
                //     // redirect: 'follow', // manual, *follow, error
                //     // referrer: 'no-referrer', // no-referrer, *client
                // }).then(response => {
                //     return response.json();
                // }).then(data => {
                //     this.tagsIndex = data;
                //     console.log(this.tagsIndex);
                // });
            }
        },
    }
</script>

<style lang="scss">
</style>
