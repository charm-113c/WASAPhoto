<template>
    <div class="sidebar-left">
        <h2>Hello, {{ username }}</h2>
        <hr>
        <input type="text" v-model="user2" @keyup="searchProfile" placeholder="Search a user">

        <div class="item-list">
            <RouterLink :to="{name: 'stream', params: {username: this.username}}">
                <svg>
                    <use href="/feather-sprite-v4.29.0.svg#home"/>
                </svg>
                <span>Your Stream</span>
            </RouterLink>
            <RouterLink :to="{name: 'myProfile', params: {username: this.username}}">
                <svg>
                    <use href="/feather-sprite-v4.29.0.svg#user"/>
                </svg>
                <span>Your Profile</span>
            </RouterLink>
            <RouterLink to="/" class="log-out">
                <svg>
                    <use href="/feather-sprite-v4.29.0.svg#log-out"/>
                </svg>
                <span>Log out</span>
            </RouterLink>
        </div>
    </div>
    <ErrorMsg v-if="showErrMsg" :msg="errorMsg" @close="closeErrMsg"></ErrorMsg>
</template>

<script>
export default {
    data() {
        return {
            username: sessionStorage.getItem("username"),
            user2: '',
            showErrMsg: false,
            errorMsg: '',
        }
    },
    methods: {
        async searchProfile(e) {
            // route to searched profile view
            // ensure there's something inside user2
            if (e.key === "Enter" && this.user2) {
                try {
                    // eliminate white spaces
                    this.user2 = this.user2.trim()
                    // route (=redirect) to searched profile view
                    await this.$router.push( {name: 'searchProfile', params: {username: this.username, user2: this.user2}} )
                    // forcing a reload because if already in searchProfile view when searching, update is not done correctly
                    window.location.reload()
                } catch (error) {
                    this.handleError(error)
                }
            }
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
    },
}
</script>

<style scoped>
.sidebar-left {
    width: 275px;
    position: fixed;
    left: 0;
    background-color: var(--strong-colour);
    height: 100%;
    border-radius: 10px;
    padding-top: 30px;
    padding-left: 10px;
    margin-left: -5px;
}
.item-list {
    display: flex;
    flex-direction: column;
    margin-top: 15px;
    align-items: baseline;
    height: 100%;
}
svg {
    stroke: var(--light-colour);
	stroke-width: 2;
	stroke-linecap: round;
	stroke-linejoin: round;
	fill: none;
    width: 30px;
    height: 40px;
    margin-right: 10px;
    margin-left: 10px;
}
span, a {
    color: var(--light-colour);
    text-decoration-line: none;
    text-decoration-color: var(--light-colour);
}
a.router-link-exact-active {
    padding: 10px;
    background-color: rgba(0,0,0,0.5);
    border-radius: 30px;
    width: 80%;
}
input {
    background-color: var(--light-colour);
    color: var(--strong-colour);
    border-radius: 15px;
    padding-left: 5px;
    padding: 5px;
    width: 95%;
}
.log-out {
    margin-top: 450px;
}
</style>