//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
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

package client

import (
	"context"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

func ApiEndpointString() string {
	return common.Getenv("API_ENDPOINT", "http://api-v2:3000/v2")
}

func ApiKeyString() string {
	return common.Getenv("API_TOKEN", "not defined")
}

func NewClient(apiEndpoint string) *api.APIClient {
	cfg := api.NewConfiguration()
	cfg.Servers = api.ServerConfigurations{{URL: apiEndpoint}}
	return api.NewAPIClient(cfg)
}

func AuthenticationContext(apiKey string) context.Context {
	return AuthenticationContextWrap(context.Background(), apiKey)
}

func AuthenticationContextWrap(ctx context.Context, apiKey string) context.Context {
	apiKeys := map[string]api.APIKey{
		"ApiKeyAuth": {Key: apiKey},
	}
	return context.WithValue(ctx, api.ContextAPIKeys, apiKeys)
}
