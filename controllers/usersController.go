package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
	"github.com/lapkomo2018/DiskordServer/validators"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
)

func Signup(c *fiber.Ctx) error {
	//get email/pass
	var body struct {
		Email    string
		Password string
	}
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	//validate email
	if validators.ValidateEmail(body.Email) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email")
	}

	//hash pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	//create user
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
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

func Login(c *fiber.Ctx) error {
	// get email/pass
	var body struct {
		Email    string
		Password string
	}
	if c.BodyParser(&body) != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to read body")
	}

	//look up user
	var user models.User
	if err := initializers.DB.First(&user, &models.User{Email: body.Email}).Error; err != nil {
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"authorization": tokenString,
	})
}

func Validate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"userEmail": user.Email,
	})
}

func GetUserFiles(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	if err := initializers.DB.Preload("Files").Find(&user).Error; err != nil {
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
