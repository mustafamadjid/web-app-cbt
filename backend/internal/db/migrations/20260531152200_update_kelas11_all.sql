-- +goose Up
-- +goose StatementBegin
WITH target_siswa AS (
    SELECT 'siswa' || LPAD(n::TEXT, 3, '0') AS username
    FROM generate_series(1, 300) AS gs(n)
),
kelas_11 AS (
    SELECT nk.id_nama_kelas
    FROM kelas k
    JOIN nama_kelas nk
        ON nk.id_kelas = k.id_kelas
       AND nk.nama_kelas = 'Kelas 11'
    WHERE k.tingkat_kelas = 11
)
UPDATE profil_siswa ps
SET
    id_nama_kelas = k11.id_nama_kelas,
    updated_at = now()
FROM pengguna p
JOIN target_siswa ts
    ON ts.username = p.username
CROSS JOIN kelas_11 k11
WHERE ps.id_pengguna = p.id_pengguna;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
WITH siswa_data AS (
    SELECT
        'siswa' || LPAD(n::TEXT, 3, '0') AS username,
        CASE (n - 1) % 3
            WHEN 0 THEN 10
            WHEN 1 THEN 11
            ELSE 12
        END AS tingkat_kelas
    FROM generate_series(1, 300) AS gs(n)
),
target_kelas AS (
    SELECT
        sd.username,
        nk.id_nama_kelas
    FROM siswa_data sd
    JOIN kelas k
        ON k.tingkat_kelas = sd.tingkat_kelas
    JOIN nama_kelas nk
        ON nk.id_kelas = k.id_kelas
       AND nk.nama_kelas = CONCAT('Kelas ', sd.tingkat_kelas)
)
UPDATE profil_siswa ps
SET
    id_nama_kelas = tk.id_nama_kelas,
    updated_at = now()
FROM pengguna p
JOIN target_kelas tk
    ON tk.username = p.username
WHERE ps.id_pengguna = p.id_pengguna;
-- +goose StatementEnd
