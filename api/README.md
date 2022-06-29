# go-eliona Api 

The go-eliona Api package provides handy methods to read and write data from and to an eliona environment using the [Eliona API](https://github.com/eliona-smart-building-assistant/eliona-api). The API is defined as OpenAPI. So the all the source code e.g. the data models can be generated. Instruction to regenerate see below. 

## Installation
To use the log package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/api"
```

To define the Api server endpoint use the `API_ENDPOINT` environment variable:

```bash
export API_ENDPOINT=https://api.eliona.io/v2
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```bash
export HTTP_PROXY=https://proxy_name:proxy_port
```

## Usage

After installation, you can access the Eliona Api.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/api"
```

## Generation ##

The OpenAPI specification can be used to generate the go source code. To do this, you can use one of the [OpenAPI Generators](https://openapi-generator.tech/). The following instruction uses the JAR-based generator.

At first, download the [Generator Jar](https://openapi-generator.tech/docs/installation#jar) file: `https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/6.0.0/openapi-generator-cli-6.0.0.jar` to any directory outside the go-eliona project. Also, Java 8 runtime at a minimum is required.

Then you can generate the go client, with the following command. Note, that you configure the package name for the generated files to `api`. The output folder is set to `api`, where the previous generated source code files are located. The `api/.openapi-generator-ignore` file defines files that should not be generated again.

```bash
java -jar openapi-generator-cli-6.0.0.jar generate \
  -g go \
  -i https://raw.githubusercontent.com/eliona-smart-building-assistant/eliona-api/develop/eliona-api-v2.yaml \
  -o ./api \
  --additional-properties=packageName=api
```