package controller

import (
	_ "embed"
	"text/template"

	"github.com/rancher/wrangler/v3/pkg/yaml"
)

//go:embed config/pod.yaml
var PodConfig string

//go:embed config/audit.tmpl.yaml
var AuditConfig string

//go:embed config/rke2.yaml
var Rk2LogConfig string

var AuditConfigTemplate = template.Must(template.New("kubeauditlogs").Parse(AuditConfig))

type AuditConfigParams struct {
	AuditLogPath string
}

func yamlMarshalString(in string) (out any, err error) {
	err = yaml.Unmarshal([]byte(in), &out)
	return
}
