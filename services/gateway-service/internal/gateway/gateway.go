package gateway

type Route struct {
	Path         string `json:"path"`
	Method       string `json:"method"`
	Backend      string `json:"backend"`
	RateLimitRPS int    `json:"rate_limit_rps"`
}

type Registry struct{ routes []Route }

func NewRegistry() *Registry { return &Registry{} }

func (reg *Registry) Register(r Route) {
	for _, e := range reg.routes {
		if e.Path == r.Path && e.Method == r.Method { return }
	}
	reg.routes = append(reg.routes, r)
}

func (reg *Registry) List() []Route { return reg.routes }
func (reg *Registry) Len() int      { return len(reg.routes) }
