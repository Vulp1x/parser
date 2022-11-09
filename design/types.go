// nolint
package design

import (
	. "goa.design/goa/v3/dsl"
)

// Creds defines the credentials to use for authenticating to service methods.
var Creds = Type("Creds", func() {
	Field(1, "jwt", String, "JWT token", func() {
		Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
	})
	Required("jwt")
})

// DatasetStatus описывает статус задачи
var DatasetStatus = Type("DatasetStatus", Int, func() {
	Enum(1, 2, 3, 4, 5, 6)
	Description(`1 - датасет только создан
	2- начали поиск блогеров
	3- успешно закончили поиска похожих блогеров
	4- начали парсинг юзеров у блогеров 
	5- успешно закончили парсинг юзеров
	6- всё сломалось `)
})

// Blogger описывает блоггера, который используется при парсинге
var Blogger = Type("Blogger", func() {
	Attribute("id", String, "", func() {
		Format(FormatUUID)
	})

	Attribute("username", String, func() {
		Description("имя аккаунта в инстаграме")
	})

	Attribute("user_id", Int64, func() {
		Description("user_id в инстаграме, -1 если неизвестен")
		Meta("struct:tag:json", "user_id")
	})

	Attribute("dataset_id", String, func() {
		Description("айди датасета, к которому принадлежит блоггер")
		Meta("struct:tag:json", "dataset_id")
		Format(FormatUUID)
	})

	Attribute("is_initial", Boolean, func() {
		Description("является ли блоггер изначально в датасете или появился при парсинге")
		Meta("struct:tag:json", "is_initial")
	})

	Required("id", "dataset_id", "username", "is_initial", "user_id")
})

// Dataset описывает рекламную кампанию
var Dataset = Type("Dataset", func() {
	Attribute("id", String, "", func() {
		Format(FormatUUID)
	})

	Attribute("bloggers", ArrayOf(Blogger))

	Attribute("status", DatasetStatus)
	Attribute("title", String, "название задачи")
	Attribute("posts_per_blogger", Int32, func() {
		Description("имена аккаунтов, для которых ищем похожих")
		Meta("struct:tag:json", "posts_per_blogger")
	})

	Attribute("liked_per_post", Int32, func() {
		Description("сколько лайкнувших для каждого поста брать")
		Meta("struct:tag:json", "liked_per_post")
	})

	Attribute("commented_per_post", Int32, func() {
		Description("сколько прокоментировааших для каждого поста брать")
		Meta("struct:tag:json", "commented_per_post")
	})

	Required("id", "bloggers", "status", "title", "posts_per_blogger", "liked_per_post", "commented_per_post")
})

var ParsingProgress = Type("ParsingProgress", func() {
	Attribute("bloggers_parsed", Int, "количество блогеров, у которых спарсили пользователей", func() {
		Meta("struct:tag:json", "bloggers_parsed")
	})

	Attribute("targets_saved", Int, "количество сохраненных доноров", func() {
		Meta("struct:tag:json", "filtered_bloggers")
	})

	Attribute("done", Boolean, "закончен ли парсинг блогеров")

	Required("bloggers_parsed", "targets_saved", "done")
})

var DatasetProgress = Type("DatasetProgress", func() {
	Attribute("bloggers", ArrayOf(Blogger), func() {
		Description("блогеры, которых уже нашли")
	})

	Attribute("initial_bloggers", Int, "количество блогеров, которые были изначально", func() {
		Meta("struct:tag:json", "initial_bloggers")
	})

	Attribute("new_bloggers", Int, "количество блогеров, которых нашли", func() {
		Meta("struct:tag:json", "new_bloggers")
	})

	Attribute("filtered_bloggers", Int, "количество блогеров, которые проходят проверку по коду региона", func() {
		Meta("struct:tag:json", "filtered_bloggers")
	})

	Attribute("done", Boolean, "закончена ли задача")

	Required("bloggers", "initial_bloggers", "new_bloggers", "filtered_bloggers", "done")
})

var Bot = Type("Bot", func() {
	Field(1, "username", String, func() {
		Description("имя аккаунта в инстаграме")
	})

	Field(2, "user_id", Int64, "количество блогеров, которые проходят проверку по коду региона", func() {
		Meta("struct:tag:json", "user_id")
	})

	Field(3, "session_id", String, "количество блогеров, которые проходят проверку по коду региона", func() {
		Meta("struct:tag:json", "session_id")
	})

	Field(4, "proxy", Proxy, "прокси для бота", func() {
	})

	Required("username", "user_id", "session_id", "proxy")
})

var Proxy = Type("Proxy", func() {
	Field(1, "host", String, func() {
		Description("имя аккаунта в инстаграме")
	})

	Field(2, "port", Int32, "количество блогеров, которые проходят проверку по коду региона")

	Field(3, "login", String, func() {
		Description("имя аккаунта в инстаграме")
	})

	Field(4, "pass", String, func() {
		Description("имя аккаунта в инстаграме")
	})

	Required("host", "port", "login", "pass")
})
