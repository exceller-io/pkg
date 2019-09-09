# http
HTTP package contains common code when implementing http servers in Go

## Using this package
If you are using Go Modules you can just reference the package directly in your code and and run go mod download or tidy, this would update your go.mod file 

## Http package
This package contains common methods helpful when implementing http web servers using Go

### Starting the server
http package uses gorilla mux router which is a widely used HTTP router in the go community. You can simply setup your route collection and instantiate a new instance of the server and call Start method to start the server. See a simple example below


```
package main

import (
    "net/http"
     ws "github.com/appsbyram/pkg/http"

)

func main() {
    routes := ws.Routes{
        ws.Route{
            "Home",
            "GET",
            "/".
            homeHandler,
        },
    }

    srv := ws.NewServer("8080", false, "", "", routes)
    srv.Start()
}

func homeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //Handle 
    }
}
```

### Reading from Request and Response Payloads Writing to Response
If you end up writing microservices in Go you typically find yourself having to read some data primarily JSON, Yaml and XML from request and Unmarshall into a struct, another use case here is you might call another service and that service sends back from response which you need to read and Unmarshall into a Go struct. Http package contains an interface called Payload that allows you to do just that without having to write the same code over and over again in various services.

#### Reading from Request
Here is a simple example that demonstrates how to read data from request and unmarshall to a Go Struct. Lets imagine we have a handler function called postHandler for one of the routes that our server supports. 

```
type SomeModel struct {
    //fields
}

func postHandler() http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        var model SomeModel 

        p := ws.NewPayload()
        err := p.ReadRequest(ws.ContentTypeJSON, &model, r)

        //Handle error if required
    }
}
```
#### Reading from Response
Here is a simple example that demonstrates how to read data from response and unmarshall to Go struct. Lets imagine we have a handler function call for one of the routes that our server supports which invokes another service and receives a response and we need to Unmarshall the response payload into a Go struct

```
type SomeModel struct {
    //fields
}

func postHandler() http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        var model SomeModel 

        //setup request

        //setup http client

        //make the call
        resp, err := client.Do(req)
        
        //Handle error if required

        //read response
        p := ws.NewPayload()
        err := p.ReadResponse(ws.ContentTypeJSON, &model, resp)
        //Handle error if required
    }
}
```

#### Write to Response
Another common thing you often find yourself especially when implementing APIs you'll need to write some JSON, YAML or XML data to response. Simple example below demonstrates just that. Lets image we have a Get Handler for one of the routes supported by server

```
type SomeModel struct {
    //fields
}

func getHandler() http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        var model SomeModel

        //retrieve data from DB

        //Setup your model

        //Write to response
        p := ws.NewPayload()
        err := p.WriteResponse(ws.ContentTypeJSON, http.StatusOK, &mode, w)
    }
}

```
