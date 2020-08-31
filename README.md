我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

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
curl -L https://github.com/jenssegers/captain/releases/download/0.3.3/captain-osx > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

#### Linux

```
curl -L https://github.com/jenssegers/captain/releases/download/0.3.3/captain-linux > /usr/local/bin/captain && chmod +x /usr/local/bin/captain
```

#### Windows (untested)

Download `captain.exe` via https://github.com/jenssegers/captain/releases/download/0.3.3/captain.exe

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
