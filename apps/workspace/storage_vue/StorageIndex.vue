<template>
    <div id="data">
        <b>Мій каталог (теми, сиріч теґи, мітки)</b>

        <br>

        <TagsIndex v-bind:tags="tags"/>

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
                tags: [],
            };
        },
        methods: {
            getTags() {
                fetch(cfg.tagsEp, {
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
                    this.tags = data;
                    console.log(this.tags);
                });
            }
        },
    }
</script>

<style lang="scss">
    #data {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
