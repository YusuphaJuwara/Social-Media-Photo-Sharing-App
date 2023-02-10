<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">New Post</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>

    <!--form @submit.prevent="submitForm"-->
    <div class="mb-3">
      <label for="picture">Picture:</label>
      <input type="file" id="picture" @change="uploadPicture" ref="pictureInput" accept="image/png"/>
      <div v-if="pictureError" class="error">Please select a PNG file</div>
    </div>
    <div class="mb-3">
      <label for="caption" class="form-label">Caption:</label>
      <textarea id="caption" class="form-control" v-model="form.caption" @mouseover="showHint('caption')"></textarea>
    </div>
    <div class="mb-3">
      <label for="hashtags" class="form-label">Hashtags:</label>
      <textarea id="hashtags" class="form-control" v-model="form.hashtags" @mouseover="showHint('hashtags')"></textarea>
    </div>
    <div>
      <button v-if="!loading" :disabled="!isValid" type="button" class="btn btn-primary" @click="uploadPhoto">
      Create Post
      </button>
    </div>
    <!--/form-->
    <div v-if="hint" class="hint">{{ hint }}</div>
  </div>
</template>

<style>
  .error {
    color: red;
  }

  .hint {
    color: gray;
  }
</style>

<script>
export default {
  data() {
    return {
      token: '',
      userid: '',
      errormsg: null,
			loading: false,
      form: {
        picture: null,
        caption: '',
        hashtags: '',
      },
      pictureError: false,
      hint: '',
    };
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
    uploadPicture(event) {
      const file = event.target.files[0];
      if (file && file.type === 'image/png') {
        this.form.picture = file;
        this.pictureError = false;
      } else {
        this.pictureError = true;
        this.$refs.pictureInput.value = '';
      }
    },
    showHint(field) {
      switch (field) {
        case 'caption':
          this.hint = 'Any char including new lines of at least one alphabet';
          break;
        case 'hashtags':
          this.hint = 'Each hashtag must contain at least one alphabet and zero or more other alphanumeric characters in [a-zA-Z0-9]. And the hashtags are delimited by a comma.';
          break;
      }
    },
    async uploadPhoto() {
      this.loading = true;
      this.errormsg = null;

      try {
        const hashtags = this.form.hashtags.split(',').map(tag => tag.trim());

        const formData = new FormData();
        formData.append('photo', this.form.picture);
        formData.append('caption', this.form.caption);
        formData.append('hashtags', hashtags);

        const response = await axios.post("/users/"+this.userid+"/posts/", formData, {
          headers: {
            'Authorization': 'Bearer ' + this.token,
            'Content-Type': 'multipart/form-data'
          }
        });
        this.form.picture = null;
        this.form.caption = '';
        this.form.hashtags = '';

        this.$router.push("/"+this.userid+"/stream/");
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },

  },
  mounted() {
    this.userid = localStorage.getItem('userid')
    this.token = localStorage.getItem('token')
  },
};
</script>
