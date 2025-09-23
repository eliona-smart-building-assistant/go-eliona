# go-eliona Frontend
The go-eliona Frontend package provides functions and data structures to handle Eliona frontend.

## Installation
To use the assetLikes package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/v2/frontend"
```

## Get frontend environment

To retrieve environment information from frontend requests, integrate the `NewEnvironmentHandler`, which embeds this information into the current `context`.

```go
err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"),
    utilshttp.NewCORSEnabledHandler(
        frontend.NewEnvironmentHandler(
            apiserver.NewRouter(
                // add controllers here
            ),
        ),
    ),
)
```
To use environment information in your controller logic, utilize the `GetEnvironment` function.

```go
env := frontend.GetEnvironment(ctx)
fmt.Printf("Environment: %v", env)
```

This function returns assetLike structure like this.

```json
{
"aud": "https://eliona.io/api",
"exp": 1706022910,
"iss": "https://eliona.io",
"role": "api",
"cust_id": "1",
"proj_id": "2",
"role_id": "3",
"user_id": "4",
"entitlements": "user"
}
```