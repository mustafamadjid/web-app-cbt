-- +goose Up
-- +goose StatementBegin

-- 1) Kalau sudah ada data > 1 row:
--    - kalau ada id=1, keep itu
--    - kalau tidak ada id=1, ambil row terbaru, jadikan id=1
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM profil_sekolah WHERE id_profil = 1) THEN
    DELETE FROM profil_sekolah WHERE id_profil <> 1;
  ELSE
    IF EXISTS (SELECT 1 FROM profil_sekolah) THEN
      -- pilih 1 row terbaru
      WITH chosen AS (
        SELECT id_profil
        FROM profil_sekolah
        ORDER BY updated_at DESC, id_profil DESC
        LIMIT 1
      )
      UPDATE profil_sekolah
      SET id_profil = 1
      WHERE id_profil IN (SELECT id_profil FROM chosen);

      -- hapus sisanya
      DELETE FROM profil_sekolah WHERE id_profil <> 1;
    END IF;
  END IF;
END $$;

-- 2) Seed id=1 kalau tabel kosong (tapi wajib isi kolom NOT NULL)
--    Kamu HARUS ganti nilai default ini ke yang masuk akal untuk sistemmu.
INSERT INTO profil_sekolah (
  id_profil,
  email_sekolah,
  no_telp_sekolah,
  kepala_sekolah,
  waka_sekolah,
  nama_sekolah,
  alamat_sekolah,
  logo_sekolah
) VALUES (
  1,
  'change-me@example.com',
  'change-me',
  'change-me',
  'change-me',
  'change-me',
  'change-me',
  NULL
)
ON CONFLICT (id_profil) DO NOTHING;

-- 3) Paksa id_profil cuma boleh 1
ALTER TABLE profil_sekolah
  ADD CONSTRAINT profil_sekolah_id_must_be_one CHECK (id_profil = 1);

-- 4) Bikin insert default mengarah ke id=1 (opsional, tapi enak)
ALTER TABLE profil_sekolah
  ALTER COLUMN id_profil SET DEFAULT 1;

-- (opsional tapi rapi) pastikan sequence BIGSERIAL tidak bikin konflik di masa depan
-- kalau suatu saat default sequence kepanggil, dia bisa generate 1 lagi dan conflict.
-- Jadi kita set sequence start > 1.
SELECT setval(pg_get_serial_sequence('profil_sekolah','id_profil'), 2, false);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_sekolah
  ALTER COLUMN id_profil DROP DEFAULT;

ALTER TABLE profil_sekolah
  DROP CONSTRAINT IF EXISTS profil_sekolah_id_must_be_one;

-- (Down tidak mengembalikan data yang sudah dihapus saat Up)
-- +goose StatementEnd
