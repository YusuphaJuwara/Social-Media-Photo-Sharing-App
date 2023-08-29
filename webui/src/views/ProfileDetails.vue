<style scoped>
.error {
  color: red;
}

.hint {
  color: gray;
}

.card-with-bg {
  margin-bottom: 20px;
  background-color:beige;
}

.user-icon {
	width: 500px;
	height: 500px;
  vertical-align: middle;
	border-radius: 50%;
	object-fit: cover;
	display: inline-block;
}

/* .btn-primary[disabled] {
  background-color: rgb(223, 82, 35);
  cursor: not-allowed;
} */
</style>

<script>

export default {
  props: ['userid'],
	data() {
		return {
      errormsg: null,
      loading: false,
      loggedUid: '',
      user: {},
      profPic: null,
      profilePicsFollower: [],
      profilePicsFollowing: [],
      profilePicsBan: [],
      followers: [],
      followings: [],
      userFollowers: [],
      userFollowings: [],
      userBans: [],
			isValid: false,
			form: {
        picture: null,
				pictureError: false,
      },
      showChangeBtn: false,
      fol: false,

      banUsers: [],
      ban: false,
      private: false,

      isUsernameValid: false,
      username: '',
      profname: '',
      profmsg: '',
      gender: null,
      bdate: null,

      highlightProfile: false,

    }
	},

	computed: {
    isValid() {
      if (!this.form.picture) {
        return false;
      }
      if (this.form.picture.type !== 'image/png') {
        return false;
      }
      return true;
    },
  },

  methods: {
    async userProfile() {
			this.highlightProfile = !this.highlightProfile

			if (this.loggedUid===this.userid){
				this.$router.push("/profile");
				console.log("userProfile(): this.loggedUid===this.userid. \nuserid: "+this.userid+"\nloggedUid: "+this.loggedUid)
			} else {
				this.$router.push("/"+this.userid+"/profile");
				console.log("userProfile(): this.loggedUid===this.userid. \nuserid: "+this.userid+"\nloggedUid: "+this.loggedUid);
			}
		},

    // async getUserProfile() {
    //   this.loading = true;
    //   this.errormsg = null;
    //   try {
    //     const response = await this.$axios.get("/users/"+this.userid);
    //     this.user = response.data;
    //     // this.user = await this.getUserProfile();
    //     await this.getUserProfilePicture();
    //   } catch (e) {
    //     this.errormsg = e.toString();
    //   } finally {
		// 		this.loading = false;
		// 	}
    // },
		async getUserProfile() {
      console.log("ProfileDetails.vue -> getUserProfile(): userid: "+this.userid)
      this.user = await this.getProfile(this.userid);
      this.profPic = await this.getUserProfilePicture(this.userid);

      console.log("ProfileDetails.vue -> getUserProfile(): user: "+this.user)
      console.log("ProfileDetails.vue -> getUserProfile(): profPic: "+this.profPic)

      this.profname = "this.user['profile-name']"
      this.profmsg = "this.user['profile-message']"
      this.gender = "this.user['gender']"
      this.bdate = "this.user['birth_date']"
      if (this.user){
        this.profname = this.user['profile-name']
        this.profmsg = this.user['profile-message']
        this.gender = this.user['gender']
        this.bdate = this.user['birth_date']
      }

      const response = await this.$axios.get("/users/" + this.userid + "/private-profile");
      this.private = response.data ? true : false

      console.log("ProfileDetails.vue -> getUserProfile(): this.private = response.data ? true : false: "+this.private+"\nresponse.data"+response.data)

      await this.getUserFollows();

      if (this.userid===this.uid){
        await this.getBanUsers();
      }
    },

    async updateUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let v = '';
        if (this.bdate != '' && this.bdate != null && this.bdate != undefined){
          v = this.bdate//.split("-").reverse().join("-");
        }

        console.log("bdate: "+this.bdate+"\nv: "+v)

        const formData = new FormData();
        formData.append("profile-name", this.profname)
        formData.append("profile-message", this.profmsg)
        formData.append("gender", this.gender)
        formData.append("birth-date", v)

        await this.$axios.patch("/users/"+this.userid, formData);
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
        this.loading = false;
      }
    },

    async setMyUserName() {
      this.loading = true;
      this.errormsg = null;
      try {
        const formData = new FormData();
        formData.append("username", this.username)
        await this.$axios.put("/users/"+this.userid+"/username", formData);

        this.username = '';

        await this.getUserProfile();

      } catch (e) {
        this.errormsg = e.toString();
      } finally {
        this.loading = false;
      }
    },

    async getProfile(uid) {
      this.loading = true;
      this.errormsg = null;
      let response = {};
      try {
        response = await this.$axios.get("/users/"+uid);
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
      return response.data;
    },

    async changeUserProfilePicture() {
      this.loading = true;
      this.errormsg = null;
      try {
        const formData = new FormData();
        formData.append('photo', this.form.picture);

        await this.$axios.put("/users/"+this.userid+"/profile-picture", formData, {
          headers: {
            "content-type": "multipart/form-data"
          }
        });

				this.form.picture = null;

        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

    async deleteUserProfilePicture() {
      this.errormsg = null;
      if (confirm('Are you sure you want to remove the profile picture? Note that even if you delete your profile picture, it will be set to the default profile picture')) {
        this.loading = true;
        try {
          await this.$axios.delete("/users/"+this.userid+"/profile-picture");

          // refresh the user profile
          await this.getUserProfile();
        } catch (e) {
          this.errormsg = e.toString();
        } finally {
          this.loading = false;
        }
      }
    },

    async deleteUser() {
      this.errormsg = null;
      if (confirm('Are you sure you want to delete your account permanently?')) {
        this.loading = true;
        try {
          await this.$axios.delete("/users/"+this.userid);
          
          // invalidate the data ...
          localStorage.removeItem('token');
          localStorage.removeItem('userid');

          this.$router.push("/login");
        } catch (e) {
          this.errormsg = e.toString();
        } finally {
          this.loading = false;
        }
      }
    },

    async getUserProfilePicture(uid) {
      this.loading = true;
      this.errormsg = null;
      let photo = null;
      try {
        let response = await this.$axios.get("/users/"+uid+"/profile-picture", {responseType: "blob"});
        photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
      } catch (e) {
        this.errormsg = "getUserProfilePicture(uid): \nuid: "+uid+"\nerror"+e.toString();
      }finally {
        this.loading = false;
      }
      return photo;
    }, 

    async getPrivate() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/users/" + this.userid + "/private-profile");
        const p = response.data ? 'Private' : 'Public'
        alert('Your profile is set to ' + p)
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
        this.loading = false;
      }
    },

    async setPrivate() {
      this.errormsg = null;
      if (confirm('Are you sure you want to set it to private?')) {
        this.loading = true;
        try {
          await this.$axios.put("/users/" + this.userid + "/private-profile");

          // refresh ...
          await this.getUserProfile();
        } catch (e) {
          this.errormsg = e.toString();
        } finally {
          this.loading = false;
        }
      }
    },
    
    async setPublic() {
      this.errormsg = null;
      if (confirm('Are you sure you want to set it public?')) {
        this.loading = true;
        try {
          await this.$axios.delete("/users/" + this.userid + "/private-profile");

          // refresh ...
          await this.getUserProfile();
        } catch (e) {
          this.errormsg = e.toString();
        } finally {
          this.loading = false;
        }
      }
    },

    async getFollowProfilePicture() {
			this.loading = true;
			this.errormsg = null;
			try {
        this.profilePicsFollower = [];
        for (const follower of this.followers) {
          console.log("getFollowProfilePicture(): this.followers: "+this.followers)
          const fUid = follower
          // const response = await this.$axios.get(`/users/"${fUid}/profile-picture`, {responseType: "blob"});
          // const photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
          console.log("getFollowProfilePicture(): before call to this.getUserProfilePicture(fUid). fUid: "+fUid)
          const photo = await this.getUserProfilePicture(fUid)
          this.profilePicsFollower.push(photo);
          console.log("getFollowProfilePicture(): \nbanned['user-id']: "+fUid);
          console.log("getFollowProfilePicture(): \nphoto: "+photo);
          console.log("getFollowProfilePicture(): \nthis.profilePicsFollower.push(photo): "+this.profilePicsFollower);
        }

        this.profilePicsFollowing = [];
        for (const following of this.followings) {
          console.log("getFollowProfilePicture(): this.followings: "+this.followings)
          const fUid = following;
          // const response = await this.$axios.get(`/users/"${fUid}/profile-picture`, {responseType: "blob"});
          // const photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
          console.log("getFollowProfilePicture(): before call to this.getUserProfilePicture(fUid). fUid: "+fUid)
          const photo = await this.getUserProfilePicture(fUid)
          this.profilePicsFollowing.push(photo);
          console.log("getFollowProfilePicture(): \nbanned['user-id']: "+fUid);
          console.log("getFollowProfilePicture(): \nphoto: "+photo);
          console.log("getFollowProfilePicture(): \nthis.profilePicsFollowing.push(photo): "+this.profilePicsFollowing);
        }
			} catch (e) {
				this.errormsg = "getFollowProfilePicture(): \nerror: "+e.toString();
			} finally {
				this.loading = false;
			}
		},

    async getFollowProfiles(){

			this.userFollowers = [];
			for (const follower of this.followers) {
				const fUid = follower
				const uDet = await this.getProfile(fUid);
				this.userFollowers.push(uDet['profile-name'])
			}

			this.userFollowings = [];
			for (const following of this.followings) {
				const fUid = following
				const uDet = await this.getProfile(fUid);
				this.userFollowings.push(uDet['profile-name'])
			}
		},

    async getBanProfilePicture() {
			this.loading = true;
			this.errormsg = null;
			try {
        this.profilePicsBan = [];
        for (const banned of this.banUsers) {
          const fUid = banned
          // const response = await this.$axios.get(`/users/"${fUid}/profile-picture`, {responseType: "blob"});
          // const photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
          const photo = await this.getUserProfilePicture(fUid)
          this.profilePicsBan.push(photo);
          console.log("getBanProfilePicture(): \nbanned['user-id']: "+fUid);
          console.log("getBanProfilePicture(): \nphoto: "+photo);
          console.log("getBanProfilePicture(): \nthis.profilePicsBan.push(photo): "+this.profilePicsBan);
        }
			} catch (e) {
				this.errormsg = "getBanProfilePicture(): \nerror: "+e.toString();
			} finally {
				this.loading = false;
			}
		},

    async getBanProfiles(){
			this.userBans = [];
      if (this.banUsers != null){
        for (const banned of this.banUsers) {
            const fUid = banned
          const uDet = await this.getProfile(fUid);
          this.userBans.push(uDet['profile-name'])
        }
      }
		},

    async getUserFollows() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/users/"+this.userid+"/follow");

        if (response.data['followers-array'] != null){
          this.followers = response.data['followers-array'];
        }
        if (response.data['followings-array'] != null){
          this.followings = response.data['followings-array'];
        }

        console.log("response: "+response+"\nfollowers: "+this.followers+"\nfollowings: "+this.followings)

        await this.getFollowProfilePicture();
        await this.getFollowProfiles();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

    async followUser() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.put("/users/"+this.loggedUid+"/follow/"+this.userid);

        // refresh...
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

    async unfollowUser() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/"+this.loggedUid+"/follow/"+this.userid);

        // refresh...
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

    async getBanUsers() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/users/"+this.userid+"/ban");
        console.log("getBanUsers(): response.data: "+response.data)

        this.banUsers = [];
        if (response.data != null) {
          this.banUsers = response.data
          await this.getBanProfilePicture();
          await this.getBanProfiles();
        }

        console.log("getBanUsers(): banUsers: "+this.banUsers)

      } catch (e) {
        this.errormsg = "getBanUsers():\nerror: "+e.toString()+"\nthis.userid: "+this.userid;
      } finally {
				this.loading = false;
			}
    },

    async banUser() {
      this.errormsg = null;
      if (confirm('Are you sure you want to ban the user?')) {
        this.loading = true;
        try {
          await this.$axios.put("/users/"+this.loggedUid+"/ban/"+this.userid);

          // refresh...
          // await this.getBanUsers();
          await this.getUserProfile();
        } catch (e) {
          this.errormsg = e.toString();
        } finally {
          this.loading = false;
        }
      }
    },

    async unbanUser() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/"+this.loggedUid+"/ban/"+this.userid);

        // refresh...
        // await this.getBanUsers();
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
          this.loading = false;
      }
    },

    uploadPicture(event) {
      const file = event.target.files[0];
      if (file && file.type === 'image/png') {
        this.form.picture = file;
        this.form.pictureError = false;
      } else {
        this.form.pictureError = true;
        this.$refs.pictureInput.value = '';
      }
    },

		checkImage() {
      this.selectedImage = this.$refs.image.files.length > 0;
    },

    toggleFollow() {
      if (this.fol) {
        this.unfollowUser()
      } else {
        this.followUser()
      }
      this.fol = !this.fol
	  },

    toggleBan() {
      if (this.ban) {
        this.unbanUser()
      } else {
        this.banUser()
      }
      this.ban = !this.ban
	  },

    toggleP() {
      if (this.private) {
        this.setPublic()
      } else {
        this.setPrivate()
      }
      this.private = !this.private
	  },

    validateUsername() {
      const regExp = /^[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$/;
      this.isUsernameValid = regExp.test(this.username);
		},

    openIt(p){
			this.showChangeBtn = false;

      if (p === "C"){
        this.showChangeBtn = true;
      }
    },

    async closeIt(p){
      this.showChangeBtn = false;

      if (p === "C"){
        await this.changeUserProfilePicture();
      }
    }

	},

  created(){
    // Currently logged in user
    this.loggedUid = localStorage.getItem('userid')
    console.log("LoggedUid: "+this.loggedUid)
  },

	async mounted(){
		await this.getUserProfile();

    await this.getUserFollows();
	},

}
</script>

<template>
<div>

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	<LoadingSpinner v-if="loading"></LoadingSpinner>

  <div class="card m-1" v-if="!loading">
    <div class="card-title m-1 text-center" style="font-size: 120px; color: chocolate;">
      {{ user["profile-name"] }}
    </div>
    <div class="row" v-if="loggedUid===userid">
      <div class="col-md-2 m-1 d-flex justify-content-center">
        <button type="button" class="btn btn-danger" @click="deleteUser">
        Delete Account
        </button>
      </div>
      <div class="col-md-3 m-1 d-flex justify-content-center">
        <button type="button" class="btn btn-primary" @click="getPrivate">
        See if your profile is set to private or public
        </button>
        </div>
      <div class="col-md-2 m-1 d-flex justify-content-center">
        <button class="btn btn-primary me-3" type="button" @click="toggleP()">
          {{ !private ? 'Set Private':'Set Public' }}
        </button>
      </div>
      <div class="col-md-2 m-1 d-flex justify-content-center">
        <button type="button" class="btn btn-primary me-3" data-bs-toggle="modal" data-bs-target="#staticBackdropUsername">
          Modify Username
        </button>
      </div>
      <div class="col-md-2 m-1 d-flex justify-content-center">
        <button type="button" class="btn btn-primary me-3" data-bs-toggle="modal" data-bs-target="#staticBackdropProfDetails">
          Modify Your Details
        </button>
      </div>
    </div>
    <div class="card-header">
      <div class="row" v-if="showChangeBtn===false && loggedUid===userid">
        <div class="col-md-6 d-flex justify-content-center">
          <button type="button" class="btn btn-primary" @click="openIt('C')">
          Change Profile Pic
          </button>
        </div>
        <div class="col-md-6 d-flex justify-content-center">
          <button type="button" class="btn btn-primary" @click="deleteUserProfilePicture">
          Delete Profile Pic
          </button>
        </div>
      </div>
      <div class="d-flex justify-content-center" v-if="showChangeBtn===true">
        <label for="picture">Picture:</label>
        <input type="file" id="picture" @change="uploadPicture" ref="pictureInput" accept="image/png"/>
        <div v-if="form.pictureError" class="error">Please select a PNG file</div>
      </div>
      <div class="row" v-if="showChangeBtn===true">
        <div class="col-md-6 d-flex justify-content-center">
          <button :disabled="!isValid" type="button" class="btn btn-primary" @click="closeIt('C')">
          Change
          </button>
        </div>
        <div class="col-md-6 d-flex justify-content-center">
          <button type="button" class="btn btn-primary" @click="showChangeBtn=false">
          Cancel
          </button>
        </div>
      </div>
    </div>
    <div class="card-body card-with-bg">
      <div class="container d-flex justify-content-center m-1">
        <img class="card-img-top user-icon" :src="profPic" alt="oops, can't load photo!" />
      </div>
    </div>
    <div class="row">
      <div class="col-md-6 d-flex justify-content-center">
        <div v-if="gender != '' && gender != null">
          <div class="card-text" 
          style="color: rgb(71, 181, 255);font-size:xx-large;">
            Gender: {{ gender }}
          </div>
        </div>
        <div v-if="gender == '' || gender == null">
          <div class="card-text" 
          style="color: rgb(71, 181, 255);font-size:xx-large;"> 
            Gender not set 
          </div>
        </div>
      </div>
      <div class="col-md-6 d-flex justify-content-center">
        <div class="card-text" style="color: rgb(71, 181, 255);font-size:xx-large;">
          Birthdate: {{ bdate }}
        </div>
      </div>
    </div>
    <div class="m-5 p-2">
      <p class="card-text" 
      style="font-weight:bolder;font-size: larger;">
        Profile Message: {{ user['profile-message'] }}
      </p>
    </div>
    <div class="card-footer justify-content-between d-grid gap-2 d-md-block">
      <button type="button" class="btn btn-primary me-3" data-bs-toggle="modal" data-bs-target="#staticBackdropFollower">Followers {{ user['follower-count'] }}</button>

      <button type="button" class="btn btn-primary me-3" data-bs-toggle="modal" data-bs-target="#staticBackdropFollowing">Followings {{ user['following-count'] }}</button>

      <button v-if="loggedUid != userid" class="btn btn-primary me-3" type="button" @click="toggleFollow()">
        {{ !fol ? 'Follow':'Unfollow' }}
      </button>

      <button v-if="loggedUid===userid" type="button" class="btn btn-primary me-3" data-bs-toggle="modal" data-bs-target="#staticBackdropBan">Banned {{ banUsers.length }}</button>

      <button v-if="loggedUid != userid" class="btn btn-primary me-3" type="button" @click="toggleBan()">
        {{ !ban ? 'Ban':'Remove Ban' }}
      </button>
    </div>

  </div>

  <!-- Modal for followers-->
  <div class="modal fade" id="staticBackdropFollower" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
      aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="staticBackdropLabel">Users who followed you</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body" v-for="(uid, idx) in followers" :key="uid">
            <div class="card mb-3 p-2">
              <div class="card-text m-2">UID: {{ uid }}</div>
              <div class="card-text m-2">profilePicsFollower[idx]: {{ profilePicsFollower[idx] }}</div>
              <div class="card-text m-2">userFollowers[idx]: {{ userFollowers[idx] }}</div>
              <div class="card-header">
                <!-- <LinkToUserProfile
                  :profpic="profilePicsFollower[idx]" 
                  :userprofname="userFollowers[idx]"
                  :uid="uid">
                </LinkToUserProfile> -->
                <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
                    @click="userProfile()">
                  <div class="d-flex align-items-center">
                    <img class="imgThumbNail me-2 mb-2" :src="profilePicsFollower[idx]" alt="Opps! error" />
                    <h5 class="user-name ms-2 mb-2" :class="{ 'highlighted': highlightProfile }">{{ userFollowers[idx] }}
                    </h5>
                  </div>
                </a>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
  </div>

  <!-- Modal for followings-->
  <div class="modal fade" id="staticBackdropFollowing" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
      aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="staticBackdropLabel">Users whom you followed</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body" v-for="(uid, idx) in followings" :key="uid">
            <div class="card mb-3 p-2">
              <div class="card-header">
                <!-- <LinkToUserProfile
                  :profpic="profilePicsFollowing[idx]" 
                  :userprofname="userFollowings[idx]"
                  :uid="uid">
                </LinkToUserProfile> -->
                <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
                    @click="userProfile()">
                  <div class="d-flex align-items-center">
                    <img class="imgThumbNail me-2 mb-2" :src="profilePicsFollowing[idx]" alt="Opps! error" />
                    <h5 class="user-name ms-2 mb-2" :class="{ 'highlighted': highlightProfile }">{{ userFollowings[idx] }}
                    </h5>
                  </div>
                </a>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
  </div>

  <!-- Modal for ban-->
  <div class="modal fade" id="staticBackdropBan" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
      aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="staticBackdropLabel">Users whom you banned</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body" v-for="(uid, idx) in banUsers" :key="uid">
            <div class="card mb-3 p-2">
              <div class="card-header">
                <!-- <LinkToUserProfile
                  :profpic="profilePicsBan[idx]" 
                  :userprofname="userBans[idx]"
                  :uid="uid">
                </LinkToUserProfile> -->
                <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
                    @click="userProfile()">
                  <div class="d-flex align-items-center">
                    <img class="imgThumbNail me-2 mb-2" :src="profilePicsBan[idx]" alt="Opps! error" />
                    <h5 class="user-name ms-2 mb-2" :class="{ 'highlighted': highlightProfile }">{{ userBans[idx] }}
                    </h5>
                  </div>
                </a>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
  </div>

  <!-- Modal: modify user details -->
  <div class="modal fade" id="staticBackdropProfDetails" data-bs-backdrop="static" data-bs-keyboard="false"
  tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <h1 class="modal-title fs-5" id="staticBackdropLabel">
            Modify Your Profile Details. At least one property must be set.
          </h1>
          <button type="button" class="btn-close" data-bs-dismiss="modal"
            aria-label="Close">
          </button>
        </div>
        <div class="modal-body">
          <label for="profname">Profile Name: </label> <br>
          <input type="text" id="profname" v-model="profname" placeholder="Enrico204" title="'^(?=.*?[a-zA-Z]).{8,20}$' # 8 to 20 chars of at least 1 alph. No new line char"/>
          <br><br>
          <label for="profmsg">Profile Message: </label> <br>
          <textarea class="form-control" id="profmsg" v-model="profmsg" rows="5" placeholder="This is my profile msg."></textarea>
          <br><br>
          <label for="gender">Gender: </label>
          <select id="gender" name="gender" v-model="gender">
            <option value="">Select Gender</option>
            <option value="Male">Male</option>
            <option value="Female">Female</option>
          </select>
          <br><br>
          <label for="birthdate">Birthdate: </label>
          <input type="date" id="bdate" name="bdate" v-model="bdate">
        </div>
        <div class="modal-footer">
          <div class="row">
            <div class="col-md-6 d-flex justify-content-center">
              <button type="button" class="btn btn-primary" data-bs-dismiss="modal" @click="updateUserProfile()">
                Save
              </button>
            </div>
            <div class="col-md-6 d-flex justify-content-center">
              <button type="button" class="btn btn-secondary"
                data-bs-dismiss="modal">Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal: change username -->
  <div class="modal fade" id="staticBackdropUsername" data-bs-backdrop="static" data-bs-keyboard="false"
    tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <h1 class="modal-title fs-5" id="staticBackdropLabel">Modify username.</h1>
          <button type="button" class="btn-close" data-bs-dismiss="modal"
            aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <div style="color: tomato;">Note that the username is not displayed for security reasons, but the Profile Name is the one displayed. It is used for login as per the project requirements. See "Modify Your Details" and the OpenAPI yaml file</div>
          <label for="loginid">Username: </label>
          <input type="text" id="loginid" v-model="username" @input="validateUsername" placeholder="Enrico204" title="The username contains at least one alphabet and 7 or more other alphanumeric characters in [a-zA-Z0-9] to have a minLength of 8"/>
        </div>
        <div class="modal-footer">
          <div class="row">
            <div class="col-md-6 d-flex justify-content-center">
              <button type="button" :disabled="!isUsernameValid" class="btn btn-primary" 
              data-bs-dismiss="modal" @click="setMyUserName()">
                Save Changes
              </button>
            </div>
            <div class="col-md-6 d-flex justify-content-center">
              <button type="button" class="btn btn-secondary"
                data-bs-dismiss="modal">Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  
</div>
</template>

