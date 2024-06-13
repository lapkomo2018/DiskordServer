package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"time"
)

type UserService interface {
	Create(email, password string) (*core.User, error)
	Login(email, password string) (token string, err error)
	LoadFiles(user *core.User) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary Signup
// @Tags user
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body signupInput true "acc info"
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /user/signup [post]

type signupInput struct {
	Email    string
	Password string
}

func (u *UserHandler) Signup(c *fiber.Ctx) error {
	//get email/pass
	var body signupInput
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	if _, err := u.userService.Create(body.Email, body.Password); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//respond
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// @Summary Login
// @Tags user
// @Description login account
// @ID login-account
// @Accept json
// @Produce json
// @Param input body loginInput true "acc info"
// @Success 200 {object} loginOutput
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /user/login [post]

type loginInput struct {
	Email    string
	Password string
}
type loginOutput struct {
	Authorization string
}

func (u *UserHandler) Login(c *fiber.Ctx) error {
	// get email/pass
	var body loginInput
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	var token string
	var err error
	if token, err = u.userService.Login(body.Email, body.Password); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//setting authorization cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HTTPOnly: true,
		SameSite: "Lax",
	})
	//respond
	return c.Status(fiber.StatusOK).JSON(loginOutput{
		Authorization: token,
	})
}

func (u *UserHandler) Validate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(core.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"userEmail": user.Email,
	})
}

type outputFile struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Size     uint64 `json:"size"`
	IsPublic bool   `json:"isPublic"`
}

func (u *UserHandler) Files(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(core.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	if err := u.userService.LoadFiles(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load user files")
	}

	var files []outputFile
	for _, file := range user.Files {
		files = append(files, outputFile{
			Id:       file.ID,
			Name:     file.Name,
			Size:     file.Size,
			IsPublic: file.IsPublic,
		})
	}

	return c.Status(fiber.StatusOK).JSON(files)
}
