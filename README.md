# cupcake
A lightweight Golang RESTful framework.

## About
While I was useing Django to build RESTful server, I realized that I like the MVC architect and those convenient interfaces offered by this framework.

After I turned to Golang, using go-restful to build the server is quite confusing for me and the usage of other Golang web frameworks are still quite different from Django. So I started this project and try to build up a Django-like golang RESTful framework with MVC architecture from scratch. ;)

## Features
* Supports method-based routing, variables in URL paths, and regexp route patterns based on radix tree implementation
* Group control
* Supports middleware for groups
* Supports static files
* Supports template render

## Getting started

### Prerequisites
* Go
* Sqlite3: Only support Sqlite3 now, the support for other databases will come soon

### Getting Cupcake

With Go module support, simply add the following import to your code:
```
import "github.com/lz-nsc/cupcake"
```

## Usage

### HTTP server
```
cc := cupcake.New()

cc.GET("/cupcake", func(resp *cupcake.Response, req *cupcake.Request) {
		resp.String(http.StatusOK, "Welcome to cupcake!")
	})

cc.GET("/cupcake/{id}", func(resp *cupcake.Response, req *cupcake.Request) {
		id := req.Param("id")
		resp.String(http.StatusOK, fmt.Sprintf("Welcome to cupcake id[%s]!", id))
	})

cc.Run(":8080")
```
### Controller

Cupcake provides `Controller` that allows users to easily create RESTful APIs. 

The core component of any RESTful API is the resource, and the Cupcake controller simplifies the process of registering the CRUD operations for a resource on the router:
```
type User struct {
	Name string `json:"user"`
	Age  int    `json:"age"`
}
type UserController struct {
	*cupcake.BaseController
}

cc := cupcake.New()

controller := UserController{cupcake.NewBaseController(&User{})}
cc.Route("/users", controller)

cc.Run(":8080")
```

In this simple example, we used BaseController, which supports the Created(POST) and Retrieve(GET) methods by default. 

Users can customize the behavior of their controller by defining their own controller methods.

For more examples, please check the Cupcake [examples](https://github.com/lz-nsc/cupcake/tree/master/examples)

### Roadmap
- [ ] Configuration. Make it convenient for user to set up the project, include choices for orm, db, or middlewares.
- [X] Controller. Controller should be bound with a resource, and when the countroller is registered to the router, the all the CURD method for this specific resource will be registered.
- [ ] A command-line tool for creating new RESTful server project with cupcake framework.
- [ ] Serializer