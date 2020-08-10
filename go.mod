module github.com/ori-amizur/introspector

go 1.13

require (
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/jaypipes/ghw v0.6.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/openshift/assisted-service v0.0.0-20200805103259-9ca9af7cddc0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/ssgreg/journald v1.0.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/openshift/assisted-service => github.com/tsorya/assisted-service v0.0.0-20200809200312-07f2712d1a44
