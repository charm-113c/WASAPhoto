<template>

    <Sidebar></Sidebar>
    <div class="container">
        <div class="searched-profile">
            <div class="metadata">
                <span class="the-man"><h5>The one and only, the myth, the legend:</h5>  <h3>{{ user2 }}</h3></span>
                <div class="f-buttons">
                    <svg @click="toggleFollow" :class="{'followed':userFollowed}">
                        <use href="/feather-sprite-v4.29.0.svg#user-plus"/>
                    </svg>
                    <!-- <button v-if="!userFollowed" @click="followUser">Follow</button>
                    <button v-else @click="unfollowUser">Unfollow</button> -->
                    <svg @click="toggleBan" :class="{'banned': userBanned}">
                        <use href="/feather-sprite-v4.29.0.svg#user-x"/>
                    </svg>
                    <!-- <button v-if="!userBanned" @click="banUser">Ban</button>
                    <button v-else @click="unbanUser">Unban</button> -->
                </div>
                <div class="f-text">
                    <span v-if="!userFollowed">Not Followed</span>
                    <span v-else>Followed</span>

                    <span v-if="!userBanned">Ban User</span>
                    <span v-else>Unban User</span>
                </div>
                
                <span>
                    <details><summary>Followers: {{ nFollowers }}</summary> {{ followers }}</details>
                    <details><summary>Following: {{ nFollowing }}</summary> {{ following }}</details>
                </span>
                <span class="stats">Number of photos: {{ nphotos }}</span>
            </div>
            <hr>

            <PhotoContainer :photos="userPhotos"></PhotoContainer>

            <p v-if="noPhotos">{{ user2 }} hasn't uploaded any photo yet</p>
        </div>
    </div>

    <ErrorMsg v-if="showErrMsg" :msg="errorMsg" @close="closeErrMsg"></ErrorMsg>
    
</template>

<script>
import PhotoContainer from '../components/PhotoContainer.vue'

export default {
    props: ['user2'],
    data() {
        return {
            username: sessionStorage.getItem('username'),
            userPhotos: {},
            // { photoID: { 'URL': str, 'description': str, 'uploadDate': date, 'likes': int, 'comments': int} }
            nphotos: 0,
            followers: [], // list of names
            nFollowers: 0,
            following: [],
            nFollowing: 0,
            showErrMsg: false,
            errorMsg: '',
            noPhotos: false,
            userFollowed: false,
            userBanned: false,
        }
    },
    methods: {
        async getUserProfile() {
            try {
                const res = await this.$axios.get(`/users/${this.user2}/profile`, 
                    { headers: {'Authorization': sessionStorage.getItem('bearerToken'), 'requesting-user': this.username} }) 
                // parse response
                this.nphotos = res.data.Nphotos
                if (res.data.Followers !== null) {
                    this.followers = res.data.Followers.join(', ')
                    this.nFollowers = res.data.Followers.length
                    this.userFollowed = res.data.Followers.includes(this.username)
                } else {
                    this.followers = ''
                    this.nFollowers = 0
                }
                if (res.data.Following !== null) {
                    this.following = res.data.Following.join(', ')
                    this.nFollowing = res.data.Following.length
                } else {
                    this.following = ''
                    this.nFollowing = 0
                }
                // photos are a particular case
                if (res.data.Photos === null) {
                    this.noPhotos = true
                } else {
                    res.data.Photos.forEach(photo => {
                        // transform img data into data URI so as to feed it to img tag
                        let imgSRC = `data:${photo.FileExtension};base64,${photo.BinaryData}`
                        this.userPhotos[photo.Uploader + photo.PhotoID] = {
                            'src': imgSRC, 'uploader': photo.Uploader, 'description': photo.Description, 'uploadDate': new Date(photo.UploadDate).toUTCString(), 
                            'likes': photo.Likes, 'comments': photo.Comments, 'id': photo.PhotoID, 'liked': false
                        }
                        if (photo.Likers) {
                            this.userPhotos[photo.Uploader + photo.PhotoID].liked = photo.Likers.includes(this.username)
                        }
                    })
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async toggleFollow() {
            try {
                // if not followed, follow
                if (!this.userFollowed) {
                    await this.$axios.put(`/users/${this.username}/following/${this.user2}`, null,
                                {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.userFollowed = true
                    // artificially increase nfollowers and follows
                    this.nFollowers++
                    this.followers = this.followers === ''? this.username : this.followers + ', ' + this.username
                }  else { // do the opposite
                    await this.$axios.delete(`/users/${this.username}/following/${this.user2}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.userFollowed = false
                    this.nFollowers--
                    // First take away the commas if they exist
                    this.followers = this.followers.replace(', ', '')
                    // then the username itself
                    this.followers = this.followers.replace(this.username, '')
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async toggleBan() {
            try {
                if (!this.userBanned) {
                    await this.$axios.put(`/users/${this.username}/blacklist/${this.user2}`, null,
                                {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.userBanned = true
                } else {
                    await this.$axios.delete(`/users/${this.username}/blacklist/${this.user2}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.userBanned = false
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async unbanUser() {
            try {
                
            } catch (error) {
                this.handleError(error)
            }
        },
        showImage(e) {
            this.$router.push({name: 'viewImage', params: {uploader: this.user2, photoID: e.target.parentElement.id}})
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
    beforeMount() {
        this.getUserProfile()
    },
}
</script>

<style scoped>
svg {
    width: 35px;
    height: 40px;
}
svg:hover {
    cursor: pointer;
}
</style>