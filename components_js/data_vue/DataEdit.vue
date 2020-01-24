<template>
    <div id="data_item_edit">
        <table class="edit_table">
        <tr v-for="field in fields" valign="top" v-if="!field.omitEmpty || (dataItem[field.key] != '' && dataItem[field.key] != undefined)">
            <td class="title_cell">{{ field.title }}</td>

            <!-- TODO!!! v-on:change="saveTemporary() -->

            <td>
                <div v-if="field.type === 'view'">
                    {{ dataItem[field.key] }}
                </div>

                <div v-else-if="field.type === 'textarea'">
                    <textarea v-model="dataItem[field.key]" v-bind:placeholder="' ' + (field.placeholder || field.title)" class="edit_field" v-bind:rows="field.lines || 2"/>
                </div>

                <div v-else-if="field.type === 'tags'">
                    <vue-tags-input
                         v-model="tag"
                         :tags="dataItem[field.key]"
                         :placeholder="field.placeholder || field.title"
                         @tags-changed="newTags => dataItem[field.key] = newTags"
                    />
                </div>

                <div v-else>
                    <input v-model="dataItem[field.key]" v-bind:placeholder="'   ' + (field.placeholder || field.title)" class="edit_field" />
                </div>
            </td>
        </tr>
        </table>
    </div>
</template>

<script>
    import e  from '../elements';

    const fields = [
        {key: "Key",      title: "ключ запису",           type: "view",     omitEmpty: true },
        {key: "Title",    title: "заголовок"},             // , placeholder: "введіть заголовок запису"
        {key: "URL",      title: "URL",},
        {key: "_tags",    title: "теґи",                  type: "tags"},
        {key: "Summary",  title: "короткий зміст, анонс", type: "textarea", lines: 3},
        {key: "_content", title: "сам запис"            , type: "textarea", lines: 30},
    ];

    export default {
        name: 'DataEdit',
        data: () => {
            return {
                tag: '',
                fields,
            };
        },
        props: ["dataItem"],
        methods: {
        },
    };
</script>

<style lang="scss">
    .edit_table {
        width: 100%;
    }

    .title_cell {
        /*min-width: 300px;*/
        padding-right: 10px;
        width: 20%;
    }

    .edit_field {
        /*min-width: 300px;*/
        width: 100%;
    }

    /*.vue-tags-input {*/
    /*    width: 100%;*/
    /*    padding: 0px;*/
    /*    align: 0px;*/
    /*}*/

    /*.vue-tags-input ::-webkit-input-placeholder {*/
    /*    padding: 0px;*/
    /*    align: 0px;*/
    /*}*/

    /*.vue-tags-input ::-moz-placeholder {*/
    /*    padding: 0;*/
    /*    align: 0px;*/
    /*}*/

    /*.vue-tags-input :-ms-input-placeholder {*/
    /*    padding: 0;*/
    /*    align: 0px;*/
    /*}*/

    /*.vue-tags-input :-moz-placeholder {*/
    /*    padding: 0;*/
    /*    align: 0px;*/
    /*}*/

</style>


<!--<template>-->
<!--    <div>-->
<!--        <vue-tags-input-->
<!--                v-model="tag"-->
<!--                :tags="tags"-->
<!--                @tags-changed="newTags => tags = newTags"-->
<!--        />-->
<!--    </div>-->
<!--</template>-->
<!--<script>-->
<!--    import VueTagsInput from '@johmun/vue-tags-input';-->

<!--    export default {-->
<!--        components: {-->
<!--            VueTagsInput,-->
<!--        },-->
<!--        data() {-->
<!--            return {-->
<!--                tag: '',-->
<!--                tags: [],-->
<!--            };-->
<!--        },-->
<!--    };-->
<!--</script>-->