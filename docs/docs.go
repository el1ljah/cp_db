// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/brands/{BRAND_ID}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "brands"
                ],
                "summary": "Get an information about one brand",
                "parameters": [
                    {
                        "type": "string",
                        "description": "BRAND_ID",
                        "name": "BRAND_ID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Brand"
                        }
                    },
                    "400": {
                        "description": "f u 400",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "f u 400",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "f u ",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Brand": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "logo": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.337",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Clothes store 👚",
	Description:      "WEB лабы. Как я хочу спать.....",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}