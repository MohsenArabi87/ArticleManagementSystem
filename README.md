# ArticleManagementSystem
A sample web API using go language

To run the server simply type "go run ." command in terminal
The server will start on default address: 0.0.0.0:9090

To register a user POST a request via 0.0.0.0:9090\signup with json body like below
     
        {
            "Email": "mohy66@gmail.com",
            "password": "123"
        }

to sign in with the user also POST a request via 0.0.0.0:9090\login with json body like below
       
        {
            "Email": "mohy66@gmail.com",
            "password": "123"
        }

This webserver uses JSON Web Token (JWT)authentication, so after successfuly logging in you will recieve an access_token and refresh_token related to your username in response.

You should pass the Access Token in each request Authorization Header with "Bearer Token" foramt to the server so as to verify your identity, otherwise you will recieve a message : "Authentication failed. Invalid token" in response.

Note that for security reasons after a specifiec time(default:120min) this access token will expire and you need to request a new access token by calling GET request 0.0.0.0:9090\refresh-token and passing given refresh token as "Bearer Token" in the  Authorization Header

To create an article POST a request via 0.0.0.0:9090\Article\Create with a json body like below

        {
            "Title": "sample title",
            "Content": "This is a sample content",
            "Tags": ["T1","T2"]
        }
in response you will recieve the full spec article :
 
       {
              "ID": "3654f2da-047c-4587-8a84-8646a7f7bee5",
              "title": "sample title",
              "content": "This is a sample content",
              "tags": [
                  "T1",
                  "T2"
              ],
              "Author": "mohy66@gmail.com",
              "createdAt": "2022-04-18T20:47:37.5466761+04:30",
              "updatedAt": "2022-04-18T20:47:37.5466761+04:30"
          }

To update an article POST a request via 0.0.0.0:9090\Article\Update with a json body like below that contains the ID and the author of the Article and other fields which needs to update

       {
            "ID": "3654f2da-047c-4587-8a84-8646a7f7bee5",
            "Title": "updated sample title",
            "Content": "This is an updated sample content",
            "author":"mohy66@gmail.com"
        }
        
note that only the author of the article can uptade or delete an article

 To delete an article call a GET request with the ID of the article via 0.0.0.0:9090\Article\Delete\3654f2da-047c-4587-8a84-8646a7f7bee5 

To fetch all existing tags call a GET request via 0.0.0.0:9090\Article\Tags

To fetch a specific article call a GET request with the ID of the article via 0.0.0.0:9090\Article\cba33430-775e-4c0b-8941-3ab6294df481

To fetch all articles with pagination call a GET request with the page number in querystring via 0.0.0.0:9090\Article?pageid=1

the default page size is 2

you can change default values in configuration

