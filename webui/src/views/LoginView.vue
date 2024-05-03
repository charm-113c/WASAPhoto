<template>
    <form @submit.prevent="doLogin">
        <label for="username">Username:</label>
        <input type="text" id="username" required v-model="username" placeholder="Username">
        <!-- <button class="submit" @click="doLogin">Login</button> -->
        <br>
        <p>Log in to WASAPhoto! Whether you're a new user or not, entering a username is all it takes!</p>
    </form>
    <ErrorMsg v-if="showErrMsg" :msg="errorMsg" @close="closeErrMsg"></ErrorMsg>
</template>

<script>
export default {
    data() {
        return {
            username: '',
            errorMsg: '',
            showErrMsg: false,
        }
    },
    methods: {
        async doLogin() {
            try {
                let res = await this.$axios.post("/login",
                    this.username, // request body
                    {headers: {'Content-Type': 'text/plain'}}) 
                // get user's token
                sessionStorage.setItem('username', this.username)
                sessionStorage.setItem('bearerToken', res.data)
                this.directToStream()
            } catch (error) {
                this.handleError(error)
            }
        },
        directToStream() {
            // before routing, we also trigger an event to tell the App to display 
            // Stream and eliminate the Login view
            this.$emit('loggedIn')
            this.$router.push({name: 'stream', params: {username: this.username}})
        },
        handleError(error) {
            if (error.response) { 
                    // check if the error is from the response
                    this.errorMsg = error.response.data
                } else if (error.request) {
                    // or from the request itself
                    this.errorMsg = error.request
                } else {
                    this.errorMsg = error.message
                }
            this.showErrMsg = true
        },
        closeErrMsg() {
            this.showErrMsg = false
        },
    }
}
</script>

<style scoped>
form {
    margin: 70px auto;
    background-color: var(--strong-colour);
    color: var(--light-colour);
    text-align: left;
    padding: 40px;
    border-radius: 5px;
}
label {
    color: var(--light-colour);
    display: inline-block;
    margin: 25px 0 15px;
    font-size: 0.7em;
    text-transform: uppercase;
    letter-spacing: 1px;
    font-weight: bold;
}
input {
    display: block;
    padding: 10px 6px;
    width: 50%;
    box-sizing: border-box;
    border: none;
    border-radius: 3px;
    border-bottom: 1px solid #2c3e50;
    color: var(--strong-colour);
}
p {
    color: var(--light-colour);
}
</style>