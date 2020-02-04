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
                    <textarea v-model="dataItem[field.key]" v-bind:rows="field.lines || 2" v-bind:placeholder="' ' + (field.placeholder || field.title)" class="edit_field" />
                </div>

                <div v-else-if="field.type === 'editor'">
                    <editor
                        api-key="no-api-key"
                        v-model="dataItem[field.key]"
                        :init="{
                            height: 500,
                            menubar: false,
                            plugins: [
                                'advlist autolink lists link image charmap print preview anchor',
                                'searchreplace visualblocks code fullscreen',
                                'insertdatetime media table paste code help wordcount'
                            ],
                            toolbar:
                                'undo redo | formatselect | bold italic backcolor | \
                                alignleft aligncenter alignright alignjustify | \
                                bullist numlist outdent indent | removeformat | table | help'
                        }"
                    ></editor>

                    <!-- <vue-editor id="aaa" v-model="dataItem[field.key]" v-bind:rows="field.lines || 2" v-bind:placeholder="' ' + (field.placeholder || field.title)" />-->
                    <!-- class="edit_field" -->
                    <!-- <froala :tag="'textarea'" v-model="dataItem[field.key]" v-bind:rows="field.lines || 2" :config="config"></froala> -->
                </div>

                <div v-else-if="field.type === 'select'">
                    <select v-model="dataItem['select.' + field.key]" class="edit_field">
                        <option v-for="v in field.values" v-bind:value="v.key">{{ v.name }}</option>
                    </select>
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
        <tr><td>
            <button v-on:click="save" class="edit_button">зберегти у нотатнику</button>
        </td></tr>
        </table>
    </div>
</template>



<script>

    // http://www.vue-tags-input.com/#/
    // https://www.froala.com/wysiwyg-editor/docs/framework-plugins/vue

    // import VueFroala from 'vue-froala-wysiwyg';
    // import {VueEditor} from "vue2-editor";
    import Editor from '@tinymce/tinymce-vue'

    import b    from '../basis';
    import Auth from '../auth_vue/Auth.vue';

    const accessKey = "_access";
    const groupsCommon = [
        {key: '', name: 'приватний запис'},
        {key: '', name: 'запис для загалу'},
    ];

    const fields = b.copy([
        {key: "Key",      title: "ключ запису",           type: "view",     omitEmpty: true },
        {key: "Title",    title: "заголовок"},             // , placeholder: "введіть заголовок запису"
        {key: accessKey,  title: "доступність",           type: "select"},
        {key: "URL",      title: "URL",},
        {key: "Summary",  title: "короткий зміст, анонс", type: "textarea", lines: 3},
        {key: "_tags",    title: "теґи",                  type: "tags"},
        {key: "_content", title: "сам запис"            , type: "editor",   lines: 30},
    ]);

    let cfg = {};

    export default {
        name: 'DataEdit',
        components: {editor: Editor},
        data: () => {
            return {
                tag: '',
                fields,
                config: {
                    events: {
                        'froalaEditor.initialized': () => {
                            console.log('FROALA initialized')
                        }
                    }
                },

            };
        },
        props: ["dataItem"],
        methods: {
            prepare(dataItem, cfgCommon) {
                cfg      = cfgCommon;
                let user = cfg.common.user || {};

                for (let f of fields) {
                    if (f.key === accessKey && user.groups instanceof Array) {
                        f.values = [...groupsCommon, ...user.groups];
                        f.values[0].key = user.Key;
                    }
                }

                if (!(dataItem instanceof Object)) {
                    dataItem = {};
                    dataItem['select.' + accessKey] = user.Key;
                }

                dataItem['select.' + accessKey] = user.Key;

                if (dataItem.Tags instanceof Array) {
                    dataItem._tags = dataItem.Tags.map(t => ({text: t.Label}));
                    delete dataItem.Tags;
                }

                if (dataItem.Data instanceof Object) {
                    dataItem._content = dataItem.Data.Content;
                    delete dataItem.Data;
                }

                return dataItem;
            },

            save(){

                let user = cfg.common.user || {};

                if ("_content" in this.dataItem) {
                    this.dataItem.Data = {
                        TypeKey: "string",
                        Content: this.dataItem._content,
                    };
                    delete this.dataItem._content;
                }

                if ("_tags" in this.dataItem) {
                    this.dataItem.Tags = this.dataItem._tags.map(
                        t => ({Label: t.text})
                    );
                    delete this.dataItem._tags;
                }

                let toSave = JSON.stringify(this.dataItem)
                console.log("TO SAVE: " + toSave);

                fetch(cfg.saveEp, {
                    method: 'POST',
                    headers: {
                        'content-type': 'application/json;charset=utf-8',
                        'authorization': user.Creds && user.Creds.jwt,
                    },
                    mode: 'cors', // no-cors, cors, *same-origin
                    body: toSave,

                }).then(response => {
                    return response.json();

                }).then(data => {
                    if (data.id) {
                        cfg.eventBus.$emit('message', "запис з id = " + data.id + " збережено");
                        this.$router.push({ name: 'NoteView',  params: { id: data.id } })

                    } else {
                        console.log(data);
                        cfg.eventBus.$emit('message', "не вдалось зберегти запис: " + data.Error);
                    }
                });
            },
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
