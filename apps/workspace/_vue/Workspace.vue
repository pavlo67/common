<template>
    <div id="workspace">

        <div>
            <b>Мій каталог</b>
        </div>

        &nbsp;

        <div v-for="item in dataItems">
            <span v-html="announce(item)" class="announce"></span><br>&nbsp;
        </div>


    </div>
</template>


<script>
  import b from '../../libraries.js/basis';

  export default {
    name: 'Workspace',
    created () {
      this.getDataItems();
    },
    data: () => {
      return {
        dataItems: [],
      };
    },
    methods: {
      announce(j) {
        if (!(typeof j === "object")) return j;

        let text =
            "[" +  b.dateStr(j.CreatedAt) + "]" +
            " &nbsp; " + j.Title +
            "&nbsp;" + "<span class=\"control\">[" +  "ред." + "]</span>" +
            "<br>" + j.Summary;

        return text;
      },

      getDataItems() {
        fetch('http://localhost:3003/workspace/v1/list', {
          method: 'GET', // *GET, POST, PUT, DELETE, etc.
          headers: {
            'Content-Type': 'application/json',
          },
          mode: 'cors', // no-cors, cors, *same-origin

          // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
          // credentials: 'same-origin', // include, *same-origin, omit
          // redirect: 'follow', // manual, *follow, error
          // referrer: 'no-referrer', // no-referrer, *client
        })
        .then(response => {
          return response.json();
        }).then(data => {
          this.dataItems = data;
          console.log(this.dataItems);
        });
      }
    },
  }
</script>

<style lang="scss">
    .announce {
        color: brown;
        font-size: small;
    }
    .control {
        color: blue;
        font-size: xx-small;
    }
    #workspace {
        padding: 0px 10px 10px 10px;
        text-align: left;
    }
</style>
