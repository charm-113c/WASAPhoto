<template>
    <div class="backdrop">
        <div class="modal">
            <p>{{ this.MsgText }}</p>
            <!-- The buttons are fixed, an 'OK' and a 'Cancel'. A click on either emits corresponding event.
            The deletion process is left to the host. We will just emit an event. So the slot will only carry 
            text. Alright. Sure. 
            Either button closes the modal -->
            <button class="ok-btn" @click="emitOK">Confirm</button>
            <button class="cancel-btn" @click="emitCancel">Cancel</button> 
        </div>
    </div>
</template>

<script>
export default {
    props: ['MsgText'],
    methods: {
        emitOK() {
            this.$emit("ok")
            this.$emit("close")
        },
        emitCancel() {
            this.$emit("close")
        }
    }
}
</script>

<style scoped>
.backdrop {
    top: 0;
    position: fixed;
    background-color: rgba(0,0,0,0.5);
    width: 100%;
    height: 100%;
}
.modal {
    width: 600px;
    height: 200px;
    background-color: var(--strong-colour);
    color: var(--light-colour);
    border-radius: 10px;
    position: fixed;
    left: 33%;
    top: 20%;
    display: grid;
    grid-template-areas: 
        'text text'
        'ok cancel';
    gap: 2px;
    grid-template-rows: 0.9fr 50px;
}
.modal p {
    grid-area: text;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 20px;
}
.modal button {
    border-radius: 15px;
    width: 100px;
    height: 30px;
    color: var(--strong-colour);
    border-color: var(--bold-colour);
    background-color: var(--light-colour);
}
.ok-btn {
    grid-area: ok;
    justify-self: right;
    margin-right: 20px;
}
.cancel-btn {
    grid-area: cancel;
    justify-self: left;
    margin-left: 20px;
}
</style>