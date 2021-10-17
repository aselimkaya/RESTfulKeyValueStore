# RESTfulKeyValueStore
Basic in-memory key value store application.

## 1) Installation
You can choose one of the following two ways to run the application:

### Run with Docker
```
git clone https://github.com/aselimkaya/RESTfulKeyValueStore.git
docker build -t keyvalstore .
docker run -p 80:80 -tid keyvalstore
```

### Run with 'go run' command
```
git clone https://github.com/aselimkaya/RESTfulKeyValueStore.git
cd src
go run main.go
```

After successfully running these commands you should see this log in terminal:  ```key-value-store-api Server started successfully at http://localhost:80```

## 2) API Reference

RESTfulKeyValueStore supports three simple operations:
* Adding key value pair
* Getting a key's value
* Flush all key value pairs

### 2.1) Adding Key Value Pair with HTTP POST
To add a key-value pair to the store, the key and value information must be provided as JSON in the body of the request.

* Sample cURL command: ```curl -X POST -d '{"key":"key1","value":"val1"}' localhost/entry```

### Server Responses
* The response from the server if the key value pair is successfully added to the store: ```{"message":"Key value pair added successfully","status":200}```
* If given key already exists: ```{"message":"Key already exists, value will be upated","status":200}```
* If a field left blank: ```{"message":"Missing field! 'key' field is required!","status":400}```

### 2.2) Fetching the Value of a Key
To find the value of a key, the key must be added to the HTTP GET request as a request parameter.
* Sample cURL command: ```curl 'http://localhost/entry?key=key1'```

### Server Responses
* If given key successfully found in store: ```{"message":"{"key":"key1","value":"key1"}","status":200}```
* If not found: ```{"message":"An error occurred while processing the data! Error: key not found","status":400}```

### 2.3) Flushing the Store
To completely delete pairs, an HTTP DELETE request must be sent to the server.
* Sample cURL command: ```curl -X DELETE localhost/entry```

### Server Responses
* If the store flushed successfully: ```{"message":"JSON file flushed successfully!","status":200```