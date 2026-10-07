### 1. Authentication

* `POST /auth/register` — Register akun baru 
* `POST /auth/login` — Login & mendapatkan access token 
* `POST /auth/logout` — Logout 

### 2. User Profile

* `GET /me` — Melihat profil user 
* `PATCH /me` — Memperbarui profil user 

### 3. Anime Discovery (Third-Party)

* `GET /animes/search` — Pencarian anime dasar 
* `GET /animes/search/advanced` — Pencarian anime tingkat lanjut / filter 
* `GET /animes/:external_id` — Detail anime dari third-party API 

### 4. Personal Library

* `POST /me/library` — Menambahkan anime ke library pribadi 
* `GET /me/library` — Melihat daftar anime di library pribadi 
* `GET /me/library/:library_anime_id` — Detail item di library pribadi 
* `PATCH /me/library/:library_anime_id` — Update status, episode, rating, atau catatan
* `DELETE /me/library/:library_anime_id` — Hapus anime dari library pribadi

### 5. Groups & Collections

* `POST /me/groups` — Membuat grup/koleksi baru 
* `GET /me/groups` — Melihat daftar grup milik user 
* `GET /me/groups/:group_id` — Melihat detail grup beserta anime di dalamnya
* `PATCH /me/groups/:group_id` — Update nama/deskripsi grup
* `DELETE /me/groups/:group_id` — Hapus grup
* `POST /me/groups/:group_id/animes` — Menambahkan anime dari library ke grup
* `DELETE /me/groups/:group_id/animes/:library_anime_id` — Hapus anime dari grup