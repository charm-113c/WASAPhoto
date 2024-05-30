<template>
    <div class="container">
        <div class="image-grid">
            <span class="img-header">{{ uploader }} | {{ uploadDate }}</span>
            <img class="img" :src="imgSRC" :alt="description"> 
            <svg class="toggle-button img-like-button" :class="{'like-button': !liked, 'unlike-button': liked}" @click="likePhoto($event,liked)">
                    <use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
            </svg>
            <span class="img-likes">{{ likes }}</span>
            <svg class="comment-button img-comment-button" :class="{'unlike-button':nComments>0, 'like-button':nComments===0}">
                    <use href="/feather-sprite-v4.29.0.svg#message-circle"/>
            </svg>
            <span class="img-comments">{{ nComments }}</span> 
            <p class="img-description">{{ description }}</p> 
        </div>

        <hr>
        
        <div class="comment-section">
            <input type="text" v-model="newComment" @keydown="postComment" placeholder="Write a comment"/>
        
            <div class="comment-container" v-for="comment in comments" :key="comment.name + comment.id" :id="comment.name +'/'+ comment.id">
                <span class="commenter">{{ comment.name }} on {{ comment.date }}</span>
                <svg class="del-button" v-if="comment.isMyComm" :id="comment.name +'|'+ comment.id" @click="showDialog">
                    <use href="/feather-sprite-v4.29.0.svg#trash-2"/>
                </svg>
                <p class="comment-text">{{ comment.text }}</p>
            </div>    
        </div>
    </div>
    <div v-if="showModal">
        <!-- the modal emits delete or cancel event upon click -->
        <PopUpMsg :MsgText="modalTxt" @close="closeModal" @ok="deleteComment"/>
    </div>
    <ErrorMsg v-if="showErrMsg" :msg="errorMsg" @close="closeErrMsg"></ErrorMsg>
</template>

<script>
export default {
    props: ['uploader', 'photoID'],
    data() {
        return {
            username: sessionStorage.getItem('username'),
            showErrMsg: false,
            errorMsg: '',
            imgSRC: null,
            likes: 0,
            liked: false,
            uploadDate: null,
            description: '',
            comments: {},
            nComments: 0,
            // allow user to post a comment
            newComment: "",
            showModal: false,
            modalTxt: "Are you sure you want to delete your comment?",
            targetCommentID: null, // remember comment targeted by user for deletion 
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
                this.uploadDate = new Date(res.data.UploadDate).toUTCString()
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
                            'text': comment.CommentText, 'date': new Date(comment.CommentDate).toUTCString(), 
                            'name': comment.CommenterName, 'id': comment.CommentID, 'isMyComm': byCurrUser
                        }
                    })
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
        async deleteComment() {
            const authorAndID = this.targetCommentID
            // note: author must === current user for the delete button to be visible in the first place
            try {
                await this.$axios.delete(`/users/${this.uploader}/photos/${this.photoID}/comments/${authorAndID[1]}`,
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken'),
                                                'requesting-user': authorAndID[0]}})
                this.refresh()
            } catch (error) {
                this.handleError(error)
            }
        },
        async likePhoto(event, liked) {
            // this handles both liking and unliking a photo
            try {
                if (!liked) {
                await this.$axios.put(`/users/${this.uploader}/photos/${this.photoID}/likes/${this.username}`, 
                                    null, // empty body
                                    {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                this.liked = true
                this.likes++
                } else {
                    await this.$axios.delete(`/users/${this.uploader}/photos/${this.photoID}/likes/${this.username}`,
                                        {headers: {'Authorization': sessionStorage.getItem('bearerToken')}})
                    this.liked = false
                    this.likes--
                }
            } catch (error) {
                this.handleError(error)
            }
        },
        showDialog(e) {
            this.showModal = true
            this.targetCommentID = e.target.parentElement.id.split(/[|/]/)
        },
        closeModal() {
            this.showModal = false
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
            location.reload()
        }
    },
    beforeMount() {
        this.getPhoto()
    }
}
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    backdrop-filter: brightness(40%);
}

.comment-button:hover {
    cursor: default;
}

hr {
    width: 90%;
    background-color: var(--strong-colour);
    margin-left: 50px;
}
</style>