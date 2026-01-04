package swagger

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/mirkobrombin/go-module-router/v2/pkg/router"
)

type Document struct {
	OpenAPI string              `json:"openapi"`
	Info    Info                `json:"info"`
	Paths   map[string]PathItem `json:"paths"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type PathItem map[string]Operation

type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

type Parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Schema   Schema `json:"schema"`
}

type Schema struct {
	Type    string   `json:"type,omitempty"`
	Minimum *float64 `json:"minimum,omitempty"`
}

type Response struct {
	Description string `json:"description"`
}

// MetaProvider is an optional interface endpoints can implement for OpenAPI metadata.
type MetaProvider interface {
	OpenAPIMeta() map[string]any
}

// Build generates OpenAPI document from registered declarative endpoints.
func Build(title, version string, endpoints ...router.Handler) ([]byte, error) {
	doc := Document{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:   title,
			Version: version,
		},
		Paths: map[string]PathItem{},
	}

	for _, ep := range endpoints {
		val := reflect.ValueOf(ep)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		typ := val.Type()

		// Extract method and path from Meta field
		var method, path string
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.Type == reflect.TypeOf(router.Meta{}) {
				method = strings.ToLower(field.Tag.Get("method"))
				path = field.Tag.Get("path")
				break
			}
		}

		if method == "" || path == "" {
			continue
		}

		pathItem, ok := doc.Paths[path]
		if !ok {
			pathItem = make(PathItem)
		}

		op := Operation{
			Responses: map[string]Response{"200": {Description: "OK"}},
		}

		// If endpoint implements MetaProvider, use its metadata
		if mp, ok := ep.(MetaProvider); ok {
			m := mp.OpenAPIMeta()

			if s, ok := m["summary"].(string); ok {
				op.Summary = s
			}
			if d, ok := m["description"].(string); ok {
				op.Description = d
			}

			if params, ok := m["parameters"].([]map[string]any); ok {
				for _, pm := range params {
					p := Parameter{
						Name:     pm["name"].(string),
						In:       pm["in"].(string),
						Required: false,
					}
					if req, ok := pm["required"].(bool); ok {
						p.Required = req
					}
					if sch, ok := pm["schema"].(map[string]any); ok {
						p.Schema.Type, _ = sch["type"].(string)
						if min, ok := sch["minimum"].(int); ok {
							f := float64(min)
							p.Schema.Minimum = &f
						}
					}
					op.Parameters = append(op.Parameters, p)
				}
			}

			if resp, ok := m["responses"].(map[int]any); ok {
				op.Responses = map[string]Response{}
				for code, desc := range resp {
					op.Responses[string(rune('0'+code/100))+string(rune('0'+(code/10)%10))+string(rune('0'+code%10))] = Response{Description: desc.(string)}
				}
			}
		}

		pathItem[method] = op
		doc.Paths[path] = pathItem
	}

	return json.MarshalIndent(doc, "", "  ")
}
