// nolint
package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("rest-api", func() {
	Title("REST api for simple route app")
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

			Attribute("title", String, func() {
				Description("название задачи")
			})

			Attribute("original_accounts", ArrayOf(String), func() {
				Description("имена аккаунтов, для которых ищем похожих")
				Meta("struct:tag:json", "original_accounts")
			})

			Required("token", "title", "original_accounts")
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

			Attribute("filter_code", Int, func() {
				Description("код региона, фильтруем аккаунты по нему")
				Example(7)
				Meta("struct:tag:json", "filter_code")
			})

			Required("token", "dataset_id")
		})

		Result(func() {
			Attribute("status", DatasetStatus)
			Attribute("dataset_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "dataset_id")
			})

			Attribute("bloggers", ArrayOf(Blogger))

			Required("dataset_id", "status", "bloggers")
		})

		HTTP(func() {
			POST("/api/datasets/{dataset_id}/start/")
			Response(StatusOK)
			Response(StatusBadRequest)
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
			POST("/api/datasets/{dataset_id}/stop/")
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

	Method("get progress", func() {
		Description("получить статус выполнения задачи по id")

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
})