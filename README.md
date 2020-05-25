# Captain

[![Build Status](https://travis-ci.org/jenssegers/captain.svg?branch=master)](https://travis-ci.org/jenssegers/captain)

Easily start and stop docker compose projects with captain, arrrrr.

<p align="center">
<img src="https://jenssegers.com/static/media/captain.png" width="250">
</p>

## Installation

Binaries can be manually downloaded from GitHub releases: https://github.com/jenssegers/captain/releases

#### OSX

```
curl -L https://github.com/jenssegers/captain/releases/download/0.3.2/captain-osx > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

#### Linux

```
curl -L https://github.com/jenssegers/captain/releases/download/0.3.2/captain-linux > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

#### Windows (untested)

Download `captain.exe` via https://github.com/jenssegers/captain/releases/download/0.3.2/captain.exe

## Usage

Captain searches for docker-compose projects in your `$HOME` folder and allows you to start and stop those projects by matching the project's directory name.

<p align="center">
<img src="https://jenssegers.com/uploads/images/captain.gif?v2">
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

Captain is smart, and does fuzzy matching:

```
captain start scrt
```

### Stopping a project

Stopping a project works similarly:

```
captain stop secret
```

### Restarting a project

Restart a project using:

```
captain restart my-secret-project
```

### Viewing project logs

View logs of a project using:

```
captain logs my-secret-project
```

### Listing projects

You can see all managable projects using:

```
captain list
```

### Stopping all containers

To quickly stop all running docker containers, use:

```
captain abandon
```

### Executing command in a service

Executing command in a running service container

```
captain exec <project> <service> <command>
captain exec my-secret-project web bash
```

Executing command as a new service container

```
captain run <project> <service> <command>
captain exec my-secret-project cli bash
```

## Tweak the running

You can tweak some behaviour through some environmental variables:

* CAPTAIN_ROOT: the starting directory from where the docker-compose are searched (default: home dir of the user)
* CAPTAIN_DEPTH: the number of subdirectory where to search the docker-compose.yml files (default: 5)
