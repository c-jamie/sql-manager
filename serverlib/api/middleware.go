package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/c-jamie/sql-manager/serverlib/log"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type contextKey string

const userContextKey = contextKey("user")

// Middleware is the interface used to control permissioning for the app
type Middleware interface {
	// Authenticate determines if a user is allowed visibility on an object
	Authenticate(c *gin.Context)
	// Authorize determines if current subject has been authorized to read an object
	Authorize(code string) gin.HandlerFunc
	// HasPermission
	HasPermission(code string, permission []string) bool
}

type middleware struct {
	AuthURL string
	DB      *sql.DB
}

func NewMiddelware(authUrl string, db *sql.DB) *middleware {
	return &middleware{DB: db, AuthURL: authUrl}
}

func (mi *middleware) Authenticate(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")

	if authorizationHeader == "" {
		c.Abort()
		return
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		mi.invalidAuthenticationTokenResponse(c)
		c.Abort()
		return
	}

	token := headerParts[1]
	if token == "" {
		mi.invalidAuthenticationTokenResponse(c)
		c.Abort()
		return
	}

	user, err := mi.getUser(token)
	log.Debug("user is ", user.Email)

	if err != nil {
		log.Debug("req error is", err)
		mi.badRequest(c, fmt.Errorf("unable to authorize: %w", err))
		return
	}

	log.Debug("user is ", user.Email)

	mi.contextSetUser(c, user)
	c.Next()
}

func (mi *middleware) Authorize(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := mi.contextGetUser(c)
		if !mi.HasPermission(code, user.Permissions) {
			mi.notPermittedResponse(c, code)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (mi *middleware) HasPermission(code string, codes []string) bool {
	for i := range codes {
		if code == codes[i] {
			return true
		}
	}
	return false

}

func (mi *middleware) getUser(token string) (*UserAccount, error) {
	url := mi.AuthURL + "/users"
	body, resp := mi.getReq(url, token)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get user; status %d; body %s", resp.StatusCode, string(body))
	}
	userJson := gjson.Get(string(body), "users")

	log.Debug("get user body: ", string(body))

	var ua UserAccount
	jsonErr := json.Unmarshal([]byte(userJson.String()), &ua)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &ua, nil
}

func (mi *middleware) getReq(url string, bearer string) ([]byte, *http.Response) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "polar-client")
	req.Header.Set("Content-Type", "application/json")

	if bearer != "" {
		log.Info("token: ", bearer)
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	client := &http.Client{}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Error("req error", getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, _ := ioutil.ReadAll(res.Body)
	return body, res
}

func (app *middleware) notPermittedResponse(c *gin.Context, code string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "code " + code + " not permitted"})
}

func (app *middleware) badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (app *middleware) contextSetUser(r *gin.Context, user *UserAccount) *gin.Context {
	r.Set(string(userContextKey), user)
	return r
}

func (app *middleware) invalidAuthenticationTokenResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing authentication token"})
}

func (app *middleware) contextGetUser(r *gin.Context) *UserAccount {
	user, ok := r.Value(string(userContextKey)).(*UserAccount)
	if !ok {
		panic("missing user value in request")
	}
	return user
}
