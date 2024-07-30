package status

type Status string

const (
	Ok       Status = "ok"
	Error    Status = "error"
	Checking Status = "checking"
)
