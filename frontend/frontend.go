//  This file is part of the eliona project.
//  Copyright © 2023 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package frontend

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/eliona-smart-building-assistant/go-utils/common"

	"github.com/golang-jwt/jwt"
)

func extractBearerToken(authHeader string) string {
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return ""
	}
	if strings.ToLower(split[0]) != "bearer" {
		return ""
	}
	return split[1]
}

func GetBearerTokenString(r *http.Request) (*string, error) {
	authHeader := r.Header.Get("Authorization")
	token := extractBearerToken(authHeader)
	if len(token) == 0 {
		cookie, err := r.Cookie("elionaAuthorization")
		if err != nil {
			return nil, err
		} else {
			return common.Ptr(fmt.Sprintf("%s", cookie.Value)), nil
		}
	}
	return common.Ptr(token), nil
}

type Environment struct {
	ProjectId    *string
	RoleId       *int32
	UserId       *int32
	Entitlements *string
}

func ParseEnvironment(r *http.Request) (*Environment, error) {
	token, err := GetBearerTokenString(r)
	if err != nil {
		return nil, err
	}
	return parseEnvironment(token)
}

func parseEnvironment(tokenString *string) (*Environment, error) {

	var env Environment
	if tokenString == nil {
		return nil, nil
	}

	// Parse the JWT token without validating its signature
	token, _, err := new(jwt.Parser).ParseUnverified(*tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("parsing JWT token: %w", err)
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	if projectId, found := claims["proj_id"].(string); found {
		env.ProjectId = common.Ptr(projectId)
	}
	if entitlements, found := claims["entitlements"].(string); found {
		env.Entitlements = common.Ptr(entitlements)
	}
	if roleIdString, found := claims["role_id"].(string); found {
		roleId, err := strconv.ParseInt(roleIdString, 10, 32)
		if err == nil {
			env.RoleId = common.Ptr(int32(roleId))
		}
	}
	if userIdString, found := claims["user_id"].(string); found {
		userId, err := strconv.ParseInt(userIdString, 10, 32)
		if err == nil {
			env.RoleId = common.Ptr(int32(userId))
		}
	}

	return &env, nil
}

type EnvironmentHandler struct {
	handler http.Handler
}

func NewEnvironmentHandler(handler http.Handler) EnvironmentHandler {
	return EnvironmentHandler{handler: handler}
}

type keyType string

const environmentKey = keyType("environment")

func (h EnvironmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	env, err := ParseEnvironment(r)
	if err != nil {
		h.handler.ServeHTTP(w, r)
	} else {
		h.handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), environmentKey, env)))
	}
}

func GetEnvironment(ctx context.Context) *Environment {
	return ctx.Value(environmentKey).(*Environment)
}
