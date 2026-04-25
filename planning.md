# Planning Deployment Docker untuk Web App CBT

Dokumen ini adalah rencana implementasi Dockerfile dan Docker Compose agar aplikasi siap deploy ke VPS Ubuntu. Target pembaca adalah agen/implementor dengan reasoning rendah, jadi setiap langkah dibuat eksplisit dan berurutan.

## 1. Kondisi Repo Saat Ini

- Backend berada di folder `backend`.
- Backend menggunakan Go module `github.com/mustafamadjid/web-app-cbt`.
- Entry point backend ada di `backend/cmd/http/main.go`.
- Backend listen di port internal `8080`.
- Backend membaca konfigurasi penting dari environment variable:
  - `ENV`
  - `BASE_URL`
  - `UPLOAD_DIR`
  - `POSTGRES_DBURL`
  - `JWT_ACCESS_SECRET`
  - `JWT_REFRESH_SECRET`
  - `TRUSTED_ORIGINS`
- Upload file backend secara default berada di `backend/cmd/http/public/uploads`.
- Path upload dapat dipindah dengan `UPLOAD_DIR`.
- Frontend berada di folder `frontend`.
- Frontend menggunakan Vite/React.
- Frontend build command: `npm run build`.
- Frontend membaca URL API dari build-time env `VITE_API_URL`.
- Database menggunakan PostgreSQL.
- Migrasi database menggunakan Goose SQL migrations di `backend/internal/db/migrations`.

## 2. Target Arsitektur Production

Di VPS Ubuntu, jalankan service berikut dengan Docker Compose:

- `frontend`: container Nginx yang menyajikan hasil build React/Vite.
- `backend`: container Go API.
- `postgres`: container PostgreSQL.
- `migrate`: one-shot container untuk menjalankan database migration.

Nginx host/server VPS dapat berjalan di luar compose sebagai reverse proxy utama. Compose tetap memakai Nginx internal untuk frontend static file.

Port yang direncanakan:

- Host `3000` -> container `frontend:80`.
- Host `8080` -> container `backend:8080`.
- Host `5432` -> container `postgres:5432`.

Catatan keamanan production:

- Jika PostgreSQL tidak perlu diakses dari luar VPS, jangan publish `5432` ke public internet. Untuk kebutuhan "FE, BE, Postgre berada di port berbeda", publish `5432` hanya ke `127.0.0.1:5432`.
- Jika Nginx host berada di VPS yang sama, publish frontend/backend ke `127.0.0.1` saja agar hanya reverse proxy host yang bisa mengakses.
- Public traffic sebaiknya hanya masuk lewat port `80` dan `443` Nginx host.

## 3. Volume Persisten

Data yang tidak boleh hilang saat image di-rebuild atau container dibuat ulang:

- Data PostgreSQL.
- File upload backend, termasuk dokumen import soal, gambar bank soal, gambar profil/pengumuman, dan dokumen pengumuman.

Volume yang harus dibuat:

- `postgres_data` untuk `/var/lib/postgresql/data`.
- `backend_uploads` untuk `/app/public/uploads`.

Backend harus dijalankan dengan:

```env
APP_DIR=/app
UPLOAD_DIR=/app/public/uploads
```

Dengan konfigurasi ini, seluruh file upload berada di volume `backend_uploads`, bukan di layer image.

## 4. File yang Perlu Dibuat

Buat file berikut:

- `.dockerignore`
- `backend/Dockerfile`
- `frontend/Dockerfile`
- `frontend/nginx.conf`
- `docker-compose.yml`
- `.env.example`

Opsional tetapi direkomendasikan:

- `deploy/nginx/web-app-cbt.conf` untuk contoh reverse proxy Nginx host.
- `backend/docker-entrypoint.sh` jika ingin backend menunggu DB atau menjalankan preflight. Untuk implementasi awal, tidak wajib karena Compose bisa memakai healthcheck dan service dependency.

## 5. Isi `.dockerignore`

Buat `.dockerignore` di root repo untuk memperkecil build context.

```dockerignore
.git
.gitignore
planning.md
random
white_box_documentation.md
white_box_documentation_v2.xlsx

**/node_modules
**/dist
**/coverage
**/.env
**/.env.*

backend/tmp
backend/bin
backend/*.exe
backend/cmd/http/public/uploads

frontend/node_modules
frontend/dist
```

Catatan:

- Jangan copy upload dari repo ke image production.
- Upload harus hidup di volume `backend_uploads`.

## 6. Backend Dockerfile

Buat `backend/Dockerfile`.

Tujuan:

- Multi-stage build.
- Stage builder memakai image Go.
- Stage runtime kecil.
- Binary dijalankan dari `/app`.
- Folder `/app/public/uploads` tersedia dan dapat di-mount volume.
- Migrasi SQL ikut tersedia di image agar service `migrate` bisa memakainya.
- Goose tersedia untuk service `migrate`.

Rencana isi:

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/web-app-cbt ./cmd/http

FROM golang:1.25-bookworm AS goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata wget \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /app/public/uploads/image /app/public/uploads/document

COPY --from=builder /out/web-app-cbt /app/web-app-cbt
COPY --from=goose /go/bin/goose /usr/local/bin/goose
COPY internal/db/migrations /app/internal/db/migrations

ENV APP_DIR=/app
ENV UPLOAD_DIR=/app/public/uploads
ENV ENV=production

EXPOSE 8080

CMD ["/app/web-app-cbt"]
```

Catatan implementasi:

- Dockerfile berada di folder `backend`, jadi build context di compose harus `./backend`.
- Karena context `./backend`, `COPY . .` hanya menyalin isi backend.
- Binary Go harus dibangun dari `./cmd/http`.
- Jika `go 1.25.7` belum tersedia di Docker Hub, gunakan tag Go terdekat yang tersedia dan kompatibel, misalnya `golang:1.25-bookworm` atau `golang:1.24-bookworm`. Setelah mengganti tag, jalankan build untuk memastikan sukses.

## 7. Frontend Dockerfile

Buat `frontend/Dockerfile`.

Tujuan:

- Build React/Vite dengan Node.
- Serve static files dengan Nginx.
- `VITE_API_URL` diisi saat build.

Rencana isi:

```dockerfile
# syntax=docker/dockerfile:1

FROM node:22-bookworm-slim AS builder

WORKDIR /app

ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.27-alpine AS runtime

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

Catatan:

- Jika API akan diakses melalui domain yang sama dengan prefix `/api`, gunakan `VITE_API_URL=https://domain.com/api`.
- Jika API akan diakses melalui subdomain/backend port langsung, gunakan `VITE_API_URL=https://api.domain.com` atau `http://domain.com:8080`.
- Karena frontend memakai cookie auth (`withCredentials: true`), lebih baik production memakai domain yang sama dengan path `/api` agar cookie dan CORS lebih sederhana.

## 8. Frontend Nginx Internal

Buat `frontend/nginx.conf`.

Tujuan:

- Serve SPA React.
- Fallback route ke `index.html`.
- Cache asset build.

Rencana isi:

```nginx
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /assets/ {
        try_files $uri =404;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location = /favicon.ico {
        try_files $uri =404;
    }
}
```

## 9. Docker Compose Production

Buat `docker-compose.yml` di root repo.

Rencana isi:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: web-app-cbt-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "127.0.0.1:${POSTGRES_PORT:-5432}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 20s

  migrate:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: web-app-cbt-migrate
    profiles:
      - migrate
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      POSTGRES_DBURL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
    command:
      [
        "goose",
        "-dir",
        "/app/internal/db/migrations",
        "postgres",
        "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable",
        "up"
      ]

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: web-app-cbt-backend
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      ENV: production
      APP_DIR: /app
      UPLOAD_DIR: /app/public/uploads
      BASE_URL: ${BACKEND_PUBLIC_URL}
      POSTGRES_DBURL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
      JWT_ACCESS_SECRET: ${JWT_ACCESS_SECRET}
      JWT_REFRESH_SECRET: ${JWT_REFRESH_SECRET}
      TRUSTED_ORIGINS: ${TRUSTED_ORIGINS}
      TZ: Asia/Jakarta
    ports:
      - "127.0.0.1:${BACKEND_PORT:-8080}:8080"
    volumes:
      - backend_uploads:/app/public/uploads
    healthcheck:
      test: ["CMD-SHELL", "wget --spider -q http://localhost:8080/auth/me || test $? -eq 8"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 20s

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        VITE_API_URL: ${VITE_API_URL}
    container_name: web-app-cbt-frontend
    restart: unless-stopped
    depends_on:
      - backend
    ports:
      - "127.0.0.1:${FRONTEND_PORT:-3000}:80"

volumes:
  postgres_data:
  backend_uploads:
```

Catatan penting:

- Healthcheck backend memakai `wget`. Karena endpoint `/auth/me` normalnya bisa mengembalikan `401`, command menganggap exit code `8` dari `wget` sebagai tanda HTTP server sudah merespons. Jika nanti ada endpoint health khusus, ganti menjadi `GET /healthz` dan cek status `200`.
- `POSTGRES_PASSWORD` dipakai di connection string URL. Gunakan password URL-safe, misalnya huruf, angka, `_`, dan `-`. Hindari karakter seperti `@`, `:`, `/`, `?`, `#`, dan `%` kecuali implementor menambahkan URL encoding.
- `migrate` memakai profile agar tidak otomatis berjalan setiap `docker compose up`.
- Jalankan migration manual setelah database hidup.
- Build backend dilakukan dua kali oleh `migrate` dan `backend`, tapi Docker cache akan reuse layer yang sama.

## 10. `.env.example`

Buat `.env.example` di root repo.

Rencana isi:

```env
# PostgreSQL
POSTGRES_DB=web_app_cbt
POSTGRES_USER=web_app_cbt
POSTGRES_PASSWORD=change-this-long-random-password
POSTGRES_PORT=5432

# Published local ports on VPS
FRONTEND_PORT=3000
BACKEND_PORT=8080

# Public URLs
# Pilihan recommended: satu domain dengan reverse proxy path /api.
APP_PUBLIC_URL=https://example.com
BACKEND_PUBLIC_URL=https://example.com/api
VITE_API_URL=https://example.com/api

# CORS + CSRF trusted origins.
# Isi origin frontend, tanpa path.
TRUSTED_ORIGINS=https://example.com

# JWT secrets. Generate minimal 32 bytes random.
JWT_ACCESS_SECRET=change-this-access-secret
JWT_REFRESH_SECRET=change-this-refresh-secret
```

Saat deploy, copy menjadi `.env`:

```bash
cp .env.example .env
nano .env
```

Generate secret di VPS:

```bash
openssl rand -base64 48
```

## 11. Reverse Proxy Nginx Host di VPS

Nginx host berjalan langsung di Ubuntu, bukan di compose. Nginx host menerima traffic public `80/443`, lalu meneruskan:

- `/` ke frontend container di `127.0.0.1:3000`.
- `/api/` ke backend container di `127.0.0.1:8080`.

Buat contoh file `deploy/nginx/web-app-cbt.conf`.

Rencana isi:

```nginx
server {
    listen 80;
    server_name example.com www.example.com;

    client_max_body_size 20m;

    location /api/ {
        proxy_pass http://127.0.0.1:8080/;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8080/uploads/;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }

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

Jika memakai HTTPS dengan Certbot:

```bash
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx
sudo cp deploy/nginx/web-app-cbt.conf /etc/nginx/sites-available/web-app-cbt.conf
sudo ln -s /etc/nginx/sites-available/web-app-cbt.conf /etc/nginx/sites-enabled/web-app-cbt.conf
sudo nginx -t
sudo systemctl reload nginx
sudo certbot --nginx -d example.com -d www.example.com
```

Setelah HTTPS aktif, pastikan `.env` memakai:

```env
APP_PUBLIC_URL=https://example.com
BACKEND_PUBLIC_URL=https://example.com/api
VITE_API_URL=https://example.com/api
TRUSTED_ORIGINS=https://example.com
```

## 12. Penyesuaian Backend untuk Prefix `/api`

Dengan konfigurasi reverse proxy:

```nginx
location /api/ {
    proxy_pass http://127.0.0.1:8080/;
}
```

Request browser ke `/api/auth/login` akan diteruskan ke backend sebagai `/auth/login`. Jadi backend tidak perlu diubah untuk mengenal prefix `/api`.

Untuk media upload:

- Frontend membentuk URL dari `VITE_API_URL + /uploads/...`.
- Jika `VITE_API_URL=https://example.com/api`, media menjadi `https://example.com/api/uploads/...`.
- Nginx `location /api/` akan meneruskan `/api/uploads/...` menjadi `/uploads/...` ke backend.
- `location /uploads/` di contoh Nginx hanya fallback jika nanti `BASE_URL` atau file lama memakai URL tanpa `/api`.

## 13. Cookie, CORS, dan CSRF

Backend memakai:

- `TRUSTED_ORIGINS` untuk CORS.
- `TRUSTED_ORIGINS` juga untuk CSRF protection.
- Cookie secure aktif jika `ENV` bukan `dev` dan bukan kosong.

Production wajib:

- Pakai HTTPS.
- Set `ENV=production`.
- Set `TRUSTED_ORIGINS` berisi origin frontend, contoh `https://example.com`.
- Jangan isi path di `TRUSTED_ORIGINS`; gunakan origin saja.

Contoh benar:

```env
TRUSTED_ORIGINS=https://example.com,https://www.example.com
```

Contoh salah:

```env
TRUSTED_ORIGINS=https://example.com/api
```

## 14. Migrasi Database

Migration tidak dijalankan otomatis oleh backend. Implementor harus menyediakan langkah manual.

Urutan deploy awal:

```bash
docker compose up -d postgres
docker compose --profile migrate run --rm migrate
docker compose up -d --build backend frontend
```

Urutan update aplikasi:

```bash
git pull
docker compose build backend frontend
docker compose --profile migrate run --rm migrate
docker compose up -d backend frontend
docker image prune -f
```

Jangan hapus volume saat update.

Perintah yang tidak boleh dipakai saat update biasa:

```bash
docker compose down -v
docker volume rm web-app-cbt_postgres_data
docker volume rm web-app-cbt_backend_uploads
```

## 15. Backup dan Restore

Backup database:

```bash
mkdir -p backups
docker compose exec -T postgres pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" > backups/web_app_cbt_$(date +%Y%m%d_%H%M%S).sql
```

Backup uploads:

```bash
mkdir -p backups
docker run --rm -v web-app-cbt_backend_uploads:/data -v "$PWD/backups:/backup" alpine tar czf /backup/uploads_$(date +%Y%m%d_%H%M%S).tar.gz -C /data .
```

Restore database ke database kosong:

```bash
docker compose exec -T postgres psql -U "$POSTGRES_USER" "$POSTGRES_DB" < backups/file.sql
```

Restore uploads:

```bash
docker run --rm -v web-app-cbt_backend_uploads:/data -v "$PWD/backups:/backup" alpine sh -c "cd /data && tar xzf /backup/uploads_file.tar.gz"
```

Catatan:

- Nama volume dapat berbeda tergantung nama folder project. Cek dengan `docker volume ls`.
- Lakukan backup sebelum migration besar atau update production.

## 16. Deployment VPS Ubuntu dari Nol

Langkah server:

```bash
sudo apt update
sudo apt install -y ca-certificates curl git nginx
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo tee /etc/apt/keyrings/docker.asc >/dev/null
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo usermod -aG docker $USER
```

Logout dan login ulang agar user bisa menjalankan Docker tanpa sudo.

Clone project:

```bash
git clone <repo-url> web-app-cbt
cd web-app-cbt
cp .env.example .env
nano .env
```

Build dan start:

```bash
docker compose build
docker compose up -d postgres
docker compose --profile migrate run --rm migrate
docker compose up -d backend frontend
docker compose ps
```

Pasang Nginx reverse proxy:

```bash
sudo cp deploy/nginx/web-app-cbt.conf /etc/nginx/sites-available/web-app-cbt.conf
sudo ln -s /etc/nginx/sites-available/web-app-cbt.conf /etc/nginx/sites-enabled/web-app-cbt.conf
sudo nginx -t
sudo systemctl reload nginx
```

Aktifkan HTTPS:

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d example.com -d www.example.com
```

## 17. Validasi Setelah Deploy

Jalankan di VPS:

```bash
docker compose ps
docker compose logs --tail=100 backend
docker compose logs --tail=100 frontend
docker compose logs --tail=100 postgres
```

Cek port lokal:

```bash
curl -I http://127.0.0.1:3000
curl -I http://127.0.0.1:8080/auth/me
```

Cek reverse proxy:

```bash
curl -I https://example.com
curl -I https://example.com/api/auth/me
```

Ekspektasi:

- Frontend mengembalikan `200`.
- `/api/auth/me` boleh mengembalikan `401` jika belum login, tetapi harus berasal dari backend, bukan `502`.
- Login dari browser berhasil.
- Upload gambar/dokumen berhasil.
- File yang sudah diupload tetap ada setelah:

```bash
docker compose build backend
docker compose up -d backend
```

## 18. Checklist Implementasi

Gunakan checklist ini saat mengerjakan.

1. Buat `.dockerignore` di root.
2. Buat `backend/Dockerfile`.
3. Buat `frontend/Dockerfile`.
4. Buat `frontend/nginx.conf`.
5. Buat `docker-compose.yml`.
6. Buat `.env.example`.
7. Buat folder `deploy/nginx`.
8. Buat `deploy/nginx/web-app-cbt.conf`.
9. Jalankan build lokal:

```bash
docker compose build
```

10. Jalankan PostgreSQL:

```bash
docker compose up -d postgres
```

11. Jalankan migration:

```bash
docker compose --profile migrate run --rm migrate
```

12. Jalankan backend dan frontend:

```bash
docker compose up -d backend frontend
```

13. Cek status:

```bash
docker compose ps
```

14. Cek log:

```bash
docker compose logs --tail=100 backend
```

15. Test frontend:

```bash
curl -I http://127.0.0.1:3000
```

16. Test backend:

```bash
curl -I http://127.0.0.1:8080/auth/me
```

17. Test upload dari aplikasi.
18. Rebuild backend dan pastikan file upload tidak hilang.
19. Dokumentasikan cara deploy/update di README jika diminta.

## 19. Acceptance Criteria

Implementasi dianggap selesai jika:

- `docker compose build` sukses.
- `docker compose up -d postgres` sukses.
- `docker compose --profile migrate run --rm migrate` sukses.
- `docker compose up -d backend frontend` sukses.
- Frontend dapat diakses dari port host `3000`.
- Backend dapat diakses dari port host `8080`.
- PostgreSQL berjalan di port host `5432` atau `127.0.0.1:5432`.
- Backend berhasil connect ke PostgreSQL.
- Login aplikasi berhasil via browser melalui reverse proxy.
- Upload file berhasil.
- File upload tetap ada setelah rebuild image backend.
- Data database tetap ada setelah rebuild image PostgreSQL atau restart container.
- Nginx host dapat proxy:
  - `/` ke frontend.
  - `/api/` ke backend.
  - `/api/uploads/` ke backend upload route.

## 20. Risiko dan Hal yang Harus Diperhatikan

- `VITE_API_URL` adalah build-time env. Jika domain berubah, frontend harus di-build ulang.
- `TRUSTED_ORIGINS` tidak boleh kosong di production jika browser memakai cookie dan request lintas origin.
- `TRUSTED_ORIGINS` harus origin saja, bukan URL dengan path.
- `ENV=production` membuat cookie secure aktif. Tanpa HTTPS, auth cookie bisa gagal.
- Jangan menjalankan `docker compose down -v` di production kecuali memang ingin menghapus seluruh data.
- Jangan simpan `.env` ke git.
- Jangan simpan file upload production ke image.
- Pastikan `client_max_body_size` Nginx minimal lebih besar dari batas upload backend:
  - Image max backend: 5 MB.
  - Document max backend: 10 MB.
  - Nginx disarankan `20m`.
- Jika ingin PostgreSQL benar-benar tidak terbuka ke luar, pakai binding `127.0.0.1:5432:5432`, bukan `5432:5432`.
