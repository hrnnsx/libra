Berikut ringkasan daftar API dari **Anime Library API** yang dikelompokkan berdasarkan modul resources:

### 1. Authentication

* `POST /auth/register` — Register akun baru (UC-01)
* `POST /auth/login` — Login & mendapatkan access token (UC-02)
* `POST /auth/logout` — Logout (UC-03)

### 2. User Profile

* `GET /me` — Melihat profil user (UC-04)
* `PATCH /me` — Memperbarui profil user (UC-05)

### 3. Anime Discovery (Third-Party)

* `GET /animes/search` — Pencarian anime dasar (UC-06)
* `GET /animes/search/advanced` — Pencarian anime tingkat lanjut / filter (UC-07)
* `GET /animes/:external_id` — Detail anime dari third-party API (UC-08)

### 4. Personal Library

* `POST /me/library` — Menambahkan anime ke library pribadi (UC-09)
* `GET /me/library` — Melihat daftar anime di library pribadi (UC-10)
* `GET /me/library/:library_anime_id` — Detail item di library pribadi (UC-11)
* `PATCH /me/library/:library_anime_id` — Update status, episode, rating, atau catatan (UC-12, UC-13, UC-14, UC-15)
* `DELETE /me/library/:library_anime_id` — Hapus anime dari library pribadi (UC-16)

### 5. Groups & Collections

* `POST /me/groups` — Membuat grup/koleksi baru (UC-17)
* `GET /me/groups` — Melihat daftar grup milik user (UC-18)
* `GET /me/groups/:group_id` — Melihat detail grup beserta anime di dalamnya (UC-19)
* `PATCH /me/groups/:group_id` — Update nama/deskripsi grup (UC-20)
* `DELETE /me/groups/:group_id` — Hapus grup (UC-21)
* `POST /me/groups/:group_id/animes` — Menambahkan anime dari library ke grup (UC-22)
* `DELETE /me/groups/:group_id/animes/:library_anime_id` — Hapus anime dari grup (UC-23)