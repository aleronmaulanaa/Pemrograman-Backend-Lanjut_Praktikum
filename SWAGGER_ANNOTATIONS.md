# Dokumentasi Anotasi Swagger - Praktikum 4 CRUD API

## Ringkasan Perubahan
Semua endpoint di proyek telah ditambahkan dengan anotasi Swagger (OpenAPI 3.0) untuk dokumentasi otomatis.

---

## 1. Authentication / Auth Service (`auth_service.go`)

### POST /login
- **Summary**: User Login
- **Description**: Authenticate user and get JWT token
- **Security**: None (endpoint publik)
- **Request**: LoginRequest (username, password)
- **Response**: JWT Token + User Info
- **Status**: 201 (Success), 400 (Invalid input), 401 (Invalid credentials), 500 (Server Error)

---

## 2. Alumni Service (`alumni_service.go`)

### GET /api/alumni
- **Summary**: Get Alumni with Pagination, Sort, and Search
- **Description**: Get list of alumni with pagination, sorting, and search capabilities
- **Security**: Bearer Token Required
- **Query Params**: page, limit, sort, order, search
- **Response**: Array of Alumni with pagination metadata

### GET /api/alumni/{id}
- **Summary**: Get Alumni by ID
- **Description**: Get alumni details by their ID
- **Security**: Bearer Token Required
- **Param**: id (path)
- **Response**: Alumni object

### POST /api/alumni
- **Summary**: Create Alumni
- **Description**: Create a new alumni record (Admin only)
- **Security**: Bearer Token Required
- **Request**: Alumni object
- **Response**: Created Alumni (Status 201)

### PUT /api/alumni/{id}
- **Summary**: Update Alumni
- **Description**: Update an existing alumni record (Admin only)
- **Security**: Bearer Token Required
- **Param**: id (path)
- **Request**: Alumni object
- **Response**: Success message

### DELETE /api/alumni/{id}
- **Summary**: Delete Alumni
- **Description**: Delete an alumni record (Admin only)
- **Security**: Bearer Token Required
- **Param**: id (path)
- **Response**: Success message

### GET /api/alumni/jumlah-angkatan
- **Summary**: Get Jumlah Alumni by Angkatan
- **Description**: Get the count of alumni grouped by year (angkatan)
- **Security**: Bearer Token Required
- **Response**: Array of {angkatan, jumlah}

### GET /api/alumni/jumlah-pekerjaan
- **Summary**: Get Alumni dengan Dua atau Lebih Pekerjaan
- **Description**: Get list of alumni who have two or more jobs
- **Security**: Bearer Token Required
- **Response**: Array of {nama, jumlah_pekerjaan}

---

## 3. Pekerjaan Service - PostgreSQL (`pekerjaan_service.go`)

### GET /api/pekerjaan
- **Summary**: Get All Pekerjaan (PostgreSQL)
- **Description**: Get list of all jobs (non-deleted)
- **Security**: Bearer Token Required
- **Response**: Array of PekerjaanAlumni

### GET /api/pekerjaan/{id}
- **Summary**: Get Pekerjaan by ID (PostgreSQL)
- **Description**: Get job details by ID
- **Security**: Bearer Token Required
- **Response**: PekerjaanAlumni object

### GET /api/pekerjaan/alumni/{alumni_id}
- **Summary**: Get Pekerjaan by Alumni ID (PostgreSQL)
- **Description**: Get jobs of alumni associated with the logged-in user
- **Security**: Bearer Token Required
- **Response**: Array of PekerjaanAlumni

### POST /api/pekerjaan
- **Summary**: Create Pekerjaan (PostgreSQL)
- **Description**: Create a new job record (Admin only)
- **Security**: Bearer Token Required
- **Request**: PekerjaanAlumni object
- **Response**: Created PekerjaanAlumni (Status 201)

### PUT /api/pekerjaan/{id}
- **Summary**: Update Pekerjaan (PostgreSQL)
- **Description**: Update an existing job record (Admin only)
- **Security**: Bearer Token Required
- **Response**: Success message

### DELETE /api/pekerjaan/{id}
- **Summary**: Soft Delete Pekerjaan with RBAC (PostgreSQL)
- **Description**: Soft delete a job record (Admin and owner only)
- **Security**: Bearer Token Required
- **RBAC**: Admin (all jobs), User (own jobs only)
- **Response**: Success message

### GET /api/pekerjaan/trash
- **Summary**: Get Trash Pekerjaan with RBAC (PostgreSQL)
- **Description**: Get list of soft-deleted jobs (Admin sees all, User sees their own)
- **Security**: Bearer Token Required
- **Response**: Array of deleted PekerjaanAlumni

### PUT /api/pekerjaan/restore/{id}
- **Summary**: Restore Pekerjaan with RBAC (PostgreSQL)
- **Description**: Restore a soft-deleted job (Admin and owner only)
- **Security**: Bearer Token Required
- **Response**: Success message

### DELETE /api/pekerjaan/hard/{id}
- **Summary**: Hard Delete Pekerjaan with RBAC (PostgreSQL)
- **Description**: Permanently delete a soft-deleted job (Admin and owner only)
- **Security**: Bearer Token Required
- **Response**: Success message

---

## 4. Pekerjaan Service - MongoDB (`pekerjaan_mongo_service.go`)

### POST /api/pekerjaan-mongo
- **Summary**: Create Pekerjaan (MongoDB)
- **Description**: Create a new job record in MongoDB (Admin only)
- **Security**: Bearer Token Required
- **Request**: PekerjaanMongo object
- **Response**: Created PekerjaanMongo (Status 201)

### GET /api/pekerjaan-mongo
- **Summary**: Get All Pekerjaan (MongoDB)
- **Description**: Get list of active jobs (Admin sees all, User sees only their own)
- **Security**: Bearer Token Required
- **Response**: Array of PekerjaanMongo

### GET /api/pekerjaan-mongo/{id}
- **Summary**: Get Pekerjaan by ID (MongoDB)
- **Description**: Get details of an active job
- **Security**: Bearer Token Required
- **Response**: PekerjaanMongo object

### DELETE /api/pekerjaan-mongo/{id}
- **Summary**: Soft Delete Pekerjaan (MongoDB)
- **Description**: Soft delete a job (Admin or job owner only)
- **Security**: Bearer Token Required
- **RBAC**: Admin (all jobs), User (own jobs only)
- **Response**: Success message

### PUT /api/pekerjaan-mongo/restore/{id}
- **Summary**: Restore Pekerjaan (MongoDB)
- **Description**: Restore a soft-deleted job (Admin or job owner only)
- **Security**: Bearer Token Required
- **Response**: Success message

### DELETE /api/pekerjaan-mongo/hard/{id}
- **Summary**: Hard Delete Pekerjaan (MongoDB)
- **Description**: Permanently delete a soft-deleted job (Admin or job owner only)
- **Security**: Bearer Token Required
- **Response**: Success message

---

## 5. Upload Service (`upload_service.go`)

### POST /api/upload
- **Summary**: Upload File
- **Description**: Upload foto or sertifikat file to the server
- **Security**: Bearer Token Required
- **Request**: FormData (file, category, optional user_id for admin)
- **Response**: Upload metadata with file info (Status 201)

### GET /api/upload
- **Summary**: Get All Uploads
- **Description**: Get list of all uploaded files (Admin sees all, User sees their own)
- **Security**: Bearer Token Required
- **Response**: Array of Upload objects

### GET /api/upload/{id}
- **Summary**: Get Upload by ID
- **Description**: Get details of an uploaded file
- **Security**: Bearer Token Required
- **Response**: Upload object

### DELETE /api/upload/{id}
- **Summary**: Delete Upload
- **Description**: Delete an uploaded file (Admin or file owner only)
- **Security**: Bearer Token Required
- **RBAC**: Admin (all files), User (own files only)
- **Response**: Success message

---

## Security Configuration

### Bearer Token (JWT)
- **Type**: API Key
- **Location**: Header
- **Header Name**: Authorization
- **Format**: "Bearer {JWT_TOKEN}"

Semua endpoint yang memerlukan autentikasi menggunakan `@Security Bearer`.

---

## Model Definitions

### Alumni
- id: int
- nama: string
- email: string
- angkatan: int
- user_id: int

### PekerjaanAlumni
- id: int
- alumni_id: int
- nama_perusahaan: string
- posisi_jabatan: string
- bidang_industri: string
- lokasi_kerja: string
- gaji_range: string
- tanggal_mulai_kerja: string
- tanggal_selesai_kerja: string
- status_pekerjaan: string
- deskripsi_pekerjaan: string
- is_deleted: nullable timestamp
- created_at: timestamp
- updated_at: timestamp

### PekerjaanMongo
- _id: ObjectID
- alumni_id: int
- nama_perusahaan: string
- posisi_jabatan: string
- bidang_industri: string
- lokasi_kerja: string
- gaji_range: string
- tanggal_mulai_kerja: string
- tanggal_selesai_kerja: string
- status_pekerjaan: string
- deskripsi_pekerjaan: string
- is_deleted: nullable timestamp
- created_at: timestamp
- updated_at: timestamp

### Upload
- _id: ObjectID
- user_id: int
- file_name: string
- original_name: string
- file_path: string
- file_size: int64
- file_type: string
- category: string (foto/sertifikat)
- created_at: timestamp

### LoginRequest
- username: string
- password: string

---

## RBAC (Role-Based Access Control)

### Roles
1. **Admin**: Full access to all resources
2. **User**: Restricted access (can only modify their own records)

### Role Enforcement
- Alumni: Create/Update/Delete (Admin only), View (Admin & User)
- Pekerjaan: Create/Update (Admin only), Soft Delete/Restore (Admin & Owner), View (Admin & User)
- Upload: Create/Delete (Admin & Owner), View (Admin & User)

---

## Testing Dengan Swagger UI

Setelah project dijalankan, akses dokumentasi Swagger di:
```
http://localhost:3000/swagger/index.html
```

Catatan: Pastikan untuk:
1. Generate docs dengan command: `swag init`
2. Import Swagger handler di main route jika belum ada

---

**Generated**: November 13, 2025
**Project**: Praktikum 4 CRUD API
**API Version**: 1.0
