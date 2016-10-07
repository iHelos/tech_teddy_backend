package user

type User struct  {
	Login string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Bears []interface{}
	Audio []interface{}
	FCMToken string `json:"fcm"`
}

type NewUser struct  {
	User
	Password1 string `json:"password1" form:"password1"`
	Password2 string `json:"password2" form:"password2"`
}

func (user NewUser) check() error{
	var usererror = NewUserError()
	if ok := IsLogin(user.Login); !ok {
		usererror.Append("login", 0)
	}
	if ok := IsEmail(user.Email); !ok {
		usererror.Append("email", 0)
	}
	if ok := IsPassword(user.Password1); !ok {
		usererror.Append("password", 0)
	}
	if user.Password1 != user.Password2{
		usererror.Append("password", 1)
	}
	if len(usererror.Messages) != 0 {return usererror}
	return nil
}
