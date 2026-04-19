-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_per_user;
DROP INDEX IF EXISTS sessions_one_unrevoked_non_admin_per_user;

CREATE UNIQUE INDEX sessions_one_unrevoked_non_admin_per_user
ON sessions (id_pengguna)
WHERE revoked_at IS NULL
  AND role IN ('SISWA', 'GURU');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_non_admin_per_user;

CREATE UNIQUE INDEX sessions_one_unrevoked_per_user
ON sessions (id_pengguna)
WHERE revoked_at IS NULL;
-- +goose StatementEnd

-- Backup schema saja (DDL)  
pg_dump -U nama_user -s -d nama_database > nama_file_schema.sql -- Backup data saja (DML)  
pg_dump -U nama_user -a -d nama_database > nama_file_data.sql -- Backup tabel tertentu saja  
pg_dump -U nama_user -t nama_tabel -d nama_database > 
nama_file_tabel.sql -- Backup seluruh database (termasuk role dan konfigurasi global)  
pg_dumpall -U nama_user > nama_file_all.sql 