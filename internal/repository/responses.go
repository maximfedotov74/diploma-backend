package repository

type UserRepoResponse struct {
	Id                    int    `json:"id"`
	ActivationAccountLink string `json:"activation_account_link"`
}
