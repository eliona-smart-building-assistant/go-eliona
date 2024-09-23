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

func ParseEnvironment(r *http.Request) (*Environment, error) {
	token, err := GetBearerTokenString(r)
	if err != nil {
		return nil, err
	}
	return parseEnvironment(token)
}

type Environment struct {
	Aud          string `json:"aud"`
	Exp          int    `json:"exp"`
	Iss          string `json:"iss"`
	Role         string `json:"role"`
	CustId       string `json:"cust_id"`
	ProjId       string `json:"proj_id"`
	RoleId       string `json:"role_id"`
	UserId       string `json:"user_id"`
	Entitlements string `json:"entitlements"`
	jwt.StandardClaims
}

func parseEnvironment(tokenString *string) (*Environment, error) {
	if tokenString == nil {
		return nil, fmt.Errorf("token string is nil")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(*tokenString, &Environment{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if token == nil || token.Claims == nil {
		return nil, fmt.Errorf("token or token claims are nil")
	}

	claims, ok := token.Claims.(*Environment)
	if !ok || claims == nil {
		return nil, fmt.Errorf("failed to parse environment claims from token or claims are nil")
	}

	return claims, nil
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
