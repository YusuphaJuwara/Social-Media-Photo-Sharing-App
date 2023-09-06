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
	width: 50%;
	height: 50%;
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
      banners: [],
      banneds: [],

      userFollowers: [],
      userFollowings: [],
      userBans: [],

			form: {
        picture: null,
				pictureError: false,
      },
      showChangeBtn: false,
      fol: false,

      ban: false,
      priv: false,

      isUsernameValid: false,
      username: '',
      profname: '',
      profmsg: '',
      gender: null,
      bdate: null,

      showBan: false, 
      showFollowing: false, 
      showFollower: false,

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
		async getUserProfile() {
      console.log("ProfileDetails.vue -> getUserProfile(): userid: "+this.userid)
      this.user = await this.getProfile(this.userid);
      this.profPic = await this.getUserProfilePicture(this.userid);

      console.log("ProfileDetails.vue -> getUserProfile(): user: "+this.user)
      console.log("ProfileDetails.vue -> getUserProfile(): profPic: "+this.profPic)

      if (this.user){
        this.profname = this.user['profile-name']
        this.profmsg = this.user['profile-message']
        this.gender = this.user['gender']
        this.bdate = this.user['birth_date']
      }

      const response = await this.$axios.get("/users/" + this.userid + "/private-profile");
      this.priv = response.data ? true : false

      console.log("ProfileDetails.vue -> getUserProfile(): \nthis.private = response.data ? true : false: "+this.priv+"\nresponse.data"+response.data)

      await this.getBanUsers();
      await this.getUserFollows();
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

          if (this.banners.includes(follower)){
            this.profilePicsFollower.push("unknown");
            console.log("getFollowProfilePicture(): follower banned: ");
          } else {
            const photo = await this.getUserProfilePicture(follower)
            this.profilePicsFollower.push(photo);
            console.log("getFollowProfilePicture(): \nphoto: "+photo);
          }

          console.log("getFollowProfilePicture(): \nthis.profilePicsFollower.push(photo): "+this.profilePicsFollower);
        }

        this.profilePicsFollowing = [];
        for (const following of this.followings) {
          console.log("getFollowProfilePicture(): this.followings: "+this.followings)

          if (this.banners.includes(following)){
            this.profilePicsFollowing.push("unknown");
            console.log("getFollowProfilePicture(): following banned: ");
          } else {
            const photo = await this.getUserProfilePicture(following)
            this.profilePicsFollowing.push(photo);
            console.log("getFollowProfilePicture(): \nphoto: "+photo);
          }

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
        if (this.banners.includes(follower)){
          this.userFollowers.push("unknown");
          console.log("getBanProfilePicture(): follower banned");
        } else {
          const uDet = await this.getProfile(follower);
          this.userFollowers.push(uDet['profile-name'])
        }
			}

			this.userFollowings = [];
			for (const following of this.followings) {
        if (this.banners.includes(following)){
          this.userFollowings.push("unknown");
          console.log("getBanProfilePicture(): following banned");
        } else {
          const uDet = await this.getProfile(following);
          this.userFollowings.push(uDet['profile-name'])
        }
			}
		},

    async getBanProfilePicture() {
			this.loading = true;
			this.errormsg = null;
			try {
        this.profilePicsBan = [];
        for (const banned of this.banneds) {

          // If they ban each other
          if (this.banners.includes(banned)){
            this.profilePicsBan.push("unknown");
            console.log("getBanProfilePicture(): they banned each other");
          } else {
            const photo = await this.getUserProfilePicture(banned)
            this.profilePicsBan.push(photo);
            console.log("getBanProfilePicture(): they didn't ban each other");
          }

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
      for (const banned of this.banneds) {
        this.userBans = [];

        // If they ban each other
        if (this.banners.includes(banned)){
          this.userBans.push("unknown");
          console.log("getBanProfilePicture(): they banned each: ");
        } else {
          const uDet = await this.getProfile(banned);
          this.userBans.push(uDet['profile-name'])
        }
      }
		},

    async getUserFollows() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/users/"+this.userid+"/follow");

        this.followers = [];
        if (response.data['followers-array'] != null){
          this.followers = response.data['followers-array'];

          // Check if the logged in user follows the user whose info is been requested.
          if (response.data['followers-array'].includes(this.loggedUid)){
						this.fol = true;
					}
        }

        this.followings = [];
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
        const response = await this.$axios.get("/users/"+this.loggedUid+"/ban");
        console.log("getBanUsers(): response.data: "+response.data)

        // Those who banned him
        this.banners = [];
        if (response.data['banners'] != null){
          this.banners = response.data['banners'];
          console.log("getBanUsers(): banners: "+this.banners);
        }

        // Those whom he banned
        this.banneds = [];
        if (response.data['banneds'] != null){
          this.banneds = response.data['banneds'];
          console.log("getBanUsers(): banneds: "+this.banneds);

          // Check if the logged in user bans the user whose info is been requested.
          if (this.banneds.includes(this.userid)){
						this.ban = true;
					}
        }

        if (this.loggedUid===this.userid){
          await this.getBanProfilePicture();
          await this.getBanProfiles();
        }

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

    async removeBanHelperFunc(uid) {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/"+this.loggedUid+"/ban/"+uid);

        // refresh...
        // await this.getBanUsers();
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
          this.loading = false;
      }
    },

    async unbanUser() {
      await this.removeBanHelperFunc(this.userid)
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
      if (this.priv) {
        this.setPublic()
      } else {
        this.setPrivate()
      }
      this.priv = !this.priv
	  },

    validateUsername() {
      const regExp = /^[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$/;
      this.isUsernameValid = regExp.test(this.username);
		},

    async openIt(p){
			this.showChangeBtn = false;
      this.showFollower = false;
      this.showFollowing = false;
      this.showBan = false;

      if (p === "C"){
        this.showChangeBtn = true;
      } else if (p === "R") {
        this.showFollower = true;
        await this.getUserFollows();
      } else if (p === "G") {
        this.showFollowing = true;
        await this.getUserFollows();
      } else if (p === "B") {
        this.showBan = true;
        await this.getBanUsers();
      }
    },

    async closeIt(p){
      this.showChangeBtn = false;
      this.showFollower = false;
      this.showFollowing = false;
      this.showBan = false;

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
    <div class="row ms-2 me-2" style="background-color: aqua;border-color: rgb(233, 6, 66);border-style: double;" v-if="loggedUid===userid">
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
          {{ !priv ? 'Set Private':'Set Public' }}
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
    <div class="card-body card-with-bg">
      <div class="container d-flex justify-content-center m-1">
        <img class="card-img-top user-icon" :src="profPic" alt="oops, can't load photo!" />
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
      <button type="button" class="btn btn-primary me-3" @click="openIt('R')">Followers {{ user['follower-count'] }}</button>
      <button type="button" class="btn btn-primary me-3" @click="openIt('G')">Followings {{ user['following-count'] }}</button>
      <button v-if="loggedUid != userid" class="btn btn-primary me-3" type="button" @click="toggleFollow()">
        {{ !fol ? 'Follow':'Unfollow' }}
      </button>
      <button type="button"  v-if="loggedUid===userid" class="btn btn-primary me-3" @click="openIt('B')">Banned {{ banneds.length }}</button>
      <button v-if="loggedUid != userid" class="btn btn-primary me-3" type="button" @click="toggleBan()">
        {{ !ban ? 'Ban':'Remove Ban' }}
      </button>
    </div>
  </div>

  <!-- Followers -->
  <div class="container mt-3 mb-2" style="border: 2px solid violet;" v-if="showFollower && followers.length>0">
    <div class="card m-2 p-2">
      <div class="card-header m-2 p-1">Users who followed you</div>
      <div v-for="(uid, idx) in followers" :key="uid">
        <div class="card-header">
          <LinkToUserProfile v-if="!banners.includes(uid)"
            :profpic="profilePicsFollower[idx]" 
            :userprofname="userFollowers[idx]"
            :uid="uid">
          </LinkToUserProfile>

          <NoLinkToUserProfile v-if="banners.includes(uid)" ></NoLinkToUserProfile>
        </div>
      </div>
    </div>
  </div>

  <!-- Followings -->
  <div class="container mt-3 mb-2" style="border: 2px solid violet;" v-if="showFollowing && followings.length>0">
    <div class="card m-2 p-2">
      <div class="card-header m-2 p-1"> Users whom you followed </div>
      <div v-for="(uid, idx) in followings" :key="uid">
        <div class="card-header">
          <LinkToUserProfile v-if="!banners.includes(uid)"
            :profpic="profilePicsFollowing[idx]" 
            :userprofname="userFollowings[idx]"
            :uid="uid">
          </LinkToUserProfile>

          <NoLinkToUserProfile v-if="banners.includes(uid)" ></NoLinkToUserProfile>
        </div>
      </div>
    </div>
  </div>

  <!-- Ban -->
  <div class="container mt-3 mb-2" style="border: 2px solid violet;" v-if="showBan && banneds.length>0">
    <div class="card m-2 p-2">
      <div class="card-header m-2 p-1"> Users whom you banned </div>
      <div v-for="(uid, idx) in banneds" :key="uid">
        <div class="card-header">
          <LinkToUserProfile v-if="!banners.includes(uid)"
            :profpic="profilePicsBan[idx]" 
            :userprofname="userBans[idx]"
            :uid="uid">
          </LinkToUserProfile>

          <NoLinkToUserProfile v-if="banners.includes(uid)" ></NoLinkToUserProfile>
          <button type="button" class="btn btn-primary m-2" @click="removeBanHelperFunc(uid)">Remove Ban</button>
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

