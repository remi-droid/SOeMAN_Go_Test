# Document management Golang API

## Prerequisites

This project provides a web service for managing document uploads and downloads using **MinIO** for object storage and **PostgreSQL** for database storage. The service is built with **Go** using the **Gin** framework and integrates seamlessly with Docker for deployment.

## Endpoints

### 1. **List Documents**
- **Description**: Retrieve a list of all stored documents.
- **URL**: `GET /list`
- **Response**:
  - `200 OK`: List of documents.
  - `500 Internal Server Error`: Error retrieving documents.
- **Example Response**:
  ```json
  {
    "documents": [
      {
        "id": 1,
        "name": "example.pdf",
        "url": "http://localhost/dl/example.pdf",
        "uploaded_at": "2024-12-01T10:00:00Z"
      }
    ]
  }
  ```

---

### 2. **Upload Document**
- **Description**: Upload a document to MinIO storage and store metadata in PostgreSQL. **The document must have a name different than all the other documents present in the database.**
- **URL**: `POST /ul`
- **Request Body**: Multipart form with a field named `document` containing the file.
- **Response**:
  - `200 OK`: Document uploaded successfully.
  - `400 Bad Request`: No document provided.
  - `500 Internal Server Error`: Error during upload or database insertion.
- **Example Response**:
  ```json
  {
    "message": "Document example.pdf uploaded successfully",
    "file_name": "example.pdf"
  }
  ```

---

### 3. **Download Document**
- **Description**: Download a specific document from MinIO storage from its filename.
- **URL**: `GET /dl/:filename`
- **Response**:
  - `200 OK`: The document is returned in the response.
  - `500 Internal Server Error`: Error retrieving the document or it doesn't exist.

---

### 4. **Clear All Data**
- **Description**: Delete all documents from PostgreSQL and MinIO storage.
- **URL**: `DELETE /clear`
- **Response**:
  - `200 OK`: All data cleared.
  - `500 Internal Server Error`: Error during deletion.
- **Example Response**:
  ```json
  {
    "message": "Database and bucket cleared successfully! 5 documents deleted."
  }
  ```

---

## MinIO Integration

### Bucket Configuration
- **Bucket Name**: `uploads`
- **MinIO Admin Console**: Available at [http://localhost:9090](http://localhost:9090)
  - **Username**: `root`
  - **Password**: `password`

### Operations
1. **Upload**: Files are stored in the `uploads` bucket.
2. **Download**: Documents are fetched from the bucket using their filename.
3. **Clear Storage**: Deletes all objects from the bucket.

---

## PostgreSQL Database

### Schema
**`Document` Table**:
| Column      | Type    | Constraints            |
|-------------|---------|------------------------|
| `id`        | `INT`   | Primary Key, Auto Increment |
| `name`      | `STRING`| Unique, Not Null       |
| `url`       | `STRING`| Unique, Not Null       |
| `uploaded_at` | `TIME`| Auto Updated           |

### Operations
1. **Insert Document Metadata**: When a document is uploaded, its metadata is stored in this table.
2. **Check Document Presence**: Ensure no duplicate documents are uploaded using the `name` field.
3. **Clear Database**: Delete all entries of the unique table.

## Deployment with Docker

### Services
1. **Upload Service**: Runs the Go application.
2. **PostgreSQL**: Stores document metadata.
3. **MinIO**: Handles object storage.

### Run the Application

Run/recompile the project with:

```shell
docker compose up --build
```

View logs with:

```shell
docker compose logs upload-service
```