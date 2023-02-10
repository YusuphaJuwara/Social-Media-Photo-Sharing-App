<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-lg-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Web and Software Architecture Photo Sharing App (WaSA Photos)</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="`/login`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Login
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink :to="`/${userid}/profile`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Profile Page
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink :to="`/${userid}/stream/`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Stream (Posts of those whom the user follows)
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink :to="`/${userid}/posts/`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								All Posts (who did not ban the user)
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink :to="`/${userid}/users/`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								All Users (who did not ban the user)
							</RouterLink>
						</li>
						<!-- <li class="nav-item">
							<RouterLink :to="`/${postid}/post`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Post
							</RouterLink>
						</li> -->
						<li class="nav-item">
							<RouterLink :to="`/create-post`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Create Post
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink :to="`/search`" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Search user by profilename or post by hashtag
							</RouterLink>
						</li>
					</ul>

					<!-- <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Secondary menu</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
								Item 1
							</RouterLink>
						</li>
					</ul> -->
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>

</template>

