-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_per_user;

CREATE UNIQUE INDEX IF NOT EXISTS sessions_one_unrevoked_non_admin_per_user
ON sessions(id_pengguna)
WHERE revoked_at IS NULL
  AND role = 'SISWA' OR role = 'GURU';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_non_admin_per_user;

CREATE UNIQUE INDEX IF NOT EXISTS sessions_one_unrevoked_per_user
ON sessions(id_pengguna)
WHERE revoked_at IS NULL;
-- +goose StatementEnd


SELECT id_buku, judul
FROM   buku
WHERE  id_buku IN (
    -- Query tengah: cari id_buku dari transaksi pelanggan Bandar Lampung
    SELECT id_buku
    FROM   transaksi
    WHERE  id_pelanggan IN (
        -- Query terdalam: cari id_pelanggan dari Bandar Lampung
        SELECT id_pelanggan
        FROM   pelanggan
        WHERE  alamat ILIKE '%Bandar Lampung%'
    )
);


