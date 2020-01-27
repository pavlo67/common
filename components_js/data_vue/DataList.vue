<template>
    <div>

        <table align="right" class="table_right">
            <tr><td>
                <span class="control">
                    [<span v-on:click="$router.push({ name: 'NoteEdit',  params: {  } })">експортувати</span>]
                </span>
            </td></tr>
        </table>

        <div v-for="dataItem in dataList" class="small">
            <!-- TODO: customize the router target -->
            <span class="data_announce">
                <span class="time">[{{ createdAt(dataItem) }}]</span>
                {{ dataItem.Title }}
                <span class="control smallest">
              &nbsp;     [<span v-on:click="$router.push({ name: 'NoteView',  params: { id: dataItem.ID } })">докладно</span>]
                    &nbsp;
                    [<span v-on:click="$router.push({ name: 'NoteEdit',  params: { id: dataItem.ID } })">ред.</span>]
                </span>
                <span v-if="dataItem.Summary">
                    <br><span v-html="dataItem.Summary" class="smaller"></span>
                </span>
                <span v-if="dataItem.Tags">
                    <br>
                    <span v-for="tag in dataItem.Tags" class="tag control smaller" >
                        [<span v-on:click="$router.push({ name: 'ListTagged',  params: { tag: tag.Label } })">{{ tag.Label }}</span>]
                        &nbsp;
                    </span>
                </span>
            </span>
            <br>&nbsp;
        </div>
    </div>
</template>


<script>
    import b from '../basis';
    import e from '../elements';
    import {createdAt} from './data';

    export default {
        name: 'DataList',
        props: ['dataList'],
        methods: {
            createdAt: createdAt,
            object:  b.object,
            href:    e.href,
        },
    }
</script>

<style lang="scss">
</style>
