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
					//this.token = response.headers.Authorization.split(" ")[1];
					//this.userid = JSON.parse(response.data);

					confirm("Auth header: "+response.headers+"['Authorization']" +".\n"+this.username +".\n"+response.data);
					console.log(response);
					// the data is globally available to all 
					// use sessionStorage if ...
					// localStorage.setItem('token', this.token);
					// localStorage.setItem('userid', this.userid);
					localStorage.setItem('token', "d73288c7-9796-4e75-9aa0-bd6147045f40");
					localStorage.setItem('userid', "fca8954a-728d-45c0-b6ba-4e1cdf2524be");

					if (statusCode == 201) { 
						confirm("You have successfully logged in! Please setup your profile details.\n" +
										"To comply with the project requirements, you can only login with the username henceforth.\n" +
										"Since the username is the only thing that can log you in, I have decided to use another name for the profile name different than the username at least to have a little bit of security because I love secure stuffs!!!.\n" +
										"The profile name is set to a default name. Change it in your personal profile page by clicking 3 vertical dots to open a dropdown for a whole lot of other stuffs too.\n" +
										"Surf your way in now ..."
										);
						this.$router.push("/"+this.userid+"/profile");
					} else {
						this.$router.push("/"+this.userid+"/stream/");
					}
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
				this.username = '';
			}
	},
	mounted(){
		this.doLogin();
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
				<button class="btn btn-primary mt-3" type="button" :disabled="!isValid" @click="doLogin">Mr.{{ username}} Button</button>
			</div>

		</div>
</template>
  
