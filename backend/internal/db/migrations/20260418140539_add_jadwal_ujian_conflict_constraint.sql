-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE jadwal_ujian
ADD CONSTRAINT excl_jadwal_ujian_ruangan_sesi_waktu_active
EXCLUDE USING gist (
    id_ruangan WITH =,
    id_sesi WITH =,
    tstzrange(waktu_mulai, waktu_selesai, '[)') WITH &&
)
WHERE (status_ujian <> 'DIBATALKAN');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
DROP CONSTRAINT IF EXISTS excl_jadwal_ujian_ruangan_sesi_waktu_active;
-- +goose StatementEnd
