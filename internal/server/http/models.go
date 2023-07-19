package internalhttp

type UserRegisterBody struct {
	FirstName  string  `json:"first_name,omitempty"`
	SecondName *string `json:"second_name,omitempty"`
	Birthdate  *string `json:"birthdate,omitempty"`
	Biography  *string `json:"biography,omitempty"`
	City       *string `json:"city,omitempty"`
	Password   string  `json:"password,omitempty"`
}

type User struct {
	Id string `json:"id,omitempty"`
	// Имя
	FirstName string `json:"first_name,omitempty"`
	// Фамилия
	SecondName *string `json:"second_name,omitempty"`
	// Возраст
	Age       *int    `json:"age,omitempty"`
	Birthdate *string `json:"birthdate,omitempty"`
	// Интересы
	Biography *string `json:"biography,omitempty"`
	// Город
	City *string `json:"city,omitempty"`
}

type InlineResponse500 struct {
	// Описание ошибки
	Message string `json:"message"`
	// Идентификатор запроса. Предназначен для более быстрого поиска проблем.
	RequestId string `json:"request_id,omitempty"`
	// Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
	Code int32 `json:"code,omitempty"`
}

type LoginBody struct {
	Id       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserIdSendBody struct {
	Text string `json:"text"`
}

type PostUpdateBody struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type PostCreateBody struct {
	Text string `json:"text"`
}

type UserLoginResponse struct {
	Token string `json:"token,omitempty"`
}

type Post struct {
	Id           string `json:"id,omitempty"`
	Text         string `json:"text,omitempty"`
	AuthorUserId string `json:"author_user_id,omitempty"`
}

type UserCreateResponse struct {
	UserId string `json:"user_id,omitempty"`
}

type DialogMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}
