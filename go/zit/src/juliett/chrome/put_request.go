package chrome

type putRequestItem struct{}

type putRequest struct {
	Deleted []string `json:"deleted"`
}
