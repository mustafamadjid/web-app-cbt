# Deploy Guide VPS Ubuntu + Docker

Panduan ini dibuat khusus untuk proyek `web-app-cbt` di repository ini. Fokusnya adalah deployment pertama ke VPS Ubuntu memakai Docker, Docker Compose, Nginx, dan SSL.

Guide ini mengikuti kondisi proyek saat ini:

- ada 4 service utama di `docker-compose.yml`: `postgres`, `migrate`, `backend`, `frontend`
- backend adalah aplikasi Go
- frontend dibuild dengan Vite lalu disajikan oleh Nginx di dalam container
- database migration dijalankan lewat service `migrate`
- data penting disimpan di Docker volume

## 1. Gambaran arsitektur

Hasil akhir yang disarankan:

1. frontend diakses dari `https://cbt.domainkamu.com`
2. API diakses dari `https://api.domainkamu.com`
3. Nginx di VPS menerima request publik
4. Nginx meneruskan request frontend ke `127.0.0.1:3000`
5. Nginx meneruskan request backend ke `127.0.0.1:8080`
6. backend berbicara ke PostgreSQL lewat network Docker internal
7. data database dan upload file tetap tersimpan walaupun container di-recreate

Kenapa saya sarankan subdomain:

- lebih mudah dipahami saat pertama deploy
- `VITE_API_URL` jadi jelas
- aturan Nginx lebih sederhana
- upload file tetap konsisten karena frontend memang memanggil host API

## 2. Hal penting yang perlu dipahami sebelum deploy

### 2.1 File env utama untuk Docker adalah `.env` di root

Untuk deployment Docker, file utama yang dibaca `docker-compose.yml` adalah:

```text
/.env
```

Bukan:

```text
/backend/.env
```

Jadi untuk deploy production, fokus edit file `.env` di root project.

### 2.2 Frontend membaca `VITE_API_URL` saat build image

Di `docker-compose.yml`, frontend dibuild dengan build arg:

```yaml
args:
  VITE_API_URL: ${VITE_API_URL}
```

Artinya:

- kalau `VITE_API_URL` berubah, frontend harus di-build ulang
- restart biasa tanpa rebuild tidak cukup

Perintah yang aman kalau env frontend berubah:

```bash
docker compose up -d --build frontend
```

### 2.3 Backend memakai `BACKEND_PUBLIC_URL` untuk membentuk URL file publik

Di compose, backend menerima:

```yaml
BASE_URL: ${BACKEND_PUBLIC_URL}
```

Nilai ini dipakai backend untuk membuat URL publik file upload. Kalau nilainya salah, gambar atau dokumen yang dihasilkan backend bisa mengarah ke alamat yang salah.

### 2.4 `APP_PUBLIC_URL` saat ini tidak dipakai langsung oleh compose

Variabel `APP_PUBLIC_URL` ada di `.env.example`, tetapi pada konfigurasi repo saat ini variabel itu tidak dipakai langsung oleh service Docker. Tetap boleh diisi untuk dokumentasi, tetapi yang paling penting untuk runtime sekarang adalah:

- `BACKEND_PUBLIC_URL`
- `VITE_API_URL`
- `TRUSTED_ORIGINS`

## 3. Prasyarat

Sebelum mulai, siapkan:

- 1 VPS Ubuntu, disarankan Ubuntu 22.04 atau 24.04
- akses SSH ke VPS
- domain aktif
- 2 subdomain yang mengarah ke IP VPS:
  - `cbt.domainkamu.com`
  - `api.domainkamu.com`

Contoh DNS:

- `A record` `cbt` -> IP VPS
- `A record` `api` -> IP VPS

Kalau DNS belum mengarah ke VPS, SSL dari Let's Encrypt tidak akan berhasil.

## 4. Struktur folder yang disarankan di VPS

Saya sarankan project diletakkan di:

```text
/opt/web-app-cbt
```

Contoh:

```text
/opt/web-app-cbt
  |- docker-compose.yml
  |- .env
  |- backend/
  |- frontend/
  `- deploy-guide.md
```

## 5. Login ke VPS dan update sistem

Masuk ke VPS:

```bash
ssh username@IP_VPS
```

Update package:

```bash
sudo apt update && sudo apt upgrade -y
```

Install package dasar:

```bash
sudo apt install -y git curl ca-certificates gnupg lsb-release nginx ufw snapd
```

## 6. Setup firewall dasar

Izinkan port yang dibutuhkan:

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

Catatan penting:

- compose saat ini mem-publish port ke `127.0.0.1`
- artinya port `3000`, `8080`, dan `5432` tidak perlu dibuka ke internet
- ini bagus untuk keamanan karena akses publik hanya lewat Nginx

## 7. Install Docker dan Docker Compose di Ubuntu

Bagian ini mengikuti alur instalasi Docker Engine via repository resmi Docker untuk Ubuntu.

### 7.1 Hapus paket Docker lama kalau ada

```bash
sudo apt remove -y docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc
```

Kalau sebagian paket tidak ada, itu normal.

### 7.2 Tambahkan repository resmi Docker

```bash
sudo apt update
sudo apt install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
```

Buat source repository Docker:

```bash
sudo tee /etc/apt/sources.list.d/docker.sources > /dev/null <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF
```

Update package index:

```bash
sudo apt update
```

### 7.3 Install Docker Engine dan Compose plugin

```bash
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### 7.4 Aktifkan Docker saat boot

```bash
sudo systemctl enable docker
sudo systemctl start docker
sudo systemctl status docker
```

### 7.5 Opsional: jalankan Docker tanpa `sudo`

```bash
sudo usermod -aG docker $USER
```

Setelah itu logout lalu login lagi ke VPS.

### 7.6 Verifikasi Docker

```bash
docker --version
docker compose version
docker run hello-world
```

Kalau command di atas jalan, instalasi Docker sudah benar.

## 8. Clone project ke VPS

Masuk ke `/opt` lalu clone repo:

```bash
cd /opt
sudo git clone <URL_REPOSITORY_KAMU> web-app-cbt
sudo chown -R $USER:$USER /opt/web-app-cbt
cd /opt/web-app-cbt
```

Kalau repository masih private, pastikan akses SSH key atau token GitHub sudah siap.

## 9. Setup `.env` production

### 9.1 Buat file env dari contoh

```bash
cp .env.example .env
```

### 9.2 Isi `.env`

Contoh untuk skenario subdomain:

```env
# PostgreSQL
POSTGRES_DB=web_app_cbt
POSTGRES_USER=web_app_cbt
POSTGRES_PASSWORD=ganti-dengan-password-random-yang-panjang
POSTGRES_PORT=5432

# Published local ports on VPS
FRONTEND_PORT=3000
BACKEND_PORT=8080

# Public URLs
APP_PUBLIC_URL=https://cbt.domainkamu.com
BACKEND_PUBLIC_URL=https://api.domainkamu.com
VITE_API_URL=https://api.domainkamu.com

# CORS + CSRF trusted origins
TRUSTED_ORIGINS=https://cbt.domainkamu.com

# JWT secrets
JWT_ACCESS_SECRET=ganti-dengan-secret-random-yang-sangat-panjang
JWT_REFRESH_SECRET=ganti-dengan-secret-random-yang-sangat-panjang
```

### 9.3 Penjelasan variabel penting

#### Database

- `POSTGRES_DB`: nama database
- `POSTGRES_USER`: user database
- `POSTGRES_PASSWORD`: password database
- `POSTGRES_PORT`: port PostgreSQL di localhost VPS

#### Port service di VPS

- `FRONTEND_PORT=3000` berarti frontend container tersedia di `127.0.0.1:3000`
- `BACKEND_PORT=8080` berarti backend container tersedia di `127.0.0.1:8080`

Karena dipublish ke `127.0.0.1`, service ini tidak langsung terbuka ke internet.

#### URL publik

- `APP_PUBLIC_URL`: dokumentasi URL frontend, tidak dipakai langsung oleh compose saat ini
- `BACKEND_PUBLIC_URL`: URL publik backend, dipakai backend untuk URL file upload
- `VITE_API_URL`: URL API yang ditanam saat build frontend

#### Origin yang diizinkan

- `TRUSTED_ORIGINS` harus diisi dengan origin frontend yang diizinkan mengakses backend

Contoh:

```env
TRUSTED_ORIGINS=https://cbt.domainkamu.com
```

Kalau ada lebih dari satu origin:

```env
TRUSTED_ORIGINS=https://cbt.domainkamu.com,https://admin.domainkamu.com
```

#### JWT secret

- wajib panjang
- wajib random
- jangan pakai secret dari development

Contoh generate secret:

```bash
openssl rand -base64 48
```

Jalankan 2 kali, satu untuk `JWT_ACCESS_SECRET`, satu untuk `JWT_REFRESH_SECRET`.

### 9.4 Hal yang wajib dihindari

Jangan deploy production dengan nilai seperti:

- `http://localhost:3000`
- `http://localhost:8080`
- password database pendek
- secret JWT lama dari development

## 10. Build dan jalankan container

Masuk ke root project:

```bash
cd /opt/web-app-cbt
```

### 10.1 Build image

```bash
docker compose build
```

Build pertama bisa lama karena:

- backend build binary Go
- frontend install dependency lalu build asset

### 10.2 Jalankan PostgreSQL dulu

```bash
docker compose up -d postgres
```

Cek status:

```bash
docker compose ps
docker compose logs -f postgres
```

Tunggu sampai postgres sehat dan tidak ada error.

### 10.3 Jalankan migration database

Service `migrate` memakai profile `migrate`. Jalankan:

```bash
docker compose --profile migrate run --rm migrate
```

Kalau selesai lalu container keluar, itu normal.

Kalau migration gagal, perbaiki dulu sebelum backend/frontend dijalankan.

### 10.4 Jalankan backend dan frontend

```bash
docker compose up -d backend frontend
```

Atau sekalian:

```bash
docker compose up -d postgres backend frontend
```

### 10.5 Verifikasi service

```bash
docker compose ps
docker compose logs --tail=100 backend
docker compose logs --tail=100 frontend
```

Kalau mau live log:

```bash
docker compose logs -f backend frontend
```

### 10.6 Test akses lokal dari VPS

Test frontend:

```bash
curl -I http://127.0.0.1:3000
```

Test backend:

```bash
curl -I http://127.0.0.1:8080/auth/me
```

Untuk `/auth/me`, status `401` tanpa login itu wajar. Yang penting backend merespons dan bukan connection refused.

## 11. Setup Nginx untuk subdomain

Kita akan buat 2 config:

- satu untuk frontend `cbt.domainkamu.com`
- satu untuk API `api.domainkamu.com`

### 11.1 Config frontend

Buat file:

```bash
sudo nano /etc/nginx/sites-available/web-app-cbt-frontend
```

Isi:

```nginx
server {
    listen 80;
    server_name cbt.domainkamu.com;

    client_max_body_size 20m;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }
}
```

### 11.2 Config backend API

Buat file:

```bash
sudo nano /etc/nginx/sites-available/web-app-cbt-api
```

Isi:

```nginx
server {
    listen 80;
    server_name api.domainkamu.com;

    client_max_body_size 20m;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }
}
```

Catatan:

- pada skenario subdomain API, semua request ke `api.domainkamu.com` diteruskan ke backend
- endpoint upload seperti `/uploads/image/...` dan `/uploads/document/...` ikut lewat sini
- jadi tidak perlu blok `location /uploads/` terpisah

### 11.3 Aktifkan site

```bash
sudo ln -s /etc/nginx/sites-available/web-app-cbt-frontend /etc/nginx/sites-enabled/
sudo ln -s /etc/nginx/sites-available/web-app-cbt-api /etc/nginx/sites-enabled/
```

Kalau default site masih aktif dan mengganggu:

```bash
sudo rm -f /etc/nginx/sites-enabled/default
```

### 11.4 Test dan reload Nginx

```bash
sudo nginx -t
sudo systemctl reload nginx
sudo systemctl status nginx
```

## 12. Setup SSL HTTPS dengan Certbot

Sebelum langkah ini:

- DNS harus sudah mengarah ke VPS
- Nginx harus aktif di port 80
- firewall harus membuka 80 dan 443

### 12.1 Install Certbot

```bash
sudo snap install core
sudo snap refresh core
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/local/bin/certbot
```

### 12.2 Buat sertifikat dan biarkan Certbot edit config Nginx

```bash
sudo certbot --nginx -d cbt.domainkamu.com -d api.domainkamu.com
```

Biasanya setelah selesai:

- HTTP akan diarahkan ke HTTPS
- config SSL Nginx akan ditambahkan otomatis

### 12.3 Test auto-renew

```bash
sudo certbot renew --dry-run
```

## 13. Verifikasi akhir

Setelah SSL selesai, cek:

### 13.1 Dari browser

- `https://cbt.domainkamu.com`
- `https://api.domainkamu.com/auth/me`

Yang diharapkan:

- frontend tampil
- API merespons
- tidak ada warning SSL

### 13.2 Dari VPS

```bash
docker compose ps
docker compose logs --tail=50 postgres
docker compose logs --tail=50 backend
docker compose logs --tail=50 frontend
sudo systemctl status nginx
```

## 14. Volume: data tersimpan di mana

Ini bagian penting supaya kamu paham data tidak hilang saat container dihapus atau di-build ulang.

### 14.1 Volume yang dipakai proyek ini

Di `docker-compose.yml`, ada 2 named volume:

```yaml
volumes:
  postgres_data:
  backend_uploads:
```

Artinya:

- `postgres_data` menyimpan data PostgreSQL
- `backend_uploads` menyimpan file upload backend

### 14.2 Lokasi mount di dalam container

Untuk `postgres`:

- volume: `postgres_data`
- mount path: `/var/lib/postgresql/data`

Untuk `backend`:

- volume: `backend_uploads`
- mount path: `/app/public/uploads`

Subfolder upload yang memang dipakai image backend saat ini antara lain:

- `/app/public/uploads/image/bank_soal`
- `/app/public/uploads/document/import_soal`
- `/app/public/uploads/document/pengumuman`

### 14.3 Lokasi sebenarnya di host Ubuntu

Karena saat ini memakai named volume Docker, data biasanya disimpan di:

```text
/var/lib/docker/volumes/
```

Nama akhirnya biasanya memakai prefix nama project Compose. Kalau folder project bernama `web-app-cbt`, biasanya hasilnya mirip:

- `web-app-cbt_postgres_data`
- `web-app-cbt_backend_uploads`

Mountpoint-nya biasanya mirip:

- `/var/lib/docker/volumes/web-app-cbt_postgres_data/_data`
- `/var/lib/docker/volumes/web-app-cbt_backend_uploads/_data`

### 14.4 Cara memastikan path volume yang benar

Lihat daftar volume:

```bash
docker volume ls
```

Lihat detail volume:

```bash
docker volume inspect web-app-cbt_postgres_data
docker volume inspect web-app-cbt_backend_uploads
```

Lihat mountpoint saja:

```bash
docker volume inspect web-app-cbt_postgres_data --format '{{ .Mountpoint }}'
docker volume inspect web-app-cbt_backend_uploads --format '{{ .Mountpoint }}'
```

Kalau nama volume berbeda, cek dulu hasil `docker volume ls` karena prefix bisa berubah tergantung nama project Compose.

### 14.5 Kapan data tetap aman

Data tetap aman saat kamu menjalankan:

- `docker compose up -d`
- `docker compose down`
- `docker compose up -d --build`

Data bisa hilang kalau kamu sengaja menghapus volume, misalnya:

```bash
docker compose down -v
```

Jadi hati-hati dengan flag `-v`.

### 14.6 Opsional: kalau ingin data ada di folder host yang benar-benar jelas

Kalau suatu saat kamu ingin lokasi data langsung terlihat di host, kamu bisa mengganti named volume menjadi bind mount. Contoh:

```yaml
services:
  postgres:
    volumes:
      - /opt/web-app-cbt/data/postgres:/var/lib/postgresql/data

  backend:
    volumes:
      - /opt/web-app-cbt/data/uploads:/app/public/uploads
```

Kalau memakai model ini, data akan terlihat langsung di:

- `/opt/web-app-cbt/data/postgres`
- `/opt/web-app-cbt/data/uploads`

Penting:

- ini bukan konfigurasi aktif repo saat ini
- kalau ingin pindah dari named volume ke bind mount, lakukan hati-hati agar data lama tidak hilang
- untuk deploy pertama, saya sarankan tetap pakai named volume dulu

## 15. Backup sederhana

Minimal pahami 2 jenis backup:

- database PostgreSQL
- file upload backend

### 15.1 Backup database

Contoh:

```bash
docker exec -t web-app-cbt-postgres pg_dump -U web_app_cbt -d web_app_cbt > backup.sql
```

Sesuaikan user dan nama database dengan isi `.env`.

### 15.2 Backup file upload

Cari dulu mountpoint volume `backend_uploads`, lalu backup folder tersebut:

```bash
sudo tar -czf backend-uploads-backup.tar.gz /var/lib/docker/volumes/web-app-cbt_backend_uploads/_data
```

## 16. Alur update aplikasi setelah deploy pertama

Kalau nanti ada perubahan code:

### 16.1 Pull source terbaru

```bash
cd /opt/web-app-cbt
git pull
```

### 16.2 Kalau backend berubah

```bash
docker compose up -d --build backend
```

### 16.3 Kalau frontend berubah

```bash
docker compose up -d --build frontend
```

### 16.4 Kalau migration bertambah

```bash
docker compose up -d postgres
docker compose --profile migrate run --rm migrate
docker compose up -d --build backend frontend
```

## 17. Troubleshooting yang paling sering

### 17.1 Frontend tampil tapi API gagal

Cek:

- `VITE_API_URL` di `.env`
- apakah frontend sudah di-build ulang
- apakah backend container hidup
- apakah `api.domainkamu.com` sudah mengarah ke VPS

Perintah:

```bash
docker compose ps
docker compose logs --tail=100 backend
docker compose up -d --build frontend
```

### 17.2 Login gagal karena cookie atau CORS

Cek:

- `TRUSTED_ORIGINS` harus berisi origin frontend yang benar
- jangan campur `http` dan `https`
- pastikan frontend memanggil `https://api.domainkamu.com`

Contoh yang benar:

```env
TRUSTED_ORIGINS=https://cbt.domainkamu.com
```

### 17.3 Upload file gagal

Cek:

- log backend
- `client_max_body_size` di Nginx
- volume `backend_uploads`

Perintah:

```bash
docker compose logs --tail=100 backend
docker volume inspect web-app-cbt_backend_uploads
```

### 17.4 Data database hilang setelah redeploy

Biasanya karena:

- volume tidak terbuat dengan benar
- project name Compose berubah
- pernah menjalankan `docker compose down -v`

Cek:

```bash
docker volume ls
docker volume inspect web-app-cbt_postgres_data
```

### 17.5 Nginx 502 Bad Gateway

Artinya Nginx hidup, tapi upstream container tidak bisa dijangkau.

Cek:

```bash
curl -I http://127.0.0.1:3000
curl -I http://127.0.0.1:8080/auth/me
docker compose ps
docker compose logs --tail=100 frontend
docker compose logs --tail=100 backend
```

## 18. Opsi alternatif: satu domain dengan path `/api`

Repository ini sudah punya contoh config di:

```text
deploy/nginx/web-app-cbt.conf
```

Model ini cocok kalau kamu ingin:

- frontend di `https://domainkamu.com`
- backend di `https://domainkamu.com/api`

Kalau memakai model ini, `.env` yang sesuai kira-kira:

```env
APP_PUBLIC_URL=https://domainkamu.com
BACKEND_PUBLIC_URL=https://domainkamu.com/api
VITE_API_URL=https://domainkamu.com/api
TRUSTED_ORIGINS=https://domainkamu.com
```

Tetapi untuk deployment pertama, saya tetap lebih menyarankan model subdomain karena lebih mudah dipahami.

## 19. Checklist singkat

Urutan ringkasnya:

1. siapkan VPS Ubuntu
2. arahkan DNS `cbt.domainkamu.com` dan `api.domainkamu.com`
3. install `git`, `nginx`, `ufw`, `snapd`
4. install Docker Engine dan Docker Compose plugin
5. clone repo ke `/opt/web-app-cbt`
6. buat `.env` production di root
7. jalankan `docker compose build`
8. jalankan `docker compose up -d postgres`
9. jalankan `docker compose --profile migrate run --rm migrate`
10. jalankan `docker compose up -d backend frontend`
11. pasang config Nginx frontend dan API
12. test `sudo nginx -t` lalu reload Nginx
13. install Certbot dan jalankan `sudo certbot --nginx -d cbt.domainkamu.com -d api.domainkamu.com`
14. cek frontend, API, log container, dan volume

## 20. Penutup

Kalau mengikuti guide ini, kamu akan mendapatkan setup production dasar yang rapi:

- port container hanya expose ke localhost VPS
- akses publik masuk lewat Nginx
- frontend dan API dipisah lewat subdomain
- database dan upload file tersimpan di Docker volume
- SSL diurus oleh Certbot

Langkah berikut yang biasanya dikerjakan setelah deploy pertama:

- backup otomatis database
- backup file upload
- CI/CD
- monitoring
- log rotation
