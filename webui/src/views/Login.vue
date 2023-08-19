<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			username: '',
			isValid: false,
			userid: null,
			token: null,
		}
	},
	methods: {
			validateUsername() {
					const regExp = /^[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$/;
					this.isValid = regExp.test(this.username);
			},
			async doLogin() {
				this.loading = true;
				this.errormsg = null;
				try {
					const formData = new FormData();
        	formData.append('username', this.username);
					let response = await this.$axios.post("/session", formData);

					let statusCode = response.status;
  				this.token = response.headers['authorization']
					this.userid = response.data;

					console.log(response);
					// the data is globally available to all 
					// use sessionStorage if ...
					localStorage.setItem('token', this.token);
					localStorage.setItem('userid', this.userid);

					if (statusCode == 201) { 
						confirm("You have successfully logged in! Please setup your profile details.\n" +
										"To comply with the project requirements, you can only login with the username henceforth.\n" +
										"Since the username is the only thing that can log you in, I found it logical to use another name for the profile name that can be different than the username to at least to have a little bit of security because I love secure stuffs!!!.\n" +
										"The profile name is set to the username as default. Change it in your personal profile page if you wish to"
										);
						//this.$router.push("/"+this.userid+"/profile");
					} //else {
					// 	this.$router.push("/"+this.userid+"/stream/");
					// }
					this.$router.push("/"+this.userid);
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
				this.username = '';
			}
	},
	mounted(){
		//this.doLogin();
	}
}
</script>
  
<style scoped>
    .btn-primary[disabled] {
        background-color: rgb(223, 82, 35);
        cursor: not-allowed;
    }
		/* .loginclass {
			display: flex; 
			align-items: center; 
			justify-content: center; 
			height: 100vh;
		} */
</style>

<template>
    <div>
      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<LoadingSpinner v-if="loading"></LoadingSpinner>

			<div class="d-grid">
				<label for="loginid">Username:</label>
				<input type="text" id="loginid" v-model="username" @input="validateUsername" placeholder="Enrico204"/>
				<button class="btn btn-primary mt-3" type="button" :disabled="!isValid" @click="doLogin">Sign in/up</button>
			</div>

		</div>
</template>
  
