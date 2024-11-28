# API DOCUMENTATION

### Daftar Endpoint Fitur Autentikasi

| No | Method | Endpoint    | Request Body                                                                                 | Deskripsi                                                                                 |
|----|--------|-------------|---------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------|
| 1  | POST   | `/register` | `{ "username": "artameviay", "first_name": "jasmine", "last_name": "arta", "phone_number": "01921", "password": "" }` | Membuat akun pengguna baru                                                               |
| 2  | POST   | `/login`    | `{ "username": "artameviay", "password": "" }`                                             | Melakukan login dan mendapatkan token autentikasi untuk akses API                         |
| 3  | POST   | `/logout`   | -                                                                                           | Menghapus token dan melakukan logout pengguna dari aplikasi                               |

