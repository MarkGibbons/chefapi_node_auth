# chefapi_node_auth

This code provides a standin for a real REST interface that checks users for access to manage chef nodes and join organizations.
See the chefapi_demo_server repository to see how this code was installed and started.

# Endpoints
-----------

        r.HandleFunc("/auth/{user}/node/{node}", authNodeCheck)
        r.HandleFunc("/auth/{user}/org/{org}", authOrgCheck)

## GET /auth/USER/node/NODENAME
==============================

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
200 - The request was valid. Authorization is returned in the JSON body.
400 - Invalid request was made

## GET /auth/USER/org/ORGNAME
============================

### Request
No body is passed

### Return
The body returned looks like this:
````json
{
  user: USERNAME
  org:  ORGNAME
  auth: true or false
}
````
Values
200 - The request was valid. Authorization is returned in the JSON body.
400 - Invalid request was made

# Links
-------
https://blog.questionable.services/article/testing-http-handlers-go/
https://github.com/quii/testing-gorillas
https://godoc.org/github.com/gorilla/mux#SetURLVars
https://github.com/gorilla/mux
# Links
-------
https://blog.questionable.services/article/testing-http-handlers-go/
https://github.com/quii/testing-gorillas
https://godoc.org/github.com/gorilla/mux#SetURLVars
https://github.com/gorilla/mux
