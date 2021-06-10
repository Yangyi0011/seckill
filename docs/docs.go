// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/goods": {
            "post": {
                "description": "添加秒杀商品进秒杀系统",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商品管理"
                ],
                "summary": "添加商品",
                "parameters": [
                    {
                        "description": "秒杀商品传输信息",
                        "name": "goods",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.GoodsDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    }
                }
            }
        },
        "/api/goods/{id}": {
            "get": {
                "description": "通过 id 查询秒杀商品",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商品管理"
                ],
                "summary": "查询商品",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "description": "用户登录签发 JWT",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "用户",
                        "name": "loginUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LoginUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    }
                }
            }
        },
        "/api/user/logout": {
            "post": {
                "description": "用户退出登录，清除登录 token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "退出登录",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    }
                }
            }
        },
        "/api/user/register": {
            "post": {
                "description": "注册用户并保存到数据库",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "注册用户",
                        "name": "registerUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RegisterUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.GoodsDTO": {
            "type": "object",
            "required": [
                "amount",
                "endTime",
                "name",
                "price",
                "startTime",
                "stock"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "default": 0,
                    "minimum": 1
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "endTime": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "img": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "originPrice": {
                    "type": "number",
                    "default": 0,
                    "minimum": 0
                },
                "price": {
                    "type": "number",
                    "default": 0,
                    "minimum": 0
                },
                "startTime": {
                    "type": "string"
                },
                "stock": {
                    "type": "integer",
                    "default": 0,
                    "minimum": 0
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "model.LoginUser": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "123"
                },
                "username": {
                    "type": "string",
                    "example": "tom"
                }
            }
        },
        "model.RegisterUser": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "kind": {
                    "type": "integer",
                    "example": 0
                },
                "password": {
                    "type": "string",
                    "example": "123"
                },
                "username": {
                    "type": "string",
                    "example": "tom"
                }
            }
        },
        "model.Result": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 0
                },
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string",
                    "example": "响应信息"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}