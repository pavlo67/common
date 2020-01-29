<template>
    <div id="action_steps">
        <a v-on:click="activate" class="control">{{ activationText }}</a>

        <div v-if="activated">
            <p/>
            <div v-if="readyToDo">
                <a v-bind:href="data" v-on:click="deactivate" download="data.json">{{ actionButtonText }}</a>
<!--                <button v-on:click="action" class="do_it">{{ actionButtonText }}</button>-->
            </div>
            <div v-else>
                {{ preparationText }}
            </div>
        </div>
    </div>
</template>

<script>
    // import e  from '../elements';

    export default {
        name: 'ActionSteps',
        props: ["activationText", "activationAction", "preparationText", "actionButtonText"],

        data: () => {
            return {
                activated: false,
                readyToDo: false,
                data: '',
            };
        },

        methods: {
            activate() {
                this.activated = true;
                this.readyToDo = false;
                this.activationAction(data => {
                    this.data = "data:application-json;base64," + btoa(unescape(encodeURIComponent(data)));
                    console.log("RETURNED FROM ACTION...");
                    this.readyToDo = true;
                });
            },

            deactivate() {
                console.log("SAVED: ", this.data);
                this.activated = false;
                this.readyToDo = false;

                return true;
            }
        }

    };
</script>

<style lang="scss">
    #action_steps {
        padding: 10px;
    }
</style>
