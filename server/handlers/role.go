package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SergeyCherepiuk/session-auth/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RolesHandler struct {
	pdb *gorm.DB
}

func NewRolesHandler(pdb *gorm.DB) *RolesHandler {
	return &RolesHandler{pdb: pdb}
}

type rolesRequestBody struct {
	Name string `json:"name"`
}

func (body rolesRequestBody) validate() error {
	if strings.TrimSpace(body.Name) == "" {
		return errors.New("role name is empty")
	}
	return nil
}

func (handler RolesHandler) GetById(c *fiber.Ctx) error {
	roleId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	role := models.Role{}
	r := handler.pdb.First(&role, roleId)
	if r.Error != nil {
		return r.Error
	}

	return c.JSON(role)
}

func (handler RolesHandler) GetAll(c *fiber.Ctx) error {
	roles := []models.Role{}
	r := handler.pdb.Find(&roles)
	if r.Error != nil {
		return r.Error
	}

	if r.RowsAffected < 1 {
		c.Status(fiber.StatusNoContent)
	} else {
		c.Status(fiber.StatusOK)
	}
	return c.JSON(roles)
}

func (handler RolesHandler) Create(c *fiber.Ctx) error {
	var body rolesRequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	err = body.validate()
	if err != nil {
		return err
	}

	role := models.Role{Name: body.Name}
	r := handler.pdb.Create(&role)
	if r.Error != nil {
		return r.Error
	}

	return c.SendString("Role has been created successfully")
}

func (handler RolesHandler) Update(c *fiber.Ctx) error {
	var body rolesRequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	err = body.validate()
	if err != nil {
		return err
	}

	roleId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	var role models.Role
	r := handler.pdb.First(&role, roleId)
	if r.Error != nil {
		return r.Error
	}

	role.Name = body.Name

	r = handler.pdb.Save(&role)
	if r.Error != nil {
		return r.Error
	}

	return c.SendString("Role has been updated successfully")
}

func (handler RolesHandler) Delete(c *fiber.Ctx) error {
	roleId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	r := handler.pdb.Delete(&models.Role{}, roleId)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected < 1 {
		return errors.New("record not found")
	}

	return c.SendString("Role has been deleted successfully")
}

func (handler RolesHandler) DeleteAll(c *fiber.Ctx) error {
	r := handler.pdb.Where("1 = 1").Delete(&models.Role{})
	if r.Error != nil {
		return r.Error
	}

	return c.SendString(fmt.Sprintf("%d role(s) has/have been deleted", r.RowsAffected))
}
