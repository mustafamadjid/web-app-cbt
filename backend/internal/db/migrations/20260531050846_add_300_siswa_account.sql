-- +goose Up
-- +goose StatementBegin
WITH siswa_data AS (
    SELECT
        n,
        'siswa' || LPAD(n::TEXT, 3, '0') AS username,
        'Siswa ' || LPAD(n::TEXT, 3, '0') AS nama_lengkap,
        CASE WHEN n % 2 = 0 THEN 2 ELSE 1 END AS jenis_kelamin,
        'siswa' || LPAD(n::TEXT, 3, '0') || '@email.com' AS email,
        '08' || LPAD((8100000000 + n)::TEXT, 10, '0') AS no_hp,
        '990' || LPAD(n::TEXT, 7, '0') AS nisn,
        ((n - 1) % 36) + 1 AS no_absen,
        CASE (n - 1) % 3
            WHEN 0 THEN 10
            WHEN 1 THEN 11
            ELSE 12
        END AS tingkat_kelas,
        CASE (n - 1) % 3
            WHEN 0 THEN 2026
            WHEN 1 THEN 2025
            ELSE 2024
        END AS angkatan,
        CASE (n - 1) % 6
            WHEN 0 THEN 'Bandar Lampung'
            WHEN 1 THEN 'Metro'
            WHEN 2 THEN 'Pringsewu'
            WHEN 3 THEN 'Kotabumi'
            WHEN 4 THEN 'Kalianda'
            ELSE 'Menggala'
        END AS tempat_lahir,
        (DATE '2008-01-01' + ((n - 1) % 1095) * INTERVAL '1 day')::DATE AS tanggal_lahir
    FROM generate_series(1, 300) AS gs(n)
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
        sd.username,
        '$2a$10$axS84TtJGIl1OOjXCvI02O/IDMy1hyoeEA.HRNPyjLHGRhBEg7oiK',
        sd.email,
        sd.no_hp,
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
    ON sd.username = up.username
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
    SELECT 'siswa' || LPAD(n::TEXT, 3, '0')
    FROM generate_series(1, 300) AS gs(n)
)
DELETE FROM profil_siswa
WHERE id_pengguna IN (
    SELECT p.id_pengguna
    FROM pengguna p
    JOIN usernames u ON u.username = p.username
);

WITH usernames(username) AS (
    SELECT 'siswa' || LPAD(n::TEXT, 3, '0')
    FROM generate_series(1, 300) AS gs(n)
)
DELETE FROM pengguna
WHERE username IN (SELECT username FROM usernames);
-- +goose StatementEnd
