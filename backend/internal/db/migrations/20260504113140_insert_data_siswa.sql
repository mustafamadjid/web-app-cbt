-- +goose Up
-- +goose StatementBegin
WITH siswa_data (
    no_absen,
    nama_lengkap,
    tingkat_kelas,
    jenis_kelamin,
    nisn,
    angkatan,
    tempat_lahir,
    tanggal_lahir
) AS (
    VALUES
        (1, 'Ahmad Zahran Riyandi', 10, 1, '0097280820', 2025, 'Bandar Lampung', DATE '2010-02-14'),
        (2, 'LUQMAN TASBIH RAMADHAN', 10, 1, '0109179220', 2025, 'Bandar Lampung', DATE '2010-05-23'),
        (3, 'Luthfi Rifatul Azzaki', 10, 1, '0104218005', 2025, 'Bandar Lampung', DATE '2010-08-07'),
        (4, 'M. Faiz Rihlinal Arifin', 10, 1, '3106188521', 2025, 'Bandar Lampung', DATE '2010-11-19'),
        (5, 'MUHAMMAD ZAKY IMANI', 10, 1, '0106234553', 2025, 'Bandar Lampung', DATE '2010-03-31'),
        (6, 'Nada Syafiyah', 10, 2, '0096370386', 2025, 'Bandar Lampung', DATE '2010-07-12'),
        (7, 'ULWAN FADHIL', 10, 1, '0104581493', 2025, 'Bandar Lampung', DATE '2010-09-28'),
        (8, 'ZASKIA ALIYA PUTRI', 10, 2, '3104727874', 2025, 'Bandar Lampung', DATE '2010-12-04'),
        (9, 'Sabrina Jihania', 10, 2, '0092136891', 2025, 'Bandar Lampung', DATE '2010-01-17'),
        (10, 'M Fatih Putra Dhafian', 10, 1, '0106300360', 2025, 'Bandar Lampung', DATE '2010-04-09'),
        (11, 'ARIQ AZAM SYAUQI', 10, 1, '3105434767', 2025, 'Bandar Lampung', DATE '2010-06-25'),
        (12, 'MUHAMMAD FANANI HUDHABITH', 10, 1, '0104424133', 2025, 'Bandar Lampung', DATE '2010-10-11'),
        (1, 'AMIR HASAN ERLANGGA', 11, 1, '0096839397', 2024, 'Bandar Lampung', DATE '2009-02-08'),
        (2, 'GHUFRAN MELSANDRI', 11, 1, '0084972663', 2024, 'Bandar Lampung', DATE '2009-05-16'),
        (3, 'MUHAMMAD JOAN ARIFA', 11, 1, '0098951684', 2024, 'Bandar Lampung', DATE '2009-08-29'),
        (4, 'Muhammad Satrio Aji', 11, 1, '0097136942', 2024, 'Bandar Lampung', DATE '2009-11-03'),
        (5, 'NAUFAL AZIZ RAMADHAN', 11, 1, '0088305862', 2024, 'Bandar Lampung', DATE '2009-03-21')
),
upserted_pengguna AS (
    INSERT INTO pengguna (
        foto,
        nama_lengkap,
        jenis_kelamin,
        username,
        password,
        email,
        no_hp,
        id_role,
        status_akun
    )
    SELECT
        NULL,
        sd.nama_lengkap,
        sd.jenis_kelamin,
        sd.nisn,
        '$2a$10$ieWQj7R9PXw4Q1UxDenBhedXUSMk1w2P7lt6cNIIPcqVIXQ8XTwci',
        NULL,
        NULL,
        r.id_role,
        'AKTIF'
    FROM siswa_data sd
    CROSS JOIN role r
    WHERE r.nama_role = 'SISWA'
    ON CONFLICT (username) DO UPDATE SET
        nama_lengkap = EXCLUDED.nama_lengkap,
        jenis_kelamin = EXCLUDED.jenis_kelamin,
        password = EXCLUDED.password,
        email = EXCLUDED.email,
        no_hp = EXCLUDED.no_hp,
        id_role = EXCLUDED.id_role,
        status_akun = EXCLUDED.status_akun,
        updated_at = now()
    RETURNING id_pengguna, username
)
INSERT INTO profil_siswa (
    id_pengguna,
    id_nama_kelas,
    nisn,
    no_absen,
    angkatan,
    tempat_lahir,
    tanggal_lahir
)
SELECT
    up.id_pengguna,
    nk.id_nama_kelas,
    sd.nisn,
    sd.no_absen,
    sd.angkatan,
    sd.tempat_lahir,
    sd.tanggal_lahir
FROM upserted_pengguna up
JOIN siswa_data sd
    ON sd.nisn = up.username
JOIN kelas k
    ON k.tingkat_kelas = sd.tingkat_kelas
JOIN nama_kelas nk
    ON nk.id_kelas = k.id_kelas
   AND nk.nama_kelas = CONCAT('Kelas ', sd.tingkat_kelas)
ON CONFLICT (id_pengguna) DO UPDATE SET
    id_nama_kelas = EXCLUDED.id_nama_kelas,
    nisn = EXCLUDED.nisn,
    no_absen = EXCLUDED.no_absen,
    angkatan = EXCLUDED.angkatan,
    tempat_lahir = EXCLUDED.tempat_lahir,
    tanggal_lahir = EXCLUDED.tanggal_lahir,
    updated_at = now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
WITH usernames(username) AS (
    VALUES
        ('0097280820'),
        ('0109179220'),
        ('0104218005'),
        ('3106188521'),
        ('0106234553'),
        ('0096370386'),
        ('0104581493'),
        ('3104727874'),
        ('0092136891'),
        ('0106300360'),
        ('3105434767'),
        ('0104424133'),
        ('0096839397'),
        ('0084972663'),
        ('0098951684'),
        ('0097136942'),
        ('0088305862')
)
DELETE FROM profil_siswa
WHERE id_pengguna IN (
    SELECT p.id_pengguna
    FROM pengguna p
    JOIN usernames u ON u.username = p.username
);

WITH usernames(username) AS (
    VALUES
        ('0097280820'),
        ('0109179220'),
        ('0104218005'),
        ('3106188521'),
        ('0106234553'),
        ('0096370386'),
        ('0104581493'),
        ('3104727874'),
        ('0092136891'),
        ('0106300360'),
        ('3105434767'),
        ('0104424133'),
        ('0096839397'),
        ('0084972663'),
        ('0098951684'),
        ('0097136942'),
        ('0088305862')
)
DELETE FROM pengguna
WHERE username IN (SELECT username FROM usernames);
-- +goose StatementEnd
