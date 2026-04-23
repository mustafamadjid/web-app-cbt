-- +goose Up
-- +goose StatementBegin
INSERT INTO role (nama_role) VALUES ('ADMIN'), ('GURU'), ('SISWA');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM profil_guru
WHERE id_pengguna IN (
    SELECT id_pengguna
    FROM pengguna
    WHERE id_role IN (
        SELECT id_role
        FROM role
        WHERE nama_role IN ('ADMIN', 'GURU', 'SISWA')
    )
);

DELETE FROM profil_siswa
WHERE id_pengguna IN (
    SELECT id_pengguna
    FROM pengguna
    WHERE id_role IN (
        SELECT id_role
        FROM role
        WHERE nama_role IN ('ADMIN', 'GURU', 'SISWA')
    )
);

DELETE FROM pengguna
WHERE id_role IN (
    SELECT id_role
    FROM role
    WHERE nama_role IN ('ADMIN', 'GURU', 'SISWA')
);

DELETE FROM role;
-- +goose StatementEnd
