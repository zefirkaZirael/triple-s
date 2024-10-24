Simple Storage Service (Triple-S)
Overview
The Simple Storage Service (Triple-S) is a RESTful API designed for managing storage containers known as "buckets," following the Amazon S3 paradigm. It allows users to create, list, and delete buckets, as well as upload, retrieve, and delete objects within those buckets. The service adheres to S3 specifications and responds in XML format.

Features
Bucket Management
Create a Bucket

Endpoint: PUT /{BucketName}
Functionality: Validates the bucket name and creates a new bucket if the name is unique and meets the naming conventions.
Response: Returns a 200 OK status with bucket details or appropriate error messages (400 Bad Request, 409 Conflict).
List All Buckets

Endpoint: GET /
Functionality: Retrieves and lists all existing buckets with their metadata.
Response: Returns a 200 OK status with an XML response of bucket information.
Delete a Bucket

Endpoint: DELETE /{BucketName}
Functionality: Checks if a bucket exists and is empty, then deletes it.
Response: Returns a 204 No Content status or relevant error messages (404 Not Found, 409 Conflict).
Object Operations
Upload a New Object

Endpoint: PUT /{BucketName}/{ObjectKey}
Functionality: Uploads a file to a specified bucket, overwriting any existing file with the same key.
Response: Returns a 200 OK status or an error message if the upload fails.
Retrieve an Object

Endpoint: GET /{BucketName}/{ObjectKey}
Functionality: Retrieves the specified object from a bucket.
Response: Returns the object data with appropriate headers or an error message if the object does not exist.
Delete an Object

Endpoint: DELETE /{BucketName}/{ObjectKey}
Functionality: Deletes a specified object from a bucket.
Response: Returns a 204 No Content status or an error message if the object does not exist.