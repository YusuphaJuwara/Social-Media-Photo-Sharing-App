<style scoped>
.card {
	margin-bottom: 20px;
}

.imgThumbNail {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
}

.user-icon {
	width: 50%;
	height: 50%;
  vertical-align: middle;
	object-fit: cover;
	display: inline-block;
}
.user-name {
	color: #333;
	cursor: pointer;
	display: inline-block;
	vertical-align: middle;
}

.highlighted {
	color: blue;
}
</style>

<script>
export default {
  props: ['postid'],

	data() {
		return {
      userid: '',
			hashtags: '',
			addhashtags: '',
			caption: '',
			highlightProfile: false,
			errormsg: null,
			loading: false,
			post: {},
      singlePhoto: null,
			uProfPic: null,
			like: false,
			commentCount: 0,
			likedBy: [],
			messagetext: '',
			user: {},
			userLikes: [],
			userComments: [],
			banners: [],
      banneds: [],

			profilePicsComment: [],
			profilePicsLike: [],
			comments: [],
			
			// Get the like-count and the user IDs who liked the post.
			likes: {},

			showCaption: false,
			showHashtag: false,
			showComment: false,
			showLikedBy: false,


		}
	},
  methods: {
		async getLikeAndCommentProfiles(){

			this.userComments = [];
			for (const comment of this.comments) {
				const cUid = comment['user-id']
				if (this.banners.includes(cUid)){
					this.userComments.push("unknown");
					console.log("getLikeAndCommentProfiles(): banned");
				} else {
					const uDet = await this.getUserProfile(cUid);
					this.userComments.push(uDet['profile-name'])
					console.log("getLikeAndCommentProfiles(): \nuDet: "+uDet);
				}

				console.log("getLikeAndCommentProfiles(): \nthis.profilePicsComment.push(photo): "+this.userComments);
			}

			this.userLikes = [];
			for (const lUid of this.likes['user-ids']) {
				console.log("getLikeAndCommentProfiles(): this.likes['user-ids']: "+this.likes['user-ids'])

				if (this.banners.includes(lUid)){
					this.userLikes.push("unknown");
					console.log("getLikeAndCommentProfiles(): banned");
				} else {
					const uDet = await this.getUserProfile(lUid);
					this.userLikes.push(uDet['profile-name'])
					console.log("getLikeAndCommentProfiles(): \nuDet: "+uDet);
				}

				console.log("getLikeAndCommentProfiles(): \nthis.profilePicsComment.push(photo): "+this.userLikes);
			}
		},

		async getUserProfile(uid) {
      this.loading = true;
      this.errormsg = null;
			let uDet = {};
      try {
        let response = await this.$axios.get("/users/"+uid);
        uDet = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
			return uDet;
    },

    async getPhoto() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/posts/"+this.postid);
				this.post = response.data;
				this.caption = this.post['caption'];

				this.commentCount = this.post['comment-count']

				this.hashtags = ''
				const htgs = this.post['hashtags']
				if (htgs){
					for (const htg of htgs) {
						if (htg != ''){
							this.hashtags += htg + ', ';
						}
					}
				}
				this.hashtags = this.hashtags.slice(0, -2);

			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

    async modifyCaption() {
			this.loading = true;
			this.errormsg = null;
			try {
				const formData = new FormData();
        formData.append('message', this.caption );

				await this.$axios.put(`/users/${this.userid}/posts/${this.postid}`, formData);

				//this.post['caption'] = this.caption

				// refresh ...
				await this.getPhoto();
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

    async addHashtag() {
			this.loading = true;
			this.errormsg = null;
			try {
				const h = this.addhashtags.split(',').map(tag => tag.trim());
				for (const hashtag of h) {
					await this.$axios.put("/users/"+this.userid+"/posts/"+this.postid+"/hashtags/"+hashtag)
				};

				this.addhashtags = '';

				// refresh ...
				await this.getPhoto();
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

    async deleteHashtag() {
			this.loading = true;
			this.errormsg = null;
			try {
				const h = this.addhashtags.split(',').map(tag => tag.trim());
				for (const hashtag of h) {
					await this.$axios.delete("/users/"+this.userid+"/posts/"+this.postid+"/hashtags/"+hashtag)
				};

				this.addhashtags = ''

				// refresh ...
				await this.getPhoto();
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

    async deletePhoto() {
			this.errormsg = null;
			if (confirm('Are you sure you want to delete this post?')) {
				this.loading = true;
				try {
					await this.$axios.delete("/users/" +this.userid+"/posts/" + this.postid);

					this.$emit('postDeleted');
					// refresh ...
					// await this.getPhoto();
				} catch (e) {
					this.errormsg = e.toString();
				} finally {
					this.loading = false;
				}
			}
		},

    async getSinglePhoto(pUid) {
			this.loading = true;
			this.errormsg = null;
			let photo = '';
			try {
				let response = await this.$axios.get("/photos/" + pUid, {responseType: "blob"});
				photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });

				// refresh...
				// await this.getPhoto();
			} catch (e) {
				this.errormsg = "getSinglePhoto(pUid) error: ------ "+e.toString();
			} finally {
				this.loading = false;
			}
			return photo;
		},

    async getUserProfilePicture(getOne) {
			this.loading = true;
			this.errormsg = null;
			try {

				if (getOne===true){

					this.profilePicsComment = [];
					for (const comment of this.comments) {
						const cUid = comment['user-id']
						if (this.banners.includes(cUid)){
							this.profilePicsComment.push("unknown");
							console.log("getCommentProfilePicture(): banned: ");
						} else {
							const response = await this.$axios.get(`/users/"${cUid}/profile-picture`, {responseType: "blob"});
							const photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
							this.profilePicsComment.push(photo);
							console.log("getCommentProfilePicture(): \nphoto: "+photo);
						}

						console.log("getCommentProfilePicture(): \nthis.profilePicsComment.push(photo): "+this.profilePicsComment);
					}

					this.profilePicsLike = [];
					for (const lUid of this.likes['user-ids']) {
						console.log("getLikesProfilePicture(): this.likes[user-ids]: "+this.likes['user-ids'])

						if (this.banners.includes(lUid)){
							this.profilePicsLike.push("unknown");
							console.log("getLikesProfilePicture(): banned: ");
						} else {
							const response = await this.$axios.get(`/users/"${lUid}/profile-picture`, {responseType: "blob"});
							const photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
							this.profilePicsLike.push(photo);
							console.log("getLikesProfilePicture(): \nphoto: "+photo);
						}

						console.log("getLikesProfilePicture(): \nthis.profilePicsComment.push(photo): "+this.profilePicsLike);
					}

				} else {
					const response = await this.$axios.get(`/users/"${this.user['user-id']}/profile-picture`, {responseType: "blob"})
					this.uProfPic = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
					
				}
			} catch (e) {
				this.errormsg = "---getLikesProfilePicture(): --- "+e.toString();
			} finally {
				this.loading = false;
			}
		},

		async getPhotoComments() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/posts/" + this.postid + "/comments/");
				this.comments = response.data;

				// refresh to get updated commenters and likers profPics
				await this.getBanUsers();
				await this.getUserProfilePicture(true);
				await this.getLikeAndCommentProfiles();
				await this.getPhoto();
			} catch (e) {
				this.errormsg = "getPhotoComments() error: "+e.toString();
			} finally {
				this.loading = false;
			}
		},

		async commentPhoto() {
			this.loading = true;
			this.errormsg = null;
			try {
				const formData = new FormData();
        formData.append('message', this.messagetext);

				await this.$axios.post("/posts/"+this.postid+"/comments/", formData);

				this.messagetext = '';

				// refresh...
				await this.getPhotoComments();
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

		async uncommentPhoto(commentid) {
			this.errormsg = null;
			let response = null;
			if (confirm('Are you sure you want to delete this comment?')) {
				this.loading = true;
				try {
					await this.$axios.delete("/comments/" + commentid);
					// refresh...
					await this.getPhotoComments();
				} catch (e) {
					this.errormsg = e.toString();
				} finally {
					this.loading = false;
				}
			}
		},

		async getLikes() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/posts/" + this.postid + "/likes/");
				this.likes = response.data;
				if (this.likes['user-ids'] == null){
					this.likes['user-ids'] = [];
					this.likes['like-count'] = 0;
					this.like = false;
					console.log("getLikes(): likes['user-ids'] is null: "+this.likes['user-ids'])
				} else {
					// Check if user's ID is in likes['user-ids']
					this.like = this.likes['user-ids'].includes(this.userid);
					console.log("getLikes(): likes['user-ids'] is not null. like: "+this.like+"\nlikes['user-ids']: "+this.likes['user-ids'])
				}

				// refresh to get updated commenters and likers profPics
				await this.getBanUsers();
				await this.getUserProfilePicture(true);
				await this.getLikeAndCommentProfiles();

			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

		async likePhoto() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.put("/posts/"+this.postid+"/likes/"+this.userid);

				await this.getLikes();
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

		async unlikePhoto() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/posts/" + this.postid + "/likes/" + this.userid);

				await this.getLikes();
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
        console.log("Post.vue -> getBanUsers(): response.data: "+response.data)

        // Those who banned him
        this.banners = [];
        if (response.data['banners'] != null){
          this.banners = response.data['banners'];
          console.log("Post.vue -> getBanUsers(): banners: "+this.banners);
        }

        // Those whom he banned
        this.banneds = [];
        if (response.data['banneds'] != null){
          this.banneds = response.data['banneds'];
          console.log("Post.vue -> getBanUsers(): banneds: "+this.banneds);
        }
      } catch (e) {
        this.errormsg = "Post.vue -> getBanUsers():\nerror: "+e.toString()+"\nthis.userid: "+this.userid;
      } finally {
				this.loading = false;
			}
    },

		toggleLike() {
			if (this.like) {
				this.unlikePhoto()
			} else {
				this.likePhoto()
			}
			this.like = !this.like
		},

		async openIt(p){
			this.showCaption = false
			this.showHashtag = false
			this.showComment = false
			this.showLikedBy = false

			if (p === "caption"){
				this.showCaption = true;
			} else if (p === "hashtag"){
				this.showHashtag = true;
				this.addhashtags = '';
			} else if (p === "comment"){
				this.showComment = true;
				await this.getPhotoComments();
			} else if (p === "likedby"){
				this.showLikedBy = true;
				await this.getLikes();
			}
		},

		async closeIt(p){
			this.showCaption = false
			this.showHashtag = false
			this.showComment = false
			this.showLikedBy = false

			if (p === "caption"){
				await this.modifyCaption();
			} else if (p === "hashtagA"){
				await this.addHashtag()
				this.addhashtags = '';
			} else if (p === "hashtagD"){
				await this.deleteHashtag()
				this.addhashtags = '';
			} else if (p === "comment"){
				if (this.messagetext != ''){
					await this.commentPhoto();
				}
			// } else if (p === "likedby"){
			// 	await this.getLikes();
			}
		},

  },
  created(){
    // Initialize variables here
    this.userid = localStorage.getItem('userid');

  },
  async mounted(){
    await this.getPhoto();

		this.user = await this.getUserProfile(this.post['user-id']);

    this.singlePhoto = await this.getSinglePhoto(this.post['photo-id']);
	
		await this.getPhotoComments();
		await this.getLikes();

		// commenters and likers profPics
		await this.getUserProfilePicture(true);

		// get uProfPic
		await this.getUserProfilePicture(false);

  },
}
</script>

<template>
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<LoadingSpinner v-if="loading"></LoadingSpinner>
		
		<div class="card m-5" style="border: 4px solid red;" v-if="!loading" >

      <div class="card-header">
				<div class="container mt-2 mb-2">
					<div class="row ">
						<div class="col-md-6">
							<LinkToUserProfile
							:profpic="uProfPic" 
							:userprofname="user['profile-name']"
							:uid="post['user-id']">
							</LinkToUserProfile>
						</div>
						<div class="col-md-3"></div>
						<div class="col-md-3 d-flex justify-content-end" v-if="userid===post['user-id']">
							<button class="btn btn-primary" type="button" @click="deletePhoto" > Delete Post</button>
						</div>
					</div>
				</div>
			</div>

			<div class="card-body">
        <small class="text-muted"> {{ post['date-time']}} </small> <br>

				<!-- Caption -->
				<div class="container m-3">
					<div class="row ">
						<div class="col-md-9">
							<div v-if="caption != ''">
								<div>Caption:</div>
								<div class="card-text"> {{ caption  }} </div>
							</div>
							<div v-if="caption  == ''">
								<div class="card-text"> No Captions </div>
							</div>
						</div>
						<div class="col-md-3 d-flex justify-content-end" v-if="userid===post['user-id'] && showCaption===false">
							<button type="button" class="btn btn-primary m-2 p-2" @click="openIt('caption')">Modify Caption</button>
						</div>
					</div>
				</div>
				<div class="container mt-2 mb-2" v-if="showCaption">
					<div class="row ">
						<div class="col-md-9">
							<textarea class="form-control" style="border: 1px solid rgb(255, 217, 47);" id="caption" v-model="caption" rows="2"
								placeholder="Edit here">
							</textarea>
						</div>
						<div class="col-md-3 d-flex justify-content-end">
							<button type="button" class="btn btn-primary m-2 p-2" @click="closeIt('caption')" >Send</button>
						</div>
					</div>
				</div>
				
				<!-- Hashtag -->
				<div class="container m-3">
					<div class="row ">
						<div class="col-md-6">
        			<div v-if="hashtags != ''">
								<div>Hashtags:</div>
								<div class="card-text"> {{ hashtags }} </div>
							</div>
							<div v-if="hashtags == ''">
								<div class="card-text"> Empty Hashtags </div>
							</div>
						</div>
						<div class="col-md-6 d-flex justify-content-end" v-if="userid===post['user-id'] && showHashtag===false">
							<button type="button" class="btn btn-primary m-2 p-2" @click="openIt('hashtag')">Modify #s</button>
						</div>
					</div>
				</div>
				<div class="container mt-2 mb-2" v-if="showHashtag">
					<div class="row ">
						<div class="col-md-7">
							<textarea class="form-control" style="border: 1px solid rgb(255, 217, 47);" id="hashtag" v-model="addhashtags" rows="2"
								placeholder="Edit here">
							</textarea>
						</div>
						<div class="col-md-5 d-flex justify-content-end">
							<button type="button" class="btn btn-primary m-1" @click="closeIt('hashtagA')" >Add #s</button>
							<button type="button" class="btn btn-danger m-1" @click="closeIt('hashtagD')" >Del #s</button>
						</div>
					</div>
				</div>
			
        <div class="container d-flex justify-content-center m-1">
        	<img class="card-img-top user-icon" :src="singlePhoto" alt="oops, can't load photo!" />
        </div>
      </div>

			<div class="card-footer">
				<div class="container mt-3 mb-3">
					<div class="row justify-content-center">
						<div class="col-md-6">
							<button type="button" class="btn btn-primary me-3" @click="openIt('likedby')">Liked by</button>
							<button class="btn btn-primary" type="button" @click="toggleLike()" :key="like">{{ !like ? 'Like':'Unlike' }} {{ likes['like-count'] }}</button>
						</div>
						<div class="col-md-6">
							<button type="button" class="btn btn-primary" @click="openIt('comment')">Comment
								{{ post['comment-count'] }}
							</button>
						</div>
					</div>
				</div>
      </div>
			
		</div>

		<!-- Liked By -->
		<div class="container mt-3 mb-2" style="border: 2px solid violet;" v-if="showLikedBy && likes['like-count'] != 0">
			<div class="card m-2 p-2">
				<div class="card-header m-2 p-1">Users who liked the post</div>
				<div v-for="(likeUid, idx) in likes['user-ids']" :key="likeUid">
					<div class="card-header">
						<LinkToUserProfile v-if="!banners.includes(likeUid)"
							:profpic="profilePicsLike[idx]" 
							:userprofname="userLikes[idx]"
							:uid="likeUid">
						</LinkToUserProfile>

						<NoLinkToUserProfile v-if="banners.includes(likeUid)" ></NoLinkToUserProfile>
					</div>
				</div>
			</div>
		</div>

		<!-- Comments -->
		<div class="container mt-3 mb-2" style="border: 2px solid blue;" v-if="showComment">

			<div class="card m-1 p-1">

				<div class="card-header m-2 p-1">Comments</div>

				<div class="container mt-3 mb-3">
					<div class="row ">
						<div class="col-md-8">
							<textarea class="form-control" style="border: 1px solid rgb(255, 217, 47);" id="messagetext" v-model="messagetext" rows="2"
								placeholder="Edit here"></textarea>
						</div>
						<div class="col-md-2 d-flex justify-content-end p-1">
							<button type="button" class="btn btn-primary" @click="commentPhoto()" >Send</button>
						</div>
						<div class="col-md-2 d-flex justify-content-end p-1">
							<button type="button" class="btn btn-danger" @click="closeIt('comment')">Close</button>
						</div>
					</div>
				</div>

				<div v-for="(comment, idx) in comments" :key="comment['comment-id']">
					<div class="card-header mt-1 p-1">
						<div class="row">
							<div class="col-md-9">
								<LinkToUserProfile v-if="!banners.includes(comment['user-id'])"
									:profpic="profilePicsComment[idx]" 
									:userprofname="userComments[idx]"
									:uid="comment['user-id']">
								</LinkToUserProfile>

								<NoLinkToUserProfile v-if="banners.includes(comment['user-id'])" ></NoLinkToUserProfile>
							</div>
							<div class="col-md-3 justify-content-end" v-if="userid===comment['user-id'] || userid==post['user-id']">
								<button @click="uncommentPhoto(comment['comment-id'])">Delete</button>
							</div>
						</div>
					</div>
					<div class="card-body m-1 p-1">
						<small class="text-muted"> {{ comment['date-time']}} </small> <br>
						<div class="card-text"> {{ comment['message'] }} </div>
					</div>
				</div>
				<div class="d-flex justify-content-end m-1">
					<button type="button" class="btn btn-danger" @click="closeIt('comment')">Close</button>
				</div>
			</div>
		</div>

	</div>
</template>
