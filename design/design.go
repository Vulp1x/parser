// nolint
package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("parser-api", func() {
	Title("api for parsing instagram")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the signin endpoint. Supports scopes "api:read" and "api:write".`)
	Scope("driver", "Read-only access")
	Scope("admin", "Read and write access")
})

var _ = Service("datasets_service", func() {
	Description("сервис для создания, редактирования и работы с задачами (рекламными компаниями)")

	Error("unauthorized", String, "Credentials are invalid")
	Error("bad request", String, "Invalid request")
	Error("internal error", String, "internal error")
	Error("dataset not found", String, "Not found")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("bad request", StatusBadRequest)
		Response("dataset not found", StatusNotFound)
		Response("internal error", StatusInternalServerError)
	})

	Method("create dataset draft", func() {
		Description("создать драфт задачи")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Token("type", DatasetType)

			Required("token", "type")
		})

		Result(String, func() {
			Description("dataset_id для созданной задачи")
			Format(FormatUUID)
		})

		HTTP(func() {
			POST("/api/datasets/draft/")
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("update dataset", func() {
		Description(`обновить информацию о задаче. Не меняет статус задачи, можно вызывать сколько угодно раз.
			Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала остановить выполнение.`)

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи, которую хотим обновить")
				Meta("struct:tag:json", "dataset_id")
			})

			Attribute("original_accounts", ArrayOf(String), func() {
				Description("имена аккаунтов, для которых ищем похожих")
				Meta("struct:tag:json", "original_accounts")
			})

			Attribute("posts_per_blogger", UInt, func() {
				Description("имена аккаунтов, для которых ищем похожих")
				Meta("struct:tag:json", "posts_per_blogger")
			})

			Attribute("liked_per_post", UInt, func() {
				Description("сколько лайкнувших для каждого поста брать")
				Meta("struct:tag:json", "liked_per_post")
			})

			Attribute("commented_per_post", UInt, func() {
				Description("сколько прокоментировааших для каждого поста брать")
				Meta("struct:tag:json", "commented_per_post")
			})

			Attribute("phone_code", Int32, func() {
				Description("код региона, по которому будем сортировать")
				Example(7)
				Minimum(1)
				Maximum(1000)
				Meta("struct:tag:json", "phone_code")
			})

			Attribute("title", String, func() {
				Description("название задачи")
			})

			Attribute("subscribers_per_blogger", UInt, func() {
				Description("сколько подписчиков для каждого блоггера брать")
				Meta("struct:tag:json", "subscribers_per_blogger")
			})

			Required("token", "dataset_id")
		})

		Result(Dataset)

		HTTP(func() {
			PUT("/api/datasets/{dataset_id}/")
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("find similar", func() {
		Description("начать выполнение задачи ")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Required("token", "dataset_id")
		})

		Result(Dataset)

		HTTP(func() {
			POST("/api/datasets/{dataset_id}/start/")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("get progress", func() {
		Description("получить статус выполнения поиска похожих аккаунтов по айди датасета")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Required("token", "dataset_id")
		})

		Result(DatasetProgress)

		HTTP(func() {
			GET("/api/datasets/{dataset_id}/progress/")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("parse dataset", func() {
		Description("получить базу доноров для выбранных блогеров")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Required("token", "dataset_id")
		})

		Result(func() {
			Attribute("status", DatasetStatus)
			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Required("dataset_id", "status")
		})

		HTTP(func() {
			POST("/api/datasets/{dataset_id}/parse/")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("get dataset", func() {
		Description("получить задачу по id")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Required("token", "dataset_id")
		})

		Result(Dataset)

		HTTP(func() {
			GET("/api/datasets/{dataset_id}/")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("get parsing progress", func() {
		Description("получить статус выполнения парсинга аккаунтов ")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Required("token", "dataset_id")
		})

		Result(ParsingProgress)

		HTTP(func() {
			GET("/api/datasets/{dataset_id}/parsing_progress/")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("download targets", func() {
		Description("получить всех пользователей, которых спарсили в задаче")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Attribute("format", Int, func() {
				Enum(1, 2, 3)
				Description(`1- только user_id, 2- только username, 3 - и то и другое`)
				Default(3)
			})

			Required("token", "dataset_id", "format")
		})

		Result(ArrayOf(String))

		HTTP(func() {
			GET("/api/datasets/{dataset_id}/download/")
			Param("format:format")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("list datasets", func() {
		Description("получить все задачи для текущего пользователя")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Required("token")
		})

		Result(ArrayOf(Dataset))

		HTTP(func() {
			GET("/api/datasets/")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
			Response(StatusInternalServerError)
		})
	})

	Method("upload files", func() {
		Description("загрузить файл с ботами, прокси")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("dataset_id", String, func() {
				Description("id задачи, в которую загружаем пользователей/прокси")
				Meta("struct:tag:json", "dataset_id")
			})

			Attribute("proxies_filename", String, func() {
				Meta("struct:tag:json", "proxies_filename")
			})

			Attribute("bots_filename", String, func() {
				Meta("struct:tag:json", "bots_filename")
			})

			Attribute("bots", ArrayOf(BotAccountRecord), "список ботов")
			Attribute("proxies", ArrayOf(ProxyRecord), "список проксей для использования")

			Required("token", "dataset_id", "bots", "proxies", "proxies_filename", "bots_filename")
		})

		Result(func() {
			Attribute("upload_errors", ArrayOf(UploadError), func() {
				Description("ошибки, которые возникли при загрузке файлов")
				Meta("struct:tag:json", "upload_errors")
			})

			Required("upload_errors")
		})

		HTTP(func() {
			POST("/api/datasets/{dataset_id}/upload/")
			MultipartRequest()
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})
})
