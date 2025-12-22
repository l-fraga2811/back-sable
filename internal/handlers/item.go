package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/models"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

type ItemHandler struct {
	client *supabase.Client
}

func NewItemHandler(client *supabase.Client) *ItemHandler {
	return &ItemHandler{
		client: client,
	}
}

func (h *ItemHandler) Create(c fiber.Ctx) error {
	userID, accessToken, ok := h.requireAuth(c)
	if !ok {
		return nil
	}

	var req models.CreateItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data: " + err.Error()})
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	payload := supabase.CreateItemPayload{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := h.client.CreateItem(c.Context(), accessToken, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating item"})
	}

	item, err := h.mapItemRow(created)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing item"})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ItemHandler) GetAll(c fiber.Ctx) error {
	_, accessToken, ok := h.requireAuth(c)
	if !ok {
		return nil
	}

	rows, err := h.client.ListItems(c.Context(), accessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching items"})
	}

	items := make([]models.Item, 0, len(rows))
	for _, row := range rows {
		item, err := h.mapItemRow(row)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing items"})
		}
		items = append(items, item)
	}

	return c.JSON(items)
}

func (h *ItemHandler) GetByID(c fiber.Ctx) error {
	userID, accessToken, ok := h.requireAuth(c)
	if !ok {
		return nil
	}

	itemID := c.Params("id")
	row, found, err := h.client.GetItemByID(c.Context(), accessToken, itemID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching item"})
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	item, err := h.mapItemRow(row)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing item"})
	}

	if item.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to access this item"})
	}

	return c.JSON(item)
}

func (h *ItemHandler) Update(c fiber.Ctx) error {
	userID, accessToken, ok := h.requireAuth(c)
	if !ok {
		return nil
	}

	itemID := c.Params("id")
	row, found, err := h.client.GetItemByID(c.Context(), accessToken, itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching item"})
	}
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	item, err := h.mapItemRow(row)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing item"})
	}
	if item.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to update this item"})
	}

	var req models.UpdateItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data: " + err.Error()})
	}

	payload := supabase.UpdateItemPayload{}
	if req.Title != "" {
		payload.Title = &req.Title
	}
	if req.Description != "" {
		payload.Description = &req.Description
	}
	if req.Price != 0 {
		payload.Price = &req.Price
	}
	// Note: Boolean zero value is false, so we might need a pointer in request struct to distinguish between explicit false and zero value
	// For simplicity, assuming Completed is updated only if explicitly set in a real scenario or passed as pointer
	payload.Completed = &req.Completed

	updatedRow, updated, err := h.client.UpdateItem(c.Context(), accessToken, itemID, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating item"})
	}
	if !updated {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	updatedItem, err := h.mapItemRow(updatedRow)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing item"})
	}

	return c.JSON(updatedItem)
}

func (h *ItemHandler) Delete(c fiber.Ctx) error {
	_, accessToken, ok := h.requireAuth(c)
	if !ok {
		return nil
	}

	itemID := c.Params("id")
	deleted, err := h.client.DeleteItem(c.Context(), accessToken, itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting item"})
	}
	if !deleted {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	return c.JSON(fiber.Map{"message": "Item deleted successfully"})
}

func (h *ItemHandler) requireAuth(c fiber.Ctx) (string, string, bool) {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		return "", "", false
	}

	token, ok := c.Locals("token").(string)
	if !ok || token == "" {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		return "", "", false
	}

	return userID, token, true
}

func (h *ItemHandler) mapItemRow(row supabase.ItemRow) (models.Item, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return models.Item{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	if err != nil {
		return models.Item{}, err
	}

	return models.Item{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Price:       row.Price,
		Completed:   row.Completed,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UserID:      row.UserID,
	}, nil
}
