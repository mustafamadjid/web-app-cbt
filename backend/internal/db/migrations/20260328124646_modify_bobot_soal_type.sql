-- +goose Up
-- +goose StatementBegin
ALTER TABLE isi_soal
    ALTER COLUMN bobot_soal TYPE DECIMAL(5,2) USING bobot_soal::DECIMAL(5,2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE isi_soal
    ALTER COLUMN bobot_soal TYPE INTEGER USING bobot_soal::INTEGER;
-- +goose StatementEnd
