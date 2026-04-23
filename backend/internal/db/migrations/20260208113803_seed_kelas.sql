-- +goose Up
-- +goose StatementBegin
INSERT INTO kelas (tingkat_kelas)
VALUES (10),(11),(12)
ON CONFLICT DO NOTHING;

WITH k AS (
    SELECT id_kelas,tingkat_kelas FROM kelas
    WHERE tingkat_kelas BETWEEN 10 AND 12
)

INSERT INTO nama_kelas (id_kelas,nama_kelas)
SELECT id_kelas,CONCAT('Kelas ',tingkat_kelas) FROM k
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM profil_siswa
WHERE id_nama_kelas IN (SELECT id_nama_kelas FROM nama_kelas)
   OR id_kelas IN (SELECT id_kelas FROM kelas);

DELETE FROM nama_kelas;
DELETE FROM kelas;

-- +goose StatementEnd
