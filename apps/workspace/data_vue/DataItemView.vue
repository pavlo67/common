<template>

    <div id="data_item_view">

        <div class="title">{{ itemToImport.Title }}</div>
        <br>
        <div>
            {{ itemToImport.Summary }}
            <br> <span v-html="href(itemToImport.URL)" class="href small"></span>
        </div>
        <span @mouseover="showDetails" @mouseleave="hideDetails" :id=itemId(itemToImport) class="small">
            подробиці...
            <div class="data_summary" :id=itemId(itemToImport,true)>{{ details(itemToImport) }}</div>
        </span>

        <br>
        &nbsp;

    </div>
</template>


<script>
    import e  from '../../components.js/elements';
    import sh from '../../components.js/show_hide/show_hide';

    let itemPostfix = "_details";
    let showHide = sh.NewShowHide(itemPostfix);

    export default {
        name: 'DataItemView',
        props: ['itemToImport'],
        methods: {
            href: e.href,
            itemId(item, postfixed) {
                return "item_to_import_" + item.ID + (postfixed ? itemPostfix : "");
            },

            showDetails(ev) {
                showHide.showContent(ev);
            },

            hideDetails(ev) {
                showHide.hideContent(ev);
            },

            details(item) {
                let itemCopy = {};
                for (let k in item) {
                    if (!["Title", "Summary", "URL"].includes(k)) {
                        itemCopy[k] = item[k];
                    }
                }

                return itemCopy;
            }
        }

    };
</script>

<style lang="scss">
    .data_summary {
        color: black;
        background-color: #97c9be;
        padding: 10px;
        position: absolute;
        display: none;
    }
    #data_item_view {
    }
</style>
