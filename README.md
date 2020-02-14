# CRUD_API_GOLANG
This project was born because only 1 reason: I'm learning Golang and trying to apply anything to this project.
Summary: A user can create/view/update/delete his tasks. A user can only view his tasks.

## What I have learned
1. CRUD functions with Mysql
2. Middleware to setHeader & Authentication (Using JWT)
3. Apply singleton to create database instance
4. Retrieve params from POST request (without using Bind a struct) 

## API Document
## Group Task
CRUD functions
#### 1. Create a new task
Method: `POST`

Route: `/task`

+ Request (application/json)
    + Headers
        - Authorization: Bearer [token]
    + Body
        - title: "Homework" (string required)
+ Response 200
    + Body
        ```json
        {
            "error_code": 200,
            "data": {
                "id": 1
            }
        }
        ```
+ Response 401
    + Body
        ```json
        {
            "error_code": 401,
            "error_message": "title is required"
        }
        ```
#### 2. Show all tasks
Method: `GET`

Route: `/task`

#### 3. View a task
Method: `GET`

Route: `/task/:id`
#### 4. Edit a task
Method: `PUT`

Route: `/task/:id`

+ Params:
`id` - `int`: Task ID
    
+ Request (application/json)
    + Headers
        - Authorization: Bearer [token]
    + Body
        + title: "Homework" (string)
        + done: 1 (int)
+ Response 200
#### 5. Delete a task
Method: `DELETE`

Route: `/task/:id`

## Group Authenticate
#### Login
Hard code with username/password to get token for testing

Method: `POST`

Route: `/login`

+ Request (application/json)
    + Body
        + username: "chris" (string)
        + password: "123456" (string)
+ Response 200
    + Body
        ```json
      {
          "message": "You were log in!",
          "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiY2hyaXMiLCJleHAiOjE1ODEzMzA1MTAsImp0aSI6IjIifQ.B7IgJBaotq7Y-is0Ba64fW36yzrpcTrwS00wzcWBRzzJDenyelNu23hjLU1H9FKZF30iBH_JycJaAcQW2fch8Q"
      }
      ```

