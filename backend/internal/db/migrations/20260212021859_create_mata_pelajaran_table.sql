-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS mata_pelajaran(
    id_mapel BIGSERIAL PRIMARY KEY,
    id_kelas BIGINT NOT NULL,
    kode_mapel VARCHAR(20) NOT NULL,
    nama_mapel VARCHAR(100) NOT NULL,
    deskripsi TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_mapel_kelas
        FOREIGN KEY (id_kelas)
        REFERENCES kelas(id_kelas)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mata_pelajaran
-- +goose StatementEnd
