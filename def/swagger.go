package def

type SwaggerOps interface {
	SetSummary(title string)
	SetDescription(description string)
}

type swaggerImpl struct {
	mi *MethodInfo
}

func (s *swaggerImpl) SetSummary(title string) {
	s.mi.KV.Store("swagger.summary", title)
}

func (s *swaggerImpl) SetDescription(description string) {
	s.mi.KV.Store("swagger.description", description)
}
