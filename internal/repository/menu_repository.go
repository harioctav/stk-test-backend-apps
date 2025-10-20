package repository

import (
	"fmt"

	"stk-technical-test-api/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new menu repository instance
func NewMenuRepository(db *gorm.DB) domain.MenuRepository {
	return &menuRepository{
		db: db,
	}
}

func (r *menuRepository) Create(menu *domain.Menu) error {
	// Generate UUID
	menu.UUID = uuid.New().String()

	// Calculate level based on parent
	if menu.ParentID != nil {
		parent, err := r.FindByID(*menu.ParentID)
		if err != nil {
			return fmt.Errorf("parent menu not found: %w", err)
		}
		menu.Level = parent.Level + 1
	} else {
		menu.Level = 0
	}

	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *domain.Menu) error {
	// Recalculate level if parent changed
	if menu.ParentID != nil {
		parent, err := r.FindByID(*menu.ParentID)
		if err != nil {
			return fmt.Errorf("parent menu not found: %w", err)
		}
		menu.Level = parent.Level + 1
	} else {
		menu.Level = 0
	}

	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id int64) error {
	// Check if menu has children
	var count int64
	r.db.Model(&domain.Menu{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete menu with children")
	}

	return r.db.Delete(&domain.Menu{}, id).Error
}

func (r *menuRepository) FindByID(id int64) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindByUUID(uuid string) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.db.Where("uuid = ?", uuid).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindAll() ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.db.Order("order_index ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByParentID(parentID *int64) ([]domain.Menu, error) {
	var menus []domain.Menu
	query := r.db.Order("order_index ASC, id ASC")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindRootMenus() ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.db.Where("parent_id IS NULL").
		Order("order_index ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindHierarchical() ([]domain.Menu, error) {
	var rootMenus []domain.Menu

	// Get all root menus (parent_id IS NULL)
	err := r.db.Where("parent_id IS NULL").
		Order("order_index ASC, id ASC").
		Find(&rootMenus).Error

	if err != nil {
		return nil, err
	}

	// Load children recursively
	for i := range rootMenus {
		r.loadChildren(&rootMenus[i])
	}

	return rootMenus, nil
}

func (r *menuRepository) loadChildren(menu *domain.Menu) {
	var children []domain.Menu
	r.db.Where("parent_id = ?", menu.ID).
		Order("order_index ASC, id ASC").
		Find(&children)

	if len(children) > 0 {
		for i := range children {
			r.loadChildren(&children[i])
		}
		menu.Children = children
	}
}

func (r *menuRepository) FindHierarchicalByRootID(rootID int64) ([]domain.Menu, error) {
	var rootMenu domain.Menu

	// Get the specific root menu
	err := r.db.First(&rootMenu, rootID).Error
	if err != nil {
		return nil, err
	}

	// Load children recursively
	r.loadChildren(&rootMenu)

	return []domain.Menu{rootMenu}, nil
}

func (r *menuRepository) FindDetailByID(id int64) (*domain.MenuDetail, error) {
	var menu domain.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}

	detail := &domain.MenuDetail{
		Menu:  menu,
		Depth: menu.Level,
	}

	// Get parent data if exists
	if menu.ParentID != nil {
		var parent domain.Menu
		err := r.db.First(&parent, *menu.ParentID).Error
		if err == nil {
			detail.ParentData = &domain.MenuParentInfo{
				ID:   parent.ID,
				UUID: parent.UUID,
				Name: parent.Name,
				Code: parent.Code,
			}
		}
	}

	return detail, nil
}

func (r *menuRepository) FindChildrenByParentID(parentID int64) ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.db.Where("parent_id = ?", parentID).
		Order("order_index ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

