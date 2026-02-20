-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pengumuman (
    id_pengumuman BIGSERIAL PRIMARY KEY,
    id_pengguna BIGINT NOT NULL,
    judul_pengumuman VARCHAR(100) NOT NULL,
    isi_pengumuman TEXT NOT NULL,
    tanggal_rilis_pengumuman DATE NOT NULL,
    dokumen_pengumuman VARCHAR(255),

    CONSTRAINT fk_id_pengguna
        FOREIGN KEY (id_pengguna)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);


CREATE INDEX IF NOT EXISTS idx_pengumuman_id_pengguna ON pengumuman(id_pengguna);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pengumuman;
-- +goose StatementEnd
