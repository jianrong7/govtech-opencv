package dto

type RegisterReq struct {
	Students []string `json:"students"`
	Teacher  string   `json:"teacher"`
}

type SuspendStudentReq struct {
	Student string `json:"student"`
}

type RetrieveForNotificationsReq struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}
