package fiberHandlers

import (
	"fmt"
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/usecases"

	"github.com/gofiber/fiber/v2"
)

type FirebaseApi struct {
	ApplicationLogic usecases.Applications
	FirebaseWebKey   string
}

type emailVer struct {
	RequestType string `json:"requestType,omitempty"`
	IdToken     string `json:"idToken,omitempty"`
	Email       string `json:"email,omitempty"`
}

type passwordReset struct {
	RequestType string `json:"requestType,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (f *FirebaseApi) SignUp(c *fiber.Ctx) error {

	var signUpUser presenter.SignUpBody

	if err := c.BodyParser(&signUpUser); err != nil {
		fmt.Println(err)
	}

	// Create the user on Firebase
	user, err := f.ApplicationLogic.UserApp.UserCreate(signUpUser.Email,
		signUpUser.Password, signUpUser.DisplayName)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Create Custom token for the user with a specific UID
	token, err := f.ApplicationLogic.UserApp.UserCreateCustomToken(user.UID)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Request a ID token for Firebase BASED on the custom token
	userIdToken, err := f.ApplicationLogic.UserApp.UserRequestIdToken(
		f.FirebaseWebKey, token)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Sent Email verification for every new user created
	emailVerificationResp, err := f.ApplicationLogic.UserApp.
		UserSendVerificationEmail(f.FirebaseWebKey, userIdToken.IdToken)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"error":   emailVerificationResp.Error,
		})
	}

	// Create User in our database
	_, err = f.ApplicationLogic.UserApp.CreateUser(user.UID, user.Email,
		user.DisplayName, "normal")
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	userApiReturn := presenter.ConvertUserToUserApiReturn(user.Email,
		user.DisplayName)

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userApiReturn,
		"message":  "User registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (f *FirebaseApi) ForgotPassword(c *fiber.Ctx) error {

	var passwordResetEmail presenter.ForgotPasswordBody

	if err := c.BodyParser(&passwordResetEmail); err != nil {
		fmt.Println(err)
	}
	fmt.Println(passwordResetEmail)

	// Send Email to reset password
	emailForgotPassResp, err := f.ApplicationLogic.UserApp.
		UserSendForgotPasswordEmail(f.FirebaseWebKey, passwordResetEmail.Email)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"error":   emailForgotPassResp.Error,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":        true,
		"passwordUpdate": emailForgotPassResp,
		"message":        "The email for password reset was successfully sent",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (f *FirebaseApi) DeleteUser(c *fiber.Ctx) error {

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	// Delete User from Firebase
	deletedUser, err := f.ApplicationLogic.UserApp.UserDelete(userId.String())
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f.ApplicationLogic.UserApp.DeleteUser(userId.String())

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": deletedUser,
		"message":  "The user was deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (f *FirebaseApi) UpdateUserInfo(c *fiber.Ctx) error {

	var userInfoUpdate presenter.SignUpBody

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if err := c.BodyParser(&userInfoUpdate); err != nil {
		fmt.Println(err)
	}

	userUpdated, err := f.ApplicationLogic.UserApp.UserUpdateInfo(userId.String(),
		userInfoUpdate.Email, userInfoUpdate.Password, userInfoUpdate.DisplayName)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f.ApplicationLogic.UserApp.UpdateUser(userId.String(), userUpdated.Email,
		userUpdated.DisplayName)

	userApiReturn := presenter.ConvertUserToUserApiReturn(userUpdated.Email,
		userInfoUpdate.DisplayName)

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userApiReturn,
		"message":  "The user information was updated successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
