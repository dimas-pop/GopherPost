# GopherPost - RESTful API Service 🚀

![Go Version](https://img.shields.io/badge/Go-1.25.3+-00ADD8?style=flat&logo=go)
![Database](https://img.shields.io/badge/PostgreSQL-18.1+-316192?style=flat&logo=postgresql)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=flat&logo=swagger)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
[![Live Demo](https://img.shields.io/badge/Demo-Try%20API%20Here-blue?style=for-the-badge&logo=swagger)](https://gopherpost-production.up.railway.app/swagger/index.html)

**GopherPost** adalah layanan REST API backend modern yang dibangun menggunakan Golang. Proyek ini mensimulasikan sistem manajemen konten (Headless CMS) dengan fitur autentikasi yang aman, manajemen relasi database, dan dokumentasi API interaktif.

Proyek ini dibuat untuk mendemonstrasikan implementasi standar industri dalam pengembangan backend, termasuk: **Clean Architecture**, **Secure Authentication**, **Performance Optimization**, dan **Testing**.

## ✨ Fitur Utama

* 🔐 **Secure Authentication**: Sistem Login & Register menggunakan **JWT (JSON Web Token)** dan Hashing Password dengan **Bcrypt**.
* 📝 **CRUD Operations**: Manajemen User, Post, dan Comment yang lengkap.
* 🛡️ **Middleware Security**: Proteksi endpoint privat dan validasi kepemilikan data (Authorization).
* 🚀 **Performance**: Implementasi **Pagination** untuk efisiensi data loading.
* 📄 **Interactive Docs**: Dokumentasi API otomatis menggunakan **Swagger UI**.
* 🗄️ **Relational Database**: Desain skema PostgreSQL yang ternormalisasi (Foreign Keys & Cascading).
* 🔍 **Observability**: Structured Logging menggunakan `slog` (JSON format).

## 🛠️ Teknologi yang Digunakan

* **Language**: Go (Golang)
* **Router**: Gorilla Mux
* **Database**: PostgreSQL
* **Driver**: Jackc/pgx (High performance driver)
* **Auth**: Golang-JWT
* **Docs**: Swaggo (Swagger)
* **Config**: Godotenv

## 🚀 Cara Menjalankan (Local)

1.  **Clone Repository**
    ```bash
    git clone [https://github.com/dimas-pop/GopherPost.git](https://github.com/dimas-pop/GopherPost.git)
    cd GopherPost
    ```

2.  **Setup Database**
    Buat database PostgreSQL:
    ```sql
    CREATE DATABASE gopherpost_db;
    ```

3.  **Environment Variables**
    Duplikasi file `.env.example` menjadi `.env` dan isi kredensial Anda:
    ```bash
    cp .env.example .env
    ```

4.  **Jalankan Aplikasi**
    ```bash
    go mod tidy
    go run main.go
    ```

5.  **Akses Dokumentasi**
    Buka browser dan kunjungi: `http://localhost:8080/swagger/index.html`

## 🧪 Testing

Proyek ini dilengkapi dengan **Unit Test** hanya untuk testing fungsi di password.go.

```bash
# Menjalankan test di folder handlers
go test ./handlers -v