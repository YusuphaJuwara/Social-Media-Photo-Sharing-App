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
</style>

<script>
export default {
	data: function () {
		return {
			token: localStorage.getItem('token'),
			userid: localStorage.getItem('userid'),
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
			messagetext: ''
		}
	},
	methods: {
		// load() {
		// 	return load
		// },
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
    async getMyStream() {
      this.getPhotos();
		},
		async getPhotos() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/posts/");
				this.posts = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
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
				// await this.getMyStream();
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
				// await this.getMyStream();
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
				// await this.getMyStream();
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
					await this.getMyStream();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
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
				// await this.getMyStream();
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
				// await this.getMyStream();
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
				// await this.getMyStream();
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
			return profileName;
		},
		userProfile(userid) {
			this.highlightProfile = !this.highlightProfile
			this.$router.push("/" + userid + "/profile");
			// this.$route.params.id
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
				// await this.getMyStream();
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
				// await this.getMyStream();
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
					// await this.getMyStream();
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
				// await this.getMyStream();
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
			return response.data['user-ids'], response.data['like-count'];
		},
	},
	async likePhoto(postid) {
		// this.loading = true;
		this.errormsg = null;
		try {
			await this.$axios.put("/posts/" + postid + "/likes/" + this.userid);

			// refresh...
			// await this.getMyStream();
		} catch (e) {
			this.errormsg = e.toString();
		}
		// this.loading = false;
	},
	async unlikePhoto(postid) {
		// this.loading = true;
		this.errormsg = null;
		try {
			await this.$axios.delete("/posts/" + postid + "/likes/" + this.userid);

			// refresh...
			// await this.getMyStream();
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
	beforeCreate() {
    // Initialize variables here
    this.userid = localStorage.getItem('userid');
    this.token = localStorage.getItem('token');
  },
	mounted() {
		this.getMyStream()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">{{ getProfileName(userid) }}'s Photo Feed</h1>
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
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="getMyStream">
						Refresh
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newPost">
						Create New Post
					</button>
				</div>
			</div>
		</div>

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
						<img class="imgThumbNail" :src="getUserProfilePicture(post['user-id'])" />
						<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ getProfileName(post['user-id']) }}
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

			<!-- Incomplete-->
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
						{{ !like ? 'Like':'Unlike' }} {{ getLikes(post['post-id'])[1] }}
					</button>
					<button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop4">Comment
						{{ commentCount }}</button>
					<!--button class="btn btn-primary" type="button" @click="showComments=true">
						Comment {{ commentCount }}</button-->
				</div>
			</div>

			<!--Modal v-if="showComments" @close="showComments=false" :comments="comments" :is-post-owner="users.isPostOwner"></Modal-->

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
										<img class="imgThumbNail" :src="getUserProfilePicture(l['user-id'])" />
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

							<!--div class="card" v-for="comment in getPhotoComments(post['post-id'])" :key="comment[comment-id]"-->
							
							<div class="card-header">
								<div class="header-left">
									<a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
										@click="userProfile(comment['user-id'])">
										<img class="imgThumbNail" :src="getUserProfilePicture(comment['user-id'])" />
										<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{
											getProfileName(comment['user-id']) }}</h5>
									</a>
								</div>
								<div class="header-right">
									<button v-if="userid==comment['user-id'] || userid==post['user-id']" data-bs-dismiss="modal"
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
