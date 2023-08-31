<style>
</style>

<script>

export default {
	data() {
		return {
			errormsg: null,
			loading: false,
      // userid: '',
      users: [],
      userProfilePics: [],

		}
	},
  methods: {
    async getAllUsers() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/");
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

    async getUserProfilePicture() {
			this.loading = true;
			this.errormsg = null;
			try {
        this.userProfilePics = [];
        for (const user of this.users) {
          var uid = user['user-id']
          const response = await this.$axios.get(`/users/"${uid}/profile-picture`, {responseType: "blob"});
          const photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
          this.userProfilePics.push(photo);
        }
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},


	},
	// created(){
	// 	this.userid = localStorage.getItem('userid')
	// },
	async mounted(){
		await this.getAllUsers();
    await this.getUserProfilePicture();

	}

}
</script>
<template>
<div>

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	<LoadingSpinner v-if="loading"></LoadingSpinner>

  <div class="card m-2" style="border: 1px solid red;" v-if="!loading" v-for="(user, idx) in users" >
    <div class="card-header m-2">
      <LinkToUserProfile
      :profpic="userProfilePics[idx]" 
      :userprofname=" user['profile-name']"
      :uid="user['user-id']">
      </LinkToUserProfile>
    </div>
    <div class="card-body">
      <div class="row">
        <div class="col-md-4 d-flex justify-content-center">Post Count: {{ user['post-count'] }}</div>
        <div class="col-md-4 d-flex justify-content-center">Follower Count: {{ user['follower-count'] }}</div>
        <div class="col-md-4 d-flex justify-content-center">Following Count: {{ user['following-count'] }}</div>
      </div>
    </div>
  </div>

</div>
</template>