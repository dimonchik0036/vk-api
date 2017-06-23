package vkapi

type Dialog struct {
	Unread     int64    `json:"unread"`
	Message    *Message `json:"message"`
	InRead     int64    `json:"in_read"`
	OutRead    int64    `json:"out_read"`
	RealOffset int64    `json:"real_offset"`
}

type Message struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	FromId    int64  `json:"from_id"`
	Date      int64  `json:"date"`
	ReadState int    `json:"read_state"`
	Out       int    `json:"out"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	/*Geo       *Geo {
		type (string) — тип места;
		coordinates (string) — координаты места;
		place (object) — описание места (если оно добавлено), объект с полями:
		id (integer) — идентификатор места (если назначено);
		title (string) — название места (если назначено);
		latitude (number) — географическая широта;
		longitude (number) — географическая долгота;
		created (integer) — дата создания (если назначено);
		icon (string) — URL изображения-иконки;
		country (string) — название страны;
		city (string) — название города;
	} `json:"geo"`*/

	/*Attachments *[]Attachments `json:"attachments"`*/
	FwdMessages *[]Message `json:"fwd_messages"`
	Emoji       int        `json:"emoji"`
	Important   int        `json:"important"`
	Deleted     int        `json:"deleted"`
	RandomId    int64      `json:"random_id"`

	ChatId     int64   `json:"chat_id"`
	ChatActive []int64 `json:"chat_active"`
	/*PushSettings *PushSettings { настройки уведомлений для беседы, если они есть.	} `json:"push_settings"`*/
	UsersCount int    `json:"users_count"`
	AdminId    int64  `json:"admin_id"`
	Action     string `json:"action"`
	/*string	тип действия (если это служебное сообщение). Возможные значения:

	  chat_photo_update — обновлена фотография беседы;
	  chat_photo_remove — удалена фотография беседы;
	  chat_create — создана беседа;
	  chat_title_update — обновлено название беседы;
	  chat_invite_user — приглашен пользователь;
	  chat_kick_user — исключен пользователь.*/

	ActionMid   int64  `json:"action_mid"`   /*идентификатор пользователя (если > 0) или email (если < 0), которого пригласили или исключили (для служебных сообщений с action = chat_invite_user или chat_kick_user). */
	ActionEmail string `json:"action_email"` /*email, который пригласили или исключили (для служебных сообщений с action = chat_invite_user или chat_kick_user и отрицательным action_mid). */
	ActionText  string `json:"action_text"`  /*название беседы (для служебных сообщений с action = chat_create или chat_title_update). */
	Photo50     string `json:"photo_50"`
	Photo100    string `json:"photo_100"`
	Photo200    string `json:"photo_200"`
}
