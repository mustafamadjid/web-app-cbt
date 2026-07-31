# Identifikasi Unit Test — Folder `backend/internal/core/service`

## Ringkasan Total Skenario

| Metode Pengujian | Total Skenario |
|---|---|
| **Basis Path Testing** (`Path N ->` dan `PN ->`) | **295** |
| **Branch Coverage** (`Branch N ->`) | **200** |
| **Total** | **495** |

---

## Basis Path Testing (295 Skenario)

Skenario diidentifikasi dari nama test case berformat `Path N ->` atau bentuk ringkas `PN ->`, misalnya `P1 ->`.

| Fungsi Test | Skenario | Service/Modul |
|---|---|---|
| `TestParseMarkers` | 13 | import_soal/parser (rich_content) |
| `TestAttemptUjianService_BasisPath` | 13 | ujian/attempt_ujian |
| `TestAuthServiceLoginBasisPath` | 12 | auth_service |
| `TestUpdateUjianService_BasisPathValidation` | 10 | ujian/ujian_penjadwalan/update |
| `TestUpdateBankSoalService_BasisPath` | 10 | bank_soal/update |
| `TestCreateSiswaBasisPath` | 10 | user/create |
| `TestUpdateAttemptUjianService_BasisPath` | 9 | ujian/attempt_ujian |
| `TestValidateParsedSoal` | 9 | import_soal/parser (validasi) |
| `TestGradingUjianService_BasisPath` | 9 | ujian/grading |
| `TestAuthServiceRefreshAccessTokenBasisPath` | 8 | auth_service |
| `TestSaveJawabanUjianService_BasisPath` | 8 | ujian/siswa_ujian/save_jawaban |
| `TestGetSiswaService_ListSiswa` | 7 | user/get |
| `TestCreateGuruBasisPath` | 7 | user/create |
| `TestUpdateGuruBasisPath` | 7 | user/update |
| `TestUpdateSiswaBasisPath` | 7 | user/update |
| `TestCreateSesiService_BasisPath` | 6 | sesi/create |
| `TestSiswaUpdateAttemptUjianService_BasisPath` | 6 | ujian/attempt_ujian |
| `TestGradingUjianWorkerService_BasisPath` | 6 | ujian/grading |
| `TestDeleteUserService_Delete` | 5 | user/delete |
| `TestEssayGradingUjianService_BasisPath` | 5 | ujian/grading/essay_grading |
| `TestUpdateSesiService_BasisPath` | 5 | sesi/update |
| `TestCreateJobService` | 5 | import_soal/create_job |
| `TestGetGuruService_ListGuru` | 5 | user/get |
| `TestExtractParagraphs` | 5 | import_soal/parser (docx) |
| `TestSubmitUjianService_BasisPath` | 5 | ujian/attempt_ujian/submit_ujian |
| `TestListSoalUjianSiswaService_BasisPath` | 4 | ujian/siswa_ujian/list_soal |
| `TestListUjianSiswaService_BasisPath` | 4 | ujian/siswa_ujian/list_ujian |
| `TestCreateBankSoalService_BasisPath` | 4 | bank_soal/create |
| `TestGetActiveAttemptUjianService_BasisPath` | 4 | ujian/attempt_ujian/active_attempt |
| `TestListSoalUjianService_BasisPath` | 4 | ujian/soal_ujian |
| `TestGetUjianService_GetAllUjian_BasisPath` | 4 | ujian/ujian_penjadwalan/get |
| `TestHasilJawabanUjianService_BasisPath` | 3 | ujian/siswa_ujian/hasil_jawaban |
| `TestResetPasswordService_ResetPasswordService` | 3 | user/reset_password |
| `TestDeleteUjianService_BasisPath` | 3 | ujian/ujian_penjadwalan/delete |
| `TestGetUjianService_GetUjianByID_BasisPath` | 3 | ujian/ujian_penjadwalan/get |
| `TestGetWaktuSelesaiService_BasisPath` | 3 | ujian/siswa_ujian/waktu_selesai |
| `TestCompositeGradingUjianExecutor_BasisPath` | 3 | ujian/grading |
| `TestDeleteSesiService_BasisPath` | 3 | sesi/delete |
| `TestGetSesiService_BasisPath` | 3 | sesi/get |
| `TestGetSesiByIdService_BasisPath` | 3 | sesi/get |
| `TestGetJobsByBankSoal` | 3 | import_soal/get_job |
| `TestAuthServiceAdminRevokingSessionBasisPath` | 3 | auth_service |
| `TestDeleteBankSoalService_BasisPath` | 3 | bank_soal/delete |
| `TestGetJobByID` | 3 | import_soal/get_job |
| `TestGradingUjianService_UpsertingToStatistikSoal_BasisPath` | 3 | ujian/grading |
| `TestGradingUjianService_TotalScore_BasisPath` | 3 | ujian/grading |
| `TestAuthServiceLogoutBasisPath` | 3 | auth_service |
| `TestDeleteAttemptUjianService_BasisPath` | 3 | ujian/attempt_ujian |
| `TestGetAttemptUjianService_BasisPath` | 3 | ujian/attempt_ujian |
| `TestGetSesiByKodeService_BasisPath` | 3 | sesi/get |
| `TestListUjianEssayUngradedService_BasisPath` | 3 | ujian/grading/list_ujian_essay_ungraded |
| `TestPesertaUjianSubmittedService_BasisPath` | 3 | ujian/attempt_ujian/list_peserta_submitted |
| `TestGetSiswaService_FindProfilSiswaByID` | 2 | user/get |
| `TestUpdateUser_DeleteFileErrorBasisPath` | 2 | user/update |
| `TestGetGuruService_FindProfilGuruByID` | 2 | user/get |
| `TestUpdateUjianService_BasisPath` | 2 | ujian/ujian_penjadwalan/update |
| `TestGetAllActiveSessionService_BasisPath` | 2 | sesi/get |
| `TestExtractImageFiles` | 2 | import_soal/parser (docx) |
| `TestDeleteUserService_DeleteMany` | 2 | user/delete |
| `TestCreateUjianService_BasisPath` | 2 | ujian/ujian_penjadwalan/create |

---

## Branch Coverage (200 Skenario)

Skenario diidentifikasi dari nama test case berformat `Branch N ->`.

| Fungsi Test | Skenario | Service/Modul |
|---|---|---|
| `TestParseMarkersFromContent_Branches` | 23 | import_soal/parser (rich_content) |
| `TestCreateUjianService_BranchCoverageValidation` | 13 | ujian/ujian_penjadwalan/create |
| `TestUpdatePengumumanService` | 12 | pengumuman/update |
| `TestUpdateProfilSekolahService` | 12 | profil_sekolah/update |
| `TestUpdateMapelService_BranchCoverage` | 10 | mata_pelajaran/update |
| `TestImportVersionService_BranchCoverage` | 6 | import_soal/import_version |
| `TestGetJobService_BranchCoverage` | 6 | import_soal/get_job |
| `TestCreateBankSoalService_BranchCoverage` | 6 | bank_soal/create |
| `TestUpdateKelasService_UpdateNamaKelas` | 6 | kelas/update |
| `TestCreatePengumumanService` | 5 | pengumuman/create |
| `TestDeletePengumumanService` | 5 | pengumuman/delete |
| `TestGetMapelService_BranchCoverage` | 5 | mata_pelajaran/get |
| `TestCreateMapelService_BranchCoverage` | 5 | mata_pelajaran/create |
| `TestWorkerProcessJobs_BranchCoverage` | 5 | import_soal/worker |
| `TestCreateJobService_BranchCoverage` | 5 | import_soal/create_job |
| `TestDeleteMapelService_BranchCoverage` | 4 | mata_pelajaran/delete |
| `TestUpdateRuangUjianService` | 4 | ruang_ujian/update |
| `TestCreateAktivitasUserService` | 4 | aktivitas_user |
| `TestGetMapelByIdService_BranchCoverage` | 4 | mata_pelajaran/get |
| `TestDeleteFileService_BranchCoverage` | 4 | delete_file_system |
| `TestCreateKelasService_CreateTingkatKelas` | 4 | kelas/create |
| `TestCreateKelasService_CreateNamaKelas` | 4 | kelas/create |
| `TestGetBankSoalService_FilterBranchCoverage` | 3 | bank_soal/get |
| `TestCreateRuangUjianService` | 3 | ruang_ujian/create |
| `TestDeleteRuangUjianService` | 3 | ruang_ujian/delete |
| `TestStatistikUjianService_BasisPath` | 3 | ujian/statistik_ujian |
| `TestGradingStatistikUjianService_BasisPath` | 3 | ujian/grading/statistik_ujian |
| `TestGetKelasService_GetFullKelas` | 3 | kelas/get |
| `TestListUjianSelesaiSiswaService_BranchCoverage` | 3 | ujian/siswa_ujian |
| `TestGetBankSoalUploadedService_FilterBranchCoverage` | 3 | bank_soal/get |
| `TestGetAktivitasUserService` | 2 | aktivitas_user |
| `TestGetRuangUjianByKode` | 2 | ruang_ujian/get |
| `TestGetRuangUjianById` | 2 | ruang_ujian/get |
| `TestGetProfilSekolahService` | 2 | profil_sekolah/get |
| `TestDeleteKelasService_DeleteNamaKelas` | 2 | kelas/delete |
| `TestValidateDate` | 2 | pengumuman |
| `TestDashboardServiceGetDashboardStatistik` | 2 | dashboard |
| `TestGetKelasService_GetKelasById` | 2 | kelas/get |
| `TestDeleteBankSoalService_BranchCoverage` | 2 | bank_soal/delete |
| `TestGetPengumumanByIdService` | 2 | pengumuman/get |
| `TestGetPengumumanActiveService` | 1 | pengumuman/get |
| `TestGetPengumumanIncomingService` | 1 | pengumuman/get |
| `TestGetRuangUjian` | 1 | ruang_ujian/get |
| `TestGetPengumumanNonActiveService` | 1 | pengumuman/get |

---

## Keterangan

- **Basis Path Testing**: Skenario dinamai `Path N -> ...` atau `PN -> ...` (contoh: `P1 -> ...`), menguji tiap jalur eksekusi independen dalam flowgraph (setiap path yang mungkin dari start sampai end).
- **Branch Coverage**: Skenario dinamai `Branch N -> ...`, menguji tiap cabang/keputusan (branch) dari struktur kontrol (if/else, switch, error handling).
- File *integration test* (folder `integration_test/`) tidak termasuk dalam perhitungan di atas.
