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

var BloggersProgress = Type("BloggersProgress", func() {
	Attribute("user_name", String, "имя пользователя бота", func() {
		Meta("struct:tag:json", "user_name")
	})
	Attribute("posts_count", Int, "количество выложенных постов", func() {
		Meta("struct:tag:json", "posts_count")
	})

	Attribute("status", Int, "текущий статус бота, будут ли выкладываться посты")

	Required("user_name", "posts_count", "status")
})

var DatasetProgress = Type("DatasetProgress", func() {
	Attribute("bots_progresses", MapOf(String, BloggersProgress), func() {
		Description("результат работы по каждому боту, ключ- имя бота")
		Meta("struct:tag:json", "bots_progresses")
	})

	Attribute("targets_notified", Int, "количество аккаунтов, которых упомянули в постах", func() {
		Meta("struct:tag:json", "targets_notified")
	})
	Attribute("targets_failed", Int, "количество аккаунтов, которых не получилось упомянуть, при перезапуске задачи будут использованы заново", func() {
		Meta("struct:tag:json", "targets_failed")
	})

	Attribute("targets_waiting", Int, "количество аккаунтов, которых не выбрали для постов", func() {
		Meta("struct:tag:json", "targets_waiting")
	})

	Attribute("targets_waiting", Int, "количество аккаунтов, которых не выбрали для постов", func() {
		Meta("struct:tag:json", "targets_waiting")
	})

	Attribute("done", Boolean, "закончена ли задача")

	Required("bots_progresses", "targets_notified", "targets_failed", "targets_waiting", "done")

})
