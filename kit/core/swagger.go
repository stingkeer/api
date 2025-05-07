package core

import (
	"fmt"

	"gitee.com/fast_api/api/def"
)

var (
	_ def.SwaggerOps = (*swaggerImpl)(nil)
)

type swaggerImpl struct {
	mi *def.MethodInfo
	SwaggerSecurit
}

type SwaggerSecurit struct {
	Ops []def.Option
}

// SecuritApiHeader implements def.SwaggerOps.
func (s *SwaggerSecurit) SecuritApiHeader(securityTag string, headerName string) {
	for i := 0; i < len(s.Ops); i++ {
		s.Ops[i].StoreKV("swagger.securit", SecuritApiHeader(securityTag, headerName))
	}

}

// SecuritCookie implements def.SwaggerOps.
func (s *SwaggerSecurit) SecuritCookie(securityTag string, cookieName string) {
	for i := 0; i < len(s.Ops); i++ {
		s.Ops[i].StoreKV("swagger.securit", SecuritCookie(securityTag, cookieName))
	}
}

// SecuritJwt implements def.SwaggerOps.
func (s *SwaggerSecurit) SecuritJwt(securityTag string) {
	for i := 0; i < len(s.Ops); i++ {
		s.Ops[i].StoreKV("swagger.securit", SecuritJwt(securityTag))
	}
}

// SetParameterDescription implements def.SwaggerOps.
func (s *swaggerImpl) SetParameterDescription(name string, description string) {
	s.mi.KV.Store(fmt.Sprintf("swagger.parameter.%s", name), description)
}

func (s *swaggerImpl) SetSummary(title string) {
	s.mi.KV.Store("swagger.summary", title)
}

func (s *swaggerImpl) SetTag(tag string) {
	s.mi.KV.Store("swagger.tag", tag)
}

func (s *swaggerImpl) SetDescription(description string) {
	s.mi.KV.Store("swagger.description", description)
}

// Security Scheme Object
type SecurityObject struct {
	Typ          string `json:"type,omitempty"`
	In           string `json:"in,omitempty"`
	Name         string `json:"name,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerformat,omitempty"`
}

func SecuritCookie(securityTag string, cookieName string) map[string]*SecurityObject {
	return map[string]*SecurityObject{
		securityTag: {
			Typ:  "apiKey",
			In:   "cookie",
			Name: cookieName,
		},
	}
}

func SecuritApiHeader(securityTag string, headerName string) map[string]*SecurityObject {
	return map[string]*SecurityObject{
		securityTag: {
			Typ:  "apiKey",
			In:   "header",
			Name: headerName,
		},
	}
}

func SecuritJwt(securityTag string) map[string]*SecurityObject {
	return map[string]*SecurityObject{
		securityTag: {
			Typ:          "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
}
