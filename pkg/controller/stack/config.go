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

//go:embed config/rke2journald.tmpl.yaml
var Rke2JournalConfig string

//go:embed config/k3sjournald.tmpl.yaml
var K3sJournalConfig string

var AuditConfigTemplate = template.Must(template.New("kubeauditlogs").Parse(AuditConfig))
var K3sJournalConfigTemplate = template.Must(template.New("k3sjournald").Parse(K3sJournalConfig))
var RK2JournalConfigTemplate = template.Must(template.New("rke2journald").Parse(Rke2JournalConfig))

type AuditConfigParams struct {
	AuditLogPath string
	StorageId    string
}

type K3sJournalParams struct {
	JournalPath string
	StorageId   string
}

type RKE2JournalParams struct {
	JournalPath string
	StorageId   string
}

func yamlMarshalString(in string) (out any, err error) {
	err = yaml.Unmarshal([]byte(in), &out)
	return
}
