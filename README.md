# Captain

Easily start and stop docker compose projects with captain, arrrrr.

<p align="center">
<img src="https://jenssegers.com/uploads/images/captain.png" width="250">
</p>

## Installation (OSX)

Install `captain` on your machine with:

```
curl https://raw.githubusercontent.com/jenssegers/captain/master/captain > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

## Usage

Captain searches for docker-compose projects in your `$HOME` folder and allows you to start and stop those projects by passing a part of the parent directory name.

<p align="center">
<img src="https://jenssegers.com/uploads/images/captain.gif">
</p>

### Starting a project

If I have a folder called `my-secret-project` that contains a `docker-compose.yml` file, I can start that project using:

```
captain start my-secret-project
```

Captain will also do partial matching of the project name, so that you can also use:

```
captain start secret
```

### Stopping a project

Stopping a project works similarly:

```
captain stop secret
```

### Listing projects

You can see all managable projects using:

```
captain list
```

### Updating captain

Update your local captain version with:

```
captain update
```
