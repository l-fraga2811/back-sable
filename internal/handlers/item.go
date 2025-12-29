package handlers

import (
    "github.com/gofiber/fiber/v3"
    "github.com/l-fraga2811/back-sable/internal/models"
    "github.com/l-fraga2811/back-sable/internal/repository"
)

type ItemHandler struct {
    itemRepo repository.ItemRepository
}

func NewItemHandler(itemRepo repository.ItemRepository) *ItemHandler {
    return &ItemHandler{
        itemRepo: itemRepo,
    }
}

func (h *ItemHandler) Create(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    var req models.CreateItemRequest
    if err := c.Bind().Body(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data: " + err.Error()})
    }

    item := &models.Item{
        UserID:      userID,
        Title:       req.Title,
        Description: req.Description,
        Price:       req.Price,
        Completed:   false,
    }

    if err := h.itemRepo.Create(item); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating item"})
    }

    return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ItemHandler) GetAll(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    items, err := h.itemRepo.GetAll(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching items"})
    }

    return c.JSON(items)
}

func (h *ItemHandler) GetByID(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    itemID := c.Params("id")
    item, err := h.itemRepo.GetByID(itemID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
    }

    if item.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to access this item"})
    }

    return c.JSON(item)
}

func (h *ItemHandler) Update(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    itemID := c.Params("id")
    item, err := h.itemRepo.GetByID(itemID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
    }

    if item.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to update this item"})
    }

    var req models.UpdateItemRequest
    if err := c.Bind().Body(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data: " + err.Error()})
    }

    if req.Title != "" {
        item.Title = req.Title
    }
    if req.Description != "" {
        item.Description = req.Description
    }
    if req.Price != 0 {
        item.Price = req.Price
    }
    item.Completed = req.Completed

    if err := h.itemRepo.Update(item); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating item"})
    }

    return c.JSON(item)
}

func (h *ItemHandler) Delete(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    itemID := c.Params("id")
    item, err := h.itemRepo.GetByID(itemID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
    }

    if item.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to delete this item"})
    }

    if err := h.itemRepo.Delete(itemID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting item"})
    }

    return c.JSON(fiber.Map{"message": "Item deleted successfully"})
}

func (h *ItemHandler) requireAuth(c fiber.Ctx) string {
    userID, ok := c.Locals("userID").(string)
    if !ok || userID == "" {
        c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
        return ""
    }
    return userID
}