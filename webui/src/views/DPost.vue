
<script>
export default {
  props: {
    postIds: Array,
  },
  components: {
    // Import and register other components if needed
  },
  data() {
    return {
      posts: [],
    };
  },
  async mounted() {
    // Fetch posts using axios and this.postIds
    // Assign fetched posts to this.posts
  },
};
</script>

<template>
  <div>
    <!-- Render user's posts here -->
    <div v-for="postId in postIds" :key="postId">
      <!-- Render each post using postId -->
    </div>
  </div>
</template>
