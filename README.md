# Ripple
a lightweight web framework for Go(base on [Echo](https://github.com/labstack/echo))

## Which features to include in a framework

A framework has 3 parts. A router receiving a request and directing it to a handler, a middleware system to add reusable pieces of software before and after the handler, and the handler processing the request and writing the response.
- Router
- Middlewares
- Handler processing

<img src="https://raw.githubusercontent.com/bmbstack/ripple/master/screenshots/framework.png" width="720" height="576" />

Middlewares handle:
- error/panic
- logging
- security
- sessions
- cookies
- body parsing

## Features
* [x] MySQL and Foundation database support
* [x] Modular (you can choose which components to use)
* [x] Middleware support, compatible Middleware works out of the box
* [x] Lightweight
* [x] Multiple configuration files support (currently ```.env```)

## Overview
`ripple` is a lightweight framework. It is based on the principles of simplicity, relevance and elegance.
- Simplicity. The design is simple, easy to understand and doesn't introduce many layers between you and the standard library. It is a goal of the project that users should be able to understand the whole framework in a single day.
- Relevance. `ripple` doesn't assume anything. We focus on things that matter, this way we are able to ensure easy maintenance and keep the system well-organized, well-planned and sweet.
- Elegance. `ripple` uses golang best practises. We are not afraid of heights, it's just that we need a parachute in our backpack. The source code is heavily documented, any functionality should be well explained and well tested.

## Project structure
This is the structure of the `rippleApp` list application that will showcase how you can build web apps with `ripple`:

```shell
.
├── config.json.example
├── controllers
│   ├── router.go
│   └── v1
│       ├── base.go
│       ├── citys.go
│       └── home.go
├── forum
├── frontend
│   ├── static
│   │   ├── css
│   │   └── js
│   └── templates
│       └── home
├── helper
│   └── helper.go
├── logger
│   └── logger.go
├── main.go
├── models
│   ├── one
│   │   ├── base_model.go
│   │   ├── bmb_city.go
│   │   └── user.go
│   └── two
│       ├── base_model.go
│       ├── bmb_city.go
│       └── user.go
└── scripts
    ├── commands.go
    ├── const.go
    ├── init.go
    └── server.go

12 directories, 19 files

```

## Configurations

ripple support .env configurations files. In our rippleApp app, we put the configuration files in the current project directory. I have included all three formats for clarity, you can be just fine with either one of them.

ripple searches for a file named .env in the config directory. The first to be found is the one to be used.

This is the content of ```.env.example``` file:

```shell
{
  "debugOn": true,
  "domain": "127.0.0.1:8090",
  "static": "frontend/static",
  "templates": "frontend/templates",
  "databases": [
    {
      "alias": "one",
      "dialect": "mysql",
      "host": "127.0.0.1",
      "port": 3306,
      "name": "one",
      "username": "root",
      "password": "123456"
    },
    {
      "alias": "two",
      "dialect": "mysql",
      "host": "127.0.0.1",
      "port": 3306,
      "name": "two",
      "username": "root",
      "password": "123456"
    }
  ],
  "caches": [
    {
      "alias": "one",
      "section": "one",
      "adapter": "redis",
      "host": "127.0.0.1",
      "port": 6379,
      "password": "123456"
    },
    {
      "alias": "two",
      "section": "two",
      "adapter": "redis",
      "host": "127.0.0.1",
      "port": 6379,
      "password": "123456"
    }
  ]
}

```
You can override the values from the config file by setting environment variables.
The names of the environment variables are shown below (with their details)

setting               | details
----------------------|----------------
debugOn               | you can open application debug mode when developing
domain                | the vps ip address or domain
static                | Static serves static files from the directory
templates             | directory to look for templates
databases             | the database list
alias                 | the database alias
dialect               | the database dialect
host                  | the database Host
port                  | the database port
name                  | the database name
user                  | the database user
password              | the database password
caches                | the cache list
alias                 | the cache alias
section               | the cache section
adapter               | the cache adapter
host                  | the cache host
port                  | the cache port
password              | the cache password


## Models
ripple uses the [gorm](https://gorm.io/gorm) library as its Object Relational Mapper, so you won't need to learn anything fancy. In our rippleApp app, we need to define a User model that will be used to store our todo details.

In the file models/user.go we define our rippleApp model like this
```go
package models

import (
	"github.com/bmbstack/ripple"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login       string `sql:"size:255;not null"`
	Password    string `sql:"size:255;not null"`
	Email       string `sql:"size:255"`
	Avatar      string `sql:"size:255"`
	Github      string `gorm:"column:github"`
	QQ          string `gorm:"column:qq"`
	Weibo       string `gorm:"column:weibo"`
	Weixin      string `gorm:"column:weixin"`
	HomePage    string `gorm:"column:home_page"`
	Tagline     string
	Description string
	Location    string
}

func init() {
	ripple.RegisterModels(&User{})
}

```
Notice that we need to register our model by calling ripple.RegisterModels(&User{}) in the init function otherwise ripple won't be aware of the model.

ripple will automatically create the table users if it doesn't exist.

Don't be confused by the schema tag, I just added them since we will use the schema package to decode form values(this has nothing to do with ripple, you can use whatever form library you fancy.)

## Controllers

ripple controllers are structs that implement the Controller interface. To help make ripple usable, Structs must implement the Controller interface.

```go
type Controller interface {
    Path() string
}
```
This creates a new Echo group at the Controller#Path, in our example /posts, with all the defined actions.
```shell
 GET /posts     => #ActionIndex
 GET /posts/:id => #ActionShow
```

Our rippleApp Controller is in the controllers/home.go

```go
package controllers

import (
	"net/http"
	"github.com/bmbstack/ripple"
	"github.com/labstack/echo"
)

type HomeController struct {
	Index echo.HandlerFunc `controller:"GET /"`
	Html echo.HandlerFunc `controller:"GET html"`
	String  echo.HandlerFunc `controller:"GET string"`
}

func init() {
	ripple.RegisterController(&HomeController{})
}

func (this HomeController) Path() string {
	return "/"
}

func (this HomeController) ActionIndex(ctx echo.Context) error {
	ctx.Render(http.StatusOK, "home/index.html", map[string]interface{}{
		"title": "Hello, forum is a Ripple application ",
	})

	return nil
}

func (this HomeController) ActionHtml(ctx echo.Context) error {
	ctx.Render(http.StatusOK, "home/html.html", map[string]interface{}{
		"title": "Hello, this is a html template",
	})

	return nil
}

func (this HomeController) ActionString(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, this is a string")
}

```
Note that we registered our controller by calling ripple.RegisterController($HomeController{}) in the init function so as to make ripple aware of our controller. See Routing section below for more explanation of what the controller is doing.

## Templates
ripple templates are golang templates(use [pongo2](https://github.com/flosch/pongo2)). This is the content of frontend/templates/home/index.html:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title></title>
</head>
<body>
<h3>{{title}}</h3>

<h3>Routes</h3>
<ul>
    <li>
        <a href="http://localhost:8090/html" target="_blank">Html template Page</a>
    </li>
    <li>
        <a href="http://localhost:8090/string" target="_blank">String Page</a>
    </li>
</ul>
</body>
</html>
```

## The server.go file

```go
package scripts

import (
	_ "github.com/bmbstack/ripple/fixtures/forum/controllers"
	_ "github.com/bmbstack/ripple/fixtures/forum/models"
	"github.com/bmbstack/ripple/fixtures/forum/logger"
	"github.com/bmbstack/ripple"
)

// Server commands
func GetServerCommands(db string) []string {
	commands := make([]string, 2)
	switch db {
	case "mysql":
		commands = append(commands, "/usr/local/bin/mysqld_safe &")
		commands = append(commands, "sleep 10s")
	}
	return commands
}

// Run server
func RunServer() {
	logger.Logger.Info("Run server ....")

	ripple.Run()
}

```

## Installation
```shell
$ curl https://glide.sh/get | sh
$ go get github.com/bmbstack/ripple
$ go get github.com/bmbstack/ripple/cmd/ripple
$ ripple new rippleApp
$ cd $GOPATH/src/rippleApp
$ go run main.go server
```
	
Then, Open the url:	[http://127.0.0.1:8090](http://127.0.0.1:8090)
