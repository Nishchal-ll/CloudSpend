package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	"github.com/Nishchal-ll/CloudSpend/models"
	"github.com/gin-gonic/gin"
)

const sessionCookieName = "cloudspend_session"

var secretKey = []byte(getSecretKey())

func getSecretKey() string {
	k := os.Getenv("SESSION_SECRET")
	if k == "" {
		k = "cloudspend-secure-secret-key-2026-production"
	}
	return k
}

// createSessionToken signs user email
func createSessionToken(email string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(email))
	sig := hex.EncodeToString(h.Sum(nil))
	return email + ":" + sig
}

// verifySessionToken verifies signature and extracts email
func verifySessionToken(token string) (string, bool) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return "", false
	}
	email, sig := parts[0], parts[1]

	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(email))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return email, true
	}
	return "", false
}

// RenderLogin renders the login page
func RenderLogin(c *gin.Context) {
	if cookie, err := c.Cookie(sessionCookieName); err == nil {
		if _, ok := verifySessionToken(cookie); ok {
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "CloudSpend — Login",
	})
}

// LoginHandler processes login form or JSON API
func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid email and password."})
		return
	}

	user, err := models.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
		return
	}

	// Set session cookie
	token := createSessionToken(user.Email)
	c.SetCookie(sessionCookieName, token, 86400*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"redirect": "/dashboard",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// LogoutHandler clears session and redirects
func LogoutHandler(c *gin.Context) {
	c.SetCookie(sessionCookieName, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// AuthRequiredMiddleware ensures request is authenticated
func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(sessionCookieName)
		if err != nil {
			handleUnauthenticated(c)
			return
		}

		email, ok := verifySessionToken(cookie)
		if !ok {
			handleUnauthenticated(c)
			return
		}

		user, err := models.GetUserByEmail(email)
		if err != nil {
			handleUnauthenticated(c)
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}

func handleUnauthenticated(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required. Please login."})
	} else {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
