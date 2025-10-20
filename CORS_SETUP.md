# 🌐 CORS Configuration Guide

## ❓ Apa itu CORS?

**CORS (Cross-Origin Resource Sharing)** adalah security feature di browser yang membatasi request dari origin yang berbeda.

### Contoh CORS Block:

```
Frontend: http://localhost:3000
Backend:  http://localhost:8080
```

Browser akan **block** karena beda port (3000 vs 8080) = different origins!

**Error yang muncul:**

```
Access to XMLHttpRequest at 'http://localhost:8080/api/menus/hierarchy'
from origin 'http://localhost:3000' has been blocked by CORS policy:
No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

---

## ✅ Solusi yang Sudah Diimplementasi

### 1. Install CORS Middleware ✅

```bash
go get github.com/gin-contrib/cors
```

### 2. Configure CORS ✅

File: `cmd/api/main.go`

```go
import (
    "github.com/gin-contrib/cors"
    "time"
)

// In setupRouter function
router.Use(cors.New(cors.Config{
    AllowOrigins:     cfg.CORS.AllowedOrigins,
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

### 3. Environment Configuration ✅

File: `.env`

```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001,http://localhost:5173
```

---

## 🎯 Configuration Explained

### AllowOrigins

```go
cfg.CORS.AllowedOrigins
```

**Sources from `.env`:**

- `http://localhost:3000` → React default port
- `http://localhost:3001` → Alternative React port
- `http://localhost:5173` → Vite default port

**Production example:**

```env
ALLOWED_ORIGINS=https://yourapp.com,https://www.yourapp.com,https://admin.yourapp.com
```

---

### AllowMethods

```go
[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
```

- All CRUD operations
- `OPTIONS` for preflight requests (browser sends this first)

---

### AllowHeaders

```go
[]string{"Origin", "Content-Type", "Accept", "Authorization"}
```

Headers yang frontend bisa kirim:

- `Origin` → Browser auto-send
- `Content-Type` → For JSON requests
- `Accept` → For response format
- `Authorization` → For JWT/Bearer tokens

---

### ExposeHeaders

```go
[]string{"Content-Length"}
```

Headers yang frontend bisa baca dari response

---

### AllowCredentials

```go
true
```

Allow frontend kirim:

- Cookies
- Authorization headers
- TLS client certificates

**Important:** Jika `true`, `AllowOrigins` tidak boleh `*` (wildcard)

---

### MaxAge

```go
12 * time.Hour
```

Browser cache preflight request selama 12 jam

- Reduce OPTIONS requests
- Better performance

---

## 🚀 How to Use

### Step 1: Update `.env` (if needed)

Add your frontend URL:

```env
ALLOWED_ORIGINS=http://localhost:3000,https://myapp.com
```

**Multiple origins:** Separate with comma (NO SPACES!)

---

### Step 2: Restart Server

```bash
go run cmd/api/main.go
```

---

### Step 3: Test from Frontend

#### **React Example:**

```javascript
// src/api/menuService.js
import axios from "axios";

const API_BASE_URL = "http://localhost:8080/api";

export const getMenuHierarchy = async () => {
	try {
		const response = await axios.get(`${API_BASE_URL}/menus/hierarchy`);
		return response.data;
	} catch (error) {
		console.error("Error fetching menu hierarchy:", error);
		throw error;
	}
};

export const createMenu = async (menuData) => {
	try {
		const response = await axios.post(`${API_BASE_URL}/menus`, menuData);
		return response.data;
	} catch (error) {
		console.error("Error creating menu:", error);
		throw error;
	}
};
```

#### **Vue Example:**

```javascript
// src/services/menuService.js
import axios from "axios";

const apiClient = axios.create({
	baseURL: "http://localhost:8080/api",
	withCredentials: true, // if using cookies
	headers: {
		"Content-Type": "application/json",
	},
});

export default {
	getMenuHierarchy() {
		return apiClient.get("/menus/hierarchy");
	},
	createMenu(menu) {
		return apiClient.post("/menus", menu);
	},
	updateMenu(id, menu) {
		return apiClient.put(`/menus/${id}`, menu);
	},
	deleteMenu(id) {
		return apiClient.delete(`/menus/${id}`);
	},
};
```

#### **Vanilla JavaScript (Fetch API):**

```javascript
// Get menu hierarchy
fetch("http://localhost:8080/api/menus/hierarchy")
	.then((response) => response.json())
	.then((data) => {
		console.log("Menu hierarchy:", data);
		// Render your tree view here
	})
	.catch((error) => console.error("Error:", error));

// Create new menu
fetch("http://localhost:8080/api/menus", {
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify({
		name: "New Menu",
		code: "new_menu",
		order_index: 1,
		is_active: true,
	}),
})
	.then((response) => response.json())
	.then((data) => console.log("Created:", data))
	.catch((error) => console.error("Error:", error));
```

---

## 🔒 Security Best Practices

### Development

```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

✅ Specific origins only

### Production

```env
ALLOWED_ORIGINS=https://yourapp.com,https://www.yourapp.com
```

✅ HTTPS only
✅ No wildcards (\*)
✅ Specific domains

### ❌ NEVER DO THIS in Production:

```go
AllowOrigins: []string{"*"}  // DON'T!
```

This allows **ANY** website to access your API!

---

## 🐛 Troubleshooting

### Error: "CORS policy: No 'Access-Control-Allow-Origin'"

**Solution:**

1. Check `.env` has correct `ALLOWED_ORIGINS`
2. Restart server after changing `.env`
3. Clear browser cache
4. Check frontend URL matches exactly (including protocol & port)

---

### Error: "CORS policy: Request header not allowed"

**Solution:**
Add header to `AllowHeaders` in `main.go`:

```go
AllowHeaders: []string{
    "Origin",
    "Content-Type",
    "Accept",
    "Authorization",
    "X-Custom-Header",  // Add your header
}
```

---

### Error: "CORS policy: Method not allowed"

**Solution:**
Add method to `AllowMethods`:

```go
AllowMethods: []string{
    "GET",
    "POST",
    "PUT",
    "DELETE",
    "PATCH",    // Add if needed
    "OPTIONS"
}
```

---

### Preflight Request Failed

**What is Preflight?**
Browser sends OPTIONS request first to check if actual request is allowed.

**Solution:**

- Make sure `OPTIONS` is in `AllowMethods` ✅ (already configured)
- Check `MaxAge` is set ✅ (already configured)

---

## 📊 How CORS Works

### Simple Request (GET, POST with simple headers)

```
1. Frontend → Backend: GET /api/menus
2. Backend → Frontend: Response + CORS headers
3. Browser: ✅ Allow or ❌ Block
```

### Preflight Request (PUT, DELETE, custom headers)

```
1. Frontend → Backend: OPTIONS /api/menus (preflight)
2. Backend → Frontend: CORS headers
3. Browser checks: ✅ OK → Continue
4. Frontend → Backend: Actual request (PUT/DELETE)
5. Backend → Frontend: Response
```

---

## 🎯 Advanced Configuration

### Allow All Origins (Development Only!)

File: `internal/config/config.go`

```go
func getCORSOrigins() []string {
    env := getEnv("APP_ENV", "development")
    if env == "development" {
        return []string{"*"}  // Allow all in dev
    }
    origins := getEnv("ALLOWED_ORIGINS", "")
    return strings.Split(origins, ",")
}
```

⚠️ **Warning:** Only use `*` for local development!

---

### Dynamic Origins (Database-driven)

```go
func (c *Config) GetAllowedOrigins() []string {
    // Fetch from database or config service
    // This allows runtime configuration
    return fetchOriginsFromDB()
}
```

---

### Environment-Specific Origins

**.env.development:**

```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

**.env.production:**

```env
ALLOWED_ORIGINS=https://app.example.com,https://www.example.com
```

---

## ✅ Verification Checklist

After setup, verify:

- [ ] `.env` has `ALLOWED_ORIGINS` configured
- [ ] Server restarted after config changes
- [ ] Frontend can fetch data without CORS error
- [ ] POST/PUT/DELETE requests work
- [ ] Network tab shows CORS headers in response
- [ ] Preflight (OPTIONS) requests succeed

**Check CORS headers in browser DevTools:**

```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization
Access-Control-Allow-Credentials: true
```

---

## 📚 Additional Resources

- [MDN CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [gin-contrib/cors GitHub](https://github.com/gin-contrib/cors)
- [Understanding CORS](https://www.codecademy.com/article/what-is-cors)

---

## 🎉 Summary

✅ CORS middleware installed
✅ Configuration via `.env` file
✅ Support multiple origins
✅ All HTTP methods allowed
✅ Credentials enabled
✅ Production-ready setup

Your API is now ready for frontend integration! 🚀
