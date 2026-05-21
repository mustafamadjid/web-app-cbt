# Black Box Selenium SMAFI CBT

Suite ini mengikuti `planning.md` dan ditempatkan mandiri di root repo pada folder `blackbox_selenium`.

## Prasyarat

- Backend, database, dan frontend sudah berjalan.
- Frontend default: `http://localhost:5173`.
- Backend default: `http://localhost:8080`.
- Chrome tersedia di mesin lokal.

## Instalasi

```bash
cd blackbox_selenium
npm install
copy .env.e2e.example .env.e2e
```

Sesuaikan kredensial di `.env.e2e` dengan akun database E2E, bukan data produksi.

## Menjalankan Test

```bash
npm run test:e2e
npm run test:e2e:auth
npm run test:e2e:master
npm run test:e2e:exam
```

Spec destructive hanya berjalan jika `E2E_ALLOW_DESTRUCTIVE=true`:

```bash
npm run test:e2e:destructive
```

Screenshot gagal tersimpan di `artifacts/screenshots`, sedangkan download di `artifacts/downloads`.

## Status Implementasi

- Semua 144 ID skenario dari rencana dimasukkan ke spec.
- Skenario autentikasi utama sudah berisi langkah Selenium nyata.
- Skenario lain dibuat sebagai kerangka `it.skip` dengan alasan eksplisit sampai data seed E2E dan selector final tersedia.
- Helper, page object, dan fixture disiapkan agar skenario bisa dinaikkan bertahap menjadi executable tanpa mengubah struktur suite.
