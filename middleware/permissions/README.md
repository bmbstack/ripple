#Permissions [![Build Status](https://travis-ci.org/xyproto/permissions2.svg?branch=master)](https://travis-ci.org/xyproto/permissions2) [![Build Status](https://drone.io/github.com/xyproto/permissions2/status.png)](https://drone.io/github.com/xyproto/permissions2/latest) [![GoDoc](https://godoc.org/github.com/xyproto/permissions2?status.svg)](http://godoc.org/github.com/xyproto/permissions2)

Middleware for keeping track of users, login states and permissions.

Online API Documentation
------------------------

[godoc.org](http://godoc.org/github.com/xyproto/permissions2)


Features and limitations
------------------------

* Uses secure cookies and stores user information in a Redis database. 
* Suitable for running a local Redis server, registering/confirming users and managing public/user/admin pages.
* Also supports connecting to remote Redis servers.
* Does not support SQL databases. For MariaDB/MySQL support, look into [permissionsql](https://github.com/xyproto/permissionsql).
* For Bolt database support, look into [permissionbolt](https://github.com/xyproto/permissionbolt).
* Supports registration and confirmation via generated confirmation codes.
* Tries to keep things simple.
* Only supports *public*, *user* and *admin* permissions out of the box, but offers functionality for implementing more fine grained permissions, if so desired.
* The default permissions can be cleared with the `Clear()` function.
* Supports [Negroni](https://github.com/codegangsta/negroni), [Martini](https://github.com/go-martini/martini), [Gin](https://github.com/gin-gonic/gin), [Macaron](https://github.com/Unknwon/macaron) and [Echo](https://github.com/labstack/echo).
* Should also work with other frameworks, since the standard http.HandlerFunc is used everywhere.

Example for [Negroni](https://github.com/codegangsta/negroni)
--------------------
~~~ go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/xyproto/permissions2"
)

func main() {
	n := negroni.Classic()
	mux := http.NewServeMux()

	// New permissions middleware
	perm := permissions.New()

	// Blank slate, no default permissions
	//perm.Clear()

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Has user bob: %v\n", userstate.HasUser("bob"))
		fmt.Fprintf(w, "Logged in on server: %v\n", userstate.IsLoggedIn("bob"))
		fmt.Fprintf(w, "Is confirmed: %v\n", userstate.IsConfirmed("bob"))
		fmt.Fprintf(w, "Username stored in cookies (or blank): %v\n", userstate.Username(req))
		fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *user rights*: %v\n", userstate.UserRights(req))
		fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *admin rights*: %v\n", userstate.AdminRights(req))
		fmt.Fprintf(w, "\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
		userstate.AddUser("bob", "hunter1", "bob@zombo.com")
		fmt.Fprintf(w, "User bob was created: %v\n", userstate.HasUser("bob"))
	})

	mux.HandleFunc("/confirm", func(w http.ResponseWriter, req *http.Request) {
		userstate.MarkConfirmed("bob")
		fmt.Fprintf(w, "User bob was confirmed: %v\n", userstate.IsConfirmed("bob"))
	})

	mux.HandleFunc("/remove", func(w http.ResponseWriter, req *http.Request) {
		userstate.RemoveUser("bob")
		fmt.Fprintf(w, "User bob was removed: %v\n", !userstate.HasUser("bob"))
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		userstate.Login(w, "bob")
		fmt.Fprintf(w, "bob is now logged in: %v\n", userstate.IsLoggedIn("bob"))
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
		userstate.Logout("bob")
		fmt.Fprintf(w, "bob is now logged out: %v\n", !userstate.IsLoggedIn("bob"))
	})

	mux.HandleFunc("/makeadmin", func(w http.ResponseWriter, req *http.Request) {
		userstate.SetAdminStatus("bob")
		fmt.Fprintf(w, "bob is now administrator: %v\n", userstate.IsAdmin("bob"))
	})

	mux.HandleFunc("/clear", func(w http.ResponseWriter, req *http.Request) {
		userstate.ClearCookie(w)
		fmt.Fprintf(w, "Clearing cookie")
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "user page that only logged in users must see!")
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "super secret information that only logged in administrators must see!\n\n")
		if usernames, err := userstate.AllUsernames(); err == nil {
			fmt.Fprintf(w, "list of all users: "+strings.Join(usernames, ", "))
		}
	})

	// Custom handler for when permissions are denied
	perm.SetDenyFunction(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Permission denied!", http.StatusForbidden)
	})

	// Enable the permissions middleware
	n.Use(perm)

	// Use mux for routing, this goes last
	n.UseHandler(mux)

	// Serve
	n.Run(":3000")
}
~~~

Example for [Martini](https://github.com/go-martini/martini)
--------------------
~~~ go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-martini/martini"
	"github.com/xyproto/permissions2"
)

func main() {
	m := martini.Classic()

	// New permissions middleware
	perm := permissions.New()

	// Blank slate, no default permissions
	//perm.Clear()

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	m.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Has user bob: %v\n", userstate.HasUser("bob"))
		fmt.Fprintf(w, "Logged in on server: %v\n", userstate.IsLoggedIn("bob"))
		fmt.Fprintf(w, "Is confirmed: %v\n", userstate.IsConfirmed("bob"))
		fmt.Fprintf(w, "Username stored in cookies (or blank): %v\n", userstate.Username(req))
		fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *user rights*: %v\n", userstate.UserRights(req))
		fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *admin rights*: %v\n", userstate.AdminRights(req))
		fmt.Fprintf(w, "\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
	})

	m.Get("/register", func(w http.ResponseWriter) {
		userstate.AddUser("bob", "hunter1", "bob@zombo.com")
		fmt.Fprintf(w, "User bob was created: %v\n", userstate.HasUser("bob"))
	})

	m.Get("/confirm", func(w http.ResponseWriter) {
		userstate.MarkConfirmed("bob")
		fmt.Fprintf(w, "User bob was confirmed: %v\n", userstate.IsConfirmed("bob"))
	})

	m.Get("/remove", func(w http.ResponseWriter) {
		userstate.RemoveUser("bob")
		fmt.Fprintf(w, "User bob was removed: %v\n", !userstate.HasUser("bob"))
	})

	m.Get("/login", func(w http.ResponseWriter) {
		userstate.Login(w, "bob")
		fmt.Fprintf(w, "bob is now logged in: %v\n", userstate.IsLoggedIn("bob"))
	})

	m.Get("/logout", func(w http.ResponseWriter) {
		userstate.Logout("bob")
		fmt.Fprintf(w, "bob is now logged out: %v\n", !userstate.IsLoggedIn("bob"))
	})

	m.Get("/makeadmin", func(w http.ResponseWriter) {
		userstate.SetAdminStatus("bob")
		fmt.Fprintf(w, "bob is now administrator: %v\n", userstate.IsAdmin("bob"))
	})

	m.Get("/clear", func(w http.ResponseWriter) {
		userstate.ClearCookie(w)
		fmt.Fprintf(w, "Clearing cookie")
	})

	m.Get("/data", func(w http.ResponseWriter) {
		fmt.Fprintf(w, "user page that only logged in users must see!")
	})

	m.Get("/admin", func(w http.ResponseWriter) {
		fmt.Fprintf(w, "super secret information that only logged in administrators must see!\n\n")
		if usernames, err := userstate.AllUsernames(); err == nil {
			fmt.Fprintf(w, "list of all users: "+strings.Join(usernames, ", "))
		}
	})

	// Set up a middleware handler for Martini, with a custom "permission denied" message.
	permissionHandler := func(w http.ResponseWriter, req *http.Request, c martini.Context) {
		// Check if the user has the right admin/user rights
		if perm.Rejected(w, req) {
			// Deny the request
			http.Error(w, "Permission denied!", http.StatusForbidden)
			// Reject the request by not calling the next handler below
			return
		}
		// Call the next middleware handler
		c.Next()
	}

	// Enable the permissions middleware
	m.Use(permissionHandler)

	// Serve
	m.Run()
}
~~~

Example for [Gin](https://github.com/gin-gonic/gin)
--------------------
~~~ go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyproto/permissions2"
)

func main() {
	g := gin.New()

	// New permissions middleware
	perm := permissions.New()

	// Blank slate, no default permissions
	//perm.Clear()

	// Set up a middleware handler for Gin, with a custom "permission denied" message.
	permissionHandler := func(c *gin.Context) {
		// Check if the user has the right admin/user rights
		if perm.Rejected(c.Writer, c.Request) {
			// Deny the request, don't call other middleware handlers
			c.AbortWithStatus(http.StatusForbidden)
			fmt.Fprint(c.Writer, "Permission denied!")
			return
		}
		// Call the next middleware handler
		c.Next()
	}

	// Logging middleware
	g.Use(gin.Logger())

	// Enable the permissions middleware, must come before recovery
	g.Use(permissionHandler)

	// Recovery middleware
	g.Use(gin.Recovery())

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	g.GET("/", func(c *gin.Context) {
		msg := ""
		msg += fmt.Sprintf("Has user bob: %v\n", userstate.HasUser("bob"))
		msg += fmt.Sprintf("Logged in on server: %v\n", userstate.IsLoggedIn("bob"))
		msg += fmt.Sprintf("Is confirmed: %v\n", userstate.IsConfirmed("bob"))
		msg += fmt.Sprintf("Username stored in cookies (or blank): %v\n", userstate.Username(c.Request))
		msg += fmt.Sprintf("Current user is logged in, has a valid cookie and *user rights*: %v\n", userstate.UserRights(c.Request))
		msg += fmt.Sprintf("Current user is logged in, has a valid cookie and *admin rights*: %v\n", userstate.AdminRights(c.Request))
		msg += fmt.Sprintln("\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
		c.String(http.StatusOK, msg)
	})

	g.GET("/register", func(c *gin.Context) {
		userstate.AddUser("bob", "hunter1", "bob@zombo.com")
		c.String(http.StatusOK, fmt.Sprintf("User bob was created: %v\n", userstate.HasUser("bob")))
	})

	g.GET("/confirm", func(c *gin.Context) {
		userstate.MarkConfirmed("bob")
		c.String(http.StatusOK, fmt.Sprintf("User bob was confirmed: %v\n", userstate.IsConfirmed("bob")))
	})

	g.GET("/remove", func(c *gin.Context) {
		userstate.RemoveUser("bob")
		c.String(http.StatusOK, fmt.Sprintf("User bob was removed: %v\n", !userstate.HasUser("bob")))
	})

	g.GET("/login", func(c *gin.Context) {
		// Headers will be written, for storing a cookie
		userstate.Login(c.Writer, "bob")
		c.String(http.StatusOK, fmt.Sprintf("bob is now logged in: %v\n", userstate.IsLoggedIn("bob")))
	})

	g.GET("/logout", func(c *gin.Context) {
		userstate.Logout("bob")
		c.String(http.StatusOK, fmt.Sprintf("bob is now logged out: %v\n", !userstate.IsLoggedIn("bob")))
	})

	g.GET("/makeadmin", func(c *gin.Context) {
		userstate.SetAdminStatus("bob")
		c.String(http.StatusOK, fmt.Sprintf("bob is now administrator: %v\n", userstate.IsAdmin("bob")))
	})

	g.GET("/clear", func(c *gin.Context) {
		userstate.ClearCookie(c.Writer)
		c.String(http.StatusOK, "Clearing cookie")
	})

	g.GET("/data", func(c *gin.Context) {
		c.String(http.StatusOK, "user page that only logged in users must see!")
	})

	g.GET("/admin", func(c *gin.Context) {
		c.String(http.StatusOK, "super secret information that only logged in administrators must see!\n\n")
		if usernames, err := userstate.AllUsernames(); err == nil {
			c.String(http.StatusOK, "list of all users: "+strings.Join(usernames, ", "))
		}
	})

	// Serve
	g.Run(":3000")
}
~~~

Example for [Macaron](https://github.com/Unknwon/macaron)
--------------------
~~~ go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Unknwon/macaron"
	"github.com/xyproto/permissions2"
)

func main() {
	m := macaron.Classic()

	// New permissions middleware
	perm := permissions.New()

	// Blank slate, no default permissions
	//perm.Clear()

	// Logging middleware
	m.Use(macaron.Logger())

	// Renderer middleware
	m.Use(macaron.Renderer())

	// Set up a middleware handler for Macaron, with a custom "permission denied" message.
	permissionHandler := func(ctx *macaron.Context) {
		// Check if the user has the right admin/user rights
		if perm.Rejected(ctx.Resp, ctx.Req.Request) {
			fmt.Fprintf(ctx.Resp, "Permission denied!")
			// Deny the request
			ctx.Error(http.StatusForbidden)
			// Don't call other middleware handlers
			return
		}
		// Call the next middleware handler
		ctx.Next()
	}

	// Enable the permissions middleware, must come before recovery
	m.Use(permissionHandler)

	// Recovery middleware
	m.Use(macaron.Recovery())

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	m.Get("/", func(ctx *macaron.Context) string {
		msg := ""
		msg += fmt.Sprintf("Has user bob: %v\n", userstate.HasUser("bob"))
		msg += fmt.Sprintf("Logged in on server: %v\n", userstate.IsLoggedIn("bob"))
		msg += fmt.Sprintf("Is confirmed: %v\n", userstate.IsConfirmed("bob"))
		msg += fmt.Sprintf("Username stored in cookies (or blank): %v\n", userstate.Username(ctx.Req.Request))
		msg += fmt.Sprintf("Current user is logged in, has a valid cookie and *user rights*: %v\n", userstate.UserRights(ctx.Req.Request))
		msg += fmt.Sprintf("Current user is logged in, has a valid cookie and *admin rights*: %v\n", userstate.AdminRights(ctx.Req.Request))
		msg += fmt.Sprintln("\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
		return msg
	})

	m.Get("/register", func(ctx *macaron.Context) string {
		userstate.AddUser("bob", "hunter1", "bob@zombo.com")
		return fmt.Sprintf("User bob was created: %v\n", userstate.HasUser("bob"))
	})

	m.Get("/confirm", func(ctx *macaron.Context) string {
		userstate.MarkConfirmed("bob")
		return fmt.Sprintf("User bob was confirmed: %v\n", userstate.IsConfirmed("bob"))
	})

	m.Get("/remove", func(ctx *macaron.Context) string {
		userstate.RemoveUser("bob")
		return fmt.Sprintf("User bob was removed: %v\n", !userstate.HasUser("bob"))
	})

	m.Get("/login", func(ctx *macaron.Context) string {
		// Headers will be written, for storing a cookie
		userstate.Login(ctx.Resp, "bob")
		return fmt.Sprintf("bob is now logged in: %v\n", userstate.IsLoggedIn("bob"))
	})

	m.Get("/logout", func(ctx *macaron.Context) string {
		userstate.Logout("bob")
		return fmt.Sprintf("bob is now logged out: %v\n", !userstate.IsLoggedIn("bob"))
	})

	m.Get("/makeadmin", func(ctx *macaron.Context) string {
		userstate.SetAdminStatus("bob")
		return fmt.Sprintf("bob is now administrator: %v\n", userstate.IsAdmin("bob"))
	})

	m.Get("/clear", func(ctx *macaron.Context) string {
		userstate.ClearCookie(ctx.Resp)
		return "Clearing cookie"
	})

	m.Get("/data", func(ctx *macaron.Context) string {
		return "user page that only logged in users must see!"
	})

	m.Get("/admin", func(ctx *macaron.Context) {
		fmt.Fprintf(ctx.Resp, "super secret information that only logged in administrators must see!\n\n")
		if usernames, err := userstate.AllUsernames(); err == nil {
			fmt.Fprintf(ctx.Resp, "list of all users: "+strings.Join(usernames, ", "))
		}
	})

	// Serve
	m.Run(3000)
}
~~~

Example for [Echo](https://github.com/labstack/echo)
--------------------
~~~ go
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/xyproto/permissions2"
)

func main() {
	e := echo.New()

	// New permissions middleware
	perm := permissions.New()

	// Blank slate, no default permissions
	//perm.Clear()

	// Set up a middleware handler for Echo, with a custom "permission denied" message.
	permissionHandler := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Check if the user has the right admin/user rights
			if perm.Rejected(c.Response().Writer(), c.Request()) {
				// Deny the request
				return echo.NewHTTPError(http.StatusForbidden, "Permission denied!")
			}
			// Continue the chain of middleware
			return next(c)
		}
	}

	// Logging middleware
	e.Use(middleware.Logger())

	// Enable the permissions middleware, must come before recovery
	e.Use(permissionHandler)

	// Recovery middleware
	e.Use(middleware.Recover())

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	e.Get("/", func(c *echo.Context) error {
		var buf bytes.Buffer
		b2s := map[bool]string{false: "false", true: "true"}
		buf.WriteString("Has user bob: " + b2s[userstate.HasUser("bob")] + "\n")
		buf.WriteString("Logged in on server: " + b2s[userstate.IsLoggedIn("bob")] + "\n")
		buf.WriteString("Is confirmed: " + b2s[userstate.IsConfirmed("bob")] + "\n")
		buf.WriteString("Username stored in cookies (or blank): " + userstate.Username(c.Request()) + "\n")
		buf.WriteString("Current user is logged in, has a valid cookie and *user rights*: " + b2s[userstate.UserRights(c.Request())] + "\n")
		buf.WriteString("Current user is logged in, has a valid cookie and *admin rights*: " + b2s[userstate.AdminRights(c.Request())] + "\n")
		buf.WriteString("\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
		return c.String(http.StatusOK, buf.String())
	})

	e.Get("/register", func(c *echo.Context) error {
		userstate.AddUser("bob", "hunter1", "bob@zombo.com")
		return c.String(http.StatusOK, fmt.Sprintf("User bob was created: %v\n", userstate.HasUser("bob")))
	})

	e.Get("/confirm", func(c *echo.Context) error {
		userstate.MarkConfirmed("bob")
		return c.String(http.StatusOK, fmt.Sprintf("User bob was confirmed: %v\n", userstate.IsConfirmed("bob")))
	})

	e.Get("/remove", func(c *echo.Context) error {
		userstate.RemoveUser("bob")
		return c.String(http.StatusOK, fmt.Sprintf("User bob was removed: %v\n", !userstate.HasUser("bob")))
	})

	e.Get("/login", func(c *echo.Context) error {
		// Headers will be written, for storing a cookie
		userstate.Login(c.Response().Writer(), "bob")
		return c.String(http.StatusOK, fmt.Sprintf("bob is now logged in: %v\n", userstate.IsLoggedIn("bob")))
	})

	e.Get("/logout", func(c *echo.Context) error {
		userstate.Logout("bob")
		return c.String(http.StatusOK, fmt.Sprintf("bob is now logged out: %v\n", !userstate.IsLoggedIn("bob")))
	})

	e.Get("/makeadmin", func(c *echo.Context) error {
		userstate.SetAdminStatus("bob")
		return c.String(http.StatusOK, fmt.Sprintf("bob is now administrator: %v\n", userstate.IsAdmin("bob")))
	})

	e.Get("/clear", func(c *echo.Context) error {
		userstate.ClearCookie(c.Response().Writer())
		return c.String(http.StatusOK, "Clearing cookie")
	})

	e.Get("/data", func(c *echo.Context) error {
		return c.String(http.StatusOK, "user page that only logged in users must see!")
	})

	e.Get("/admin", func(c *echo.Context) error {
		var buf bytes.Buffer
		buf.WriteString("super secret information that only logged in administrators must see!\n\n")
		if usernames, err := userstate.AllUsernames(); err == nil {
			buf.WriteString("list of all users: " + strings.Join(usernames, ", "))
		}
		return c.String(http.StatusOK, buf.String())
	})

	// Serve
	e.Run(":3000")
}
~~~


Default permissions
-------------------

* The */admin* path prefix has admin rights by default.
* These path prefixes have user rights by default: */repo* and */data*
* These path prefixes are public by default: */*, */login*, */register*, */style*, */img*, */js*, */favicon.ico*, */robots.txt* and */sitemap_index.xml*

The default permissions can be cleared with the `Clear()` function.


Password hashing
----------------

* bcrypt is used by default for hashing passwords. sha256 is also supported.
* By default, all new password will be hashed with bcrypt.
* For backwards compatibility, old password hashes with the length of a sha256 hash will be checked with sha256. To disable this behavior, and only ever use bcrypt, add this line: `userstate.SetPasswordAlgo("bcrypt")`


Coding style
------------

* The code shall always be formatted with `go fmt`.


General information
-------------------

* Version: 2.2
* License: MIT
* Alexander F Rødseth

