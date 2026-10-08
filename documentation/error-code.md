
| Code                      | HTTP Status | Keterangan                               |
| ------------------------- | ----------: | ---------------------------------------- |
| `INVALID_REQUEST`         |         400 | Request body/parameter tidak valid       |
| `USERNAME_EXISTS`         |         409 | Username sudah digunakan                 |
| `EMAIL_EXISTS`            |         409 | Email sudah digunakan                    |
| `INVALID_CREDENTIALS`     |         401 | Email atau password salah                |
| `TOKEN_GENERATION_FAILED` |         500 | Gagal membuat JWT                        |
| `INTERNAL_SERVER_ERROR`   |         500 | Internal server error                    |
| `MISSING_AUTHORIZATION`   |         401 | Authorization header tidak ada           |
| `INVALID_AUTHORIZATION`   |         401 | Format Authorization tidak valid         |
| `INVALID_TOKEN`           |         401 | JWT tidak valid atau expired             |
| `EXTERNAL_API_ERROR`      |         502 | Gagal berkomunikasi dengan API eksternal |
