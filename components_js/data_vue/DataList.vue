<template>
    <div>
        <div v-for="item in dataList" class="small">
            <span class="data_announce">
                <span class="time">[{{ dateStr(object(item.Status).CreatedAt) }}]</span> {{ item.Title }}
                 &nbsp; <span class="control" v-on:click="$router.push({ name: 'StorageItem',  params: { id: item.ID } })">[докладно][ред]</span>
                <br><span v-html="item.Summary" class="smaller"></span>
            </span>
            <br>&nbsp;
        </div>
    </div>
</template>


<script>
    import b from '../basis';
    import e from '../elements';

    export default {
        name: 'DataList',
        props: ['dataList'],
        methods: {
            object:  b.object,
            dateStr: b.dateStr,
            href:    e.href,
        },
    }
</script>

<style lang="scss">
</style>
