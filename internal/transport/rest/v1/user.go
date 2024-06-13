package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"github.com/lapkomo2018/DiskordServer/pkg/validation"
	"time"
)

func (h *Handler) initUserRouters(api fiber.Router) {
	user := api.Group("/user")
	{
		user.Post("/signup", h.userSignup)
		user.Post("/login", h.userLogin)

		authorized := user.Group("", h.userIdentify)
		{
			authorized.Get("/validate", userValidate)
			authorized.Get("/files", h.userFiles)
		}
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

func (h *Handler) userSignup(c *fiber.Ctx) error {
	//get email/pass
	var body signupInput
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	if err := validation.ValidateEmail(body.Email); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := validation.ValidatePassword(body.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := h.userService.Create(body.Email, body.Password); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

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
	Token string `json:"token"`
}

func (h *Handler) userLogin(c *fiber.Ctx) error {
	// get email/pass
	var body loginInput
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	if err := validation.ValidateEmail(body.Email); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := validation.ValidatePassword(body.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var token string
	var err error
	if token, err = h.userService.Login(body.Email, body.Password); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//setting authorization cookie
	c.Cookie(&fiber.Cookie{
		Name:     AuthorizationCookie,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(loginOutput{token})
}

func userValidate(c *fiber.Ctx) error {
	user, ok := c.Locals(UserLocals).(core.User)
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

func (h *Handler) userFiles(c *fiber.Ctx) error {
	user, ok := c.Locals(UserLocals).(core.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	if err := h.userService.LoadFiles(&user); err != nil {
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
