<template>
    <ErrorMsg v-if="errorMsg" :msg="errorMsg" :code="errorCode"></ErrorMsg>
    Welcome to your profile, {{ username }}. Here, you can do stuff. Stuff includes:
    <ul>
        <li><details><summary>Followers: {{ nFollowers }}</summary> <p>{{ followers }}</p></details></li>
        <li><details><summary>Following: {{ nFollowing }}</summary> <p>{{ following }}</p></details></li>
        <li>Uploading new photos</li>
        <li>If we're feeling frisky, deleting photos and comments</li>
        <li>
            <input type="text" v-model="newUsername">
            <button @click="changeUsername">Change your username</button>
        </li>
    </ul>
    Priority goes to uploading your own photos, so here you go: 
    <form>
        <label for="images">Click to select images</label> <br>
        <input type="file" id="images" accept="image/png, image/jpeg" @change="selectImage" required/>
        <br>
        <label for="description">Description [optional]</label> <br>
        <input type="text" id="description" v-model="description">
        <div class="submit">
            <button @click="uploadPhotos">Upload</button>
        </div>
    </form>
    <div v-for="photo in this.myPhotos" :id="photo.ID">
        <img :src="photo.src" :alt="photo.description" @click="showImage">
        <button @click="deletePhoto">Delete Photo</button> <br>
        <span>{{ photo.uploadDate }} - Likes: {{ photo.likes }} - Comments: {{ photo.comments }}</span>
        <button v-if="!photo.liked" @click="likePhoto">Like</button>
        <button v-else @click="unlikePhoto">Unlike</button>
        <p>{{ photo.description }}</p>
    </div>

</template>

<script>
export default {
    data() {
        return {
            username: sessionStorage.getItem('username'),
            // we'll show user's own profile:
            myPhotos: {}, 
            noPhotos: true,
            nphotos: 0,
            followers: [], // list of names
            nFollowers: 0,
            following: [],
            nFollowing: 0,
            // and we'll allow them to post photos
            selectedFile: null,
            description: '',
            errorCode: null,
            errorMsg: '',
            // or change username
            newUsername: '',
        }
    },
    methods: {
        selectImage(event) {
            this.selectedFile = event.target.files[0]
            this.previewImgURL = URL.createObjectURL(this.selectedFile)
        },
        async uploadPhotos() {
            // we create a formData to send
            let formData = new FormData()
            // we extract the necessary metadata
            formData.append('Description', this.description)
            formData.append('UploadDate', new Date().toISOString())
            formData.append('UploadedImage', this.selectedFile)
            try {
                await this.$axios.post(`/users/${this.username}/photos`, formData,
                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                // upon success show message and void inputs
                // then update the current stream with new the photo
                this.refresh()
                } catch (error) {
                this.handleError(error)
            }
            
        },
        async getMyPhotos() {
            try {
                // show user's own photos in their profile
                const res = await this.$axios.get(`/users/${this.username}/profile`,
                    {headers: {'Authorization': sessionStorage.getItem('bearerToken'), 'requesting-user': this.username}})
                this.nphotos = res.data.Nphotos
                this.followers = res.data.Followers 
                this.nFollowers = this.followers !== null ? this.followers.length : 0
                this.following = res.data.Following 
                this.nFollowing = this.following !== null ? this.following.length : 0
                // as usual, photos need special treatment
                if (res.data.Photos === null) {
                        this.noPhotos = true
                } else {
                    res.data.Photos.forEach(photo => {
                        // transform img data into data URI so as to feed it to img tag
                        let imgSRC = `data:${photo.FileExtension};base64,${photo.BinaryData}`
                        this.myPhotos[photo.Uploader + photo.PhotoID] = {
                            'src': imgSRC, 'description': photo.Description, 'uploadDate': new Date(photo.UploadDate), 
                            'likes': photo.Likes, 'comments': photo.Comments, 'ID': photo.PhotoID, 'liked': false
                        }
                        // set 'liked' to true if userrname is liking current photo
                        if (photo.Likers) {
                            this.myPhotos[photo.Uploader + photo.PhotoID].liked = photo.Likers.includes(this.username)
                        }
                    })
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async changeUsername() {
            try {
                if (this.newUsername) {
                    // trim white spaces
                    this.newUsername = this.newUsername.trim()
                    // make request to change
                    await this.$axios.put(`/users/${this.username}/username`,
                        this.newUsername,
                        {headers: {'Authorization': sessionStorage.getItem('bearerToken'), 'Content-Type': 'text/plain'}})
                    this.username = this.newUsername
                    sessionStorage.setItem('username', this.newUsername)
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async deletePhoto(e) {
            try {
                await this.$axios.delete(`/users/${this.username}/photos/${e.target.parentElement.id}`,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.refresh()
            } catch (error) {
                this.handleError(error)
            }
        },
        showImage(e) {
            this.$router.push({name: 'viewImage', params: {uploader: this.username, photoID: e.target.parentElement.id}})
        },
        async likePhoto(e) {
            // turns out that if I don't await a request I might fail to catch incoming errors
            try {
                // elemID is uploaderName+'/'+photoID, we want them separated
                const pID = e.target.parentElement.id
                await this.$axios.put(`/users/${this.username}/photos/${pID}/likes/${this.username}`, 
                                    null, // empty body
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.myPhotos[this.username + pID].liked = true
                this.myPhotos[this.username + pID].likes++
            } catch (error) {
                this.handleError(error)
            }
        },
        async unlikePhoto(e) {
            try {
                const pID = e.target.parentElement.id
                await this.$axios.delete(`/users/${this.username}/photos/${pID}/likes/${this.username}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.myPhotos[this.username + pID].liked = false
                this.myPhotos[this.username + pID].likes--
            } catch (error) {
                this.handleError(error)
            }
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
        },
        refresh() {
            // refresh really means update showed photos
            // we're forcing a re-mount of the component
            window.location.reload()
        }
    },
    beforeMount() {
        this.getMyPhotos()
    },
}
</script>