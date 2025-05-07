package def

type SwaggerSecurity interface {
	SecuritCookie(name string, cookieName string)
	SecuritApiHeader(securityTag string, headerName string)
	SecuritJwt(name string)
}

type SwaggerOps interface {
	SwaggerSecurity
	SetSummary(title string)
	SetTag(tag string)
	SetDescription(description string)
	SetParameterDescription(name, description string)
}
