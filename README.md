# Captain

Easily start and stop docker compose projects with captain, arrrrr.

<p align="center">
<img src="https://jenssegers.com/uploads/images/captain.png" width="250">
</p>

## Installation

Install `captain` on your machine:

### OSX

```
curl -L https://github.com/jenssegers/captain/releases/download/0.2.0/captain-osx > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

### Linux

```
curl -L https://github.com/jenssegers/captain/releases/download/0.2.0/captain-linux > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

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
