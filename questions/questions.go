package questions

import (
	"errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"regexp"
)

var email = survey.Question{
	Name: "email",
	Prompt: &survey.Input{
		Message: "Enter Email :",
		Help:    "Please provide your email for registration",
	},
	Validate: func(val interface{}) error {
		paaswdReg := regexp.MustCompile(`^\w+@(gmail|yahoo)\.[a-zA-Z]{2,3}$`)
		if str, ok := val.(string); !ok || !paaswdReg.Match([]byte(str)) {
			return errors.New("Please enter valid email")
		}
		return nil
	},
}

var random = survey.Question{
	Name: "random",
	Prompt: &survey.Input{
		Message: "Please enter random string from your email:",
	},
	Validate: survey.Required,
}

var password = survey.Question{
	Name: "password",
	Prompt: &survey.Password{
		Message: "Enter Password:",
	},
	Validate: func(val interface{}) error {
		paaswdReg := regexp.MustCompile(`([a-zA-Z\d!@#$%^&*]+)`)
		if str, ok := val.(string); !ok || !paaswdReg.Match([]byte(str)) || len(str) < 7 {
			return errors.New("Password must be longer than 6 and must cotaines [1-9a-zA-Z] and [! @ # $ % ^ & *]")
		}
		return nil
	},
}

var currPassword = survey.Question{
	Name: "curr_password",
	Prompt: &survey.Password{
		Message: "Enter Current Password:",
	},
	Validate: func(val interface{}) error {
		paaswdReg := regexp.MustCompile(`([a-zA-Z\d!@#$%^&*]+)`)
		if str, ok := val.(string); !ok || !paaswdReg.Match([]byte(str)) || len(str) < 7 {
			return errors.New("Password must be longer than 6 and must cotaines [1-9a-zA-Z] and [! @ # $ % ^ & *]")
		}
		return nil
	},
}

var newPassword = survey.Question{
	Name: "new_password",
	Prompt: &survey.Password{
		Message: "Enter New Password:",
	},
	Validate: func(val interface{}) error {
		paaswdReg := regexp.MustCompile(`([a-zA-Z\d!@#$%^&*]+)`)
		if str, ok := val.(string); !ok || !paaswdReg.Match([]byte(str)) || len(str) < 7 {
			return errors.New("Password must be longer than 6 and must cotaines [1-9a-zA-Z] and [! @ # $ % ^ & *]")
		}
		return nil
	},
}

var name = survey.Question{
	Name: "name",
	Prompt: &survey.Input{
		Message: "Enter your name :",
		Help:    "Please provide your name for registration",
	},
	Validate: survey.Required,
}

var RegisterUserInit = append([]*survey.Question{}, &email)

var RegisterUserFinish = append([]*survey.Question{}, &random, &name, &password)

var LoginUser = append([]*survey.Question{}, &email, &password)

var ChangePassword = append([]*survey.Question{}, &currPassword, &newPassword)
