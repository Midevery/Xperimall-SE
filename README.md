# 🏬 Xperimall

**Xperimall** adalah aplikasi direktori mall digital yang dikembangkan secara internal.  
Aplikasi ini menyediakan fitur seperti daftar tenant, promo, perencana aktivitas, dan lainnya.

Proyek ini terdiri dari dua bagian utama:

- **Frontend**: React Native menggunakan Expo (SDK 52)
- **Backend**: Golang + MySQL

---

## 🔗 Repository

GitHub: [https://github.com/Midevery/Xperimall-SE.git](https://github.com/Midevery/Xperimall-SE.git)

---

## 📦 Langkah Awal

```bash
git clone https://github.com/Midevery/Xperimall-SE.git
cd xperimall-se
```

---

## 🖥️ Backend (Golang)

### 📁 Lokasi
`/XperimallBackend`

### ▶️ Cara Menjalankan

Pastikan sudah menginstall **Go** dan **MySQL**.

```bash
cd XperimallBackend
go run main.go
```

> 💡 **Catatan**: Pastikan konfigurasi koneksi database sudah sesuai di dalam kode (`main.go` atau file konfigurasi yang digunakan).

---

## 📱 Frontend (React Native + Expo)

### 📁 Lokasi
Root directory (`/`)

### 🧩 Install Dependencies

```bash
npm install
```

### ▶️ Jalankan App

> ❗ **Catatan penting**:  
> Aplikasi ini menggunakan **Expo SDK 52**, yang **tidak kompatibel dengan Expo Go**.  
> **Wajib menggunakan Android Studio Emulator atau iOS Simulator.**

Pastikan emulator aktif, lalu jalankan:

```bash
npx expo start
```

Kemudian pilih opsi:

- ✅ **Run on Android emulator**
- ✅ **Run on iOS simulator** (hanya untuk Mac)

---

