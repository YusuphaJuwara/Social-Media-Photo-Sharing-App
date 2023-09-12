# Wasa Photo Sharing App

This is the [Web and Software Architecture (WaSA)](http://gamificationlab.uniroma1.it/en/wasa/) course project. You need permissions and the wasa virtual machine to build and run this project. The permissions and the wasa vm are aslo needed to build and run docker containers for the project (Scroll down to read how to build and run docker containers).

## Read the project specification
[Project Specs](http://gamificationlab.uniroma1.it/notes/Project.pdf)

### Summary:
A user has:
- Profile page containing:
	- Profile picture
	- Change or delete the profile picture
	- Change user details (profile name, birthdate, gender, profile msg)
	- Modify username (the profile name is, as default, set to the username)
	- Delete account and all the user's data
	- Set profile to public or private
	- All posts of the user in reversed chronological order
	- User follwers, followings, banned users, and their counts
	- Etc

- Post containing:
	- A view component called in "Stream Posts", "All Posts", and "Profile"
	- Photo(required)
	- Hashtags
	- Captions
	- Date and Time posted
	- Likes count and Likers' profile names and profile pics (little, round thumbnail)
	- Comment count and Commenters' profile names and profile pics (little, round thumbnail)
	- Like the post and unlike
	- Comment on the post (even multiple times), and only the commenter or the user on whose post the comment is placed can delete it
	- Etc

- Stream Posts (Feed):
	- This calls the post component
	- It contains only posts of the user himself and those whom he followed (Provided not banned by them)
	- Etc

- All Posts:
	- This also calls the post component
	- It contains :
		- All posts of the user himself
		- Those whom he followed (Provided not banned by them)
		- Those whom he did not follow but did not set their profiles to private (Provided not banned by them)
	- Essentially, this allows the user to see other users' posts even if he didn't follow them.
	- He can click on their profiles and follow them if he wants
	- Etc

All users containing:
	- All the users whom he followed (Provided not banned by them)
	- All the users whom he did not follow but did not set their profiles to private (Provided not banned by them)
	- The user can view other user profiles and follow them, ban them, etc
	- Etc 

Search:
	- Search a user by his profile name (not username -> the project requirements)
	- Search a post by one of its hashtags
	- The search checks if there are user profile names and/or post hashtags that correspond to the serch term
	- It returns two arrays: one for the users if any, the other for the posts if any

Login:
	- A user must login before making any further requests
	- If its the first time, the user is loged in
	- Else, his account is created (read the project equirements about login)
	- In any case, a seesion token and a user id are returned; session id as an "authorization" header, user id as a json object
	- Any subsequent request must also send the session token, else, Unauthorized error
	- It routes the user to the profile page for first time users, else to the Stream profile (Feed)

Log out:
	- This logs out the user if he is currently logged in and rounds the user ro the log in.

## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities
	* `service/api` contains an example of an API server
	* `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is an example of a web frontend in Vue.js; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

Other project files include:
* `open-npm.sh` starts a new (temporary) container using `node:lts` image for safe web frontend development (you don't want to use `npm` in your system, do you?)

## Go vendoring

This project uses [Go Vendoring](https://go.dev/ref/mod#vendoring). You must use `go mod vendor` after changing some dependency (`go get` or `go mod tidy`) and add all files under `vendor/` directory in your commit.

For more information about vendoring:

* https://go.dev/ref/mod#vendoring
* https://www.ardanlabs.com/blog/2020/04/modules-06-vendoring.html

## Node/NPM vendoring

This repository contains the `webui/node_modules` directory with all dependencies for Vue.JS. You should commit the content of that directory and both `package.json` and `package-lock.json`.

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```shell
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run dev
```

## How to build container images

### Backend

```sh
$ docker build -t wasa-photos-backend:latest -f Dockerfile.backend .
```

### Frontend

```sh
$ docker build -t wasa-photos-frontend:latest -f Dockerfile.frontend .
```

## How to run container images

### Backend

```sh
$ docker run -it --rm -p 3000:3000 wasa-photos-backend:latest
```

### Frontend

```
$ docker run -it --rm -p 8081:80 wasa-photos-frontend:latest
```

## Known issues

### Apple M1 / ARM: `failed to load config from`...

If you use Apple M1/M2 hardware, or other ARM CPUs, you may encounter an error message saying that `esbuild` (or some other tool) has been built for another platform.

If so, you can fix issuing these commands **only the first time**:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm install
exit
# Now you can continue as indicated in "How to build/run"
```

**Use these instructions only if you get an error. Do not use it if your build is OK**.

## License

See [LICENSE](LICENSE).

## Use golangci-lint and go fmt to format code for best practice

E.g., ```golangci-lint run -E go fumpt```
enables fumpt with the -E tag
Use the full part to run this. E.g., 
/Social-Media-Photo-Sharing-App$ ```/home/wasa/go/bin/golangci-lint run -E revive```

E.g., ```go fmt ./...``` or ```gofmt -w .```
reformats code
