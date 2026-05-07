-- +goose Up
-- +goose StatementBegin
WITH siswa_credentials (
    nisn,
    username,
    password_hash
) AS (
    VALUES
        ('0097280820', 'ahmad_zahran', '$2a$10$KQ.mVWCYWSRJ/v2/5GSj8O54d64SVXdAqCMCzC.3EBoYcA1sZIXR6'),
        ('0109179220', 'luqman_tasbih', '$2a$10$jBtrQkEVdhDkqp/FCDnIzOxEcw4iARpT2BNJ6XLn8ZJj8miGlwUaa'),
        ('0104218005', 'luthfi_rifatul', '$2a$10$nt8SqSC7cNB/yDL/WG8M2OhtVrXRFJIgCHzKBKmtdwv6Hcivl79Hu'),
        ('3106188521', 'm_faiz', '$2a$10$BF7wFD0RXgDcJ76gfwcYzeZuyUJuiQaAkjihVSE4.sVQ7yXEsgKDO'),
        ('0106234553', 'muhammad_zaky', '$2a$10$igS7DCa4qgcApXyNXqC62uKJ5SPN82PnBGImXgEZ57jEbI5r9CrwK'),
        ('0096370386', 'nada_syafiyah', '$2a$10$nTodifpy3Fie8amBT19Oo.BjZ5XVUx9OqwAljVbcJasSBwJp3/4Hy'),
        ('0104581493', 'ulwan_fadhil', '$2a$10$QIwY5wb6GWBbr3qWM3iJYOdRm2pw1gyEZwZsQa.ZCWRwZr6KostxC'),
        ('3104727874', 'zaskia_aliya', '$2a$10$saEVeIeuYr.SLN4icG5OGOopLxykB6kqhvrO5coKaW6C8.YzFXx2e'),
        ('0092136891', 'sabrina_jihania', '$2a$10$S58wy2OwK1jiPe0MzwVwR./wFamIdOzLKTJTTi8cQZj.9KM9kalNW'),
        ('0106300360', 'm_fatih', '$2a$10$8YmR/Lvv7wgrG/DL4zN33ePH/CDd2vJLz50sxqaWNXxLedQgAYTiu'),
        ('3105434767', 'ariq_azam', '$2a$10$QEVyqv.Sj4ot59.hDSbd6OHRvaujA8F3xth9EjdW6KRqzvD62WCUW'),
        ('0104424133', 'muhammad_fanani', '$2a$10$8Y14Y4PQEMK41R7/W.RRveyQJ6ti/mMqco0WZPmHFQZOp620SbB5q'),
        ('0096839397', 'amir_hasan', '$2a$10$lupUJ6Fg7HGNex0zfVQKDeX/TbSGcAcE8lv5F5Lp2sKKfKBaYoFeq'),
        ('0084972663', 'ghufran_melsandri', '$2a$10$nzTICxOiqtyY8Q/QUjrMw.AhODMi3yxrSFvbIh6.BjSEOMSlHuYhG'),
        ('0098951684', 'muhammad_joan', '$2a$10$X/R80CWAdZ71SY32MU0TuuzplPuCw0FcxiPTUd/oKdGrBhqKbPaKK'),
        ('0097136942', 'muhammad_satrio', '$2a$10$xzMjDLL42RbbAyVP74QnyunksSfqqcv7pODncRj1FQoqc7pj68aQC'),
        ('0088305862', 'naufal_aziz', '$2a$10$05ARJKfKEx/oChJNTWYu5.8vk4huv89bsFZsuFwQdxkYbQU4fuDHC')
)
UPDATE pengguna p
SET
    username = sc.username,
    password = sc.password_hash,
    updated_at = now()
FROM profil_siswa ps
JOIN siswa_credentials sc
    ON sc.nisn = ps.nisn
WHERE p.id_pengguna = ps.id_pengguna;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
WITH siswa_credentials(nisn) AS (
    VALUES
        ('0097280820'),
        ('0109179220'),
        ('0104218005'),
        ('3106188521'),
        ('0106234553'),
        ('0096370386'),
        ('0104581493'),
        ('3104727874'),
        ('0092136891'),
        ('0106300360'),
        ('3105434767'),
        ('0104424133'),
        ('0096839397'),
        ('0084972663'),
        ('0098951684'),
        ('0097136942'),
        ('0088305862')
)
UPDATE pengguna p
SET
    username = sc.nisn,
    password = '$2a$10$ieWQj7R9PXw4Q1UxDenBhedXUSMk1w2P7lt6cNIIPcqVIXQ8XTwci',
    updated_at = now()
FROM profil_siswa ps
JOIN siswa_credentials sc
    ON sc.nisn = ps.nisn
WHERE p.id_pengguna = ps.id_pengguna;
-- +goose StatementEnd
