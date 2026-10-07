package domain_health

type Health struct {
	postgres bool
}

func NewHealth(postgres bool) Health {
	return Health{
		postgres: postgres,
	}
}

func (h *Health) IsPostgres() bool { return h.postgres }
