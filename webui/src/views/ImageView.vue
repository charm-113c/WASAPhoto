<template>
    For likes: the server returns a list of people who like the photo. `liked` will then be true if the user is in the list<br>
    <img :src="imgSRC" :alt="description"> <br>
    <span>{{ uploader }} - {{ uploadDate }} - Likes: {{ likes }} - Comments: {{ nComments }}</span> <br>
    <button v-if="!liked" @click="likePhoto">Like</button>
    <button v-else @click="unlikePhoto">Unlike</button>
    <p>{{ description }}</p> <br>
    <input type="text" v-model="newComment" @keydown="postComment"/>
    
    <div v-for="comment in comments" :id="comment.name +'/'+ comment.id">
        <span>{{ comment.date }} - by {{ comment.name }}</span> 
        <button v-if="comment.isMyComm" @click="deleteComment">Delete</button>
        <p>{{ comment.text }}</p>
    </div>

    <ErrorMsg v-if="errorMsg" :msg="errorMsg" :code="errorCode"></ErrorMsg>
</template>

<script>
export default {
    props: ['uploader', 'photoID'],
    data() {
        return {
            username: sessionStorage.getItem('username'),
            errorCode: null,
            errorMsg: '',
            imgSRC: null,
            likes: 0,
            liked: false,
            uploadDate: null,
            description: '',
            comments: {},
            nComments: 0,
            // allow user to post a comment
            newComment: 'Got a comment to make?',
        }
    },
    methods: {
        async getPhoto() {
            try {
                const res = await this.$axios.get(`/users/${this.uploader}/photos/${this.photoID}`,
                                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken'),
                                                                'requesting-user': this.username},}) 
                this.imgSRC = `data:${res.data.FileExtension};base64,${res.data.BinaryData}`
                this.likes = res.data.Likes 
                this.uploadDate = new Date(res.data.UploadDate) 
                this.description = res.data.Description
                if (res.data.Likers) {
                    this.liked = res.data.Likers.includes(this.username)
                }
                // comments have a lot of information to them, if they exist
                if (res.data.Comments !== null) {
                    this.nComments = res.data.Comments.length
                    res.data.Comments.forEach(comment => {
                        let byCurrUser = comment.CommenterName === this.username
                        this.comments[comment.CommenterName + comment.CommentID] = {
                            'text': comment.CommentText, 'date': new Date(comment.CommentDate), 
                            'name': comment.CommenterName, 'id': comment.CommentID, 'isMyComm': byCurrUser
                        }
                    });
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        async postComment(e) {
            if (e.key === 'Enter' && this.newComment) {
                try {
                    // post the comment
                    await this.$axios.post(`/users/${this.uploader}/photos/${this.photoID}/comments`, this.newComment,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken'),
                                                    'Content-Type': 'text/plain', 
                                                    'commenter-username': this.username,
                                                    'upload-date': new Date().toISOString()}})
                    this.refresh()
                } catch (error) {
                    this.handleError(error)
                }
            }
        },
        async deleteComment(e) {
            const authorAndID = e.target.parentElement.id.split('/')
            // note: author must === current user for the delete button to be visible in the first place
            console.log(authorAndID);
            try {
                await this.$axios.delete(`/users/${this.uploader}/photos/${this.photoID}/comments/${authorAndID[1]}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken'),
                                                'requesting-user': authorAndID[0]}})
                this.refresh()
            } catch (error) {
                this.handleError(error)
            }
        },
        async likePhoto() {
            try {
                await this.$axios.put(`/users/${this.uploader}/photos/${this.photoID}/likes/${this.username}`,
                                        null,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.liked = true
                this.likes++
            } catch (error) {
                this.handleError(error)
            }
        },
        async unlikePhoto() {
            try {
                await this.$axios.delete(`/users/${this.uploader}/photos/${this.photoID}/likes/${this.username}`,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.liked = false
                this.likes--
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
            location.reload()
        }
    },
    beforeMount() {
        this.getPhoto()
    }
}
</script>