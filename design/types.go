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

// Dataset описывает рекламную кампанию
var Dataset = Type("Dataset", func() {
	Attribute("id", String, "", func() {
		Format(FormatUUID)
	})

	Attribute("text_template", String, func() {
		Meta("struct:tag:json", "text_template")
		Description("описание под постом")
	})

	Attribute("post_images", ArrayOf(String), "список base64 строк картинок", func() {
		Meta("struct:tag:json", "post_images")
	})

	Attribute("landing_accounts", ArrayOf(String), func() {
		Description("имена аккаунтов, на которых ведем трафик")
		Meta("struct:tag:json", "landing_accounts")
	})

	Attribute("bot_names", ArrayOf(String), func() {
		Description("имена для аккаунтов-ботов")
		Meta("struct:tag:json", "bot_names")
	})

	Attribute("bot_last_names", ArrayOf(String), func() {
		Description("фамилии для аккаунтов-ботов")
		Meta("struct:tag:json", "bot_last_names")
	})

	Attribute("bot_images", ArrayOf(String), func() {
		Description("аватарки для ботов")
		Meta("struct:tag:json", "bot_images")
	})

	Attribute("bot_urls", ArrayOf(String), func() {
		Description("ссылки для описания у ботов")
		Meta("struct:tag:json", "bot_urls")
	})
	Attribute("status", DatasetStatus)
	Attribute("title", String, "название задачи")

	Attribute("bots_num", Int, "количество ботов в задаче", func() {
		Meta("struct:tag:json", "bots_num")
	})
	Attribute("residential_proxies_num", Int, "количество резидентских прокси в задаче", func() {
		Meta("struct:tag:json", "residential_proxies_num")
	})
	Attribute("cheap_proxies_num", Int, "количество дешёвых прокси в задаче", func() {
		Meta("struct:tag:json", "cheap_proxies_num")
	})

	Attribute("targets_num", Int, "количество целевых пользователей в задаче", func() {
		Meta("struct:tag:json", "targets_num")
	})

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

	Attribute("follow_targets", Boolean, func() {
		Description("нужно ли подписываться на аккаунты")
		Meta("struct:tag:json", "follow_targets")
	})

	Attribute("need_photo_tags", Boolean, func() {
		Description("делать отметки на фотографии")
		Meta("struct:tag:json", "need_photo_tags")
	})

	Attribute("per_post_sleep_seconds", UInt, func() {
		Description("делать отметки на фотографии")
		Meta("struct:tag:json", "per_post_sleep_seconds")
	})

	Attribute("photo_tags_delay_seconds", UInt, func() {
		Description("задержка перед проставлением отметок")
		Meta("struct:tag:json", "photo_tags_delay_seconds")
	})

	Attribute("posts_per_bot", UInt, func() {
		Description("количество постов для каждого бота")
		Meta("struct:tag:json", "posts_per_bot")
	})

	Attribute("targets_per_post", UInt, func() {
		Description("количество упоминаний под каждым постом")
		Meta("struct:tag:json", "targets_per_post")
	})

	Required("id", "text_template", "post_images", "status", "title", "bots_num", "residential_proxies_num",
		"cheap_proxies_num", "targets_num", "bot_images", "landing_accounts", "bot_names", "bot_last_names", "bot_urls",
		"targets_per_post", "posts_per_bot", "photo_tags_delay_seconds", "per_post_sleep_seconds", "need_photo_tags", "follow_targets",
	)
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

var Blogger = Type("Blogger", func() {
	Attribute("login", String, func() {
		Description("логин блогера")
	})

	Attribute("user_id", String, func() {
		Description("user_id блогера")
	})

	Required("login", "user_id")
})
