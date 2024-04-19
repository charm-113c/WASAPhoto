<template>

    <!-- Here we display the profile of the searched user. Thankfully, backend provides us with the data, we just fill in.
    This page is only accessed upopn being routed to it. That's a plus.   -->
    <!-- First we show general user data  -->
    <div v-if="!errorMsg">
        <div class="section">
            <p>{{ user2 }} - Number of photos: {{ nphotos }} - Followers: {{ nFollowers }} - Following: {{ nFollowing }}</p>
            <button v-if="!userFollowed" @click="followUser">Follow</button>
            <button v-else @click="unfollowUser">Unfollow</button>
            <button v-if="!userBanned" @click="banUser">Ban</button>
            <button v-else @click="unbanUser">Unban</button>

            <p>Urgh... Style Followers and Following so that they display the corresponding list when clicked</p>  
        </div> 
        <!-- Then the photos, if they exist -->
        <div v-for="photo in userPhotos" v-if="!noPhotos" :id="photo.ID">
            <img :src="photo.src" :alt="photo.description" @click="showImage"/> <br>
            <span>{{ photo.uploadDate }} - Likes: {{ photo.likes }} - Comments: {{ photo.comments }}</span>
            <button v-if="!photo.liked" @click="likePhoto">Like</button>
            <button v-else @click="unlikePhoto">Unlike</button>
            <p>{{ photo.description }}</p>
            <hr>
        </div>

        <p v-if="noPhotos">{{ user2 }} hasn't uploaded any photo yet</p>
    </div>

    <ErrorMsg v-if="errorMsg" :msg="errorMsg" :code="errorCode"></ErrorMsg>
    
</template>

<script>
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
            errorCode: null,
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
                this.followers = res.data.Followers
                this.nFollowers = this.followers !== null? this.followers.length : 0
                this.following = res.data.Following
                this.nFollowing = this.following !== null? this.following.length : 0
                // photos are a particular case
                if (res.data.Photos === null) {
                    this.noPhotos = true
                } else {
                    res.data.Photos.forEach(photo => {
                        // transform img data into data URI so as to feed it to img tag
                        let imgSRC = `data:${photo.FileExtension};base64,${photo.BinaryData}`
                        this.userPhotos[photo.Uploader + photo.PhotoID] = {
                            'src': imgSRC, 'description': photo.Description, 'uploadDate': new Date(photo.UploadDate), 
                            'likes': photo.Likes, 'comments': photo.Comments, 'ID': photo.PhotoID, 'liked': false
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
        async followUser() {
            try {
                await this.$axios.put(`/users/${this.username}/following/${this.user2}`, null,
                                {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.userFollowed = true
            } catch (error) {
                this.handleError(error)
            }
        },
        async unfollowUser() {
            try {
                await this.$axios.delete(`/users/${this.username}/following/${this.user2}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.userFollowed = false
            } catch (error) {
                this.handleError
            }
        },
        async banUser() {
            try {
                await this.$axios.put(`/users/${this.username}/blacklist/${this.user2}`, null,
                                {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.userBanned = true
            } catch (error) {
                this.handleError(error)
            }
        },
        async unbanUser() {
            try {
                await this.$axios.delete(`/users/${this.username}/blacklist/${this.user2}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.userBanned = false
            } catch (error) {
                this.handleError(error)
            }
        },
        async likePhoto(e) {
            try {
                // elemID is uploaderName+'/'+photoID, we want them separated
                const pID = e.target.parentElement.id
                await this.$axios.put(`/users/${this.user2}/photos/${pID}/likes/${this.username}`, 
                                    null, // empty body
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.userPhotos[this.user2 + pID].liked = true
                this.userPhotos[this.user2 + pID].likes++
            } catch (error) {
                this.handleError(error)
            }
        },
        async unlikePhoto(e) {
            try {
                const pID = e.target.parentElement.id
                await this.$axios.delete(`/users/${this.user2}/photos/${pID}/likes/${this.username}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.userPhotos[this.user2 + pID].liked = false
                this.userPhotos[this.user2 + pID].likes--
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
                    this.errorCode = error.response.status
                    this.errorMsg = error.response.data
                } else if (error.request) {
                    // or from the request itself
                    this.errorCode = 400
                    this.errorMsg = error.request
                } else {
                    this.errorCode = 500
                    this.errorMsg = error.message
                }
        }
    },
    beforeMount() {
        this.getUserProfile()
    },
}
</script>

<style scoped>
p {
    color: black
}
</style>