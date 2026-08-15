| Komponen | Spesifikasi / Teknologi | Keterangan |
|---|---|---|
| Nama Sistem | Sistem Ujian Sekolah Berbasis Web | Sistem yang diuji pada pengujian performa |
| Jenis Aplikasi | Web Application | Aplikasi diakses melalui browser |
| Frontend | ReactJS | Antarmuka pengguna |
| Backend | Golang | Penyedia REST API |
| Database | PostgreSQL 16 | Penyimpanan data pengguna, soal, ujian, jawaban, dan nilai |
| Reverse Proxy | Traefik v3 | Mengatur routing HTTP/HTTPS |
| Deployment | Docker Compose | Menjalankan setiap service dalam container |
| Server | 4 vCPU, 4 GB RAM | Lingkungan server pengujian |
| Sistem Operasi | Ubuntu Server | Sistem operasi pada server |
| Tool Pengujian | k6 / Locust | Tool untuk simulasi virtual user |
| Target Pengguna | Siswa | Pengujian difokuskan pada alur pengerjaan ujian |