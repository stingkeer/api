package def

type SwaggerOps interface {
	SetSummary(title string)
	SetDescription(description string)
	SetParameterDescription(name, description string)
}
