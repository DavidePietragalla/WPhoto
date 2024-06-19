<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,

			nickname: localStorage.getItem("nickname"),
			userId: localStorage.getItem("token"),

			url: "",
			followers: new Array(),
			followings: new Array(),
			posts: new Array(),
			
			stream: new Array(),

			researchNick: "",
			researchUsers: null,
			userFound: true,
		}
	},
	methods: {

		async getProfile(){
            try {
                let response = await this.$axios.get(
                    "/user/"+this.userId+"?nickname="+this.nickname,
                    {headers: {
						Authorization: "Bearer " + localStorage.getItem("token")
					}})
                if (response.data.followers != null) {
                    this.followers = response.data.followers;
                }
                if (response.data.following != null) {
                    this.followings = response.data.following;
                }
                if (response.data.posts != null) {
                    this.posts = response.data.posts;
                }
            } catch (e) {
				this.errormsg = e.toString();
			}
        },

		async uploadPhoto(){
			if (this.url == "") return;
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/user/"+this.userId+"/posts",
				{urlImage: this.url},
				{headers: {
						Authorization: "Bearer " + localStorage.getItem("token"),
					}},
				);
				this.url = response.data.urlImage
			} catch (e) {
				this.errormsg = e.toString();
			}
			let urlBox = document.getElementById("urlBox");
			urlBox.value = "";
            this.loading = false;
		},

		async getStream(){
		this.loading = true;
		this.errormsg = null;
		try {
			let response = await this.$axios.get("/user/"+this.userId+"/stream",
			{headers: {
					Authorization: "Bearer " + localStorage.getItem("token"),
				}},
			);
			this.stream = response.data;
			for (let p of this.stream){
				if (p.comments == null)
					p.comments = new Array();
				if (p.likes == null)
					p.likes = new Array();
			}
		} catch (e) {
			this.errormsg = e.toString();
		}
		this.loading = false;
		},

		async search(){
			if (this.researchNick == "") return;
			this.loading = true;
			this.errormsg = null;
			this.userFound = true;
			try {
				let response = await this.$axios.get("/user?nickname="+this.researchNick,
				{headers: {
						Authorization: "Bearer " + localStorage.getItem("token"),
					}},
				);
				this.researchUsers = response.data
				if (0 >= this.researchUsers.length) {
					this.userFound = false;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			let searchBox = document.getElementById("searchBox");
			searchBox.value = ""
			this.loading = false
		},

		async watchProfile(x, y){
			localStorage.setItem("watchingId", x);
			localStorage.setItem("watchingNick", y);
			this.$router.replace("/profile");
		}

	},
	mounted() {
		this.getProfile();
		this.getStream();
	}
}
</script>

<template>
	<div class="HtopBar">
		<h1>WASAPHOTO</h1>
	</div>

	<div class="Hspace">
		<div class="HtableSpace">
			<div class="HhSpace"></div>

			<div class="Hprofile">
				<ul class="HinfoProfile">
					<li @click="watchProfile(this.userId, this.nickname)"><h3>{{this.nickname}}</h3></li>
					<li>{{"Followers: "+this.followers.length}}</li>
					<li>{{"Followings: "+this.followings.length}}</li>
					<li>{{"Posts: "+this.posts.length}}</li>
				</ul>
			</div>

			<div class="HhSpace"></div>

			<div class="Hstream">
				<ul class="HlistPost">
					<li v-for="p in stream" v-bind:key="p.post_id"><Post :postId=p.post_id :userOwner=p.owner :date=p.date :urlPost="'data:image/*;base64,'+p.urlImage" :comments=p.comments :likes=p.likes></Post></li>
				</ul>
			</div>

			<div class="HhSpace"></div>

			<div class="Hsearch">
				<div class="HsearchBar">
				<input id="searchBox" type="text" placeholder="Search" required v-model="researchNick"><button class="HsearchButton" type="submit" @click="search"></button>
			    </div>
				<ul class="HlistUsers">
					<li v-if="this.userFound==false"><div>Utente non trovato :(</div></li>
					<li v-for="u in researchUsers" v-bind:key="u.user_id" @click="watchProfile(u.user_id, u.nickname)">{{ u.nickname }}</li>
				</ul>
			</div>

			<div class="HhSpace"></div>
		</div>
	</div>
</template>

<style>
.HtopBar {
    width: 100%;
    color: #FFF;
    background-image: linear-gradient(to bottom right, #0F006E, #B5068F, #CD6700, #FEF9D5);
    overflow: hidden;
    background-size: 250% 800%;
    animation: waveanimation 5s ease infinite;
    text-indent: 1em;
	padding-top: 10px;
}

@keyframes waveanimation {
    0%{
        background-position: 0% 50%;
    }
    50%{
        background-position: 100% 50%;
    }
    100%{
        background-position: 0% 50%;
    }
}

.Hspace {
	display: flex;
	justify-content: flex-start;
    margin: 0;
	padding-top: 10px;
    background-image: url('./src/spazio1.png');
	min-height: 100vh;
	height:fit-content;
}

.HhSpace {
    width: 2%;
    background-color: none;
}

.Hspace .HtableSpace {
	display: flex;
    width: 100%;
    background-color: none;
}

.Hprofile {
	padding-top: 10px;
	font-size: 15px;
	height:fit-content;
	justify-content: flex-start;
    width: 24%;
    background-color: #FFF;
    border-radius: 40px;
    text-indent: 0em;
}

.Hprofile .HinfoProfile {
	list-style-type: none;
}

.HinfoProfile .HurlImage{
	display: flex;
}

.HurlImage .bx{
	font-size: 2em;
}

.HurlImage .HuploadPhoto{
	background: #FFF;
	border: none;
}

.HurlImage input{
	width: 90%;
	text-indent: 0.5em;
	outline: none;
	border: none;
	border: 1px solid rgba(0, 0, 0, 2);
	border-radius: 40px;
}

.Hstream {
	padding-top: 10px;
	height:fit-content;
    width: 44%;
    background-color: #FFF;
    border-radius: 40px;
    text-indent: 0em;
	margin-bottom: 10px;
}

.Hstream .HlistPost {
	margin: 0;
	padding:0;
	list-style-type: none;
}

.Hsearch {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	text-align: center;
	padding-top: 20px;
	height:fit-content;
    width: 24%;
	font-size: 15px;
    background-color: #FFF;
    border-radius: 40px;
    text-indent: 0em;
}

.Hsearch .HsearchBar{
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
}

.Hsearch .HsearchButton{
	width:27px;
	height: 27px; 
	margin-right: 10px;
	background-image: url('./src/search.png');
	background-size: cover;
	background-position: center;
	border: none;
	border-radius: 50%;
}

.Hsearch .HlistUsers{
	top:40px;
	width: 90%;
	text-align: left;
	list-style-type:square;
	
}

.Hsearch input {
	margin-left: 10px;
	font-size: 15px;
	height: 30px;
	width: 90%;
	text-indent: 0.5em;
	outline: none;
	border: none;
	border: 1px solid rgba(0, 0, 0, 0.5);
	border-radius: 40px;
	margin-right: 4px;
}

</style>