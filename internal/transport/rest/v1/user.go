package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"net/http"
	"time"
)

func (h *Handler) initUserRouters(api *echo.Group) {
	user := api.Group("/user")
	{
		user.POST("/signup", h.userSignup)
		user.POST("/login", h.userLogin)

		authorized := user.Group("", h.userIdentify)
		{
			authorized.GET("/validate", userValidate)
			authorized.GET("/files", h.userFiles)
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
	Email    string `json:"email" form:"email" validate:"email" binding:"required"`
	Password string `json:"password" form:"password" validate:"min=8" binding:"required"`
}

func (h *Handler) userSignup(c echo.Context) error {
	//get email/pass
	var body signupInput
	if c.Bind(&body) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read body")
	}

	if err := h.validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := h.userService.Create(body.Email, body.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{})
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

type (
	loginInput struct {
		Email    string `json:"email" form:"email" validate:"email" binding:"required"`
		Password string `json:"password" form:"password" validate:"min=8" binding:"required"`
	}

	loginOutput struct {
		Token string `json:"token" form:"token"`
	}
)

func (h *Handler) userLogin(c echo.Context) error {
	var input loginInput
	if c.Bind(&input) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read input")
	}

	if err := h.validator.Struct(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var token string
	var err error
	if token, err = h.userService.Login(input.Email, input.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     AuthorizationCookie,
		Value:    token,
		Expires:  time.Now().Add(CookieExpireTTL),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, loginOutput{token})
}

type userValidateOutput struct {
	UserEmail string `json:"userEmail" form:"userEmail"`
}

func userValidate(c echo.Context) error {
	user, ok := c.Get(UserLocals).(core.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse user")
	}

	return c.JSON(http.StatusOK, userValidateOutput{UserEmail: user.Email})
}

type userFilesOutputFile struct {
	Id       uint   `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Size     uint64 `json:"size" form:"size"`
	IsPublic bool   `json:"isPublic" form:"isPublic"`
}

func (h *Handler) userFiles(c echo.Context) error {
	user, ok := c.Get(UserLocals).(core.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse user")
	}

	if err := h.userService.LoadFiles(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load user files")
	}

	var files []userFilesOutputFile
	for _, file := range user.Files {
		files = append(files, userFilesOutputFile{
			Id:       file.ID,
			Name:     file.Name,
			Size:     file.Size,
			IsPublic: file.IsPublic,
		})
	}

	return c.JSON(http.StatusOK, files)
}
