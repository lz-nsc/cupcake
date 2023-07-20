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
### To Do
* Configuration. Make it convenient for user to set up the project, include choices for orm, db, or middlewares.
* Controller. Controller should be bound with a resource, and when the countroller is registered to the router, the all the CURD method for this specific resource will be registered.
* A command-line tool for creating new RESTful server project with cupcake framework.
