<template>
    <!-- Note: we're assuming photos to be non-empty. It's the upstream view that needs to check that before calling component -->
    <div class="photo-container" v-for="photo in vphotos" :key="photo.uploader + photo.id" :id="photo.uploader+'/'+photo.id">
        <span class="header">{{ photo.uploader }} | {{ photo.uploadDate }} </span>
        <img class="photo" :src="photo.src" :alt="photo.description" @click="showImage"/> 

        <svg class="toggle-button" :class="{'like-button': !photo.liked, 'unlike-button': photo.liked}" :id="photo.uploader+'|'+photo.id" @click="likePhoto($event, photo.liked)">
            <use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
        </svg>
        <span class="like-number" :class="{'liked-number': photo.liked}">{{ photo.likes }}</span>
        <svg class="comment-button" :class="{'unlike-button': photo.comments>0, 'like-button': photo.comments===0}" @click="showImage" :id="photo.uploader+'+'+photo.id">
            <use href="/feather-sprite-v4.29.0.svg#message-circle"/>
        </svg>
        <span class="comment-number">{{ photo.comments }}</span>

        <p class="description">{{ photo.description }}</p>
        <hr class="line">
    </div>
    <ErrorMsg v-if="showError" :msg="errorTxt" @close="closeErrMsg"></ErrorMsg>
</template>

<script>
export default {
    props: ['photos'],
    data() {
        return {
            username: sessionStorage.getItem("username"),
            vphotos: this.photos,
            // props can be used as data, but cannot be modified
            // so to modify data, we'll only use it as initial value, which we are free to play with
            showError: false,
            errorTxt: '',
        }
    },
    methods: {
        async likePhoto(event, liked) {
            // this handles both liking and unliking a photo
            try {
                let uploaderAndID = event.target.parentElement.id.split(/[|/]/)
                // elemID is uploaderName+'/'+photoID, we want them separated
                // Why the strange split? It's black magic, because svg's behaviour on click is near unpredictable. 
                // For one click the target is the svg, for another it's use. If you click a little outside 
                // all hell breaks loose and nothing works. It'd take me a paragraph to explain why this works, 
                // so chalk it up to black magic.
                if (!liked) {
                await this.$axios.put(`/users/${uploaderAndID[0]}/photos/${uploaderAndID[1]}/likes/${this.username}`, 
                                    null, // empty body
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.vphotos[uploaderAndID[0]+uploaderAndID[1]].liked = true
                this.vphotos[uploaderAndID[0]+uploaderAndID[1]].likes++
                } else {
                    await this.$axios.delete(`/users/${uploaderAndID[0]}/photos/${uploaderAndID[1]}/likes/${this.username}`,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.vphotos[uploaderAndID[0]+uploaderAndID[1]].liked = false
                    this.vphotos[uploaderAndID[0]+uploaderAndID[1]].likes--
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        showImage(e) {
            let uper = e.target.parentElement.id.split(/[+/]/)[0]
            let pID = e.target.parentElement.id.split(/[+/]/)[1]
            this.$router.push({name: 'viewImage', params: {uploader: uper, photoID: pID}})
        },
        handleError(error) {
            // How do we handle errors? By displaying a closable modal. A generic one saying there was an error 
            // loading some images.
            this.showError = true 
        },
        closeErrMsg() {
            this.showError = false
        }
    },
}
</script>