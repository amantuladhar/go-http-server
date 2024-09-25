package zhttp

type Request struct {
	Method      string
	Path        string
	HttpVersion string
	PathParam   map[string]string
}
