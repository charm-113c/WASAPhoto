<template>

    <Sidebar></Sidebar>

    <div class="container">
        <div class="profile">

            <h4 class="top-box">
                Welcome to your profile, {{ username }}. Here, you can do stuff. <br> Stuff includes:
            </h4>
            <ul>
                <li>See your follows:</li>
                <li><details><summary>Followers: {{ nFollowers }}</summary>{{ followers }}</details></li>
                <li><details><summary>Following: {{ nFollowing }}</summary> {{ following }}</details></li>
                <li>
                    Change your username: <br>
                    <input class="new-name" type="text" v-model="newUsername" placeholder="New username" @keydown="changeUsername">
                </li>
            </ul>
            <p></p>
            You currently have {{ Object.keys(myPhotos).length }} photo(s).
            Feel free to upload some, or delete uploaded ones: 
            <form>
                <label for="images">Select a photo to upload</label> <br>
                <input type="file" id="images" accept="image/png, image/jpeg" @change="selectImage" required/>
                <br>
                <label for="description">Description (optional):</label> <br>
                <input class="desc" type="text" id="description" placeholder="Enter description" v-model="description">
                <div class="submit">
                    <button @click.prevent="uploadPhotos">Upload</button>
                </div>
            </form>
            <hr>
            <div class="profile-photos" v-for="photo in this.myPhotos" :id="photo.id">
                <img :src="photo.src" :alt="photo.description" @click="showImage">
                <!-- <button class="del-btn" @click="deletePhoto">Delete Photo</button> <br> -->
                <svg class="del-btn" :id="photo.id + '.'" @click="showDialog">
                    <use href="/feather-sprite-v4.29.0.svg#trash-2"/>
                </svg>
                <span class="metadata">{{ photo.uploadDate }}</span>
                <svg class="like" :class="{'not-liked': !photo.liked, 'liked': photo.liked}" :id="photo.id + '_'" @click="likePhoto($event, photo.liked)">
                    <use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
                </svg>
                <span class="nlikes">{{ photo.likes }}</span>
                <!-- <button class="photo-btn" v-if="!photo.liked" @click="likePhoto">Like</button>
                <button class="photo-btn" v-else @click="unlikePhoto">Unlike</button> -->
                <svg class="comment" :class="{'commented': photo.comments > 0}" :id="photo.id + '-'" @click="showImage">
                    <use href="/feather-sprite-v4.29.0.svg#message-circle"/>
                </svg>
                <span class="ncomments">{{ photo.comments }}</span>
                <p class="photo-desc">{{ photo.description }}</p>
                <hr class="line">
            </div>

        </div>
    </div>

    <div v-if="showModal">
        <!-- the modal emits delete or cancel event upon click -->
        <PopUpMsg :MsgText="modalTxt" @close="closeModal" @ok="deletePhoto"/>
    </div>

    <ErrorMsg v-if="showErrMsg" :msg="errorMsg" @close="closeErrMsg"></ErrorMsg>
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
            followers: '', // list of names
            nFollowers: 0,
            following: '',
            nFollowing: 0,
            // and we'll allow them to post photos
            selectedFile: null,
            description: '',
            showErrMsg: false,
            errorMsg: '',
            // or change username
            newUsername: '',
            showModal: false,
            modalTxt: 'Are you sure you want to delete your photo?',
            targetPhotoID: '',
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
                if (res.data.Followers !== null) {
                    this.followers = res.data.Followers.join(', ')
                    this.nFollowers = res.data.Followers.length
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
                // as usual, photos need special treatment
                if (res.data.Photos === null) {
                        this.noPhotos = true
                } else {
                    res.data.Photos.forEach(photo => {
                        // transform img data into data URI so as to feed it to img tag
                        let imgSRC = `data:${photo.FileExtension};base64,${photo.BinaryData}`
                        this.myPhotos[photo.Uploader + photo.PhotoID] = {
                            'src': imgSRC, 'uploader': this.username, 'description': photo.Description, 'uploadDate': new Date(photo.UploadDate).toUTCString(), 
                            'likes': photo.Likes, 'comments': photo.Comments, 'id': photo.PhotoID, 'liked': false
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
        async likePhoto(e, liked) {
            try {
                // get photoID
                let pID = e.target.parentElement.id.replace('_', '')
                if (!liked) {
                    await this.$axios.put(`/users/${this.username}/photos/${pID}/likes/${this.username}`,
                                            null,
                                            {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.myPhotos[this.username + pID].liked = true
                    this.myPhotos[this.username + pID].likes++
                } else {
                    await this.$axios.delete(`/users/${this.username}/photos/${pID}/likes/${this.username}`,
                                            {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.myPhotos[this.username + pID].liked = false
                    this.myPhotos[this.username + pID].likes--
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async changeUsername(e) {
            try {
                if (this.newUsername && e.key === 'Enter') {
                    // trim white spaces
                    this.newUsername = this.newUsername.trim()
                    // make request to change
                    await this.$axios.put(`/users/${this.username}/username`,
                        this.newUsername,
                        {headers: {'Authorization': sessionStorage.getItem('bearerToken'), 'Content-Type': 'text/plain'}})
                    this.username = this.newUsername
                    sessionStorage.setItem('username', this.newUsername)
                    this.refresh()
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        showDialog(e) {
            this.showModal = true
            this.targetPhotoID = e.target.parentElement.id.replace('.', '')
        },
        closeModal() {
            this.showModal = false
        },
        async deletePhoto() {
            try {
                await this.$axios.delete(`/users/${this.username}/photos/${this.targetPhotoID}`,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.refresh()
            } catch (error) {
                this.handleError(error)
            }
        },
        showImage(e) {
            let pID = e.target.parentElement.id.replace('-', '')
            this.$router.push({name: 'viewImage', params: {uploader: this.username, photoID: pID}})
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

<style scoped>

.container {
    margin-right: 0px;
}

.profile {
    display: flex;
    flex-direction: column;
    align-items: left;
    margin-left: 40px;
    width: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    padding-right: 20px;
    padding-left: 20px;
}

.top-box {
    align-self: center;
    width: 70%;
    background-color: var(--strong-colour);
    text-align: center;
    padding: 10px 20px 10px 20px;
    border-radius: 50px;
    margin-top: 10px;
}

ul {
    list-style-type: none;
}

input {
    height: 30px;
}

input::file-selector-button {
    background-color: transparent;
    color: var(--bold-colour);
    border-color: var(--bold-colour);
    width: 100px;
    border-radius: 15px;
}
input::file-selector-button:hover {
    cursor: pointer;
} 

.new-name {
    background-color: var(--light-colour);
    color: var(--strong-colour);
    border-radius: 20px;
    border: 0px;
    padding-left: 15px;
    width: 40%;
}

.desc {
    border: 0px;
    background-color: var(--light-colour);
    color: var(--strong-colour);
    width: 50%;
    border-radius: 20px;
    padding-left: 15px;
}

.submit button{
    margin-top: 5px;
    background-color: transparent;
    color: var(--bold-colour);
    border-color: var(--bold-colour);
    width: 100px;
    border-radius: 15px;
}

.profile-photos {
    text-align: center;
    display: grid;
    width: 100%;
    gap: 3px;
    grid-template-areas: 
        'photo photo photo photo del-button metadata'
        'photo photo photo photo empty description'
        'like nlikes comment ncomments empty description'
        'line line line line line line';
    grid-template-rows: 40px auto 35px 15px;
    grid-template-columns: 25% 5% 6% 24% 40px auto;
}
.profile-photos img {
    grid-area: photo;
    max-width: 100%;
    max-height: 100%;
    justify-self: center;
}
.del-btn {
    grid-area: del-button;
    stroke: var(--bold-colour);
	stroke-width: 1;
	stroke-linecap: round;
	stroke-linejoin: round;
    width: 30px;
	height: 30px;
}
.metadata {
    grid-area: metadata;
    align-self: center;
}
svg {
    stroke: var(--bold-colour);
	stroke-width: 1;
	stroke-linecap: round;
	stroke-linejoin: round;
    width: 28px;
	height: 30px;
}
.like {
    grid-area: like;
    justify-self: right;
}
.liked {
    fill: var(--strong-colour);
}
.nlikes {
    grid-area: nlikes;
    justify-self: left;
    align-self: center;
}
.comment {
    grid-area: comment;
    justify-self: right;
}
.commented {
    fill: var(--strong-colour);
}
.ncomments {
    grid-row: ncomments;
    justify-self: left;
    align-self: center;
}
.photo-desc {
    grid-area: description;
}
.line {
    grid-area: line;
    align-self: center;
    width: 61%;
    margin-top: 5px;
}
svg:hover, img:hover {
    cursor: pointer;
}
</style>