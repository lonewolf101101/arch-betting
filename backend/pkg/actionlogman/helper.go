package actionlogman

type Filter struct {
	Actions     []string
	CustomerIDs []int
	RefIDs      []int
	CustomerID  int
	RefID       int
	Action      string
}

type AdminActionFilter struct {
	Actions []string
}
