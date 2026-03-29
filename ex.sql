CREATE TABLE buku ( 
    id_buku    VARCHAR(10)   NOT NULL, 
    judul      VARCHAR(150)  NOT NULL,  
    kategori   VARCHAR(50), 
    id_penulis VARCHAR(10), 
    harga      NUMERIC(10,2) NOT NULL, 
    stok       INTEGER DEFAULT 50, 
    PRIMARY KEY (id_buku), 
    FOREIGN KEY (id_penulis) REFERENCES penulis(id_penulis) 
        ON DELETE RESTRICT   -- mencegah hapus penulis jika masih ada buku 
        ON UPDATE CASCADE    -- update id_penulis otomatis jika berubah 
); 

-- setelah muncul postgres=# 
CREATE DATABASE db_toko_buku_NIM; 
 -- lalu masuk ke database yang telah dibuat 
\c db_toko_buku_NIM;


CREATE TABLE penulis ( 
    id_penulis  VARCHAR(10)  NOT NULL, 
    nama        VARCHAR(100) NOT NULL, 
    asal_daerah TEXT,              -- TEXT lebih fleksibel dari VARCHAR di PostgreSQL 
    PRIMARY KEY (id_penulis) 
); 

-- Verifikasi: lihat struktur tabel 
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'penulis';

-- Kalau menggunakan CMD atau Powershell, jalankan psql terlebih dahulu. 
psql -U postgres 
 -- setelah muncul postgres=# 
CREATE DATABASE db_toko_buku_NIM; 
 -- lalu masuk ke database yang telah dibuat 
\c db_toko_buku_NIM 

CREATE TABLE buku ( 
    id_buku    VARCHAR(10)   NOT NULL, 
    judul      VARCHAR(150)  NOT NULL,  
    kategori   VARCHAR(50), 
    id_penulis VARCHAR(10), 
    harga      NUMERIC(10,2) NOT NULL, 
    stok       INTEGER DEFAULT 50, 
    PRIMARY KEY (id_buku), 
    FOREIGN KEY (id_penulis) REFERENCES penulis(id_penulis) 
        ON DELETE RESTRICT   -- mencegah hapus penulis jika masih ada buku 
        ON UPDATE CASCADE    -- update id_penulis otomatis jika berubah 
); 