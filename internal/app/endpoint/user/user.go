package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"github.com/lapkomo2018/DiskordServer/pkg/email"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

type User struct {
	db *gorm.DB
}

func New(db *gorm.DB) *User {
	return &User{
		db: db,
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

func (u *User) Signup(c *fiber.Ctx) error {
	//get email/pass
	var body signupInput
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	//validate email
	if email.Validate(body.Email) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email")
	}

	//hash pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	//create user
	user := model.User{
		Email:    body.Email,
		Password: string(hash),
	}

	if err := u.db.Create(&user).Error; err != nil {
		errString := err.Error()
		if strings.Contains(errString, "SQLSTATE 23505") {
			errString = "Email already registered"
		} else if strings.Contains(errString, "SQLSTATE") {
			errString = "Failed to create user due to a database error"
		}
		return fiber.NewError(fiber.StatusInternalServerError, errString)
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

func (u *User) Login(c *fiber.Ctx) error {
	// get email/pass
	var body loginInput
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	//look up user
	var user model.User
	if err := u.db.First(&user, &model.User{Email: body.Email}).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User not found")
	}

	//compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Incorrect password")
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create token")
	}

	//setting authorization cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HTTPOnly: true,
		SameSite: "Lax",
	})
	//respond
	return c.Status(fiber.StatusOK).JSON(loginOutput{
		Authorization: tokenString,
	})
}

func (u *User) Validate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"userEmail": user.Email,
	})
}

func (u *User) Files(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	if err := u.db.Preload("Files").Find(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load user files")
	}

	type returnedFile struct {
		Id       uint   `json:"id"`
		Name     string `json:"name"`
		Size     uint64 `json:"size"`
		IsPublic bool   `json:"isPublic"`
	}
	var files []returnedFile
	for _, file := range user.Files {
		files = append(files, returnedFile{
			Id:       file.ID,
			Name:     file.Name,
			Size:     file.Size,
			IsPublic: file.IsPublic,
		})
	}

	return c.Status(fiber.StatusOK).JSON(files)
}
