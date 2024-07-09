package httpcaller

type CallerOptions struct {
	DefaultHeaders      map[string]string
	BaseSuccessResponse map[string]interface{}
}

type CallOption struct {
	Header    map[string]string
	PathParam map[string]string
}
