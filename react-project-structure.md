# React Project Structure

Berikut dokumentasi struktur folder React berdasarkan susunan folder pada `src/`:
[text](https://dev.to/ziontutorial/best-project-structure-reactjs-project-22ef)

```
src/
├── assets/
│   ├── images/
│   └── styles/
│       └── global.css
├── components/
│   ├── Common/
│   │   ├── Button.js
│   │   └── Input.js
│   ├── Feature1/
│   │   ├── Feature1Component1.js
│   │   └── Feature1Component2.js
│   └── Feature2/
│       ├── Feature2Component1.js
│       └── Feature2Component2.js
├── containers/
│   ├── Feature1Container.js
│   └── Feature2Container.js
├── contexts/
│   └── AuthContext.js
├── hooks/
│   ├── useFetch.js
│   └── useAuth.js
├── services/
│   ├── ApiService.js
│   └── AuthService.js
├── redux/
│   ├── actions/
│   ├── reducers/
│   └── store/
├── routes/
│   └── AppRouter.js
├── utils/
│   └── helpers.js
├── App.js
└── index.js
```

## Penjelasan Per Folder

### ✅ 1. `assets/`

Tempat file statis seperti gambar, ikon, dan file style global.

### ✅ 2. `components/`

Komponen UI yang reusable dan presentasional (tanpa logic kompleks). Bisa dikelompokkan per fitur atau berdasarkan kegunaan umum.

### ✅ 3. `containers/`

Komponen yang menghubungkan UI dengan data/state, misalnya container yang memanggil API atau Redux.

### ✅ 4. `contexts/`

Untuk React Context API (global state yang lebih ringan dibanding Redux).

### ✅ 5. `hooks/`

Custom React hooks yang bisa dipakai ulang.

### ✅ 6. `services/`

Logika bisnis eksternal seperti API calls atau integrasi layanan luar.

### ✅ 7. `redux/`

Semua file Redux, seperti:

- `actions`
- `reducers`
- konfigurasi `store`

Ini cocok kalau kamu pakai Redux.

### ✅ 8. `routes/`

Tempat definisi route (biasanya dengan React Router).

### ✅ 9. `utils/`

Fungsi utilitas yang dipakai di berbagai bagian app (helper functions).

### 📌 10. `App.js` dan `index.js`

- `App.js` → struktur utama aplikasi
- `index.js` → titik entry, mount React ke DOM

## 🔄 Opsional: Folder `layouts/`

Kalau kamu punya layout yang dipakai banyak halaman (misalnya header, footer, sidebar), artikel itu juga menyarankan memasukkan folder `layouts/`:

```
layouts/
├── MainLayout/
│   ├── Header.js
│   ├── Footer.js
│   └── Navigation.js
```

Ini membantu memisahkan struktur halaman dari konten komponen.
