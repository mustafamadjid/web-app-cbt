# Planning Refactoring White Box Unit Test Layer Service

Dokumen ini adalah rencana kerja untuk menyesuaikan unit test layer service dengan teknik white box pada `revisi_modul_white_box_final.md`.

Target utama:

- Refactor isi test yang sudah ada di setiap folder `test` layer service, bukan membuat struktur folder baru.
- Semua test menggunakan table driven test.
- Semua assertion menggunakan `github.com/stretchr/testify/assert`.
- Nama skenario wajib mengikuti teknik:
  - Basis Path Testing: `"Path 1 -> skenario"`
  - Branch Coverage: `"Branch 1 -> skenario"`
  - Loop Coverage: `"Loop 1 -> skenario"`
- Test tetap fokus pada layer service. Jangan memindahkan fokus ke handler HTTP, database integration, atau adapter.

## 1. Acuan Teknik Per Folder Service

Gunakan mapping ini ketika refactor test:

| Modul | Fitur di modul | Folder service | Teknik utama |
|---|---|---|---|
| M-01 | Ujian | `backend/internal/core/service/ujian` | Basis Path Testing |
| M-02 | Manajemen Akun Pengguna | `backend/internal/core/service/user` | Basis Path Testing |
| M-03 | Autentikasi dan Otorisasi | `backend/internal/core/service/auth_service` | Basis Path Testing |
| M-04 | Manajemen Import Soal | `backend/internal/core/service/import_soal` | Branch Coverage |
| M-05 | Manajemen Pengumuman | `backend/internal/core/service/pengumuman` | Branch Coverage |
| M-06 | Manajemen Profil Sekolah | `backend/internal/core/service/profil_sekolah` | Branch Coverage |
| M-07 | Manajemen Bank Soal | `backend/internal/core/service/bank_soal` | Branch Coverage |
| M-08 | Manajemen Mata Pelajaran | `backend/internal/core/service/mata_pelajaran` | Branch Coverage |
| M-09 | Manajemen Kelas | `backend/internal/core/service/kelas` | Branch Coverage |
| M-10 | Upload Berkas | `backend/internal/core/service/delete_file_system` dan service yang menghapus/mengelola path file | Branch Coverage |

Folder service yang belum ada di modul revisi:

- `backend/internal/core/service/sesi`
- `backend/internal/core/service/ruang_ujian`
- `backend/internal/core/service/dashboard`
- `backend/internal/core/service/aktivitas_user`

Untuk folder di atas, jangan refactor dulu kecuali ada instruksi tambahan. Jika tetap harus disentuh karena dependency test, pertahankan perilaku lama dan beri catatan bahwa folder belum terpetakan di `revisi_modul_white_box_final.md`.

## 2. Aturan Umum Refactoring Test

Kerjakan hanya isi test di folder yang sudah ada:

- Pakai file test yang sudah tersedia, misalnya `backend/internal/core/service/user/create/test/create_test.go`.
- Jika file test lama namanya mengandung teknik yang salah, boleh rename file agar tidak membingungkan. Contoh: test `ujian` yang bernama `*_branch_coverage_test.go` perlu menjadi basis path atau isinya diubah menjadi basis path.
- Jangan menghapus fake repository/helper jika masih dipakai.
- Jangan mengubah source service kecuali test tidak mungkin dibuat tanpa dependency injection kecil dan perubahan itu benar-benar diperlukan.
- Unit test service memakai fake/mock repository sederhana, bukan database asli.
- Integration test di folder `integration_test` tidak menjadi prioritas refactor white box.

Pola table driven yang wajib dipakai:

```go
tests := []struct {
	name string
	// input
	// fake behavior
	// expected result
}{
	{
		name: "Path 1 -> valid input berhasil",
	},
}

for _, tc := range tests {
	t.Run(tc.name, func(t *testing.T) {
		// arrange
		// act
		// assert
	})
}
```

Assertion:

- Gunakan `assert.NoError(t, err)` untuk sukses.
- Gunakan `assert.ErrorIs(t, err, wantErr)` jika service mengembalikan sentinel error.
- Gunakan `assert.Equal(t, want, got)` untuk nilai deterministik.
- Gunakan `assert.ElementsMatch(t, want, got)` untuk hasil random yang urutan pastinya tidak boleh diasumsikan.
- Boleh tetap memakai `require` hanya untuk setup yang harus berhenti lebih awal, tetapi validasi behavior utama harus memakai `assert`.

## 3. Definisi Praktis Teknik

### 3.1 Basis Path Testing

Dipakai untuk `ujian`, `user`, dan `auth_service`.

Langkah untuk setiap fungsi service:

1. Baca source service.
2. Identifikasi node keputusan: validasi input, pemanggilan repo, transaksi, hasher, branch role, delete file, commit, rollback, error handling.
3. Buat daftar independent path dari awal fungsi sampai return.
4. Buat satu test case untuk setiap path penting.
5. Nama test harus `"Path N -> ..."` dengan urutan dari path error paling awal sampai happy path.

Minimal skenario basis path yang harus dicakup bila ada di source:

- Path input tidak valid.
- Path sanitizer mengubah input lalu validasi berhasil.
- Path validator gagal.
- Path repository pertama error.
- Path repository berikutnya error.
- Path external dependency error, misalnya hasher atau delete file.
- Path transaction begin error.
- Path commit error.
- Path rollback terpanggil saat error setelah transaksi dimulai.
- Path role khusus, misalnya guru/siswa/admin.
- Path happy path lengkap.

### 3.2 Branch Coverage

Dipakai untuk `import_soal`, `pengumuman`, `profil_sekolah`, `bank_soal`, `mata_pelajaran`, `kelas`, dan upload berkas.

Langkah untuk setiap fungsi service:

1. Tandai semua kondisi `if`, `else`, `switch`, dan guard clause.
2. Pastikan setiap kondisi punya test true dan false.
3. Jika ada banyak branch validasi field, buat case per field yang gagal.
4. Jika ada repo atau dependency yang bisa error, buat branch error dan branch sukses.
5. Nama test harus `"Branch N -> ..."`.

Minimal skenario branch coverage yang harus dicakup bila ada di source:

- Branch ID kosong/invalid vs valid.
- Branch request kosong atau field wajib kosong vs lengkap.
- Branch string perlu trim vs sudah bersih.
- Branch tanggal invalid vs valid.
- Branch data lama punya file vs tidak punya file.
- Branch file baru ada vs tidak ada.
- Branch delete file sukses vs gagal.
- Branch repo error vs repo sukses.
- Branch list/filter kosong vs filter terisi.
- Branch hasil kosong vs hasil ada.

### 3.3 Loop Coverage

Loop coverage hanya wajib untuk randomization urutan soal pada proses pengacakan soal.

Lokasi source yang harus diuji:

- `backend/internal/core/service/ujian/soal_ujian/list_soal.go`
- `backend/internal/core/service/ujian/siswa_ujian/soal_ujian/list_soal_ujian_siswa.go`

Loop yang diuji adalah loop Fisher-Yates:

```go
for i := lastElement; i > 0; i-- {
	j := rand.IntN(i + 1)
	soal[i], soal[j] = soal[j], soal[i]
}
```

Skenario loop coverage wajib:

- `"Loop 1 -> acak false tidak masuk loop"`: `acakSoal=false`, jumlah soal lebih dari 1, hasil harus sama persis dengan input.
- `"Loop 2 -> acak true dengan 0 soal tidak masuk loop"`: hasil kosong, tidak panic.
- `"Loop 3 -> acak true dengan 1 soal tidak masuk loop"`: hasil satu soal tetap sama.
- `"Loop 4 -> acak true dengan 2 soal masuk loop 1 iterasi"`: hasil memiliki elemen yang sama, panjang sama, tidak boleh hilang/duplikat.
- `"Loop 5 -> acak true dengan lebih dari 2 soal masuk loop banyak iterasi"`: hasil memiliki elemen yang sama, panjang sama, tidak boleh hilang/duplikat.

Catatan penting:

- Karena random memakai `math/rand/v2` global dan tidak deterministic, jangan assert urutan berubah.
- Assert yang benar untuk random adalah `assert.ElementsMatch(t, expected, got)` dan `assert.Len(t, got, len(expected))`.
- Untuk case non-random, boleh pakai `assert.Equal(t, expected, got)` agar memastikan urutan tidak berubah.
- Loop coverage ini adalah tambahan khusus di dalam modul `ujian`. Teknik utama modul `ujian` tetap Basis Path Testing.

## 4. Rencana Kerja Per Modul

### 4.1 M-01 Ujian - Basis Path Testing + Loop Coverage Khusus

Folder target:

- `backend/internal/core/service/ujian/ujian_penjadwalan/create/test`
- `backend/internal/core/service/ujian/ujian_penjadwalan/get/test`
- `backend/internal/core/service/ujian/ujian_penjadwalan/update/test`
- `backend/internal/core/service/ujian/ujian_penjadwalan/delete/test`
- `backend/internal/core/service/ujian/attempt_ujian/create/test`
- `backend/internal/core/service/ujian/attempt_ujian/get/test`
- `backend/internal/core/service/ujian/attempt_ujian/update/test`
- `backend/internal/core/service/ujian/attempt_ujian/delete/test`
- `backend/internal/core/service/ujian/attempt_ujian/active_attempt/test`
- `backend/internal/core/service/ujian/attempt_ujian/submit_ujian/test`
- `backend/internal/core/service/ujian/attempt_ujian/list_peserta_submitted/test`
- `backend/internal/core/service/ujian/soal_ujian/test`
- `backend/internal/core/service/ujian/statistik_ujian/test`
- `backend/internal/core/service/ujian/siswa_ujian/test`
- `backend/internal/core/service/ujian/siswa_ujian/list/test`
- `backend/internal/core/service/ujian/siswa_ujian/waktu_selesai/test`
- `backend/internal/core/service/ujian/siswa_ujian/soal_ujian/test`
- `backend/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban/test`
- `backend/internal/core/service/ujian/siswa_ujian/jawaban_ujian/get_jawaban/test`
- `backend/internal/core/service/ujian/siswa_ujian/jawaban_ujian/hasil_jawaban/test`
- `backend/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/test`
- `backend/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading/test`
- `backend/internal/core/service/ujian/siswa_ujian/grading_ujian/statistik_ujian/test`
- `backend/internal/core/service/ujian/siswa_ujian/grading_ujian/list/list_ujian_essay_ungraded/test`
- `backend/internal/core/service/ujian/siswa_ujian/grading_ujian/worker/test`

Arahan refactor:

- Semua skenario utama di folder `ujian` harus memakai nama `"Path N -> ..."`.
- Jika ada test lama bernama `"branch N -> ..."` di folder `ujian`, ubah menjadi path.
- Untuk `soal_ujian/list_soal.go` dan `siswa_ujian/soal_ujian/list_soal_ujian_siswa.go`, tambahkan table test loop coverage dengan nama `"Loop N -> ..."`.

Contoh path umum untuk service ujian:

- `"Path 1 -> id tidak valid mengembalikan ErrMissingId"`
- `"Path 2 -> validasi command gagal"`
- `"Path 3 -> repo get gagal"`
- `"Path 4 -> repo create/update/delete gagal"`
- `"Path 5 -> status attempt tidak mengizinkan aksi"`
- `"Path 6 -> waktu ujian sudah selesai"`
- `"Path 7 -> happy path berhasil"`

### 4.2 M-02 User - Basis Path Testing

Folder target:

- `backend/internal/core/service/user/create/test`
- `backend/internal/core/service/user/get/test`
- `backend/internal/core/service/user/update/test`
- `backend/internal/core/service/user/delete/test`
- `backend/internal/core/service/user/reset_password/test`

Arahan skenario:

- `create`: path validasi command, role guru, role siswa, role tidak dikenal jika ada, hasher error, begin tx error, create user error, create profil error, commit error, rollback terpanggil, happy path guru, happy path siswa.
- `get`: path id invalid, filter invalid, repo error, data tidak ditemukan jika service membedakan, happy path guru, happy path siswa, list kosong, list terisi.
- `update`: path id invalid, validasi gagal, user lama tidak ditemukan/error, delete foto lama sukses/gagal jika foto berubah, begin tx error, update user error, update profil error, commit error, happy path guru/siswa.
- `delete`: path id invalid, get user error, delete file error jika ada foto, delete repo error, happy path.
- `reset_password`: path id invalid, password invalid, hasher error, repo error, happy path.

Semua nama case harus `"Path N -> ..."`.

### 4.3 M-03 Auth Service - Basis Path Testing

Folder target:

- `backend/internal/core/service/auth_service/test`

Arahan skenario:

- Login username kosong.
- Login password kosong.
- Repo find user error.
- User tidak ditemukan jika dibedakan dari repo error.
- Password mismatch.
- Generate token/session error jika ada dependency.
- Simpan session error.
- Happy path login berhasil.
- Session validation: token kosong, token invalid, session tidak ditemukan, session expired jika ada, happy path.
- Logout: session id/token invalid, repo delete error, happy path.

Semua nama case harus `"Path N -> ..."`.

### 4.4 M-04 Import Soal - Branch Coverage

Folder target:

- `backend/internal/core/service/import_soal/create_job/test`
- `backend/internal/core/service/import_soal/get_job/test`
- `backend/internal/core/service/import_soal/import_version/test`
- `backend/internal/core/service/import_soal/parser/test`
- `backend/internal/core/service/import_soal/worker/test`

Arahan skenario:

- `create_job`: branch id bank soal invalid/valid, id user invalid/valid, file path kosong/terisi, repo error/sukses.
- `get_job`: branch id invalid/valid, repo error/sukses, data kosong/ada.
- `import_version`: branch version invalid/valid, repo error/sukses, version aktif/tidak aktif jika ada.
- `parser`: branch format dokumen invalid/valid, section soal kosong/ada, pilihan kosong/ada, jawaban benar kosong/ada, rich content text/image/table jika didukung, error media missing, parse sukses.
- `worker`: branch tidak ada job, ambil job error, file tidak bisa dibaca, parser error, import service error, update status gagal, gambar ditemukan/tidak ditemukan, happy path.

Semua nama case harus `"Branch N -> ..."`.

### 4.5 M-05 Pengumuman - Branch Coverage

Folder target:

- `backend/internal/core/service/pengumuman/create/test`
- `backend/internal/core/service/pengumuman/get/test`
- `backend/internal/core/service/pengumuman/update/test`
- `backend/internal/core/service/pengumuman/delete/test`
- `backend/internal/core/service/pengumuman/date_validation/test`

Arahan skenario:

- Create: judul kosong/valid, isi kosong/valid, tanggal invalid/valid, dokumen kosong/ada, repo error/sukses.
- Get: id invalid/valid, filter kosong/terisi, repo error/sukses, result kosong/ada.
- Update: id invalid/valid, field kosong/terisi, dokumen lama ada/tidak, dokumen baru ada/tidak, delete file error/sukses, repo error/sukses.
- Delete: id invalid/valid, pengumuman punya dokumen/tidak, delete file error/sukses, repo delete error/sukses.
- Date validation: tanggal mulai setelah tanggal selesai, tanggal sama jika valid, format kosong/valid.

Semua nama case harus `"Branch N -> ..."`.

### 4.6 M-06 Profil Sekolah - Branch Coverage

Folder target:

- `backend/internal/core/service/profil_sekolah/get/test`
- `backend/internal/core/service/profil_sekolah/update/test`

Arahan skenario:

- Get: repo error, data tidak ada jika service membedakan, data ada.
- Update: nama sekolah kosong/valid, alamat kosong/valid jika wajib, kontak/email invalid/valid jika ada, logo lama ada/tidak, logo baru ada/tidak, delete file error/sukses jika ada, repo error/sukses.

Semua nama case harus `"Branch N -> ..."`.

### 4.7 M-07 Bank Soal - Branch Coverage

Folder target:

- `backend/internal/core/service/bank_soal/create/test`
- `backend/internal/core/service/bank_soal/get/test`
- `backend/internal/core/service/bank_soal/update/test`
- `backend/internal/core/service/bank_soal/delete/test`

Arahan skenario:

- Create: nama kosong/valid, id mapel invalid/valid, id kelas invalid/valid, id pengguna invalid/valid, repo error/sukses.
- Get: id invalid/valid, filter kosong/terisi, uploaded flag true/false jika ada, repo error/sukses, result kosong/ada.
- Update: id invalid/valid, nama kosong/valid, id mapel/id kelas invalid/valid, repo error/sukses.
- Delete: id invalid/valid, repo error/sukses.

Semua nama case harus `"Branch N -> ..."`.

### 4.8 M-08 Mata Pelajaran - Branch Coverage

Folder target:

- `backend/internal/core/service/mata_pelajaran/create/test`
- `backend/internal/core/service/mata_pelajaran/get/test`
- `backend/internal/core/service/mata_pelajaran/update/test`
- `backend/internal/core/service/mata_pelajaran/delete/test`

Arahan skenario:

- Create: nama mapel kosong/valid, kode mapel kosong/valid jika ada, repo error/sukses.
- Get: id invalid/valid, filter kosong/terisi, repo error/sukses, result kosong/ada.
- Update: id invalid/valid, field kosong/valid, repo error/sukses.
- Delete: id invalid/valid, repo error/sukses.

Semua nama case harus `"Branch N -> ..."`.

### 4.9 M-09 Kelas - Branch Coverage

Folder target:

- `backend/internal/core/service/kelas/create/test`
- `backend/internal/core/service/kelas/get/test`
- `backend/internal/core/service/kelas/update/test`
- `backend/internal/core/service/kelas/delete/test`

Arahan skenario:

- Create: nama kelas kosong/valid, tingkat/jurusan invalid/valid jika ada, repo error/sukses.
- Get: id invalid/valid, filter kosong/terisi, repo error/sukses, result kosong/ada.
- Update: id invalid/valid, field kosong/valid, repo error/sukses.
- Delete: id invalid/valid, repo error/sukses.

Semua nama case harus `"Branch N -> ..."`.

### 4.10 M-10 Upload Berkas - Branch Coverage

Folder target utama:

- `backend/internal/core/service/delete_file_system/test`

Folder terkait file pada service lain tetap mengikuti modul masing-masing:

- `user` tetap basis path walaupun ada delete foto.
- `pengumuman` tetap branch coverage.
- `profil_sekolah` tetap branch coverage.
- `import_soal/worker` tetap branch coverage.

Arahan skenario untuk `delete_file_system`:

- `"Branch 1 -> path traversal ditolak"`
- `"Branch 2 -> filepath.Rel error jika bisa disimulasikan"`
- `"Branch 3 -> file tidak ada mengembalikan error"`
- `"Branch 4 -> file valid berhasil dihapus"`
- `"Branch 5 -> path dengan prefix uploads dibersihkan dengan benar"`

Catatan:

- Gunakan `t.TempDir()` untuk filesystem test.
- Jangan memakai path absolut berbahaya.
- Pastikan test tidak menghapus file di luar temporary directory.

## 5. Checklist Eksekusi

Kerjakan bertahap agar mudah diverifikasi:

1. Pastikan semua test service yang disentuh memakai table driven test.
2. Refactor `ujian` ke Basis Path Testing.
3. Tambahkan Loop Coverage untuk randomization soal.
4. Refactor `user` ke Basis Path Testing.
5. Refactor `auth_service` ke Basis Path Testing.
6. Refactor `import_soal` ke Branch Coverage.
7. Refactor `pengumuman` ke Branch Coverage.
8. Refactor `profil_sekolah` ke Branch Coverage.
9. Refactor `bank_soal` ke Branch Coverage.
10. Refactor `mata_pelajaran` ke Branch Coverage.
11. Refactor `kelas` ke Branch Coverage.
12. Refactor `delete_file_system` ke Branch Coverage.
13. Jalankan test per folder kecil dulu, lalu seluruh service.

## 6. Perintah Verifikasi

Jalankan dari folder `backend`.

Per modul:

```powershell
go test ./internal/core/service/ujian/...
go test ./internal/core/service/user/...
go test ./internal/core/service/auth_service/...
go test ./internal/core/service/import_soal/...
go test ./internal/core/service/pengumuman/...
go test ./internal/core/service/profil_sekolah/...
go test ./internal/core/service/bank_soal/...
go test ./internal/core/service/mata_pelajaran/...
go test ./internal/core/service/kelas/...
go test ./internal/core/service/delete_file_system/...
```

Seluruh service:

```powershell
go test ./internal/core/service/...
```

Opsional coverage:

```powershell
go test ./internal/core/service/... -cover
```

## 7. Definition of Done

Refactor dianggap selesai jika:

- Semua folder yang terpetakan di modul revisi sudah memakai teknik yang sesuai.
- Tidak ada test di `ujian`, `user`, atau `auth_service` yang skenario utamanya masih bernama `"Branch N -> ..."`.
- Tidak ada test di `import_soal`, `pengumuman`, `profil_sekolah`, `bank_soal`, `mata_pelajaran`, `kelas`, atau `delete_file_system` yang skenario utamanya masih bernama `"Path N -> ..."`, kecuali helper lama yang belum menjadi target dan tidak menguji service utama.
- Randomization urutan soal punya Loop Coverage dengan `"Loop N -> ..."` pada dua service pengacakan.
- Semua test yang disentuh berbentuk table driven.
- Semua assertion behavior utama memakai `assert`.
- `go test ./internal/core/service/...` lulus dari folder `backend`.

