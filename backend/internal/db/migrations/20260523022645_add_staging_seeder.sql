-- +goose Up
INSERT INTO pengguna (foto, nama_lengkap, jenis_kelamin, username, password, email, no_hp, id_role, status_akun)
VALUES (NULL, 'Administrator', 1, 'superadmin', '$2a$10$gBRZEtuWGIycsVSBFkXkAOgucHFl2NDcXn3kB.7Ss2P7HWlTTln9e', 'superadmin@example.com', '081267668224', 1, 'AKTIF')
ON CONFLICT (username) DO NOTHING;

-- +goose Down
DELETE FROM pengguna WHERE username = 'superadmin';
