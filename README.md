# Triple-S: A Simple Storage Service

Triple-S (Simple Storage Service) is a lightweight cloud storage system inspired by Amazon S3. It provides basic functionalities for bucket and object management through a RESTful API, allowing users to create buckets, upload files, and retrieve them.

## Features
- Bucket management (create, list, delete buckets)
- Object operations (upload, retrieve, delete objects)
- RESTful API with XML responses
- Data persistence using CSV for metadata storage
- Configurable server port and storage directory
- Error handling for invalid requests and conflicts

## Installation
Ensure you have [Go](https://go.dev/) installed, then clone this repository:

```sh
$ git clone https://github.com/yourusername/triple-s.git
$ cd triple-s
```

Build the project:

```sh
$ go build -o triple-s .
```

## Usage
Run the server with the desired port and storage directory:

```sh
$ ./triple-s --port 8080 --dir ./data
```

To display help information:

```sh
$ ./triple-s --help
```

## API Endpoints
### Bucket Management
#### Create a Bucket
- **Method:** `PUT`
- **Endpoint:** `/{BucketName}`
- **Response:**
  - `200 OK` on success
  - `400 Bad Request` for invalid names
  - `409 Conflict` if the bucket already exists

#### List Buckets
- **Method:** `GET`
- **Endpoint:** `/`
- **Response:**
  - `200 OK` with an XML list of buckets

#### Delete a Bucket
- **Method:** `DELETE`
- **Endpoint:** `/{BucketName}`
- **Response:**
  - `204 No Content` on success
  - `404 Not Found` if bucket doesn’t exist
  - `409 Conflict` if the bucket is not empty

### Object Operations
#### Upload an Object
- **Method:** `PUT`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Headers:**
  - `Content-Type: <mime-type>`
  - `Content-Length: <size>`
- **Response:**
  - `200 OK` on success
  - `404 Not Found` if the bucket doesn’t exist

#### Retrieve an Object
- **Method:** `GET`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Response:**
  - `200 OK` with file content
  - `404 Not Found` if object doesn’t exist

#### Delete an Object
- **Method:** `DELETE`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Response:**
  - `204 No Content` on success
  - `404 Not Found` if object doesn’t exist

## Requirements
- Go 1.18+
- Only standard Go packages are allowed

## Contribution
Feel free to fork this repository and submit pull requests for improvements!

## License
This project is licensed under the MIT License.

