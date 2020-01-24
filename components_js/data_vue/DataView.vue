<template>

    <div id="data_item_view">

        <div v-if="showTitle">
            {{ dataItem.Title }}
            <br>&nbsp;
        </div>

        <div v-if="dataItem.Summary">
            <span v-html="dataItem.Summary"></span>
            <br>
        </div>

        <span v-html="href(dataItem.URL)" class="href smaller"></span>
        <br>

        <span @mouseover="showDetails" @mouseleave="hideDetails" :id=itemId(dataItem) class="smaller">
            подробиці...
            <div class="data_summary" :id=itemId(dataItem,true)>{{ details(dataItem) }}</div>
        </span>

        <br>
        &nbsp;

    </div>
</template>


<script>
    import e  from '../elements';
    import sh from '../show_hide/show_hide';

    let itemPostfix = "_details";
    let showHide = sh.NewShowHide(itemPostfix);

    export default {
        name: 'DataItemView',
        props: ["dataItem", "showTitle"],
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
</style>
