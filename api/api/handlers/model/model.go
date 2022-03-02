package model

type DeleteOrGet struct {
	Id string
}

type TaskReq struct {
	Assignee_id string `json:"assignee_id"`
	Title       string `json:"title"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
}

type TaskRes struct {
	Id          string `json:"id"`
	Assignee_id string `json:"assignee_id"`
	Title       string `json:"title"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateTask struct {
	Id          string `json:"id"`
	Assignee_id string `json:"assignee_id"`
	Title       string `json:"title"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
}

type Update1 struct {
	Id       string
	Assignee string
	Title    string
	Deadline string
	Status   string
}

type Tasks struct {
	Task []TaskRes
}

type Mess struct {
	Message string
}
type ListTasks struct {
	Page  int64
	Limit int64
}

type NewUser struct {
	Id            string
	First_name    string
	Last_name     string
	Email         string
	Profile_photo string
	Location      []Location
	Phone         []Phone
}

type Phone struct {
	Phone_number string
}
type Location struct {
	Location string
}

type UserRes struct {
	Id            string
	First_name    string
	Last_name     string
	Email         string
	Profile_photo string
	Location      []Location
	Phone         []Phone
}
