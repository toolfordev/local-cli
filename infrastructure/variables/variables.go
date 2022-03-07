package variables

import (
	"github.com/go-resty/resty/v2"
	encryptedModels "github.com/toolfordev/local-api-encrypted-variables/models"
	globalModels "github.com/toolfordev/local-api-global-variables/models"
)

func GetGlobal() (variables []globalModels.GlobalVariable, err error) {
	client := resty.New()
	_, err = client.
		R().
		SetResult(variables).
		SetHeader("Accept", "application/json").
		Get("http://localhost:14001/global-variables")
	return
}

func GetEncrypted() (variables []encryptedModels.EncryptedVariable, err error) {
	client := resty.New()
	_, err = client.
		R().
		SetResult(variables).
		SetHeader("Accept", "application/json").
		Get("http://localhost:14002/encrypted-variables")
	return
}

func UnlockEncrypted(password string) (err error) {
	client := resty.New()
	_, err = client.
		R().
		SetBody(encryptedModels.Password{Value: password}).
		SetHeader("Accept", "application/json").
		Patch("http://localhost:14002/password/unlock")
	return
}

func LockEncrypted() (err error) {
	client := resty.New()
	_, err = client.
		R().
		SetHeader("Accept", "application/json").
		Patch("http://localhost:14002/password/lock")
	return
}

func SetPasswordEncrypted(password string) (err error) {
	client := resty.New()
	_, err = client.
		R().
		SetBody(encryptedModels.Password{Value: password}).
		SetHeader("Accept", "application/json").
		Patch("http://localhost:14002/password/set")
	return
}
