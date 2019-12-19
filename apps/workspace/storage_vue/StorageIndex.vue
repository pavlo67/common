<template>
    <div id="storage_index">
        <b>Мій каталог (теми, сиріч теґи, мітки)</b>

        <br>&nbsp;

        <TagsIndex v-bind:tagsIndex="tagsIndex"/>

    </div>
</template>


<script>
    import b       from '../../components.js/basis';
    import { cfg } from './init';

    export default {
        title: 'мій каталог',
        created () {
            this.getTags();
        },
        data: () => {
            return {
                tagsIndex: [],
            };
        },
        methods: {
            getTags() {
                fetch(cfg.tagsEp + "?key=storage", {
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
                }).then(data => {
                    this.tagsIndex = data;
                    console.log(this.tagsIndex);
                });
            }
        },
    }
</script>

<style lang="scss">
</style>
