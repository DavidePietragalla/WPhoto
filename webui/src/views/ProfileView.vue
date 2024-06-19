<script>
export default {
	data: function() {
		return {
            errormsg: null,
            loading: false,

            nickname: localStorage.getItem("nickname"),
	        userId: localStorage.getItem("token"),
            watchingId: localStorage.getItem("watchingId"),
            watchingNick: localStorage.getItem("watchingNick"),

            url: null,
            newNick: "",
			followed: false,
            followers: new Array(),
			followings: new Array(),
			posts: new Array(),

            researchNick: "",
			researchUsers: null,
			userFound: true,
        }
    },
    methods: {
		async ban(){
			this.loading = true;
			try {
				await this.$axios.put("/user/"+this.userId+"/banned_users?nickname="+this.watchingNick,
				{headers: {
					Authorization: "Bearer " + localStorage.getItem("token"),
				}});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading=false;
			location.reload();
		},

		async unban(){
			this.loading = true;
			try {
				await this.$axios.delete("/user/"+this.userId+"/banned_users?nickname="+this.watchingNick,
				{headers: {
					Authorization: "Bearer " + localStorage.getItem("token"),
				}});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading=false;
			location.reload();
		},

        async follow(){
            this.loading = true;
	        this.errormsg = null;
			if(!this.followed){
				try {
					await this.$axios.put("/user/"+this.userId+"/followers?nickname="+this.watchingNick,
					{headers: {
						Authorization: "Bearer " + localStorage.getItem("token"),
					}},);
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.followed = true;
			}
			else if(this.followed){
				try {
					await this.$axios.delete("/user/"+this.userId+"/followers?nickname="+this.watchingNick,
					{headers: {
						Authorization: "Bearer " + localStorage.getItem("token"),
					}});
				} catch (e) {
					this.errormsg = e.toString();
				}
				
				this.followed = false;
			}
			this.loading = false;
			location.reload();
			},

        async uploadPhoto(){
			this.loading = true;
			this.errormsg = null;
			let fileInput = document.getElementById('urlBox');
			const file = fileInput.files[0];
			const fileReader = new FileReader();

			fileReader.readAsArrayBuffer(file);

			fileReader.onload = async () => {
				try {
					await this.$axios.post("/user/" + this.userId + "/posts",
						fileReader.result,
						{headers: {
							Authorization: "Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.loading = false;
					location.reload();
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
			//-------
		},

        async getProfile(){
            try {
                let response = await this.$axios.get(
                    "/user/"+this.userId+"?nickname="+this.watchingNick,
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
                    for (let p of this.posts){
				        if (p.comments == null)
					        p.comments = new Array();
				        if (p.likes == null)
					        p.likes = new Array();
			        }
                }
            } catch (e) {
                this.errormsg = e.toString();
            }
			this.getFollow();
        },

        async changeUsername(){
			if (this.newNick == this.nickname) {
				let inputB = document.getElementById("newNickBox");
				inputB.value = "Gia in uso...";
				return
			}
            if (this.newNick == "") return;
			this.loading = true;
			this.errormsg = null;
            try {
				let response = await this.$axios.get("/user?nickname="+this.newNick,
				{headers: {
						Authorization: "Bearer " + localStorage.getItem("token"),
					}},
				);
				if (0 < response.data.length) {
					let inputB = document.getElementById("newNickBox");
					inputB.value = "Gia in uso...";
					return
				}
			    await this.$axios.put("/user/" + this.userId + "/nickname",
                {nickname: this.newNick},
				{headers: {
						Authorization: "Bearer " + localStorage.getItem("token"),
					}},
				);
			} catch (e) {
				this.errormsg = e.toString();
			}
			localStorage.setItem("nickname", this.newNick);
			localStorage.setItem("watchingNick", this.newNick)
			this.loading = false;
            location.reload();
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
			if(this.watchingId==x) {return}
			localStorage.setItem("watchingId", x);
			localStorage.setItem("watchingNick", y);
			window.location.reload();
		},

		async getFollow() {
            this.followers.forEach(f => {
                if (this.nickname == f.nickname){
                    this.followed = true;
                }
            });
    	},
	},
    mounted() {
		
		this.getProfile();
		
	}
}
</script>

<template>
    <div class="PtopBar">
		<h1>WASAPHOTO</h1>
	</div>

    <div class="Pspace">
        <div class="PtableSpace">
			<div class="PhSpace"></div>

            <div class="PwrapProfile">
			<div class="Pprofile">
				<ul class="PinfoProfile">
					<li><h3>{{this.watchingNick}}</h3></li>
					<li>{{"Followers: "+this.followers.length}}</li>
					<li>{{"Followings: "+this.followings.length}}</li>
					<li>{{"Posts: "+this.posts.length}}</li>
                    <li v-if="this.userId!=this.watchingId"><button class="Pfollow" type="submit" v-if="!this.followed" @click="follow">Follow</button><button class="Pfollow" type="submit" v-if="this.followed" @click="follow">Unfollow</button><button @click="ban" class="Pban" type="submit">Ban</button><button @click="unban" class="Punban" type="submit">Unban</button></li>
				</ul>
            </div>

			<div class="PbackProfile" v-if="this.userId!=this.watchingId">
                <ul class="PchangeNickCont">
                    <li class="PnickBar">{{this.nickname}}<label class="labelImage" for="urlBox" @click="watchProfile(this.userId, this.nickname)"><i class='bx bxs-user'></i></label></li>
                    
                </ul>
            </div>

            <div class="PchangeNick" v-if="this.userId==this.watchingId">
                <ul class="PchangeNickCont">
                    <li class="PnickBar">New photo? :o<label class="labelImage" for="urlBox"><i class='bx bx-image-add' style='color:#23aa00' ></i></label><div><input id="urlBox" type="file" accept="image/*" @change="uploadPhoto"></div></li>
                    
                </ul>
            </div>

            <div class="PchangeNick" v-if="this.userId==this.watchingId">
                <ul class="PchangeNickCont">
                    <li><h6>New username? ;)</h6></li>
                    <li class="PnickBar"><button class="PsendNick" type="submit" @click="changeUsername"><i class='bx bxs-edit-alt'></i></button><input id="newNickBox" type="text" placeholder="New Username" required v-model="newNick"></li>
                </ul>
            
			</div>
            </div>

			<div class="PhSpace"></div>

			<div class="Pstream">
				<ul class="PlistPost">
					<li v-for="p in this.posts" v-bind:key="p.post_id"><Post :postId=p.post_id :userOwner=p.owner :date=p.date :urlPost="'data:image/*;base64,'+p.urlImage" :comments=p.comments :likes=p.likes></Post></li>
				</ul>
			</div>

			<div class="PhSpace"></div>

			<div class="Psearch">
				<div class="PsearchBar">
				<input id="searchBox" type="text" placeholder="Search" required v-model="researchNick"><button class="PsearchButton" type="submit" @click="search"></button>
			    </div>
				<ul class="PlistUsers">
					<li v-if="this.userFound==false"><div>Utente non trovato :(</div></li>
					<li v-for="u in researchUsers" v-bind:key="u.user_id" @click="watchProfile(u.user_id, u.nickname)">{{ u.nickname }}</li>
				</ul>
			</div>

			<div class="PhSpace"></div>
		</div>
    </div>
</template>

<style>
.PtopBar {
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

.Pspace {
	display: flex;
	justify-content: flex-start;
    margin: 0;
	padding-top: 10px;
    background-image: url('./src/spazio1.png');
	min-height: 100vh;
	height:fit-content;
}

.PhSpace {
    width: 2%;
    background-color: none;
}

.Pspace .PtableSpace {
	display: flex;
    width: 100%;
    background-color: none;
}

.PwrapProfile{
    display: flex;
    flex-direction: column;
    width: 24%;
    background-color: transparent;
}

.PchangeNick {
    display: flex;
    margin-top: 10px;
    padding-top: 10px;
	font-size: 15px;
	height:fit-content;
    width:100%;
    background-color: #FFF;
    border-radius: 40px;
}

.PbackProfile {
	display: flex;
    margin-top: 10px;
    padding-top: 10px;
	font-size: 15px;
	height:fit-content;
    width: fit-content;
    background-color: #FFF;
    border-radius: 40px;
	padding-right: 25px;
}

.PchangeNickCont {
    list-style-type: none;
}

.PsendNick {
    background: #FFF;
	border: none;
}

.PnickBar {
    display: flex;
	align-items: center;
}

.PnickBar .bx{
    font-size: 2em;
}

.PchangeNickCont input{
    height: fit-content;
    margin-right: 10px;
    width:90%;
    text-indent: 0.5em;
	outline: none;
	border: none;
	border: 1px solid rgba(0, 0, 0, 2);
	border-radius: 40px;
}



.Pprofile {
    display: flex;
	padding-top: 10px;
	font-size: 15px;
	height:fit-content;
    width:100%;
    background-color: #FFF;
    border-radius: 40px;
    text-indent: 0em;
}

.Pfollow {
    width: 50%;
    background:  #272b8b ;
    text-align: center;
    color: #FFF;
	border: none;
    border-bottom-left-radius: 20px;
    border-top-left-radius: 20px;
    margin-top: 5px;
}

.Pban {
    width: 20%;
    background:  #880e0e ;
    text-align: center;
    color: #FFF;
	border: none;
    margin-top: 5px;
}
 
.Punban {
    width: 30%;
    background:  #c95d13 ;
    text-align: center;
    color: #FFF;
	border: none;
    border-bottom-right-radius: 20px;
    border-top-right-radius: 20px;
    margin-top: 5px;
}

.Pprofile .PinfoProfile {
	list-style-type: none;
    flex: 1;
    margin-right: 30px;
}

.Pstream {
	padding-top: 10px;
	height:fit-content;
    width: 44%;
    background-color: #FFF;
    border-radius: 40px;
    text-indent: 0em;
    margin-bottom: 10px;
}

.Pstream .PlistPost {
	margin: 0;
	padding:0;
	list-style-type: none;
}

.Psearch {
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

.Psearch .PsearchBar{
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
}

.Psearch .PsearchButton{
	width:27px;
	height: 27px; 
	margin-right: 10px;
	background-image: url('./src/search.png');
	background-size: cover;
	background-position: center;
	border: none;
	border-radius: 50%;
}

.Psearch .PlistUsers{
	top:40px;
	width: 90%;
	text-align: left;
	list-style-type:square;
	
}

.Psearch input {
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

input#urlBox{
	display:none;
}

.labelImage{
	margin-left:15px;
}


</style>