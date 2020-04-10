# chefapi_node_auth

This code provides a standin for a real REST interface that checks users for access to manage chef nodes.

# Endpoints
-----------

## GET /auth/NODENAME/user/USERNAME
===========================

### Request
No body is passed

### Return
The body returned looks like this:
````json
{
  user: USERNAME
  node: NODENAME
  auth: true or false
}
````
Values
200 - The user is allowed to modify the node
400 - Invalid request was made

# Links
-------
https://blog.questionable.services/article/testing-http-handlers-go/
https://github.com/quii/testing-gorillas
https://godoc.org/github.com/gorilla/mux#SetURLVars
https://github.com/gorilla/mux
