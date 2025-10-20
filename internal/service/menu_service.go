package service

import (
	"fmt"
	"time"

	"stk-technical-test-api/internal/domain"
)

type menuService struct {
	repo domain.MenuRepository
}

// NewMenuService creates a new menu service instance
func NewMenuService(repo domain.MenuRepository) domain.MenuService {
	return &menuService{
		repo: repo,
	}
}

func (s *menuService) CreateMenu(req *domain.CreateMenuRequest) (*domain.Menu, error) {
	// Validate parent exists if provided
	if req.ParentID != nil {
		_, err := s.repo.FindByID(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent menu not found")
		}
	}

	menu := &domain.Menu{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Route:       req.Route,
		Icon:        req.Icon,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(menu)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	return menu, nil
}

func (s *menuService) UpdateMenu(id int64, req *domain.UpdateMenuRequest) (*domain.Menu, error) {
	// Check if menu exists
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	// Validate parent exists if provided
	if req.ParentID != nil {
		// Check if trying to set itself as parent
		if *req.ParentID == id {
			return nil, fmt.Errorf("menu cannot be its own parent")
		}

		_, err := s.repo.FindByID(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent menu not found")
		}
	}

	// Update fields
	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Code = req.Code
	menu.Description = req.Description
	menu.Route = req.Route
	menu.Icon = req.Icon
	menu.OrderIndex = req.OrderIndex
	menu.IsActive = req.IsActive
	menu.UpdatedAt = time.Now()

	err = s.repo.Update(menu)
	if err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	return menu, nil
}

func (s *menuService) DeleteMenu(id int64) error {
	// Check if menu exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("menu not found")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	return nil
}

func (s *menuService) GetMenuByID(id int64) (*domain.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}
	return menu, nil
}

func (s *menuService) GetMenuByUUID(uuid string) (*domain.Menu, error) {
	menu, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}
	return menu, nil
}

func (s *menuService) GetAllMenus() ([]domain.Menu, error) {
	menus, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get menus: %w", err)
	}
	return menus, nil
}

func (s *menuService) GetRootMenus() ([]domain.Menu, error) {
	menus, err := s.repo.FindRootMenus()
	if err != nil {
		return nil, fmt.Errorf("failed to get root menus: %w", err)
	}
	return menus, nil
}

func (s *menuService) GetMenuHierarchy() ([]domain.Menu, error) {
	menus, err := s.repo.FindHierarchical()
	if err != nil {
		return nil, fmt.Errorf("failed to get menu hierarchy: %w", err)
	}
	return menus, nil
}

func (s *menuService) GetHierarchyByRootID(rootID int64) ([]domain.Menu, error) {
	// Check if menu exists and is a root menu
	menu, err := s.repo.FindByID(rootID)
	if err != nil {
		return nil, fmt.Errorf("root menu not found")
	}

	if menu.ParentID != nil {
		return nil, fmt.Errorf("menu is not a root menu")
	}

	menus, err := s.repo.FindHierarchicalByRootID(rootID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu hierarchy: %w", err)
	}
	return menus, nil
}

func (s *menuService) GetMenuDetail(id int64) (*domain.MenuDetail, error) {
	detail, err := s.repo.FindDetailByID(id)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}
	return detail, nil
}

func (s *menuService) GetChildrenByParentID(parentID int64) ([]domain.Menu, error) {
	// Validate parent exists
	_, err := s.repo.FindByID(parentID)
	if err != nil {
		return nil, fmt.Errorf("parent menu not found")
	}

	menus, err := s.repo.FindChildrenByParentID(parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}
	return menus, nil
}

