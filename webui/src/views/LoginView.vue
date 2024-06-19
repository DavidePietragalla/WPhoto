<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			nick: "",
            user: null,
		}
	},
	methods: {
		async doLogin() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/session", {
                    nickname: this.nick
                });
				this.user = response.data;
                localStorage.setItem("token", this.user.user_id)
                localStorage.setItem("nickname", this.user.nickname)
                localStorage.setItem("watchingId", this.user.user_id)
                localStorage.setItem("watchingNick", this.user.nickname)
			} catch (e) {
				this.errormsg = e.toString();
			}
            this.loading = false;
            this.$router.replace("/home");
		},
	},
	mounted() {
	}
}
</script>

<template>
    <div class ="logspace">
        <div class="logwrapper">
                <h1>WASAPHOTO</h1>

                <div class="loginputbox">
                    <input type="text" placeholder="Nickname" required v-model="nick">
                </div>
                
                <button type="submit" class="logbtn" @click="doLogin">Log In</button>

        </div>
    </div>
</template>

<style>
.logspace {
    display: flex;
    position: relative;
    justify-content: center;
    align-items: center;
    background-image: url('./src/spazio1.png');
    height:100vh;
}

.logspace .logwrapper{
    position: absolute;
    width: 340px;
    height: 300px;
    color: #FFF;
    text-align: center;
    border-radius: 40px;
    background-image: linear-gradient(to bottom right, #0F006E, #B5068F, #CD6700, #FEF9D5);
    background-size: 450% 450%;
    animation: waveanimation 5s ease infinite;
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

.logwrapper h1 {
    font-size: 36px;
    text-align: center;
    margin: 50px;
}

.logwrapper .loginputbox {
    width: 100%;
    height: 80px;
    margin: 0px;
}

.loginputbox input{
    width: 70%;
    height: 50%;
    background: transparent;
    border: none;
    outline: none;
    border: 2px solid rgba(255, 255, 255, 2);
    border-radius: 40px;
    color:#FFF;
    font-size: 18;
    text-indent: 1em;
}

.loginputbox input::placeholder{
    color:#FFF;
    font-size: 18;
    text-indent: 1em;
}

.logwrapper .logbtn {
    width: 70%;
    height: 40px;
    background: #FFF;
    border: #FFF;
    outline: none;
    border-radius: 40px;
    color:#000000;
    font-size: 18;
}
</style>