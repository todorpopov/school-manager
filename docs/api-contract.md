# API Contract — School Manager

**Backend:** Go (stdlib `net/http`)  
**Base URL:** `http://localhost:8080`  
**Date:** 2026-03-28

---

## General Conventions

### Response Envelope

Every response — success or error — is wrapped in the same JSON envelope:

```json
{
  "error": false,
  "message": "Human-readable message",
  "data": { }
}
```

| Field | Type | Description |
|---|---|---|
| `error` | `boolean` | `true` on error, `false` on success |
| `message` | `string` | Human-readable status message |
| `data` | `object \| array \| null` | Payload; `null` on errors or delete operations |

### Authentication

Protected endpoints require a session ID passed as a request header:

```
X-Session-Id: <sessionId>
```

The `sessionId` is returned upon login. There is no Bearer token — the session ID **is** the auth credential.

### Content Type

All request and response bodies are `application/json`.

---

## Error Codes

| HTTP Status | `error` code | Meaning |
|---|---|---|
| `400` | `VALIDATION_ERROR` | Field-level validation failed. `data` contains a map of `field → message` |
| `400` | `REQUEST_VALIDATION_ERROR` | Malformed request (e.g. invalid path param, missing session header) |
| `400` | `UNIQUE_VIOLATION` | Duplicate value violates a unique constraint |
| `400` | `FOREIGN_KEY_VIOLATION` | Referenced entity does not exist |
| `400` | `NOT_NULL_VIOLATION` | Required DB field is null |
| `400` | `ROLES_DO_NOT_EXIST` | One or more provided role names are unknown |
| `401` | `INVALID_CREDENTIALS` | Wrong email or password |
| `401` | `SESSION_EXPIRED` | Session has expired |
| `401` | `UNAUTHORIZED` | Valid session but insufficient role |
| `404` | `NOT_FOUND` | Resource does not exist |
| `500` | `INTERNAL_ERROR` | Unhandled server error |
| `500` | `DATABASE_ERROR` | Generic database error |

**Validation error example:**
```json
{
  "error": true,
  "message": "Validation failed during registration",
  "data": {
    "email": "Invalid email format",
    "password": "Password must be at least 8 characters"
  }
}
```

---

## Roles

The system ships with the following predefined role names (seeded in `data_migration_init_roles.sql`):

| Role name | Description |
|---|---|
| `ADMIN` | Full system access across all schools |
| `DIRECTOR` | Full access within their school |
| `TEACHER` | Manages grades and absences for own students |
| `PARENT` | Read-only access to own children's data |
| `STUDENT` | Read-only access to own data |
| `USER` | Default role assigned on self-registration |

---

## Endpoints

---

### Health Check

#### `GET /`
Public. Returns server status.

**Response `200`**
```json
{
  "error": false,
  "message": "pong",
  "data": null
}
```

---

### Auth

#### `POST /api/auth/register`
Public. Self-registration. The role is automatically set to `USER`.

**Request Body**
```json
{
  "first_name": "Ivan",
  "last_name":  "Petrov",
  "email":      "ivan@example.com",
  "password":   "secret123"
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `first_name` | `string` | ✅ | 1–255 chars |
| `last_name` | `string` | ✅ | 1–255 chars |
| `email` | `string` | ✅ | Valid email format |
| `password` | `string` | ✅ | Meets password policy |

**Response `201`**
```json
{
  "error": false,
  "message": "User registered successfully",
  "data": {
    "sessionId": "abc123"
  }
}
```

---

#### `POST /api/auth/login`
Public. Authenticate an existing user.

**Request Body**
```json
{
  "email":    "ivan@example.com",
  "password": "secret123"
}
```

| Field | Type | Required |
|---|---|---|
| `email` | `string` | ✅ |
| `password` | `string` | ✅ |

**Response `200`**
```json
{
  "error": false,
  "message": "User logged in successfully",
  "data": {
    "sessionId": "abc123"
  }
}
```

---

### Users

> All user endpoints require `X-Session-Id`.  
> Most are restricted to `ADMIN` — exceptions noted below.

#### `POST /api/user`
🔒 `ADMIN`. Create a user with explicit role assignment.

**Request Body**
```json
{
  "first_name": "Maria",
  "last_name":  "Georgieva",
  "email":      "maria@school.bg",
  "password":   "secret123",
  "roles":      ["TEACHER"]
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `first_name` | `string` | ✅ | 1–255 chars |
| `last_name` | `string` | ✅ | 1–255 chars |
| `email` | `string` | ✅ | Valid email |
| `password` | `string` | ✅ | Meets password policy |
| `roles` | `string[]` | ✅ | At least one valid role name |

**Response `201`**
```json
{
  "error": false,
  "message": "User created successfully",
  "data": {
    "user_id":    1,
    "first_name": "Maria",
    "last_name":  "Georgieva",
    "email":      "maria@school.bg",
    "roles":      ["TEACHER"]
  }
}
```

---

#### `GET /api/users`
🔒 Authenticated (any role). Returns all users.

**Response `200`**
```json
{
  "error": false,
  "message": "Users retrieved successfully",
  "data": [
    {
      "user_id":    1,
      "first_name": "Maria",
      "last_name":  "Georgieva",
      "email":      "maria@school.bg",
      "roles":      ["TEACHER"]
    }
  ]
}
```

---

#### `GET /api/user/{user_id}`
🔒 `ADMIN`.

**Path Params**

| Param | Type | Description |
|---|---|---|
| `user_id` | `integer` | User ID |

**Response `200`**
```json
{
  "error": false,
  "message": "User retrieved successfully",
  "data": {
    "user_id":    1,
    "first_name": "Maria",
    "last_name":  "Georgieva",
    "email":      "maria@school.bg",
    "roles":      ["TEACHER"]
  }
}
```

---

#### `GET /api/user/email/{email}`
🔒 `ADMIN`.

**Path Params**

| Param | Type | Description |
|---|---|---|
| `email` | `string` | URL-encoded email address |

**Response `200`** — same shape as `GET /api/user/{user_id}`.

---

#### `PUT /api/user/{user_id}`
🔒 `ADMIN`. Update user profile and roles.

**Request Body**
```json
{
  "first_name": "Maria",
  "last_name":  "Ivanova",
  "email":      "maria.new@school.bg",
  "roles":      ["TEACHER", "DIRECTOR"]
}
```

**Response `200`** — returns updated user object (same shape as create).

---

#### `PUT /api/user/password/{user_id}`
🔒 Authenticated (any role — user can change own password).

**Request Body**
```json
{
  "password": "newSecret456"
}
```

**Response `200`**
```json
{
  "error": false,
  "message": "User password updated successfully",
  "data": null
}
```

---

#### `DELETE /api/user/{user_id}`
🔒 `ADMIN`.

**Response `200`**
```json
{
  "error": false,
  "message": "User deleted successfully",
  "data": null
}
```

---

### Roles

> All role endpoints require `X-Session-Id` and the `ADMIN` role.

#### `POST /api/role`
Create a new role.

**Request Body**
```json
{
  "role_name": "TEACHER"
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `role_name` | `string` | ✅ | Must be a valid role name (uppercase letters, no spaces) |

**Response `201`**
```json
{
  "error": false,
  "message": "Role created successfully",
  "data": {
    "role_id":   3,
    "role_name": "TEACHER"
  }
}
```

---

#### `GET /api/roles`
Returns all roles.

**Response `200`**
```json
{
  "error": false,
  "message": "Roles retrieved successfully",
  "data": [
    { "role_id": 1, "role_name": "ADMIN" },
    { "role_id": 2, "role_name": "DIRECTOR" },
    { "role_id": 3, "role_name": "TEACHER" },
    { "role_id": 4, "role_name": "PARENT" },
    { "role_id": 5, "role_name": "STUDENT" }
  ]
}
```

---

#### `GET /api/role/{role_id}`
Get a single role by ID.

**Response `200`**
```json
{
  "error": false,
  "message": "Role retrieved successfully",
  "data": {
    "role_id":   3,
    "role_name": "TEACHER"
  }
}
```

---

#### `DELETE /api/role/{role_id}`

**Response `200`**
```json
{
  "error": false,
  "message": "Role deleted successfully",
  "data": null
}
```

---

## Endpoints To Be Implemented

The following endpoints are **required by the system specification** but are not yet present in the backend. They are documented here as the agreed contract for future implementation.

---

### Schools

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school` | `ADMIN` | Create a school |
| `GET` | `/api/schools` | `ADMIN`, `DIRECTOR` | List all schools |
| `GET` | `/api/school/{school_id}` | `ADMIN`, `DIRECTOR` | Get school by ID |
| `PUT` | `/api/school/{school_id}` | `ADMIN` | Update school name/address |
| `DELETE` | `/api/school/{school_id}` | `ADMIN` | Delete a school |

**School object**
```json
{
  "school_id": 1,
  "name":      "СУ Васил Левски",
  "address":   "бул. България 1, София"
}
```

---

### Classes

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school/{school_id}/class` | `ADMIN`, `DIRECTOR` | Create a class |
| `GET` | `/api/school/{school_id}/classes` | `ADMIN`, `DIRECTOR`, `TEACHER` | List classes in a school |
| `PUT` | `/api/school/{school_id}/class/{class_id}` | `ADMIN`, `DIRECTOR` | Update class |
| `DELETE` | `/api/school/{school_id}/class/{class_id}` | `ADMIN`, `DIRECTOR` | Delete class |
| `POST` | `/api/class/{class_id}/student/{student_id}` | `ADMIN`, `DIRECTOR` | Enroll student in class |
| `DELETE` | `/api/class/{class_id}/student/{student_id}` | `ADMIN`, `DIRECTOR` | Remove student from class |

**Class object**
```json
{
  "class_id":  1,
  "school_id": 1,
  "name":      "10А",
  "year":      2026
}
```

---

### Director Profile

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school/{school_id}/director` | `ADMIN` | Assign director to school |
| `GET` | `/api/school/{school_id}/director` | `ADMIN`, `DIRECTOR` | Get director for school |
| `PUT` | `/api/school/{school_id}/director/{director_id}` | `ADMIN` | Update director personal data |
| `DELETE` | `/api/school/{school_id}/director/{director_id}` | `ADMIN` | Remove director from school |

**Director profile object**
```json
{
  "director_id": 1,
  "user_id":     5,
  "school_id":   1,
  "first_name":  "Georgi",
  "last_name":   "Stoyanov",
  "email":       "director@school.bg"
}
```

---

### Teacher Profile

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school/{school_id}/teacher` | `ADMIN`, `DIRECTOR` | Add teacher to school |
| `GET` | `/api/school/{school_id}/teachers` | `ADMIN`, `DIRECTOR` | List teachers in school |
| `GET` | `/api/school/{school_id}/teacher/{teacher_id}` | `ADMIN`, `DIRECTOR` | Get teacher details |
| `PUT` | `/api/school/{school_id}/teacher/{teacher_id}` | `ADMIN`, `DIRECTOR` | Update teacher personal data + subjects |
| `DELETE` | `/api/school/{school_id}/teacher/{teacher_id}` | `ADMIN`, `DIRECTOR` | Remove teacher from school |

**Teacher profile object**
```json
{
  "teacher_id": 1,
  "user_id":    7,
  "school_id":  1,
  "first_name": "Ana",
  "last_name":  "Dimitrova",
  "email":      "ana@school.bg",
  "subjects":   ["Математика", "Информатика"]
}
```

---

### Student Profile

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school/{school_id}/student` | `ADMIN`, `DIRECTOR` | Add student to school |
| `GET` | `/api/school/{school_id}/students` | `ADMIN`, `DIRECTOR`, `TEACHER` | List students in school |
| `GET` | `/api/school/{school_id}/student/{student_id}` | `ADMIN`, `DIRECTOR`, `TEACHER`, `PARENT` | Get student details |
| `PUT` | `/api/school/{school_id}/student/{student_id}` | `ADMIN`, `DIRECTOR` | Update student personal data + class |
| `DELETE` | `/api/school/{school_id}/student/{student_id}` | `ADMIN`, `DIRECTOR` | Remove student from school |

**Student profile object**
```json
{
  "student_id": 1,
  "user_id":    10,
  "school_id":  1,
  "first_name": "Petar",
  "last_name":  "Nikolov",
  "email":      "petar@school.bg",
  "class_id":   3,
  "class_name": "10А"
}
```

---

### Parent Profile

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school/{school_id}/parent` | `ADMIN`, `DIRECTOR` | Add parent |
| `GET` | `/api/school/{school_id}/parent/{parent_id}` | `ADMIN`, `DIRECTOR`, `PARENT` | Get parent details |
| `PUT` | `/api/school/{school_id}/parent/{parent_id}` | `ADMIN`, `DIRECTOR` | Update parent data + linked children |
| `DELETE` | `/api/school/{school_id}/parent/{parent_id}` | `ADMIN`, `DIRECTOR` | Remove parent |
| `POST` | `/api/parent/{parent_id}/child/{student_id}` | `ADMIN`, `DIRECTOR` | Link child to parent |
| `DELETE` | `/api/parent/{parent_id}/child/{student_id}` | `ADMIN`, `DIRECTOR` | Unlink child from parent |

**Parent profile object**
```json
{
  "parent_id":  1,
  "user_id":    12,
  "first_name": "Elena",
  "last_name":  "Nikolova",
  "email":      "elena@example.com",
  "children":   [
    { "student_id": 1, "first_name": "Petar", "last_name": "Nikolov" }
  ]
}
```

---

### Curriculum (Учебна програма)

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/school/{school_id}/curriculum` | `ADMIN`, `DIRECTOR` | Create curriculum for a term |
| `GET` | `/api/school/{school_id}/curriculum/{curriculum_id}` | `ADMIN`, `DIRECTOR`, `TEACHER`, `STUDENT`, `PARENT` | Get curriculum |
| `PUT` | `/api/school/{school_id}/curriculum/{curriculum_id}` | `ADMIN`, `DIRECTOR` | Update curriculum |
| `DELETE` | `/api/school/{school_id}/curriculum/{curriculum_id}` | `ADMIN`, `DIRECTOR` | Delete curriculum |
| `POST` | `/api/curriculum/{curriculum_id}/entry` | `ADMIN`, `DIRECTOR` | Add subject+teacher entry to curriculum |
| `DELETE` | `/api/curriculum/{curriculum_id}/entry/{entry_id}` | `ADMIN`, `DIRECTOR` | Remove entry |

**Curriculum object**
```json
{
  "curriculum_id": 1,
  "school_id":     1,
  "class_id":      3,
  "term":          "2025-2026/1",
  "entries": [
    {
      "entry_id":   1,
      "subject":    "Математика",
      "teacher_id": 1,
      "teacher_name": "Ana Dimitrova"
    }
  ]
}
```

---

### Grades

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/grade` | `TEACHER` | Add a grade |
| `GET` | `/api/student/{student_id}/grades` | `ADMIN`, `DIRECTOR`, `TEACHER`, `PARENT`, `STUDENT` | Get grades for a student |
| `PUT` | `/api/grade/{grade_id}` | `TEACHER` | Update a grade |
| `DELETE` | `/api/grade/{grade_id}` | `TEACHER` | Delete a grade |

> Teachers may only manage grades for students they teach (enforced server-side).  
> Parents may only view grades for their own children.

**Grade object**
```json
{
  "grade_id":   1,
  "student_id": 1,
  "teacher_id": 1,
  "subject":    "Математика",
  "value":      5.50,
  "date":       "2026-03-15",
  "note":       "Контролно №2"
}
```

**Create grade request**
```json
{
  "student_id": 1,
  "subject":    "Математика",
  "value":      5.50,
  "date":       "2026-03-15",
  "note":       "Контролно №2"
}
```

---

### Absences

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/absence` | `TEACHER` | Register an absence |
| `GET` | `/api/student/{student_id}/absences` | `ADMIN`, `DIRECTOR`, `TEACHER`, `PARENT`, `STUDENT` | Get absences for a student |
| `PUT` | `/api/absence/{absence_id}` | `TEACHER` | Update absence (e.g. excused) |
| `DELETE` | `/api/absence/{absence_id}` | `TEACHER` | Delete absence |

> Same teacher-scoping and parent-scoping rules apply as for grades.

**Absence object**
```json
{
  "absence_id": 1,
  "student_id": 1,
  "teacher_id": 1,
  "subject":    "Математика",
  "date":       "2026-03-20",
  "excused":    false,
  "note":       ""
}
```

---

### Statistics

> Read-only endpoints. Available to `DIRECTOR` (own school) and `ADMIN` (all schools).

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/school/{school_id}/stats/grades` | `ADMIN`, `DIRECTOR` | Grade distribution for school |
| `GET` | `/api/school/{school_id}/stats/grades/subject/{subject}` | `ADMIN`, `DIRECTOR` | Grade distribution by subject |
| `GET` | `/api/school/{school_id}/stats/grades/teacher/{teacher_id}` | `ADMIN`, `DIRECTOR` | Grade distribution by teacher |
| `GET` | `/api/school/{school_id}/stats/absences` | `ADMIN`, `DIRECTOR` | Absence summary for school |
| `GET` | `/api/stats/grades` | `ADMIN` | Grade distribution across all schools |
| `GET` | `/api/stats/absences` | `ADMIN` | Absence summary across all schools |

**Stats response example**
```json
{
  "error": false,
  "message": "Statistics retrieved successfully",
  "data": {
    "school_id":    1,
    "subject":      "Математика",
    "total_grades": 120,
    "average":      4.75,
    "distribution": {
      "2": 3,
      "3": 10,
      "4": 28,
      "5": 45,
      "6": 34
    }
  }
}
```

---

## Data Models Summary

| Model | Key Fields |
|---|---|
| `User` | `user_id`, `first_name`, `last_name`, `email`, `roles[]` |
| `Role` | `role_id`, `role_name` |
| `School` | `school_id`, `name`, `address` |
| `Class` | `class_id`, `school_id`, `name`, `year` |
| `Director` | `director_id`, `user_id`, `school_id` |
| `Teacher` | `teacher_id`, `user_id`, `school_id`, `subjects[]` |
| `Student` | `student_id`, `user_id`, `school_id`, `class_id` |
| `Parent` | `parent_id`, `user_id`, `children[]` |
| `Curriculum` | `curriculum_id`, `school_id`, `class_id`, `term`, `entries[]` |
| `Grade` | `grade_id`, `student_id`, `teacher_id`, `subject`, `value`, `date` |
| `Absence` | `absence_id`, `student_id`, `teacher_id`, `subject`, `date`, `excused` |

