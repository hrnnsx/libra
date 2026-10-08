# Anime Library API --- Requirements Document

## 1. Overview

### 1.1 Purpose

Anime Library API adalah REST API yang menyediakan layanan perpustakaan
digital untuk anime. Sistem memungkinkan pengguna memiliki library anime
pribadi, mengatur status menonton, mencatat progress dan rating, serta
mengelompokkan anime ke dalam group atau collection.

Sistem juga menyediakan fitur pencarian anime menggunakan Third-Party
Anime API sebagai sumber data eksternal.

### 1.2 Scope

Versi awal (MVP) sistem mencakup:

-   Authentication dan pengelolaan akun dasar
-   Anime discovery dan pencarian anime
-   Advanced anime search/filter
-   Personal anime library
-   Status dan progress menonton
-   Group/collection anime
-   Integrasi Third-Party Anime API
-   Integrasi Anime Identification API
-   Privacy pada library dan group pengguna

------------------------------------------------------------------------

## 2. Actors

### 2.1 User

User adalah pengguna terdaftar yang dapat:

-   Membuat dan mengakses akun
-   Mencari anime
-   Melihat detail anime
-   Menambahkan anime ke library
-   Mengelola library pribadi
-   Membuat dan mengelola group anime

### 2.2 Third-Party Anime API

Third-Party Anime API adalah sistem eksternal yang menyediakan informasi,
pencarian, dan detail anime.

Sistem digunakan sebagai sumber data anime ketika user melakukan
discovery/search dan ketika sebuah anime akan ditambahkan ke library.

### 2.3 Anime Identification API

Anime Identification API adalah sistem eksternal yang digunakan untuk
mengidentifikasi anime berdasarkan screenshot atau gambar sebuah scene.

Sistem digunakan ketika user ingin mencari anime berdasarkan gambar.

------------------------------------------------------------------------

## 3. Functional Requirements

## 3.1 Authentication & Account

### FR-01 --- User Registration

Sistem harus memungkinkan user membuat akun menggunakan informasi dasar
seperti username, email, dan password.

### FR-02 --- User Login

Sistem harus memungkinkan user melakukan login menggunakan kredensial
yang terdaftar.

### FR-03 --- User Logout

Sistem harus memungkinkan user melakukan logout dari sesi
autentikasinya.

### FR-04 --- Manage Basic Profile

Sistem harus memungkinkan user melihat dan mengubah informasi profil
dasar miliknya.

Informasi profil dasar minimal:

-   Username
-   Email
-   Password

Sistem tidak menyediakan public profile pada MVP.

------------------------------------------------------------------------

## 3.2 Anime Discovery

### FR-05 --- Search Anime

Sistem harus memungkinkan user mencari anime berdasarkan kata kunci.

Contoh:

``` http
GET /animes/search?q=Frieren
```

Hasil pencarian diperoleh dari Third-Party Anime API.

### FR-06 --- View Anime Detail

Sistem harus memungkinkan user melihat detail anime dari hasil
pencarian.

Informasi yang dapat ditampilkan minimal:

-   External ID
-   Title
-   Synopsis
-   Image URL

### FR-07 --- Advanced Anime Search

Sistem harus menyediakan pencarian/filter anime berdasarkan beberapa
parameter yang didukung oleh Third-Party Anime API.

Parameter dapat mencakup:

-   Keyword
-   Genre
-   Status
-   Tahun
-   Sorting

Contoh:

``` http
GET /animes/search?q=Naruto&genre=action&status=finished&year=2015
```

Parameter final akan disesuaikan dengan kemampuan Third-Party Anime API
yang dipilih.

### FR-08 --- Third-Party Anime API Integration

Sistem harus mengintegrasikan Third-Party Anime API untuk memperoleh
data anime.

Third-Party API digunakan sebagai sumber data untuk:

1.  Anime search
2.  Advanced search/filter
3.  Anime detail
4.  Pengambilan data anime ketika user menambahkan anime ke library
5. Identifikasi anime berdasarkan screenshot sebuah scene

Sistem tidak menyimpan seluruh katalog anime dari Third-Party API.

------------------------------------------------------------------------

## 3.3 Anime Data

### FR-09 --- Store Selected Anime

Sistem hanya menyimpan data anime ke database lokal ketika user
menambahkan anime ke library.

Data anime minimal yang disimpan:

-   Internal ID
-   External ID
-   Title
-   Synopsis
-   Image URL
-   Created At
-   Updated At

Hasil pencarian dari Third-Party API yang belum ditambahkan ke library
tidak disimpan secara permanen.

### FR-10 --- Prevent Duplicate Anime in User Library

Satu user tidak boleh memiliki anime yang sama lebih dari satu kali
dalam library.

Constraint yang digunakan:

``` text
UNIQUE(user_id, anime_id)
```

------------------------------------------------------------------------

## 3.4 Personal Anime Library

### FR-11 --- Add Anime to Library

Sistem harus memungkinkan user menambahkan anime yang ditemukan melalui
Anime Discovery ke library pribadi.

Ketika user menambahkan anime:

1.  Sistem mengambil detail anime dari Third-Party API.
2.  Sistem menyimpan data anime jika belum tersedia di database lokal.
3.  Sistem membuat record library untuk user tersebut.

### FR-12 --- View Personal Library

Sistem harus memungkinkan user melihat seluruh anime yang terdapat dalam
library miliknya.

### FR-13 --- View Library Anime Detail

Sistem harus memungkinkan user melihat detail anime beserta informasi
personal dari library.

Informasi personal dapat meliputi:

-   Status
-   Current episode
-   Rating
-   Notes
-   Started date
-   Completed date

### FR-14 --- Update Watching Status

User harus dapat mengubah status anime dalam library.

Status yang tersedia:

-   `WATCHING`
-   `COMPLETED`
-   `PLAN_TO_WATCH`
-   `ON_HOLD`
-   `DROPPED`

### FR-15 --- Update Episode Progress

User harus dapat memperbarui progress episode anime yang terdapat dalam
library.

Contoh:

``` text
Current Episode: 17
```

### FR-16 --- Rate Anime

User harus dapat memberikan rating terhadap anime dalam library.

### FR-17 --- Add Personal Notes

User harus dapat menambahkan atau mengubah catatan pribadi terhadap
anime dalam library.

### FR-18 --- Remove Anime from Library

User harus dapat menghapus anime dari library pribadinya.

Penghapusan dari library tidak berarti menghapus data anime global dari
database apabila anime tersebut masih digunakan oleh user lain.

------------------------------------------------------------------------

## 3.5 Anime Group / Collection

### FR-19 --- Create Group

User harus dapat membuat group/collection anime pribadi.

Group minimal memiliki:

-   Name
-   Description

Contoh:

``` text
Favorite
Best Isekai
Anime 10/10
Watch Later
```

### FR-20 --- View Groups

User harus dapat melihat group yang dimilikinya.

### FR-21 --- Update Group

User harus dapat mengubah informasi group miliknya.

### FR-22 --- Delete Group

User harus dapat menghapus group miliknya.

Penghapusan group tidak menghapus anime dari library.

### FR-23 --- Add Library Anime to Group

User harus dapat memasukkan anime yang terdapat di library miliknya ke
dalam group.

### FR-24 --- Remove Anime from Group

User harus dapat menghapus anime dari group tanpa menghapus anime
tersebut dari library.

### FR-25 --- Multiple Group Membership

Satu anime dapat berada di beberapa group milik user yang sama.

Contoh:

``` text
Frieren
├── Favorite
├── Best Fantasy
└── Anime 10/10
```

------------------------------------------------------------------------

## 3.6 Privacy & Authorization

### FR-26 --- Private Library

Library user bersifat private.

User hanya dapat melihat library miliknya sendiri.

### FR-27 --- Private Groups

Group user bersifat private.

User hanya dapat melihat dan mengelola group miliknya sendiri.

### FR-28 --- Resource Ownership

User hanya dapat membuat, mengubah, atau menghapus resource yang menjadi
miliknya.

Sistem harus melakukan authorization berdasarkan `user_id`.

### FR-29 --- No Public Profile

MVP tidak menyediakan public profile, social feed, follow, like, atau
comment.

------------------------------------------------------------------------

## 4. Non-Functional Requirements

### NFR-01 --- REST API

Sistem harus dikembangkan sebagai REST API.

### NFR-02 --- JSON Response

Request dan response API menggunakan format JSON kecuali terdapat
kebutuhan khusus dari integrasi eksternal.

### NFR-03 --- Authentication

Endpoint yang membutuhkan identitas user harus dilindungi oleh mekanisme
authentication.

### NFR-04 --- Authorization

Sistem harus memastikan user hanya dapat mengakses dan memodifikasi
resource yang dimilikinya.

### NFR-05 --- Validation

Sistem harus melakukan validasi terhadap input user.

Validasi minimal mencakup:

-   Username
-   Email
-   Password
-   Rating
-   Episode progress
-   Anime status
-   Group name

### NFR-06 --- Error Handling

Sistem harus memberikan HTTP status code dan response error yang sesuai.

Contoh:

``` text
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
409 Conflict
500 Internal Server Error
```

### NFR-07 --- External API Failure Handling

Sistem harus menangani kondisi ketika Third-Party Anime API tidak
tersedia, mengalami error, atau mengembalikan response yang tidak valid.

### NFR-08 --- Data Integrity

Relasi antara user, library, anime, dan group harus dijaga melalui
database constraints dan validasi pada application layer.

------------------------------------------------------------------------

## 5. Business Rules

### BR-01 --- Anime Search Does Not Persist Data

Hasil pencarian anime dari Third-Party API tidak otomatis disimpan ke
database lokal.

### BR-02 --- Anime Is Persisted on Library Addition

Data anime disimpan secara lokal ketika user menambahkan anime ke
library.

### BR-03 --- Library Ownership

Setiap library entry harus dimiliki tepat oleh satu user.

### BR-04 --- Unique Library Entry

User tidak dapat memiliki anime yang sama lebih dari satu kali dalam
library.

### BR-05 --- Group Ownership

Setiap group harus dimiliki tepat oleh satu user.

### BR-06 --- Group Membership

Anime hanya dapat dimasukkan ke group jika anime tersebut sudah terdapat
di library user yang sama.

### BR-07 --- Group Removal Does Not Remove Library Entry

Menghapus anime dari group tidak menghapus anime dari library.

### BR-08 --- Library Removal Does Not Delete Global Anime Data

Menghapus anime dari library tidak secara otomatis menghapus record
anime dari tabel anime.

Anime hanya dapat dihapus secara fisik apabila tidak lagi digunakan dan
kebijakan data mengizinkannya.

### BR-09 --- Anime Identification Does Not Persist Data

Hasil identifikasi anime dari Anime Identification API tidak otomatis
disimpan ke database lokal.

Hasil identifikasi hanya digunakan untuk menampilkan informasi kepada user
dan dapat digunakan sebagai referensi untuk menemukan anime.

------------------------------------------------------------------------

## 6. High-Level User Flow

### 6.1 Search Anime

``` text
User
  ↓
Search Anime
  ↓
Our API
  ↓
Third-Party Anime API
  ↓
Search Results
  ↓
Anime Detail
```

### 6.2 Add Anime to Library

``` text
Anime Detail
  ↓
Add to Library
  ↓
Fetch Anime Detail
  ↓
Save Anime Locally
  ↓
Create Library Entry
  ↓
Library Updated
```

### 6.3 Manage Library

``` text
My Library
  ↓
Select Anime
  ↓
Update Status / Progress / Rating / Notes
  ↓
Save
```

### 6.4 Manage Group

``` text
My Groups
  ↓
Create Group
  ↓
Select Anime from My Library
  ↓
Add to Group
  ↓
Group Updated
```

### 6.5 Identify Anime from Screenshot

```text
User
  ↓
Upload / Provide Image URL
  ↓
Our API
  ↓
Anime Identification API
  ↓
Identification Result
  ↓
Anime / Episode / Timestamp

------------------------------------------------------------------------

## 7. Out of Scope for MVP

Fitur berikut belum termasuk dalam versi awal:

-   Public profile
-   Follow user
-   Social feed
-   Like
-   Comment
-   User review publik
-   Payment gateway
-   Subscription
-   Email notification
-   Chat
-   Anime streaming server
-   Sinkronisasi penuh seluruh katalog Third-Party Anime API

------------------------------------------------------------------------