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
	Description(`1 - задача только создана, нужно загрузить список ботов, прокси и получателей
	2- в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	3- задача готова к запуску
	4- задача запущена 
	5 - задача остановлена
	6 - задача завершена`)
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

	Required("id", "bloggers", "status", "title")
})

var DatasetFilenames = Type("DatasetFileNames", func() {
	Attribute("bots_filename", String, "название файла, из которого брали ботов", func() {
		Meta("struct:tag:json", "bots_filename")
	})
	Attribute("residential_proxies_filename", String, "название файла, из которого брали резидентские прокси", func() {
		Meta("struct:tag:json", "residential_proxies_filename")
	})
	Attribute("cheap_proxies_filename", String, "название файла, из которого брали дешёвые прокси", func() {
		Meta("struct:tag:json", "cheap_proxies_filename")
	})
	Attribute("targets_filename", String, "название файла, из которого брали целевых пользователей", func() {
		Meta("struct:tag:json", "targets_filename")
	})

	Required("bots_filename", "residential_proxies_filename", "cheap_proxies_filename", "targets_filename")
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
