<template>
    
    <div class="top-bar">
        <div class="logo">
            <h2>WASAPhoto</h2>
            <p>Keep in touch with your friends by sharing photos of special moments, thanks to WASAPhoto!</p>
        </div>
    </div>
    
    <ErrorMsg v-if="errorMsg" :msg="errorMsg" :code="errorCode"></ErrorMsg>

    <div class="container">
        <div class="sidebar-left">
            <h2>Hello, {{ username }}</h2>
            <hr>
            <input type="text" v-model="user2" @keyup="searchProfile" placeholder="Search a user">
            <ul>
                <li>This is the first element</li>
                <li>This sidebar will only be visible here</li>
                <li>Or, we can paste it back everywhere</li>
                <li>But we deal with that later</li>
                <li>
                    <RouterLink :to="{name: 'myProfile', params: {username: this.username}}">
                        <!-- Router link to own profile -->
                        My profile
                    </RouterLink>
                </li>
            </ul>
        </div>

        <!-- Display the stream only if there aren't any errors -->
        <div class="stream">
            <div class="photo-container" v-for="photo in streamPhotos" v-if="!emptyStream" :id="photo.uploader+'/'+photo.id">
                <span class="header">{{ photo.uploader }} | {{ photo.uploadDate }} </span>
                <img class="photo" :src="photo.src" :alt="photo.description" @click="showImage"/> 
                <!-- <button class="like-" v-if="!photo.liked" @click="likePhoto">Like</button> -->
                <!-- <button class="like-" v-if="photo.liked" @click="unlikePhoto">Unlike</button> -->
                <svg class="toggle-button" :class="{'like-button': !photo.liked, 'unlike-button': photo.liked}" :id="photo.uploader+'|'+photo.id" @click="likePhoto($event, photo.liked)">
                    <use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
                </svg>
            
                <!-- <svg class="unlike-button" id="unlikeButton" v-else @click="unlikePhoto">
                    <use href="/feather-sprite-v4.29.0.svg#thumbs-up" @click="dispatchClick"/>
                </svg> -->

                <span class="like-number" :class="{'liked-number': photo.liked}">{{ photo.likes }}</span>
                <!-- <button class="comment-button">Comment</button> -->
                <svg class="comment-button">
                    <use href="/feather-sprite-v4.29.0.svg#message-circle" @click="dispatchClick"/>
                </svg>
                <span class="comment-number">{{ photo.comments }}</span>
                <p class="description">{{ photo.description }}</p> 
                <hr class="line">
            </div>
            <p v-if="emptyStream">Looks like there is nothing to show yet. Follow other users and see if this changes!</p>
        </div>
        
    </div>

</template>


<script>
import { RouterLink } from 'vue-router'

export default {
    data() {
        return {
            username: sessionStorage.getItem('username'),
            // show the stream
            streamPhotos: {},
            // streamPhotos[Uploader + PhotoID] = 
            //      { 'src': str, 'uploader': str, 'description': str, 'uploadDate': date, 
            //      'likes': int, 'comments': int, 'liked': bool, 'id': int}
            errorCode: null,
            errorMsg: null,
            user2: '',
            emptyStream: false,
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
                } catch (error) {
                    this.handleError(error)
                }
            }
        },
        async getStream() {
            try {
                // make the request
                const res = await this.$axios.get(`/users/${this.username}/stream`,
                    { headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                // and handle the JSON data (photos + metadata)
                res.data.forEach(photo => {
                    // transform the JSON'ed imgFile into a data URI by encoding it to base 64
                    // effectively transforming it into a file again and allowing it to be used 
                    // as src for img tags 
                    let imgSRC = `data:${photo.FileExtension};base64,${photo.BinaryData}`
                    this.streamPhotos[photo.Uploader + photo.PhotoID] = {
                        'src': imgSRC, 'uploader': photo.Uploader, 'description': photo.Description, 'uploadDate': new Date(photo.UploadDate).toUTCString(), 
                        'likes': photo.Likes, 'comments': photo.Comments, 'liked': false, 'id':photo.PhotoID
                    }
                    // set 'liked' to true if user has liked the photo
                    if (photo.Likers) {
                        this.streamPhotos[photo.Uploader + photo.PhotoID].liked = photo.Likers.includes(this.username)
                    }
                })
            } catch (error) {
                if (error.message === 'res.data is null') {
                    this.emptyStream = true
                }
                this.handleError(error)
            }
        },
        async likePhoto(event, liked) {
            // this handles both liking and unliking a photo
            try {
                if (!liked) {
                // elemID is uploaderName+'/'+photoID, we want them separated
                let uploaderAndID = event.target.parentElement.id.split(/[|/]/)
                // Why the strange split? It's black magic, because svg's behaviour on click is near unpredictable. 
                // For one click the target is the svg, for another it's use. If you click a little outside 
                // all hell breaks loose and nothing works. It'd take me a paragraph to explain why this works, 
                // so chalk it up to black magic.
                await this.$axios.put(`/users/${uploaderAndID[0]}/photos/${uploaderAndID[1]}/likes/${this.username}`, 
                                    null, // empty body
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.streamPhotos[uploaderAndID[0]+uploaderAndID[1]].liked = true
                this.streamPhotos[uploaderAndID[0]+uploaderAndID[1]].likes++
                } else {
                    let uploaderAndID = event.target.parentElement.id.split(/[|/]/)
                    await this.$axios.delete(`/users/${uploaderAndID[0]}/photos/${uploaderAndID[1]}/likes/${this.username}`,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.streamPhotos[uploaderAndID[0]+uploaderAndID[1]].liked = false
                    this.streamPhotos[uploaderAndID[0]+uploaderAndID[1]].likes--
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        dispatchClick(id) {
            console.log(id)
        },
        // async unlikePhoto(e) {
        //     console.log(e);
        //     try {
        //         let uploaderAndID = e.target.parentElement.id.split('/')
        //         console.log(uploaderAndID)
        //         await this.$axios.delete(`/users/${uploaderAndID[0]}/photos/${uploaderAndID[1]}/likes/${this.username}`,
        //                             {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
        //         this.streamPhotos[uploaderAndID[0]+uploaderAndID[1]].liked = false
        //         this.streamPhotos[uploaderAndID[0]+uploaderAndID[1]].likes--
        //     } catch (error) {
        //         this.handleError(error)
        //     }
        // },
        showImage(e) {
            let uploaderAndID = e.target.parentElement.id.split('/')
            this.$router.push({name: 'viewImage', params: {uploader: uploaderAndID[0], photoID: uploaderAndID[1]}})
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
                }
        },
    },
    async beforeMount() {
        // get the photos before loading the page
        await this.getStream()
    },
}
</script>