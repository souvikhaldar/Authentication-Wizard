# Authentication Wizard
**Authentication Wizard** is a a backend system useful for user authentication. It consists the following features:  
1. The user can sign up as a new user. They are expected to see an email in their inbox.  
2. The user can click the link/copy and paste the link, in the email to verify the account.  
3. The user can log in only after their email address is verified.  
4. The user cannot create another account with the same email from an existing account.  
5. A hacker if committing Man-in-the-middle attacks cannot know the plain-text password.  
6. Password is not stored in plaintext.   
7. Email is checked for compliance to standard format.  


## Pre-requisites 
1. [Install Redis](https://redis.io/topics/quickstart)  
2. [Install Golang](https://golang.org/doc/install)  
3. Export the email ID and password of the user who would be sending the verification email. Eg.  
```
export AW_PASSWORD=<password>
export AW_EMAIL=<email>
```  
*NOTE* Gmail need to be allowed access to unsafe app. Refer [here](https://serverfault.com/questions/635139/how-to-fix-send-mail-authorization-failed-534-5-7-14)   



## Running the server
1. Clone this repository. `git clone https://github.com/souvikhaldar/Authentication-Wizard.git`  
2. `cd` into this repository. 
  
3. `go run main.go`  

Server should be now running on port 8192. 


## API Documentation  
**POST Sign up**  
`localhost:8192/signup`  
Make request to this endpoint to sign up the user. On success, it returns a URL which can be copy pasted on the browser for verification. Also, the same URL is mailed to the provided email address. 
```
HEADERS
Content-Type application/json
BODY raw
{
	"email_id": "example@gmail.com",
	"password": "12345678"
}

Example Request
Sign up
curl --location --request POST "localhost:8192/signup" \
  --header "Content-Type: application/json" \
  --data "{
	\"email_id\": \"example@gmail.com\",
	\"password\": \"12345678\"
}"

```   
**GET Verify**  
`localhost:8192/verify?e=example@gmail.com&t=EyrMU`  
This endpoint verifies the authenticity of the user who is trying to sign up. The endpoint is the URL which is returned as response to /signup and body of the mail sent to the user.  
```
PARAMS
example@gmail.com
tEyrMU

Example Request
Verify
curl --location --request GET "localhost:8192/verify?e=example@gmail.com&t=EyrMU"

```  
**POST Log in**  
`localhost:8192/login`  
This endpoint is used to allowing the registered user to log in, if he/she is successfully registered.  
```
HEADERS
Content-Type application/json
BODY raw
{
	"email_id": "example@gmail.com",
	"password": "12345678"
}
```
**NOTE**
Postman API Documentation can be found [here](https://documenter.getpostman.com/view/7875071/SWDze13k?version=latest)  

Video Demonstration can be found [here](https://youtu.be/dkU-kxDsuw4)  



