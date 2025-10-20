package handler

import (
	"net/http"
	"strconv"

	"stk-technical-test-api/internal/domain"
	"stk-technical-test-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	service domain.MenuService
}

func NewMenuHandler(service domain.MenuService) *MenuHandler {
	return &MenuHandler{
		service: service,
	}
}

// CreateMenu godoc
// @Summary Create a new menu
// @Description Create a new menu or submenu
// @Tags menus
// @Accept json
// @Produce json
// @Param menu body domain.CreateMenuRequest true "Menu data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req domain.CreateMenuRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	menu, err := h.service.CreateMenu(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create menu", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Menu created successfully", menu)
}

// GetMenuHierarchy godoc
// @Summary Get menu hierarchy
// @Description Get all menus in hierarchical structure
// @Tags menus
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/menus/hierarchy [get]
func (h *MenuHandler) GetMenuHierarchy(c *gin.Context) {
	menus, err := h.service.GetMenuHierarchy()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get menu hierarchy", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Menu hierarchy retrieved successfully", menus)
}

// GetAllMenus godoc
// @Summary Get all menus
// @Description Get all menus (flat list)
// @Tags menus
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/menus [get]
func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	menus, err := h.service.GetAllMenus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get menus", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Menus retrieved successfully", menus)
}

// GetMenuByID godoc
// @Summary Get menu by ID
// @Description Get a single menu by ID
// @Tags menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/menus/{id} [get]
func (h *MenuHandler) GetMenuByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid menu ID", err.Error())
		return
	}

	menu, err := h.service.GetMenuByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Menu not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Menu retrieved successfully", menu)
}

// GetMenuByUUID godoc
// @Summary Get menu by UUID
// @Description Get a single menu by UUID
// @Tags menus
// @Produce json
// @Param uuid path string true "Menu UUID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/menus/uuid/{uuid} [get]
func (h *MenuHandler) GetMenuByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		response.Error(c, http.StatusBadRequest, "Invalid menu UUID", "UUID is required")
		return
	}

	menu, err := h.service.GetMenuByUUID(uuid)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Menu not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Menu retrieved successfully", menu)
}

// UpdateMenu godoc
// @Summary Update a menu
// @Description Update an existing menu
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu body domain.UpdateMenuRequest true "Menu data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid menu ID", err.Error())
		return
	}

	var req domain.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	menu, err := h.service.UpdateMenu(id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update menu", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Menu updated successfully", menu)
}

// DeleteMenu godoc
// @Summary Delete a menu
// @Description Delete a menu by ID
// @Tags menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid menu ID", err.Error())
		return
	}

	err = h.service.DeleteMenu(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to delete menu", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Menu deleted successfully", nil)
}

