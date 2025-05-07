package todo

const (
	NotStarted string = "not started"
	Started    string = "started"
	Completed  string = "completed"
)

type Todo struct {
	ID     int    `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
	UserID int    `json:"-"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}
