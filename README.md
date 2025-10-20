# STK Technical Test API - Menu Management System

A RESTful API for hierarchical menu management built with Go, Gin, and GORM.

## 🚀 Features

- Create, Read, Update, Delete (CRUD) operations for menus
- Hierarchical menu structure (parent-child relationships)
- Get menus in flat or hierarchical format
- **UUID support** for secure and distributed identification
- **CORS enabled** for frontend integration
- Clean Architecture implementation
- RESTful API design
- Auto-generated UUID for each menu
- Multi-level hierarchy support (unlimited depth)

## 📋 Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- golang-migrate (for database migrations)

## 🛠️ Installation

1. **Clone the repository**

   ```bash
   cd stk-technical-test-api
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Setup environment variables**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` file with your database credentials:

   ```
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=stk_menu_system
   SERVER_PORT=8080
   APP_ENV=development
   ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001,http://localhost:5173
   ```

4. **Create database**

   ```bash
   mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS stk_menu_system;"
   ```

5. **Run migrations**
   ```bash
   migrate -path database/migrations -database "mysql://root:1111@tcp(localhost:3306)/stk_menu_system" up
   ```

## 🏃 Running the Application

```bash
go run cmd/api/main.go
```

The server will start at `http://localhost:8080`

## 📚 API Endpoints

### Health Check

- `GET /health` - Check if server is running

### Menu Management

| Method | Endpoint                | Description                             |
| ------ | ----------------------- | --------------------------------------- |
| GET    | `/api/menus/hierarchy`  | Get all menus in hierarchical structure |
| GET    | `/api/menus`            | Get all menus (flat list)               |
| GET    | `/api/menus/:id`        | Get a specific menu by ID               |
| GET    | `/api/menus/uuid/:uuid` | **NEW!** Get a specific menu by UUID    |
| POST   | `/api/menus`            | Create a new menu (UUID auto-generated) |
| PUT    | `/api/menus/:id`        | Update an existing menu                 |
| DELETE | `/api/menus/:id`        | Delete a menu                           |

## 📝 Request/Response Examples

### Create Menu (POST /api/menus)

**Request Body:**

```json
{
	"name": "System Management",
	"code": "system.management",
	"description": "System management module",
	"route": "/system",
	"icon": "settings",
	"order_index": 1,
	"is_active": true
}
```

**Response:**

```json
{
	"success": true,
	"message": "Menu created successfully",
	"data": {
		"id": 1,
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"parent_id": null,
		"name": "System Management",
		"code": "system.management",
		"description": "System management module",
		"route": "/system",
		"icon": "settings",
		"order_index": 1,
		"level": 0,
		"is_active": true,
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}
}
```

**Note:** UUID is automatically generated for each menu.

### Create Submenu (POST /api/menus)

**Request Body:**

```json
{
	"parent_id": 1,
	"name": "Users",
	"code": "users",
	"description": "User management",
	"route": "/system/users",
	"icon": "user",
	"order_index": 1,
	"is_active": true
}
```

### Get Menu Hierarchy (GET /api/menus/hierarchy)

**Response:**

```json
{
	"success": true,
	"message": "Menu hierarchy retrieved successfully",
	"data": [
		{
			"id": 1,
			"parent_id": null,
			"name": "System Management",
			"code": "system.management",
			"level": 0,
			"children": [
				{
					"id": 2,
					"parent_id": 1,
					"name": "Users",
					"code": "users",
					"level": 1,
					"children": []
				}
			]
		}
	]
}
```

## 🏗️ Project Structure

```
stk-technical-test-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection
│   ├── domain/
│   │   └── menu.go              # Domain models & interfaces
│   ├── repository/
│   │   └── menu_repository.go   # Data access layer
│   ├── service/
│   │   └── menu_service.go      # Business logic layer
│   └── handler/
│       └── menu_handler.go      # HTTP request handlers
├── pkg/
│   └── response/
│       └── response.go          # Standard API response
├── database/
│   └── migrations/              # Database migration files
├── .env                         # Environment variables
├── .env.example                 # Environment variables example
├── go.mod
└── README.md
```

## 🧪 Testing with cURL

### Create root menu

```bash
curl -X POST http://localhost:8080/api/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "System Management",
    "code": "system.management",
    "order_index": 1,
    "is_active": true
  }'
```

### Create submenu

```bash
curl -X POST http://localhost:8080/api/menus \
  -H "Content-Type: application/json" \
  -d '{
    "parent_id": 1,
    "name": "Users",
    "code": "users",
    "order_index": 1,
    "is_active": true
  }'
```

### Get hierarchy

```bash
curl http://localhost:8080/api/menus/hierarchy
```

### Get menu by UUID

```bash
curl http://localhost:8080/api/menus/uuid/550e8400-e29b-41d4-a716-446655440000
```

## 🔒 Database Schema

```sql
CREATE TABLE menus (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    parent_id BIGINT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE,
    description TEXT,
    route VARCHAR(255),
    icon VARCHAR(100),
    order_index INT DEFAULT 0,
    level INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,

    FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE,
    INDEX idx_uuid (uuid)
);
```

## 📦 Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [google/uuid](https://github.com/google/uuid) - UUID generation

## 📥 Import Postman Collection

1. Download `STK_Menu_API.postman_collection.json` from the project root
2. Open Postman
3. Click **Import** → Select the JSON file
4. Collection will be imported with all endpoints ready to use
5. The collection includes:
   - ✅ Complete CRUD operations
   - ✅ 19 pre-configured requests for creating full menu hierarchy
   - ✅ Health check endpoint
   - ✅ Environment variables setup

## 🎯 Quick Start with Postman

After importing the collection:

1. **Set Environment Variables** (optional):

   - `base_url`: `http://localhost:8080` (default)
   - `menu_uuid`: Copy UUID from any created menu

2. **Create Complete Menu Structure**:

   - Navigate to folder: **Menus - Create Data**
   - Run requests **01 to 19** in order
   - This will create the entire hierarchical menu structure

3. **View Results**:
   - Use **Get Menu Hierarchy** to see the tree structure
   - Use **Get All Menus** to see flat list with all fields including UUID

## 🔑 UUID Features

### Why UUID?

- **Security**: IDs are not sequential and hard to predict
- **Distributed Systems**: No collision between different servers
- **External APIs**: Safe for public-facing endpoints

### UUID Usage Examples

**Get by ID (Internal):**

```bash
GET /api/menus/1
```

**Get by UUID (External/Public):**

```bash
GET /api/menus/uuid/550e8400-e29b-41d4-a716-446655440000
```

Both endpoints return the same menu, but UUID is safer for external use.

## 🌐 CORS Configuration

This API is configured with CORS (Cross-Origin Resource Sharing) to allow frontend applications to access it.

### Allowed Origins

By default, the following origins are allowed:

- `http://localhost:3000` (React default)
- `http://localhost:3001` (Alternative React port)
- `http://localhost:5173` (Vite default)

### Configure CORS

You can configure allowed origins in the `.env` file:

```env
ALLOWED_ORIGINS=http://localhost:3000,https://yourapp.com,https://www.yourapp.com
```

**Multiple origins:** Separate with commas (no spaces)

### CORS Settings

The API allows:

- **Methods:** GET, POST, PUT, DELETE, OPTIONS
- **Headers:** Origin, Content-Type, Accept, Authorization
- **Credentials:** Enabled (for cookies/auth)
- **Max Age:** 12 hours (preflight cache)

### Frontend Integration Example

**Fetch API:**

```javascript
fetch("http://localhost:8080/api/menus/hierarchy")
	.then((response) => response.json())
	.then((data) => console.log(data))
	.catch((error) => console.error("Error:", error));
```

**Axios:**

```javascript
axios
	.get("http://localhost:8080/api/menus/hierarchy")
	.then((response) => console.log(response.data))
	.catch((error) => console.error("Error:", error));
```

**With Credentials:**

```javascript
fetch("http://localhost:8080/api/menus/hierarchy", {
	credentials: "include", // Send cookies
	headers: {
		Authorization: "Bearer your-token-here",
	},
});
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.
