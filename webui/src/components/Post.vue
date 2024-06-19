<script>
export default {
    data: function() {
		return {
            errormsg: null,
			loading: false,
            requestingId: localStorage.getItem("token"),
            nickname: "undefined",
            newComment: "",
            liked: false,
            numLikes: this.likes.length+0,
        }
    },
	props: ['postId','userOwner', 'date','urlPost', 'comments', 'likes'],
    methods:{
        async getNickname(){
            this.loading = true;
	        this.errormsg = null;
            try {
                let response = await this.$axios.get("/user/"+this.requestingId+"/nickname?requestedId="+this.userOwner,
                {headers: {
                        Authorization: "Bearer " + localStorage.getItem("token"),
                    }},
                );
                this.nickname = response.data.nickname;
            } catch (e) {
		        this.errormsg = e.toString();
	        }
	        this.loading = false;
        },

        async deletePost(p){
            this.loading = true;
	        this.errormsg = null;
            try {
                await this.$axios.delete("/user/"+this.requestingId+"/posts?postId="+p,
                {headers: {
                    Authorization: "Bearer " + localStorage.getItem("token"),
                }});
            } catch (e) {
		        this.errormsg = e.toString();
	        }
            this.loading = false;
            window.location.reload();
        },

        async postComment(){
            this.loading = true;
	        this.errormsg = null;
            try {
                await this.$axios.post("/user/"+this.requestingId+"/posts/comments?postId="+this.postId,
                {comment: this.newComment},
                {headers: {
                    Authorization: "Bearer " + localStorage.getItem("token"),
                }});
            } catch (e) {
		        this.errormsg = e.toString();
	        }
            this.loading = false;
            window.location.reload();
        },

        async getLike(){
            let strId = this.requestingId
            this.likes.forEach(l => {
                if (strId == (l.user_id+"")){
                    this.liked = true;
                }
            });
        },

        async deleteComment(c){
            this.loading = true;
	        this.errormsg = null;
            try {
                await this.$axios.delete("/user/"+this.requestingId+"/posts/comments?commentId="+c,
                {headers: {
                    Authorization: "Bearer " + localStorage.getItem("token"),
                }});
            } catch (e) {
		        this.errormsg = e.toString();
	        }
            this.loading = false;
            window.location.reload();
        },

        async setLike(){
            this.loading = true;
	        this.errormsg = null;
            if (this.likes == null) {this.likes = new Array()}
            if (this.liked) {
                try {
                    await this.$axios.delete("/user/"+this.requestingId+"/posts/likes?postId="+this.postId,
                    {headers: {
                        Authorization: "Bearer " + localStorage.getItem("token"),
                    }});
                } catch (e) {
		            this.errormsg = e.toString();
	            }
                this.numLikes -= 1;
                this.liked = false;
                
            } else if(!this.liked) {
                try {
                    await this.$axios.put("/user/"+this.requestingId+"/posts/likes?postId="+this.postId,
                    {headers: {
                        Authorization: "Bearer " + localStorage.getItem("token"),
                    }});
                } catch (e) {
		            this.errormsg = e.toString();
	            }
                this.numLikes += 1;
                this.liked = true;
                
            }
            this.loading = false;
        },

        async watchProfile(x, y){
			localStorage.setItem("watchingId", x);
			localStorage.setItem("watchingNick", y);
			this.$router.replace("/profile");
		},
    },
    mounted() {
		this.getNickname();
        this.getLike();
	}
}
</script>

<template>
<div class="CPwrapperPost" >
    <div class="CPnickOwner" v-if="this.requestingId!=this.userOwner" @click="watchProfile(this.userOwner,this.nickname)">
        <h3>{{ this.nickname }}</h3>
    </div>

    <hr class="CPseparatorPost" v-if="this.requestingId!=this.userOwner">
    
    <div class="CPphotoPost">
        <div class="CPcenterImage"></div>
        <img class="CPphoto" :src=this.urlPost alt="Error Image">
        <div class="CPcenterImage"></div>
    </div>

    <div class="CPwrapperStats">
        <div class="CPlikeDiv" @click="setLike"><h3><i class='bx bxs-star' style='color:#ffdd09'  v-if="this.liked"></i><i class='bx bx-star' v-if="!this.liked"></i>{{ this.numLikes }}</h3></div>
        <div class="CPcommentDiv"><h3><i class='bx bx-comment-detail'></i>{{ comments.length }}</h3></div>
        <div class="CPdeleteDiv" v-if="this.requestingId==this.userOwner"><h3><i class='bx bxs-trash' style='color:#cb1e1e' @click="deletePost(this.postId)"></i></h3></div>
    </div>
    <div class="CPwrapperComments">
        <div class="CPcommentBar">
            <input id="commentBox" type="text" placeholder="Comment the picture :)" required v-model="this.newComment"><button class="CPcommentButton" type="submit" @click="postComment"></button>
        </div>
        <ul class="CPlistComments">
            <li v-for="c in comments" v-bind:key="c.nickname">{{ c.nickname + ": " + c.comment + "  "}}<i class='bx bxs-trash' @click="deleteComment(c.comment_id)" v-if="c.user_id==this.requestingId"></i></li>
        </ul>
    </div>
    <hr class="CPseparatorPost">

</div>
</template>

<style>

.CPwrapperPost{
    width:100%;
    display: flex;
    align-items: start;
    flex-direction: column;
    
}

.CPwrapperPost .CPseparatorPost{
    width: 100%;
    height: 1px;
    margin-top: 0px;
}

.CPwrapperPost .CPnickOwner {
    text-indent: 2em;
}


.CPwrapperPost .CPphotoPost {
    width: 100%;
    display: flex;
}

.CPphotoPost .CPcenterImage{
    width: 2%;
}

.CPphotoPost .CPphoto {
    width: 96%;
    height: auto;
    display: block;
    border-radius: 40px;
    padding-bottom: 10px;
}

.CPwrapperStats{
    width: 100%;
    display: flex;
}

.CPwrapperStats .CPlikeDiv {
    margin-left: 10px;
    padding-right: 10px;
}

.CPwrapperStats .CPdeleteDiv {
    flex: 1;
    text-align: right;
    padding-right: 10px;
}

.CPwrapperComments {
    width: 100%;
    display: flex;
    flex-direction: column;
    margin-bottom: 10px;
}

.CPwrapperComments .CPcommentBar{
    width: 100%;
    display: flex;
    flex-direction: row;
}

.CPcommentBar .CPcommentButton{
    width:27px;
	height: 27px;
    margin-right: 10px;
    background-color: black;
	background-image: url('./src/rocket.png');
	background-size:80%;
	background-position: center;
	border: none;
	border-radius: 50%;
}

.CPcommentBar input {
    flex: 1;
    text-indent: 1em;
    margin-left: 10px;
	font-size: 15px;
	height: 30px;
	outline: none;
	border: none;
	border: 1px solid rgba(0, 0, 0, 0.5);
	border-radius: 40px;
	margin-right: 1%;
}

.CPwrapperComments .CPlistComments{
    width: 80%;
	text-align: left;
    list-style-type:square;
}






</style>