# 🖥️ UI Workflow Guide - Menu Management

## 📋 Overview

This guide explains the complete workflow for the menu management UI as shown in the interface design.

---

## **🎯 UI Workflow**

### **1. Initial Load: Get Root Menus for Dropdown**

**Endpoint:**

```
GET /api/menus/root
```

**Purpose:** Populate dropdown dengan list root menus

**Response Example:**

```json
{
	"success": true,
	"message": "Root menus retrieved successfully",
	"data": [
		{
			"id": 1,
			"uuid": "56320ee9-6af6-11ed-a7ba-f220afe5e4a9",
			"parent_id": null,
			"name": "system.management",
			"code": "system_management",
			"level": 0
		}
	]
}
```

**UI Implementation:**

```javascript
// Load root menus for dropdown
async function loadRootMenus() {
	const response = await fetch("http://localhost:8080/api/menus/root");
	const data = await response.json();

	// Populate dropdown
	const dropdown = document.getElementById("root-menu-dropdown");
	data.data.forEach((menu) => {
		const option = document.createElement("option");
		option.value = menu.id;
		option.textContent = menu.name;
		dropdown.appendChild(option);
	});
}
```

---

### **2. User Selects Root Menu: Load Hierarchy Tree**

**Endpoint:**

```
GET /api/menus/:id/hierarchy
```

**Purpose:** Load hierarchical tree untuk root menu yang dipilih

**Request Example:**

```bash
GET /api/menus/1/hierarchy
```

**Response Example:**

```json
{
	"success": true,
	"message": "Menu hierarchy retrieved successfully",
	"data": [
		{
			"id": 1,
			"uuid": "...",
			"name": "system.management",
			"level": 0,
			"children": [
				{
					"id": 2,
					"name": "System Management",
					"level": 1,
					"children": [
						{
							"id": 3,
							"name": "Systems",
							"level": 2,
							"children": [
								{
									"id": 4,
									"name": "System Code",
									"level": 3,
									"children": []
								}
							]
						}
					]
				}
			]
		}
	]
}
```

**UI Implementation:**

```javascript
// When dropdown changes
async function onRootMenuChange(rootId) {
	const response = await fetch(
		`http://localhost:8080/api/menus/${rootId}/hierarchy`
	);
	const data = await response.json();

	// Render tree view
	renderTreeView(data.data);
}

function renderTreeView(menus) {
	const treeContainer = document.getElementById("menu-tree");
	treeContainer.innerHTML = "";

	menus.forEach((menu) => {
		renderTreeNode(menu, treeContainer);
	});
}

function renderTreeNode(menu, container) {
	const node = document.createElement("div");
	node.className = "tree-node";
	node.style.paddingLeft = `${menu.level * 20}px`;
	node.innerHTML = `
    <span class="toggle">${menu.children?.length > 0 ? "▼" : ""}</span>
    <span class="menu-name" data-id="${menu.id}">${menu.name}</span>
  `;

	// Click handler untuk show detail
	node.querySelector(".menu-name").addEventListener("click", () => {
		loadMenuDetail(menu.id);
	});

	container.appendChild(node);

	// Recursive untuk children
	if (menu.children && menu.children.length > 0) {
		menu.children.forEach((child) => {
			renderTreeNode(child, container);
		});
	}
}
```

---

### **3. User Clicks Menu Item: Show Detail Panel**

**Endpoint:**

```
GET /api/menus/:id/detail
```

**Purpose:** Get detail menu termasuk parent info dan depth

**Request Example:**

```bash
GET /api/menus/4/detail
```

**Response Example:**

```json
{
	"success": true,
	"message": "Menu detail retrieved successfully",
	"data": {
		"id": 4,
		"uuid": "56320ee9-6af6-11ed-a7ba-f220afe5e4a9",
		"parent_id": 3,
		"name": "System Code",
		"code": "system_code",
		"description": null,
		"route": "/system/code",
		"icon": "code",
		"order_index": 1,
		"level": 3,
		"is_active": true,
		"created_at": "2025-10-20T18:00:00+07:00",
		"updated_at": "2025-10-20T18:00:00+07:00",
		"parent_data": {
			"id": 3,
			"uuid": "...",
			"name": "Systems",
			"code": "systems"
		},
		"depth": 3
	}
}
```

**UI Implementation:**

```javascript
async function loadMenuDetail(menuId) {
	const response = await fetch(
		`http://localhost:8080/api/menus/${menuId}/detail`
	);
	const data = await response.json();
	const menu = data.data;

	// Populate detail panel
	document.getElementById("menu-id").value = menu.uuid;
	document.getElementById("depth").value = menu.depth;
	document.getElementById("parent-data").value =
		menu.parent_data?.name || "Root";
	document.getElementById("menu-name").value = menu.name;
	document.getElementById("menu-code").value = menu.code;
	document.getElementById("menu-route").value = menu.route || "";
	document.getElementById("menu-icon").value = menu.icon || "";
	document.getElementById("order-index").value = menu.order_index;

	// Store current menu ID for update
	document.getElementById("detail-form").dataset.menuId = menu.id;
}
```

---

### **4. User Edits & Saves: Update Menu**

**Endpoint:**

```
PUT /api/menus/:id
```

**Purpose:** Update menu yang sudah di-edit

**Request Example:**

```bash
PUT /api/menus/4
Content-Type: application/json

{
  "parent_id": 3,
  "name": "System Code Updated",
  "code": "system_code",
  "description": "Updated description",
  "route": "/system/code",
  "icon": "code",
  "order_index": 1,
  "is_active": true
}
```

**UI Implementation:**

```javascript
async function saveMenuChanges(event) {
	event.preventDefault();

	const form = document.getElementById("detail-form");
	const menuId = form.dataset.menuId;

	const payload = {
		parent_id: parseInt(document.getElementById("parent-id").value),
		name: document.getElementById("menu-name").value,
		code: document.getElementById("menu-code").value,
		description: document.getElementById("menu-description").value || null,
		route: document.getElementById("menu-route").value || null,
		icon: document.getElementById("menu-icon").value || null,
		order_index: parseInt(document.getElementById("order-index").value),
		is_active: document.getElementById("is-active").checked,
	};

	const response = await fetch(`http://localhost:8080/api/menus/${menuId}`, {
		method: "PUT",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(payload),
	});

	const data = await response.json();

	if (data.success) {
		alert("Menu updated successfully!");
		// Reload tree
		const rootId = document.getElementById("root-menu-dropdown").value;
		onRootMenuChange(rootId);
	} else {
		alert("Failed to update menu: " + data.error);
	}
}
```

---

### **5. (Optional) Load Children on Expand**

**Endpoint:**

```
GET /api/menus/:id/children
```

**Purpose:** Lazy load children hanya ketika node di-expand (untuk performance pada tree besar)

**Request Example:**

```bash
GET /api/menus/3/children
```

**Response Example:**

```json
{
	"success": true,
	"message": "Children retrieved successfully",
	"data": [
		{
			"id": 4,
			"parent_id": 3,
			"name": "System Code",
			"level": 3
		},
		{
			"id": 5,
			"parent_id": 3,
			"name": "Code Registration",
			"level": 3
		}
	]
}
```

**UI Implementation:**

```javascript
async function loadChildren(parentId, containerElement) {
	const response = await fetch(
		`http://localhost:8080/api/menus/${parentId}/children`
	);
	const data = await response.json();

	data.data.forEach((child) => {
		renderTreeNode(child, containerElement);
	});
}
```

---

## **🔄 Complete Workflow Diagram**

```
┌─────────────────────────────────────────────────────────────────┐
│                        PAGE LOAD                                │
│                          ↓                                      │
│              GET /api/menus/root                                │
│              (Load root menus for dropdown)                     │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│                   USER SELECTS ROOT MENU                        │
│                          ↓                                      │
│           GET /api/menus/:id/hierarchy                          │
│           (Load full tree for selected root)                    │
│                          ↓                                      │
│              Render Tree View (Left Panel)                      │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│               USER CLICKS MENU NODE IN TREE                     │
│                          ↓                                      │
│            GET /api/menus/:id/detail                            │
│            (Get menu detail with parent info)                   │
│                          ↓                                      │
│          Show Detail Form (Right Panel)                         │
│          - Menu ID (UUID)                                       │
│          - Depth/Level                                          │
│          - Parent Data                                          │
│          - Name, Code, Route, etc.                              │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│                USER EDITS & CLICKS SAVE                         │
│                          ↓                                      │
│              PUT /api/menus/:id                                 │
│              (Update menu)                                      │
│                          ↓                                      │
│        Refresh Tree (call hierarchy again)                      │
└─────────────────────────────────────────────────────────────────┘
```

---

## **📊 API Endpoints Summary for UI**

| UI Action          | Endpoint                       | Purpose                     |
| ------------------ | ------------------------------ | --------------------------- |
| Load Dropdown      | `GET /api/menus/root`          | Get all root menus          |
| Select Root        | `GET /api/menus/:id/hierarchy` | Load tree for specific root |
| Click Menu         | `GET /api/menus/:id/detail`    | Show detail panel           |
| Save Changes       | `PUT /api/menus/:id`           | Update menu                 |
| Expand Node (lazy) | `GET /api/menus/:id/children`  | Load children on demand     |
| Create New         | `POST /api/menus`              | Add new menu/submenu        |
| Delete             | `DELETE /api/menus/:id`        | Remove menu                 |

---

## **💡 Frontend Implementation Tips**

### **1. State Management**

```javascript
const appState = {
	currentRootId: null,
	selectedMenuId: null,
	treeData: [],
	isDirty: false,
};
```

### **2. Expand/Collapse Tree Nodes**

```javascript
function toggleNode(nodeId) {
	const node = document.getElementById(`node-${nodeId}`);
	const children = node.querySelector(".children");
	const toggle = node.querySelector(".toggle");

	if (children.style.display === "none") {
		children.style.display = "block";
		toggle.textContent = "▼";
	} else {
		children.style.display = "none";
		toggle.textContent = "▶";
	}
}
```

### **3. Highlight Selected Node**

```javascript
function selectNode(menuId) {
	// Remove previous highlight
	document.querySelectorAll(".tree-node").forEach((node) => {
		node.classList.remove("selected");
	});

	// Add highlight to selected
	const selectedNode = document.querySelector(`[data-id="${menuId}"]`);
	selectedNode.closest(".tree-node").classList.add("selected");

	// Load detail
	loadMenuDetail(menuId);
}
```

### **4. Unsaved Changes Warning**

```javascript
let isDirty = false;

document.querySelectorAll("#detail-form input").forEach((input) => {
	input.addEventListener("change", () => {
		isDirty = true;
	});
});

window.addEventListener("beforeunload", (e) => {
	if (isDirty) {
		e.preventDefault();
		e.returnValue = "";
	}
});
```

---

## **🎨 Example Complete HTML Structure**

```html
<!DOCTYPE html>
<html>
	<head>
		<title>Menu Management</title>
		<style>
			.container {
				display: flex;
			}
			.left-panel {
				flex: 1;
				padding: 20px;
				border-right: 1px solid #ccc;
			}
			.right-panel {
				flex: 1;
				padding: 20px;
			}
			.tree-node {
				padding: 5px;
				cursor: pointer;
			}
			.tree-node.selected {
				background: #e3f2fd;
			}
			.children {
				margin-left: 20px;
			}
		</style>
	</head>
	<body>
		<div class="header">
			<h1>Menu Management</h1>
			<select id="root-menu-dropdown" onchange="onRootMenuChange(this.value)">
				<option value="">Select Root Menu...</option>
			</select>
		</div>

		<div class="container">
			<div class="left-panel">
				<button onclick="expandAll()">Expand All</button>
				<button onclick="collapseAll()">Collapse All</button>
				<div id="menu-tree"></div>
			</div>

			<div class="right-panel">
				<form id="detail-form" onsubmit="saveMenuChanges(event)">
					<h3>Menu Details</h3>

					<label>Menu ID</label>
					<input type="text" id="menu-id" readonly />

					<label>Depth</label>
					<input type="number" id="depth" readonly />

					<label>Parent Data</label>
					<input type="text" id="parent-data" readonly />

					<label>Name</label>
					<input type="text" id="menu-name" required />

					<label>Code</label>
					<input type="text" id="menu-code" required />

					<label>Route</label>
					<input type="text" id="menu-route" />

					<label>Icon</label>
					<input type="text" id="menu-icon" />

					<label>Order</label>
					<input type="number" id="order-index" />

					<label>
						<input type="checkbox" id="is-active" checked />
						Active
					</label>

					<button type="submit">Save</button>
				</form>
			</div>
		</div>

		<script src="app.js"></script>
	</body>
</html>
```

---

## **✅ Checklist Implementation**

- [ ] Load root menus on page load
- [ ] Render dropdown dengan root menus
- [ ] Load hierarchy when root selected
- [ ] Render tree view dengan expand/collapse
- [ ] Show detail panel when node clicked
- [ ] Populate form dengan menu data
- [ ] Handle form submission (update)
- [ ] Refresh tree after update
- [ ] Add unsaved changes warning
- [ ] Add loading indicators
- [ ] Add error handling
- [ ] Add success notifications

---

Your API is now fully ready for this UI workflow! 🎉
