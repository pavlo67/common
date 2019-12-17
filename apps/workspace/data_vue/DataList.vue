<template>
    <div>
        <div v-for="item in dataItems" class="small">
            <span v-html="announce(item)" class="data_announce"></span><br>&nbsp;
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
        props: ['tags'],
        methods: {
            href: e.href,
            announce(j) {
                if (!(typeof j === "object")) return j;

                console.log(j)

                let text =
                    "[" +  b.dateStr(b.object(j.Status).CreatedAt) + "]" +
                    " &nbsp; " + j.Title +
                    "&nbsp;" + "<span class=\"control\">[" +  "докладно" + "][" +  "ред." + "]</span>" +
                    "<br>" + j.Summary;

                return text;
            },
        },
    }
</script>

<style lang="scss">
    .data_announce {
        font-size: small;
    }
</style>
