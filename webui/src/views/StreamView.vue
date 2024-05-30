<template>
    
    <div class="top-bar">
        <div class="logo">
            <h2>WASAPhoto</h2>
            <p>Keep in touch with your friends by sharing photos of special moments, thanks to WASAPhoto!</p>
        </div>
    </div>
    
    <div class="container">
        <Sidebar></Sidebar>
        <div class="stream">
            <PhotoContainer :photos="streamPhotos"></PhotoContainer>
            <p v-if="emptyStream"> 
                <br> Looks like there's nothing to show yet. Follow others and see if this changes!
            </p>
        </div>      
    </div>

    <ErrorMsg v-if="showErrMsg" :msg="errorMsg" @close="closeErrMsg"></ErrorMsg>
</template>

<script>
export default {
    data() {
        return {
            username: sessionStorage.getItem('username'),
            // show the stream
            streamPhotos: {},
            // streamPhotos[Uploader + PhotoID] = 
            //      { 'src': str, 'uploader': str, 'description': str, 'uploadDate': date, 
            //      'likes': int, 'comments': int, 'liked': bool, 'id': int}
            showErrMsg: false,
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
                if (res.data === null) {
                    this.emptyStream = true
                } else {
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
                }
                
            } catch (error) {
                if (error.message === 'res.data is null') {
                    this.emptyStream = true
                }
                this.handleError(error)
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
    async beforeMount() {
        // get the photos before loading the page
        await this.getStream()
    },
}
</script>