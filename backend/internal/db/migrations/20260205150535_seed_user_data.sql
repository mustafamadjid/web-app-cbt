-- +goose Up
INSERT INTO pengguna (foto, nama_lengkap, jenis_kelamin, username, password, email, no_hp, id_role, status_akun)
VALUES (NULL, 'Administrator', 1, 'myadmin', '$2a$10$ieWQj7R9PXw4Q1UxDenBhedXUSMk1w2P7lt6cNIIPcqVIXQ8XTwci', 'admin@example.com', "08127234567", 1, 'AKTIF')
ON CONFLICT (username) DO NOTHING;

-- +goose Down
DELETE FROM pengguna WHERE username = 'myadmin';
