-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bank_soal_version (
    id_bank_soal_version BIGSERIAL PRIMARY KEY,
    id_bank_soal BIGINT NOT NULL,
    version_no INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by BIGINT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_bank_soal_version_bank'
          AND conrelid = 'bank_soal_version'::regclass
    ) THEN
        ALTER TABLE bank_soal_version
            ADD CONSTRAINT fk_bank_soal_version_bank
            FOREIGN KEY (id_bank_soal)
            REFERENCES bank_soal(id_bank_soal)
            ON UPDATE CASCADE
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uq_bank_soal_version'
          AND conrelid = 'bank_soal_version'::regclass
    ) THEN
        ALTER TABLE bank_soal_version
            ADD CONSTRAINT uq_bank_soal_version UNIQUE (id_bank_soal, version_no);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'ck_bank_soal_version_status'
          AND conrelid = 'bank_soal_version'::regclass
    ) THEN
        ALTER TABLE bank_soal_version
            ADD CONSTRAINT ck_bank_soal_version_status
            CHECK (status IN ('draft', 'published', 'archived'));
    END IF;
END $$;

ALTER TABLE bank_soal
    ADD COLUMN IF NOT EXISTS id_bank_soal_version_aktif BIGINT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_bank_soal_version_aktif'
          AND conrelid = 'bank_soal'::regclass
    ) THEN
        ALTER TABLE bank_soal
            ADD CONSTRAINT fk_bank_soal_version_aktif
            FOREIGN KEY (id_bank_soal_version_aktif)
            REFERENCES bank_soal_version(id_bank_soal_version)
            ON UPDATE CASCADE
            ON DELETE SET NULL;
    END IF;
END $$;

ALTER TABLE isi_soal
    ADD COLUMN IF NOT EXISTS id_bank_soal_version BIGINT NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'isi_soal'
          AND column_name = 'id_bank_soal'
    ) THEN
        INSERT INTO bank_soal_version (id_bank_soal, version_no, status, created_by)
        SELECT b.id_bank_soal, 1, 'published', b.id_pengguna
        FROM bank_soal b
        WHERE EXISTS (
            SELECT 1
            FROM isi_soal s
            WHERE s.id_bank_soal = b.id_bank_soal
        )
          AND NOT EXISTS (
            SELECT 1
            FROM bank_soal_version v
            WHERE v.id_bank_soal = b.id_bank_soal
        );

        UPDATE isi_soal s
        SET id_bank_soal_version = (
            SELECT v.id_bank_soal_version
            FROM bank_soal_version v
            WHERE v.id_bank_soal = s.id_bank_soal
            ORDER BY
                CASE v.status
                    WHEN 'published' THEN 0
                    WHEN 'draft' THEN 1
                    ELSE 2
                END,
                v.version_no DESC,
                v.id_bank_soal_version DESC
            LIMIT 1
        )
        WHERE s.id_bank_soal_version IS NULL;
    END IF;
END $$;

UPDATE bank_soal b
SET id_bank_soal_version_aktif = x.id_bank_soal_version
FROM (
    SELECT DISTINCT ON (v.id_bank_soal)
        v.id_bank_soal,
        v.id_bank_soal_version
    FROM bank_soal_version v
    ORDER BY
        v.id_bank_soal,
        CASE v.status
            WHEN 'published' THEN 0
            WHEN 'draft' THEN 1
            ELSE 2
        END,
        v.version_no DESC,
        v.id_bank_soal_version DESC
) x
WHERE b.id_bank_soal = x.id_bank_soal
  AND b.id_bank_soal_version_aktif IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_isi_soal_bank_soal_version'
          AND conrelid = 'isi_soal'::regclass
    ) THEN
        ALTER TABLE isi_soal
            ADD CONSTRAINT fk_isi_soal_bank_soal_version
            FOREIGN KEY (id_bank_soal_version)
            REFERENCES bank_soal_version(id_bank_soal_version)
            ON UPDATE CASCADE
            ON DELETE CASCADE;
    END IF;
END $$;

ALTER TABLE isi_soal
    DROP CONSTRAINT IF EXISTS fk_soal_bank_soal;

ALTER TABLE isi_soal
    DROP COLUMN IF EXISTS id_bank_soal;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM isi_soal WHERE id_bank_soal_version IS NULL) THEN
        ALTER TABLE isi_soal
            ALTER COLUMN id_bank_soal_version SET NOT NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_bank_soal_version_bank
    ON bank_soal_version(id_bank_soal);

CREATE INDEX IF NOT EXISTS idx_isi_soal_id_bank_soal_version
    ON isi_soal(id_bank_soal_version);

CREATE UNIQUE INDEX IF NOT EXISTS uq_bank_soal_version_one_published
    ON bank_soal_version(id_bank_soal)
    WHERE status = 'published';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_bank_soal_version_one_published;
DROP INDEX IF EXISTS idx_isi_soal_id_bank_soal_version;
DROP INDEX IF EXISTS idx_bank_soal_version_bank;

ALTER TABLE isi_soal
    ADD COLUMN IF NOT EXISTS id_bank_soal BIGINT NULL;

UPDATE isi_soal s
SET id_bank_soal = v.id_bank_soal
FROM bank_soal_version v
WHERE s.id_bank_soal_version = v.id_bank_soal_version
  AND s.id_bank_soal IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_soal_bank_soal'
          AND conrelid = 'isi_soal'::regclass
    ) THEN
        ALTER TABLE isi_soal
            ADD CONSTRAINT fk_soal_bank_soal
            FOREIGN KEY (id_bank_soal)
            REFERENCES bank_soal(id_bank_soal)
            ON UPDATE CASCADE
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM isi_soal WHERE id_bank_soal IS NULL) THEN
        ALTER TABLE isi_soal
            ALTER COLUMN id_bank_soal SET NOT NULL;
    END IF;
END $$;

ALTER TABLE isi_soal
    DROP CONSTRAINT IF EXISTS fk_isi_soal_bank_soal_version;

ALTER TABLE isi_soal
    DROP COLUMN IF EXISTS id_bank_soal_version;

ALTER TABLE bank_soal
    DROP CONSTRAINT IF EXISTS fk_bank_soal_version_aktif;

ALTER TABLE bank_soal
    DROP COLUMN IF EXISTS id_bank_soal_version_aktif;

DROP TABLE IF EXISTS bank_soal_version;
-- +goose StatementEnd
