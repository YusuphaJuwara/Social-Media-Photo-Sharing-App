<style scoped>
.card {
  margin-bottom: 20px;
}

.imgThumbNail {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  object-fit: cover;
}

/* .user-icon {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 8px;
	display: inline-block;
} */

.user-name {
  color: #333;
  cursor: pointer;
  display: inline-block;
  vertical-align: middle;
}

.highlighted {
  color: blue;
}

.btn-primary[disabled] {
  background-color: rgb(223, 82, 35);
  cursor: not-allowed;
}
.loginclass {
  display: flex; 
  align-items: center; 
  justify-content: center; 
  height: 100vh;
}
</style>

<script>
export default {
  data: function () {
    return {
      fol: false,
      ban: false,
      banUsers: [],
      user: null,
      userid: null,
      userID: localStorage.getItem('userid'),
      token: localStorage.getItem('token'),
      clicked: false,
      hashtags: '',
      caption: '',
      highlightProfile: false,
      showCaption: false,
      showHashtag: false,
      showComments: false,
      showLikes: false,
      errormsg: null,
      loading: false,
      posts: [],
      like: false,
      commentCount: 0,
      likedBy: [],
      messagetext: '',
      profname: user['profile-name'],
      profmsg: user['profile-message'],
      gender: user['gender'],
      bdate: user['birth-date'],
      username: '',
			isValid: false,
      selectedImage: false,
    }
  },
  methods: {
    // load() {
    //   return load
    // },
    checkImage() {
      this.selectedImage = this.$refs.image.files.length > 0;
    },
    async changeUserProfilePicture() {
      // this.loading = true;
      this.errormsg = null;
      try {
        const formData = new FormData();
        formData.append("photo", this.$refs.image.files[0]);

        await this.$axios.put("/users/"+this.userid+"/profile-picture", formData, {
          headers: {
            "Content-Type": "multipart/form-data"
          }
        });
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    validateUsername() {
      const regExp = /^[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$/;
      this.isValid = regExp.test(this.username);
			},
    async logOut() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/session");

        // invalidate the data ...
        localStorage.removeItem('token');
        localStorage.removeItem('userid');

        this.$router.push("/login");
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async getUserPhotos() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/" + this.userid+"/posts/");
        this.posts = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/"+this.userid);
        this.user = response.data;
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async updateUserProfile() {
      // this.loading = true;
      this.errormsg = null;
      try {
        let v = this.bdate.split("-").reverse().join("-");
        await this.$axios.patch("/users/"+this.userid, {
          "profile-name": this.profname,
          "profile-message": this.profmsg,
          "gender": this.gender,
          "birth-date": v,
        });
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async setMyUserName() {
      // this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.put("/users/"+this.userid+"/username", {
          "username": this.username,
        });
        // await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async userDetails(uidd) {
      this.errormsg = null;
      try {
        var response = await this.$axios.get("/users/"+uidd);
        var photo = this.getUserProfilePicture(uidd)
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      return photo, response.data['profile-name'];
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
        }
        this.loading = false;
      }
    },
    async deleteUserProfilePicture() {
      this.errormsg = null;
      if (confirm('Are you sure you want to remove the profile picture?')) {
        this.loading = true;
        try {
          await this.$axios.delete("/users/"+this.userid+"/profile-picture");

          // refresh the user profile
          await this.getUserProfile();
        } catch (e) {
          this.errormsg = e.toString();
        }
        this.loading = false;
      }
    },
    async getUserFollows() {
      // this.loading = true;
      this.errormsg = null;
      try {
        var response = await this.$axios.get("/users/"+this.userid+"/follow");
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
      return response.data['followers-array'], response.data['followings-array'];
    },
    async followUser() {
      // this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.put("/users/"+this.userid+"/follow/"+this.userID);

        // refresh...
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async unfollowUser() {
      // this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/"+this.userid+"/follow/"+this.userID);

        // refresh...
        await this.getUserProfile();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    toggleFollow() {
      if (fol) {
        this.unfollowUser()
      } else {
        this.followUser()
      }
      fol = !fol
	  },
    // toggleBan() {
    //   if (ban) {
    //     this.unbanUser()
    //   } else {
    //     this.banUser()
    //   }
    //   ban = !ban
	  // },
    async getBanUsers() {
      this.loading = true;
      this.errormsg = null;
      try {
        var response = await this.$axios.get("/users/"+this.userid+"/ban/");
        this.banUsers = response.data
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async banUser() {
      // this.loading = true;
      this.errormsg = null;
      if (confirm('Are you sure?')) {
        try {
          await this.$axios.put("/users/"+this.userid+"/ban/"+this.userID);

          // refresh...
          await this.getUserProfile();
        } catch (e) {
          this.errormsg = e.toString();
        }
      }
      // this.loading = false;
    },
    async unbanUser() {
      // this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/"+this.userid+"/ban/"+this.userID);

        // refresh...
        await this.getBanUsers();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async modifyCaption(postid) {
      // this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.put("/users/" + this.userid + "/posts/" + postid,
          {
            message: this.caption
          });

        // close if successful
        this.showCaption = false

        // refresh ...
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async addHashtag(postid) {
      // this.loading = true;
      this.errormsg = null;
      try {
        const h = this.hashtags.split(',').map(tag => tag.trim());
        for (hashtag of h) {
          await this.$axios.put("/users/" + this.userid + "/posts/" + postid + "/hashtags/" + hashtag)
        };

        // close if successful
        this.showHashtag = false

        // refresh ...
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async deleteHashtag(postid) {
      // this.loading = true;
      this.errormsg = null;
      try {
        const h = this.hashtags.split(',').map(tag => tag.trim());
        for (hashtag of h) {
          await this.$axios.delete("/users/" + this.userid + "/posts/" + postid + "/hashtags/" + hashtag)
        };

        // close if successful
        this.showHashtag = false

        // refresh ...
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      // this.loading = false;
    },
    async newPost() {
      this.$router.push("/create-post");
    },
    async deletePhoto(postid) {
      // this.loading = true;
      this.errormsg = null;
      if (confirm('Are you sure you want to delete?')) {
        this.loading = true;
        try {
          await this.$axios.delete("/users/" + this.userid + "/posts/" + postid);

          // refresh ...
          await this.getUserPhotos();
        } catch (e) {
          this.errormsg = e.toString();
        }
        this.loading = false;
      }
      // this.loading = false;
    },
    async getPrivate() {
      // this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/" + this.userid + "/private-profile");
        let p = response.data ? 'Private' : 'Public'
        alert('Your profile is set to ' + p)
      } catch (e) {
        this.errormsg = e.toString();
      }
    }
    // this.loading = false;
  },
  async setPrivate() {
    // this.loading = true;
    this.errormsg = null;
    if (confirm('Are you sure you want to set it to private?')) {
      this.loading = true;
      try {
        await this.$axios.put("/users/" + this.userid + "/private-profile");

        // refresh ...
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
    }
    // this.loading = false;
  },
  async setPublic() {
    // this.loading = true;
    this.errormsg = null;
    if (confirm('Are you sure you want to set it public?')) {
      this.loading = true;
      try {
        await this.$axios.delete("/users/" + this.userid + "/private-profile");

        // refresh ...
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
    }
    // this.loading = false;
  },
  async getSinglePhoto(photoid) {
    // this.loading = true;
    this.errormsg = null;
    let photo = ''
    try {
      let response = await this.$axios.get("/photos/" + photoid);
      photo = URL.createObjectURL(new Blob([response.data]));

      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
    return photo;
  },
  async getUserProfilePicture(userid) {
    // this.loading = true;
    this.errormsg = null;
    let photo = ''
    try {
      let response = await this.$axios.get("/users/" + userid + "/profile-picture");
      photo = URL.createObjectURL(new Blob([response.data]));

      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
    return photo;
  }, 
		async getProfileName(userid) {
    // this.loading = true;
    this.errormsg = null;
    let profileName = ''
    try {
      let response = await this.$axios.get("/users/" + userid);
      profileName = response.data['profile-name']

      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
    return profileName;
  },
  userProfile(userid) {
    this.highlightProfile = !this.highlightProfile
    this.$router.push("/" + userid + "/profile");
    // this.$route.params.userid
  },
  async getPhotoComments(postid) {
    // this.loading = true;
    this.errormsg = null;
    let response = null;
    try {
      response = await this.$axios.get("/posts/" + postid + "/comments/");
      // this.likedBy = response.data['user-ids'];
      // this.commentCount = response.data['like-count'];
      // this.showLikes = true;
      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
    return response.data;
  },
  async commentPhoto(postid, message) {
    // this.loading = true;
    this.errormsg = null;
    let response = null;
    try {
      await this.$axios.post("/posts/" + postid + "/comments/", message);
      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
  },
  async uncommentPhoto(commentid) {
    // this.loading = true;
    this.errormsg = null;
    let response = null;
    if (confirm('Are you sure you want to delete this comment?')) {
      try {
        await this.$axios.delete("/comments/" + commentid);
        // refresh...
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
    }
    // this.loading = false;
    // return response.data;
  },
  async getLikes(postid) {
    // this.loading = true;
    this.errormsg = null;
    let response = null;
    try {
      response = await this.$axios.get("/posts/" + postid + "/likes/");
      // this.likedBy = response.data['user-ids'];
      // this.commentCount = response.data['like-count'];
      // this.showLikes = true;
      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
    return response.data['user-ids'], response.data['like-count'];
  },

  async likePhoto(postid){
    // this.loading = true;
    this.errormsg = null;
    try {
      await this.$axios.put("/posts/" + postid + "/likes/" + this.userid);

      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
  },
  async unlikePhoto(postid){
    // this.loading = true;
    this.errormsg = null;
    try {
      await this.$axios.delete("/posts/" + postid + "/likes/" + this.userid);

      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    // this.loading = false;
  },
  toggleLike(postid) {
    if (like) {
      this.unlikePhoto(postid)
    } else {
      this.likePhoto(postid)
    }
    like = !like
  },
  beforeCreate(){
    this.userID = localStorage.getItem('userid')
    this.token = localStorage.getItem('token')
  },
  mounted() {
    //this.getUserPhotos()
    //this.getUserProfile()
    this.userid = this.$route.params.userid
  }
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2"> {{ this.user['profile-name'] }}'s Photo Feed</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="logOut">
            <svg class="feather">
              <use href="/feather-sprite-v4.29.0.svg#log-out" />
            </svg>
            Logout
          </button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="getUserProfile">
            Refresh Profile
          </button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="newPost">
            Create New Post
          </button>
        </div>
      </div>
    </div>
    <!-- Profile details card-->
    <div class="card">
      <div class="card-header text-center">
        <img :src="getUserProfilePicture(this.user['user-id'])" alt="Oops!" width="150rem" height="150rem" class="rounded-circle mx-auto">
        <h1 class="card-title mt-3 text-center">{{ this.user['profile-name'] }}</h1>
        <p class="card-text m-2 text-center">
          <span>Gender: {{ user['gender'] }}  </span>  
          <span> Birthdate: {{ user['birth-date'] }} </span>
        </p>
      </div>
      <div class="card-body">
        <p class="card-text">{{ this.user['profile-message'] }}</p>
      </div>
      <div class="card-footer justify-content-between d-grid gap-2 d-md-block">

        <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop10">Follower {{ this.user['follower-count'] }}</button>

        <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop11">Following {{ this.user['following-count'] }}</button>

        <button v-if="userID != userid" class="btn btn-primary" type="button" @click="toggleFollow()">
					{{ !fol ? 'Follow':'Unfollow' }}
				</button>
        <button v-if="userID != userid" class="btn btn-primary" type="button" @click="banUser()">
					'Ban'
				</button>

        <button v-if="userID == userid" type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop12">Banned Users</button>

        <div class="dropdown" v-if="userID == userid">
          <button class="btn btn-primary dropdown-toggle" type="button" data-bs-toggle="dropdown" aria-expanded="false">
            <svg class="feather">
              <use href="/feather-sprite-v4.29.0.svg#more-vertical" />
            </svg>
          </button>
          <ul class="dropdown-menu" >
            <li>
              <a class="dropdown-item" href="javascript:" @click="deleteUser()">Delete Account</a>
            </li>
            <li>
              <a class="dropdown-item" href="javascript:" @click="deleteUserProfilePicture()">Delete Profile Picture</a>
            </li>
            <li>
              <a class="dropdown-item" href="javascript:" @click="getPrivate()">See if Profile is private</a>
            </li>
            <li>
              <a class="dropdown-item" href="javascript:" @click="setPrivate()">Set Private</a>
            </li>
            <li>
              <a class="dropdown-item" href="javascript:" @click="setPublic()">Set Public</a>
            </li>

            <li>
              <a class="dropdown-item" href="javascript:" data-bs-toggle="modal" data-bs-target="#staticBackdrop20">Modify Your Profile Details</a>
            </li>
            <li>
              <a class="dropdown-item" href="javascript:" data-bs-toggle="modal" data-bs-target="#staticBackdrop21">Change Username</a>
            </li>
            <li>
              <a class="dropdown-item" href="javascript:" data-bs-toggle="modal" data-bs-target="#staticBackdrop22">Change Profile Picture</a>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <!----------->
    <!-- Dropdown for followers-->
			<div class="modal fade" id="staticBackdrop10" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
				aria-labelledby="staticBackdropLabel" aria-hidden="true">
				<div class="modal-dialog">
					<div class="modal-content">
						<div class="modal-header">
							<h1 class="modal-title fs-5" id="staticBackdropLabel">Users who followed you</h1>
							<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
						</div>
						<div class="modal-body" v-for="uid in getUserFollows()[0]" :key="uid">
							<div class="card mb-3 p-2">
								<div class="card-header header-left">
									<a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
										@click="userDetails(uid)[0]" data-bs-dismiss="modal">
										<img class="imgThumbNail" :src="getUserProfilePicture(uid)" />
										<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ this.userDetails(uid)[1] }}</h5>
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
      <!-- Dropdown for followings-->
			<div class="modal fade" id="staticBackdrop11" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
				aria-labelledby="staticBackdropLabel" aria-hidden="true">
				<div class="modal-dialog">
					<div class="modal-content">
						<div class="modal-header">
							<h1 class="modal-title fs-5" id="staticBackdropLabel">Users whom you are following</h1>
							<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
						</div>
						<div class="modal-body" v-for="uid in getUserFollows()[1]" :key="uid">
							<div class="card mb-3 p-2">
								<div class="card-header header-left">
									<a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
										@click="userProfile(uid)" data-bs-dismiss="modal">
										<img class="imgThumbNail" :src="userDetails(uid)[0]" />
										<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ this.userDetails(uid)[1] }}</h5>
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

      <!-- Dropdown for banned users-->
			<div class="modal fade" id="staticBackdrop12" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
				aria-labelledby="staticBackdropLabel" aria-hidden="true">
				<div class="modal-dialog">
					<div class="modal-content">
						<div class="modal-header">
							<h1 class="modal-title fs-5" id="staticBackdropLabel">Users whom you have banned. Click on a user to unban</h1>
							<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
						</div>
						<div class="modal-body" v-for="uid in banUsers" :key="uid">
							<div class="card mb-3 p-2">
								<div class="card-header header-left">
									<a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
										@click="unbanUser()">
										<img class="imgThumbNail" :src="userDetails(uid)[0]" />
										<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ this.userDetails(uid)[1] }}</h5>
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

    <!-- Dropdown Button item: modify user details -->
    <div class="modal fade" id="staticBackdrop20" data-bs-backdrop="static" data-bs-keyboard="false"
      tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="staticBackdropLabel">
              Modify Your Profile Details. At least one property must be set.
            </h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal"
              aria-label="Close"></button>
          </div>
          <div class="modal-body">

            <label for="profname">Profile Name: </label> <br>
            <input type="text" id="profname" v-model="profname" placeholder="Enrico204" title="'^(?=.*?[a-zA-Z]).{8,20}$' # 8 to 20 chars of at least 1 alph. No new line char"/>
            <br><br>
            <label for="profmsg">Profile Message: </label> <br>
            <textarea class="form-control" id="profmsg" v-model="profmsg" rows="20" placeholder="Edit here"></textarea>
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
            <button type="button" class="btn btn-secondary"
              data-bs-dismiss="modal">Close</button>
            <button type="button" class="btn btn-primary" data-bs-dismiss="modal" @click="updateUserProfile()">Save</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Dropdown Button item change username -->
    <div class="modal fade" id="staticBackdrop21" data-bs-backdrop="static" data-bs-keyboard="false"
      tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="staticBackdropLabel">Modify username.</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal"
              aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <label for="loginid">Username: </label>
            <input type="text" id="loginid" v-model="username" @input="validateUsername" placeholder="Enrico204" title="The username contains at least one alphabet and 7 or more other alphanumeric characters in [a-zA-Z0-9] to have a minLength of 8"/>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary"
              data-bs-dismiss="modal">Close</button>
            <button type="button" :disabled="!isValid" class="btn btn-primary" data-bs-dismiss="modal" @click="setMyUserName()">Save</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Dropdown Button item change prof pic -->
    <div class="modal fade" id="staticBackdrop22" data-bs-backdrop="static" data-bs-keyboard="false"
      tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="staticBackdropLabel">Modify Profile Picture.</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal"
              aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <label for="image">Upload Image: </label>
            <input type="file" id="image" name="image" ref="image" accept="image/png" @change="checkImage" title="Select exactly 1 png image">
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary"
              data-bs-dismiss="modal">Close</button>
            <button type="button" :disabled="!selectedImage" class="btn btn-primary" data-bs-dismiss="modal" @click="changeUserProfilePicture()">Save</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal end-->

    <!---------------------------------------->
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <LoadingSpinner v-if="loading"></LoadingSpinner>

    <div class="card" v-if="posts.length===0">
      <div class="card-body">
        <p>No posts in the database.</p>

        <a href="javascript:" class="btn btn-primary" @click="newPost">Create New Post</a>
      </div>
    </div>

    <div class="card" v-if="!loading" v-for="post in posts" :key="post['post-id']">
      <div class="card-header">
        <div class="header-left">
          <!--a href="javascript:" class="btn btn-danger" @click="deleteFountain(post.id)">Delete Comment</a-->

          <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
            @click="userProfile(post['user-id'])">
            <img class="imgThumbNail" :src="userDetails(post['user-id'])[0]" />
            <h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ userDetails(post['user-id'])[1] }}
            </h5>
          </a>
        </div>
        <div class="header-right" v-if="userid==post['user-id']">
          <div class="dropdown">
            <button class="btn btn-primary dropdown-toggle" type="button" data-bs-toggle="dropdown"
              aria-expanded="false">
              <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#more-vertical" />
              </svg>
            </button>
            <ul class="dropdown-menu">
              <li><a class="dropdown-item" href="javascript:" @click="deletePhoto(post['post-id'])">Delete</a></li>
              <li><a class="dropdown-item" href="javascript:" data-bs-toggle="modal" data-bs-target="#staticBackdrop1"
                  @click="showCaption=true">Modify Caption</a></li>
              <li><a class="dropdown-item" href="javascript:" data-bs-toggle="modal" data-bs-target="#staticBackdrop2"
                  @click="showHashtag=true">Add/Delete Hashtags</a></li>
            </ul>
          </div>
        </div>
      </div>

      // style="width:fit-content"
      <img class="card-img-top" :src="getSinglePhoto(post['photo-id'])" alt="oops, can't load photo!" />

      <div class="card-body">
        <small class="text-muted"> {{ post['date-time']}} </small> <br>
        <p class="card-text"> {{ caption }} </p>
        <p class="card-text"> Hashtags: {{ hashtags }} </p>
      </div>
      <div class="card-footer">
        <div class="d-grid gap-2 d-md-block">
          <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop3">Liked
            by</button>
          <!--button class="btn btn-primary" type="button" @click="getLikes(post['post-id'])">Liked by</button-->
          <button class="btn btn-primary" type="button" @click="toggleLike(post['post-id'])">
            {{ !like? 'Like':'Unlike' }} {{ getLikes(post['post-id'])[1] }}
          </button>
          <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop4">Comment
            {{ commentCount }}</button>
          <!--button class="btn btn-primary" type="button" @click="showComments=true">
						Comment {{ commentCount }}</button-->
        </div>
      </div>

      <!-- Dropdown for likes-->
      <div class="modal fade" id="staticBackdrop3" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h1 class="modal-title fs-5" id="staticBackdropLabel">Users who liked the post</h1>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body" v-for="uid in getLikes(post['post-id'])[0]" :key="uid">
              <div class="card mb-3 p-2">
                <div class="card-header header-left">
                  <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
                    @click="userProfile(uid)" data-bs-dismiss="modal">
                    <img class="imgThumbNail" :src="getUserProfilePicture(uid['user-id'])" />
                    <h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ getProfileName(uid) }}</h5>
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
      <!-- Dropdown for comments-->
      <div class="modal fade" id="staticBackdrop4" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h1 class="modal-title fs-5" id="staticBackdropLabel">Comments</h1>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body card" v-for="comment in getPhotoComments(post['post-id'])"
              :key="comment[comment-id]">
              <!--------------------->
              <div class="card-header">
                <div class="header-left">
                  <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
                    @click="userProfile(comment['user-id'])">
                    <img class="imgThumbNail" :src="getUserProfilePicture(comment['user-id'])" />
                    <h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{
                      getProfileName(comment['user-id'])
                    }}</h5>
                  </a>
                </div>
                <div class="header-right">
                  <button v-if="userid==comment['user-id']||userid==post['user-id']" data-bs-dismiss="modal"
                    @click="uncommentPhoto(comment['comment-id'])">&times</button>
                </div>l
              </div>
              <div class="card-body">
                <small class="text-muted"> {{ comment['date-time']}} </small> <br>
                <p class="card-text"> {{ comment.message }} </p>
              </div>
            </div>
            <div class="modal-footer">
              <div class="row m-2 p-2">
                <textarea class="form-control" id="messagetext" v-model="messagetext" cols="10"
                  placeholder="Edit here"></textarea>
                <button type="button" data-bs-dismiss="modal" class="btn btn-primary"
                  @click="commentPhoto(comment['postid-id'], messagetext)" cols="2">Send</button>
              </div>
            </div>
          </div>
        </div>

        <!-- Dropdown Button items trigger modals -->
        <div v-if="showCaption" class="modal fade" id="staticBackdrop1" data-bs-backdrop="static"
          data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
          <div class="modal-dialog modal-dialog-scrollable">
            <div class="modal-content">
              <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel">Modify Caption</h1>
                <button @click="showCaption=false" type="button" class="btn-close" data-bs-dismiss="modal"
                  aria-label="Close"></button>
              </div>
              <div class="modal-body">
                <textarea class="form-control" id="caption" v-model="caption" rows="20"
                  placeholder="Edit here"></textarea>
              </div>
              <div class="modal-footer">
                <button type="button" @click="showCaption=false" class="btn btn-secondary"
                  data-bs-dismiss="modal">Close</button>
                <button type="button" class="btn btn-primary" @click="modifyCaption(post['post-id'])">Save</button>
              </div>
            </div>
          </div>
        </div>

        <div class="modal fade" id="staticBackdrop2" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
          aria-labelledby="staticBackdropLabel" aria-hidden="true">
          <div class="modal-dialog modal-dialog-scrollable">
            <div class="modal-content">
              <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel">Add/Delete Hashtags. Add Commas if more than one
                </h1>
                <button type="button" @click="showHashtag=false" class="btn-close" data-bs-dismiss="modal"
                  aria-label="Close"></button>
              </div>
              <div class="modal-body">
                <textarea class="form-control" id="hashtags" v-model="hashtags" rows="20"
                  placeholder="Edit here"></textarea>
              </div>
              <div class="modal-footer">
                <button type="button" @click="showHashtag=false" class="btn btn-secondary"
                  data-bs-dismiss="modal">Close</button>
                <button type="button" class="btn btn-primary" @click="deleteHashtag(post['post-id'])">Delete</button>
                <button type="button" class="btn btn-primary" @click="addHashtag(post['post-id'])">Add</button>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal end-->
      </div>
    </div>
  </div>
</template>
