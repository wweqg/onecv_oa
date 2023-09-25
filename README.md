# OneCV OA Backend API

This repository contains the backend API for the OneCV OA. It provides a set of RESTful endpoints for managing teachers, students, and their relationships, as well as performing various operations related to students and teachers.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

Before running the API, ensure you have docker installed.

## Getting Started

1. Clone this repository to your local machine:

```bash
git clone https://github.com/wweqg/onecv_oa.git
cd onecv_oa
```

2. Create a .env file in the project root directory with the following database configuration, or u can use the `.env example` provided in the root folder:

```bash
DB_USERNAME=your_database_username
DB_PASSWORD=your_database_password
DB_HOST=your_database_host
DB_PORT=your_database_port
DB_NAME=onecv_oa
```

3. Run the following command to start the backend service, the API will be accessible at http://localhost:3000:

```bash
docker-compose up
```

## API Endpoints

This service provides the following endpoints:


### List all teachers
- **URL**: `/api/teachers`
- **Method**: `GET`
- **Description**: List all teachers.
- **Response**:
  - 200 OK on success

### Create a new teacher
- **URL**: `/api/teachers`
- **Method**: `POST`
- **Description**: Create a new teacher.
- **Request Body**:
  ```json
  {
    "email": "john_doe@example.com"
  }
  ```
- **Response**:
  - 201 Created on success

### Delete a teacher
- **URL**: `/api/teachers/:email`
- **Method**: `DELETE`
- **Description**: Delete an existing teacher.
- **Response**:
  - 204 on success

### List all students
- **URL**: `/api/students`
- **Method**: `GET`
- **Description**: List all students.
- **Response**:
  - 200 OK on success

### Create a new student
- **URL**: `/api/students`
- **Method**: `POST`
- **Description**: Create a new student.
- **Request Body**:
  ```json
  {
    "email": "john_doe@example.com"
  }
  ```
- **Response**:
  - 201 Created on success

### Delete a student
- **URL**: `/api/students/:email`
- **Method**: `DELETE`
- **Description**: Delete an existing student.
- **Response**:
  - 204 on success

### List all teacher-student links
- **URL**: `/api/teachers_students`
- **Method**: `GET`
- **Description**: List all teacher-student links.
- **Response**:
  - 200 OK on success

### Register students to a teacher
- **URL**: `/api/students`
- **Method**: `POST`
- **Description**: Register students to a teacher.
- **Request Body**:
  ```json
  {
    "teacher": "john_doe@example.com",
    "students": [
        "studentjon@gmail.com",
        "studenthon@gmail.com"
    ]
  }
  ```
- **Response**:
  - 204 on success

### Get common students for a list of teachers
- **URL**: `/api/commonstudents`
- **Method**: `GET`
- **Description**: Get common students for a list of teachers.
- **Response**:
  - 200 OK on success
- **Example request**:
  - /api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com

### Suspend a student
- **URL**: `/api/suspend`
- **Method**: `POST`
- **Description**: Suspend a student.
- **Request Body**:
  ```json
    {
    "student" : "studentmary@gmail.com"
    }
  ```
- **Response**:
  - 204 on success

### Retrieve students for notifications
- **URL**: `/api/retrievefornotifications`
- **Method**: `POST`
- **Description**: Retrieve students for notifications.
- **Request Body**:
  ```json
    {
    "teacher":  "teacherken@gmail.com",
    "notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
    }
  ```
- **Response**:
  - 200 on success
- **Details**:
    To receive notifications from e.g. 'teacherken@gmail.com', a student:
    •	MUST NOT be suspended,
    •	AND MUST fulfill AT LEAST ONE of the following:
    i.	is registered with “teacherken@gmail.com"
    ii.	has been @mentioned in the notification

Please refer to the API handlers in the source code for more details.