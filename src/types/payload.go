package types

//Repository -
type Repository struct {
	Name   string `json:"name"`
	SSHURL string `json:"ssh_url"`
}

//Payload -
type Payload struct {
	Repository Repository `json:"repository"`
}
