<template>
    <div id="storage_item" class="small">
        <div class="title"><b>Новий запис</b></div>
        <DataEdit v-bind:dataItem="dataItem"/>
        <button v-on:click="save" class="edit_button">зберегти у нотатнику</button>
    </div>
</template>


<script>
    import b       from '../basis';
    import { cfg } from './init';

    export default {
        title: () => 'новий запис',
        data: () => {
            return {
                dataItem: {
                    // Key:      "",
                    // Title:    "type title here",
                    // URL:      "type URL here",
                    // Summary:  "type summary here",
                    // _tags:    "type tags here",
                    // _content: "type content here",
                },
            };
        },
        methods: {
            save() {
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
                        'Content-Type': 'application/json;charset=utf-8',
                    },
                    mode: 'cors', // no-cors, cors, *same-origin
                    body: toSave,

                }).then(response => {
                    return response.json();

                }).then(data => {
                    if (data.ID > 0) {
                        cfg.eventBus.$emit('message', "запис з id = " + data.ID + " збережено");
                        this.$router.push({ name: 'StorageItem',  params: { id: data.ID } })

                    } else {
                        console.log(data);
                        cfg.eventBus.$emit('message', "не вдалось зберегти запис: " + data.Error);
                    }


                });

            }
        },
    }
</script>

<style lang="scss">
    .edit_button {
        height: 30px;
    }
</style>
