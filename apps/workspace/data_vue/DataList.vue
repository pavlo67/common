<template>
    <div>
        <div v-for="item in dataList" class="small">
            <span class="data_announce">
                 [{{ dateStr(object(item.Status).CreatedAt) }}] &nbsp; {{ item.Title }}
                 &nbsp; <span class="control" v-on:click="$router.push({ name: 'StorageItem',  params: { id: item.ID } })">[докладно][ред]</span>
                <br><span v-html="item.Summary"></span>
            </span>
            <br>&nbsp;
        </div>
    </div>
</template>


<script>
    import b from '../../components.js/basis';
    import e from '../../components.js/elements';

    // import sh from '../../components.js/show_hide/show_hide';
    // let itemPostfix = "_details";
    // let showHide = sh.NewShowHide(itemPostfix);

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
