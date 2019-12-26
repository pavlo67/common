<template>
    <div id="flow_select">
        <b>Вибір новин</b>

        <div>{{ inputData() }}</div>
    </div>
</template>


<script>
    import b  from '../../../components.js/basis';

    export default {
        name: 'Flow selector',
        created () {
        },
        data: () => {
            return {};
        },
        methods: {
            dateStr: b.dateStr,

            init(data) {
            },

            inputData() {
                console.log(6666, this.$route.params.id)
                console.log(7777, this.$route.params.event)
            },

        },
    };
</script>

<style lang="scss">
    #flow_select {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
