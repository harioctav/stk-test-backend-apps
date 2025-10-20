# 📮 Panduan Lengkap Postman - STK Menu API

## 🚀 Quick Start

### 1. Import Collection ke Postman

1. **Download file:** `STK_Menu_API.postman_collection.json`
2. **Buka Postman**
3. **Klik tombol "Import"** (kiri atas)
4. **Drag & drop** file JSON atau klik "Upload Files"
5. **Klik "Import"**

✅ Collection akan muncul dengan nama: **"STK Menu Management API"**

---

### 2. Setup Environment (Optional tapi Recommended)

1. **Create Environment:**

   - Klik icon ⚙️ (Settings) → **Environments**
   - Klik **"Create Environment"** atau **+**
   - Nama: `Local Development`

2. **Add Variables:**

   ```
   Variable         | Initial Value              | Current Value
   ---------------- | -------------------------- | --------------------------
   base_url         | http://localhost:8080      | http://localhost:8080
   menu_uuid        | (leave empty)              | (will be filled later)
   ```

3. **Save** dan **Select environment** dari dropdown (kanan atas)

---

## 📂 Struktur Collection

```
STK Menu Management API/
├── Health Check
│   └── Health Check
├── Menus - Create Data (19 requests)
│   ├── 01 - Create Root: system.management
│   ├── 02 - Create L1: System Management
│   ├── 03 - Create L2: Systems
│   ├── ... (hingga request 19)
├── Menus - Read (4 requests)
│   ├── Get All Menus (Flat)
│   ├── Get Menu Hierarchy
│   ├── Get Menu by ID
│   └── Get Menu by UUID
├── Menus - Update
│   └── Update Menu
└── Menus - Delete
    └── Delete Menu
```

---

## 🎯 Cara Menggunakan

### **SCENARIO 1: Buat Seluruh Struktur Menu**

Jalankan **request 01-19** secara berurutan di folder **"Menus - Create Data"**

#### **Method 1: Manual (Satu per Satu)**

1. Klik request **"01 - Create Root: system.management"**
2. **Klik "Send"**
3. Lihat response → catat ID yang dibuat
4. Lanjut ke request **"02 - Create L1: System Management"**
5. **Klik "Send"**
6. Ulangi hingga request **"19"**

**Tips:** Pastikan parent_id di request selanjutnya sesuai dengan ID yang sudah dibuat!

---

#### **Method 2: Using Collection Runner (Otomatis)** ⭐ Recommended!

1. **Klik kanan** pada folder **"Menus - Create Data"**
2. Pilih **"Run folder"**
3. **Collection Runner** akan terbuka
4. **Atur settings:**
   - ✅ Iterations: `1`
   - ✅ Delay: `100ms` (jeda antar request)
5. **Klik "Run Menus - Create Data"**
6. Tunggu hingga semua request selesai (19/19)

✅ **Hasilnya:** Seluruh struktur menu akan dibuat otomatis!

---

### **SCENARIO 2: Lihat Hasil (View Data)**

#### **A. View Hierarchy (Tree Structure)**

1. Buka request: **"Get Menu Hierarchy"**
2. **Klik "Send"**
3. Response akan menampilkan struktur tree:

```json
{
  "success": true,
  "message": "Menu hierarchy retrieved successfully",
  "data": [
    {
      "id": 1,
      "uuid": "abc-123...",
      "name": "system.management",
      "level": 0,
      "children": [
        {
          "id": 2,
          "name": "System Management",
          "level": 1,
          "children": [...]
        }
      ]
    }
  ]
}
```

---

#### **B. View All Menus (Flat List)**

1. Buka request: **"Get All Menus (Flat)"**
2. **Klik "Send"**
3. Response: Array semua menu tanpa nested structure

---

#### **C. View Single Menu by ID**

1. Buka request: **"Get Menu by ID"**
2. Ubah ID di URL: `/api/menus/1` (ganti 1 dengan ID yang diinginkan)
3. **Klik "Send"**

---

#### **D. View Single Menu by UUID** 🆕

1. Buka request: **"Get Menu by UUID"**
2. **Copy UUID** dari response sebelumnya
3. Paste UUID di URL: `/api/menus/uuid/550e8400-e29b-...`
4. **Klik "Send"**

**Alternative (using variable):**

- Set environment variable `menu_uuid` dengan UUID yang valid
- URL akan otomatis: `/api/menus/uuid/{{menu_uuid}}`

---

### **SCENARIO 3: Update Menu**

1. Buka request: **"Update Menu"**
2. Ubah ID di URL: `/api/menus/2` (menu yang akan diupdate)
3. Edit request body sesuai kebutuhan:

```json
{
	"parent_id": 1,
	"name": "System Management Updated",
	"code": "system_mgmt",
	"description": "New description",
	"route": "/system",
	"icon": "settings",
	"order_index": 1,
	"is_active": true
}
```

4. **Klik "Send"**

---

### **SCENARIO 4: Delete Menu**

1. Buka request: **"Delete Menu"**
2. Ubah ID di URL: `/api/menus/5` (menu yang akan dihapus)
3. **Klik "Send"**

⚠️ **Warning:**

- Tidak bisa delete menu yang punya children!
- Harus delete children dulu, baru parent

---

## 🎨 Tips & Tricks

### **1. Save Response as Example**

Setelah dapat response sukses:

1. **Klik "Save Response"** (di panel response)
2. **Pilih "Save as example"**
3. Berguna untuk dokumentasi dan referensi

---

### **2. Use Tests for Auto-Validation**

Tambahkan di tab **"Tests"** untuk auto-check:

```javascript
// Check status code
pm.test("Status code is 200 or 201", function () {
	pm.expect(pm.response.code).to.be.oneOf([200, 201]);
});

// Check success field
pm.test("Success is true", function () {
	var jsonData = pm.response.json();
	pm.expect(jsonData.success).to.eql(true);
});

// Check data exists
pm.test("Data exists", function () {
	var jsonData = pm.response.json();
	pm.expect(jsonData.data).to.exist;
});

// Auto-save UUID to environment
if (pm.response.code === 201) {
	var jsonData = pm.response.json();
	pm.environment.set("menu_uuid", jsonData.data.uuid);
	console.log("UUID saved: " + jsonData.data.uuid);
}
```

---

### **3. Pre-request Script untuk Dynamic Data**

Tambahkan di tab **"Pre-request Script"**:

```javascript
// Generate random menu name
pm.variables.set("random_name", "Menu_" + Math.floor(Math.random() * 1000));

// Get current timestamp
pm.variables.set("timestamp", new Date().toISOString());
```

Lalu gunakan di body:

```json
{
	"name": "{{random_name}}",
	"code": "menu_{{$timestamp}}"
}
```

---

### **4. Organize dengan Folders**

Sudah diorganize:

- ✅ **Health Check** → Testing koneksi
- ✅ **Create Data** → Seeding data
- ✅ **Read** → GET endpoints
- ✅ **Update** → PUT endpoints
- ✅ **Delete** → DELETE endpoints

---

### **5. Export Collection untuk Sharing**

1. Klik **"..."** di collection
2. Pilih **"Export"**
3. Format: **Collection v2.1** (recommended)
4. Share file JSON ke team

---

## 📊 Expected Results

Setelah menjalankan semua Create requests (01-19), Anda akan punya:

```
system.management (ID: 1) ← Root
└── System Management (ID: 2)
    ├── Systems (ID: 3)
    │   ├── System Code (ID: 4)
    │   │   ├── Code Registration (ID: 5)
    │   │   ├── Code Registration -2 (ID: 6)
    │   │   └── Properties (ID: 7)
    │   ├── Menus (ID: 8)
    │   │   └── Menu Registration (ID: 9)
    │   └── API List (ID: 10)
    │       ├── API Registration (ID: 11)
    │       └── API Edit (ID: 12)
    ├── Users & Groups (ID: 13)
    │   ├── Users (ID: 14)
    │   │   └── User Account Registration (ID: 15)
    │   └── Groups (ID: 16)
    │       └── User Group Registration (ID: 17)
    └── 사용자 승인 (ID: 18)
        └── 사용자 승인 상세 (ID: 19)
```

**Total:** 19 menus dengan 5 levels hierarki

---

## ❌ Troubleshooting

### **Error: "Connection refused"**

✅ **Solution:** Pastikan server sudah running (`go run cmd/api/main.go`)

### **Error: "parent menu not found"**

✅ **Solution:** Parent ID belum dibuat. Buat parent dulu!

### **Error: "duplicate entry for key code"**

✅ **Solution:** Field `code` harus unique. Ubah kodenya.

### **Error: "cannot delete menu with children"**

✅ **Solution:** Delete children dulu, baru parent.

---

## 🎯 Next Steps

1. ✅ Test semua endpoint
2. ✅ Experiment dengan update & delete
3. ✅ Buat menu custom sendiri
4. ✅ Test error scenarios
5. ✅ Share collection dengan team

---

## 📚 Additional Resources

- **Main Documentation:** `README.md`
- **API Source Code:** `cmd/api/main.go`
- **Handler Code:** `internal/handler/menu_handler.go`

---

Happy Testing! 🚀
