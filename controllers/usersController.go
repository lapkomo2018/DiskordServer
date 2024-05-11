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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	//validate email
	if validators.ValidateEmail(body.Email) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email",
		})
	}

	//hash pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errString,
		})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	//look up user
	var user models.User
	if err := initializers.DB.Where(&models.User{Email: body.Email}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	//compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect password",
		})
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create token",
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"userEmail": user.Email,
	})
}
