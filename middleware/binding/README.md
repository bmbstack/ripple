# binding
Request data binding and validation for echo.

## Getting Started


To install Binding:

    go get github.com/echo-contrib/binding

Code:

```go
package main

import (
    b "github.com/echo-contrib/binding"
    "github.com/labstack/echo"
)

type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

// Handler
func binding_test(c *echo.Context) error {
    u := &User{}
    err := c.Bind(u)
    if err != nil {
        return c.JSON(200, err.Error())
    }
    return c.JSON(200, u)
}

func main() {
    // Echo instance
    e := echo.New()
    e.SetBinder(&b.Binder{})

    // Routes
    e.Post("/", binding_test)

    // Start server
    e.Run(":1234")
}
```
