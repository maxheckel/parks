package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/maxheckel/parks/services/store"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxheckel/parks/services/authenticator"
	"net/http"
)

func Callback(ctx *fiber.Ctx) error {

	session, _ := store.Store.Get(ctx)
	fmt.Println(session.Get("state"))
	if ctx.Query("state", "") != session.Get("state") {
		return ctx.Status(400).JSON(map[string]string{
			"type":    "BAD_REQUEST",
			"message": "could not load store",
			"error":   "no store",
		})
	}
	auth, err := authenticator.New()
	if err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not load authenticator",
			"error":   err.Error(),
		})
	}
	// Exchange an authorization code for a token.
	token, err := auth.Exchange(ctx.Context(), ctx.Query("code"))
	if err != nil {

		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not exchange context",
			"error":   err.Error(),
		})
	}

	idToken, err := auth.VerifyIDToken(ctx.Context(), token)
	if err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not verify ID token ",
			"error":   err.Error(),
		})
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not get claims",
			"error":   err.Error(),
		})
	}

	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	if err := session.Save(); err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not save session",
			"error":   err.Error(),
		})
	}

	// Redirect to logged in page.
	return ctx.Redirect("/account/user", http.StatusTemporaryRedirect)
}

func Login(ctx *fiber.Ctx) error {
	auth, err := authenticator.New()
	if err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not load authenticator",
			"error":   err.Error(),
		})
	}

	state, err := generateRandomState()
	if err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not generate random state",
			"error":   err.Error(),
		})
	}

	session, err := store.Store.Get(ctx)
	if err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not get session",
			"error":   err.Error(),
		})
	}
	// Save the state inside the session.
	session.Set("state", state)
	if err := session.Save(); err != nil {

		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not save session",
			"error":   err.Error(),
		})
	}

	return ctx.Redirect(auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func User(ctx *fiber.Ctx) {

}

func Logout(ctx *fiber.Ctx) error {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not parse logout url",
			"error":   err.Error(),
		})
	}

	scheme := "http"
	//if ctx.Request().URI() != nil {
	//	scheme = "https"
	//}

	returnTo, err := url.Parse(scheme + "://" + string(ctx.Request().Host()))
	if err != nil {

		return ctx.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not parse return url",
			"error":   err.Error(),
		})
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	return ctx.Redirect(logoutUrl.String(), http.StatusTemporaryRedirect)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
