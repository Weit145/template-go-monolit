package v1

import _ "embed"

// Spec contains the OpenAPI document served by the Swagger UI.
//
//go:embed swagger.yaml
var Spec []byte
