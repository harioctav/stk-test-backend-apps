package domain

import "time"

// Menu represents the menu entity
type Menu struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID        string     `json:"uuid" gorm:"size:36;uniqueIndex;not null"`
	ParentID    *int64     `json:"parent_id" gorm:"index"`
	Name        string     `json:"name" gorm:"size:255;not null"`
	Code        string     `json:"code" gorm:"size:100;uniqueIndex"`
	Description *string    `json:"description" gorm:"type:text"`
	Route       *string    `json:"route" gorm:"size:255"`
	Icon        *string    `json:"icon" gorm:"size:100"`
	OrderIndex  int        `json:"order_index" gorm:"default:0;index"`
	Level       int        `json:"level" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedBy   *int64     `json:"created_by"`
	UpdatedBy   *int64     `json:"updated_by"`
	Children    []Menu     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// MenuDetail represents menu with parent information
type MenuDetail struct {
	Menu
	ParentData *MenuParentInfo `json:"parent_data,omitempty"`
	Depth      int             `json:"depth"`
}

// MenuParentInfo represents parent menu basic info
type MenuParentInfo struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// TableName specifies the table name for Menu
func (Menu) TableName() string {
	return "menus"
}

// CreateMenuRequest represents the request payload for creating a menu
type CreateMenuRequest struct {
	ParentID    *int64  `json:"parent_id"`
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Description *string `json:"description"`
	Route       *string `json:"route"`
	Icon        *string `json:"icon"`
	OrderIndex  int     `json:"order_index"`
	IsActive    bool    `json:"is_active"`
}

// UpdateMenuRequest represents the request payload for updating a menu
type UpdateMenuRequest struct {
	ParentID    *int64  `json:"parent_id"`
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Description *string `json:"description"`
	Route       *string `json:"route"`
	Icon        *string `json:"icon"`
	OrderIndex  int     `json:"order_index"`
	IsActive    bool    `json:"is_active"`
}

// MenuRepository defines the interface for menu data operations
type MenuRepository interface {
	Create(menu *Menu) error
	Update(menu *Menu) error
	Delete(id int64) error
	FindByID(id int64) (*Menu, error)
	FindByUUID(uuid string) (*Menu, error)
	FindAll() ([]Menu, error)
	FindByParentID(parentID *int64) ([]Menu, error)
	FindRootMenus() ([]Menu, error)
	FindHierarchical() ([]Menu, error)
	FindHierarchicalByRootID(rootID int64) ([]Menu, error)
	FindDetailByID(id int64) (*MenuDetail, error)
	FindChildrenByParentID(parentID int64) ([]Menu, error)
}

// MenuService defines the interface for menu business logic
type MenuService interface {
	CreateMenu(req *CreateMenuRequest) (*Menu, error)
	UpdateMenu(id int64, req *UpdateMenuRequest) (*Menu, error)
	DeleteMenu(id int64) error
	GetMenuByID(id int64) (*Menu, error)
	GetMenuByUUID(uuid string) (*Menu, error)
	GetAllMenus() ([]Menu, error)
	GetRootMenus() ([]Menu, error)
	GetMenuHierarchy() ([]Menu, error)
	GetHierarchyByRootID(rootID int64) ([]Menu, error)
	GetMenuDetail(id int64) (*MenuDetail, error)
	GetChildrenByParentID(parentID int64) ([]Menu, error)
}

