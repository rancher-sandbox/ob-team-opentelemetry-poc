module github.com/rancher-sandbox/ob-team-opentelemetry-poc

go 1.24.0

toolchain go1.24.1

require (
	github.com/go-viper/mapstructure/v2 v2.2.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.125.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage v0.125.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.125.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.125.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.125.0
	github.com/rancher/lasso v0.2.1
	github.com/rancher/wrangler/v3 v3.2.0
	github.com/samber/lo v1.50.0
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/collector/component v1.31.0
	go.opentelemetry.io/collector/config/configgrpc v0.125.0
	go.opentelemetry.io/collector/config/confignet v1.31.0
	go.opentelemetry.io/collector/config/configtelemetry v0.125.0
	go.opentelemetry.io/collector/exporter/debugexporter v0.125.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.125.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.33.0
	k8s.io/apimachinery v0.33.0
	k8s.io/client-go v0.33.0
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/elastic/lunes v0.1.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/evanphx/json-patch v5.9.11+incompatible // indirect
	github.com/expr-lang/expr v1.17.2 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/providers/confmap v1.0.0 // indirect
	github.com/knadh/koanf/v2 v2.2.0 // indirect
	github.com/leodido/go-syslog/v4 v4.2.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20190525184631-5f46317e436b // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.125.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.etcd.io/bbolt v1.4.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector v0.125.0 // indirect
	go.opentelemetry.io/collector/client v1.31.0 // indirect
	go.opentelemetry.io/collector/component/componentstatus v0.125.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.125.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.31.0 // indirect
	go.opentelemetry.io/collector/config/configmiddleware v0.125.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.31.0 // indirect
	go.opentelemetry.io/collector/config/configretry v1.31.0 // indirect
	go.opentelemetry.io/collector/config/configtls v1.31.0 // indirect
	go.opentelemetry.io/collector/confmap v1.31.0 // indirect
	go.opentelemetry.io/collector/consumer v1.31.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.125.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror/xconsumererror v0.125.0 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.125.0 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.125.0 // indirect
	go.opentelemetry.io/collector/exporter v0.125.0 // indirect
	go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper v0.125.0 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.125.0 // indirect
	go.opentelemetry.io/collector/extension v1.31.0 // indirect
	go.opentelemetry.io/collector/extension/extensionauth v1.31.0 // indirect
	go.opentelemetry.io/collector/extension/extensionmiddleware v0.125.0 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.125.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.31.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.125.0 // indirect
	go.opentelemetry.io/collector/pdata v1.31.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.125.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.125.0 // indirect
	go.opentelemetry.io/collector/pipeline/xpipeline v0.125.0 // indirect
	go.opentelemetry.io/collector/receiver v1.31.0 // indirect
	go.opentelemetry.io/collector/receiver/receiverhelper v0.125.0 // indirect
	go.opentelemetry.io/collector/semconv v0.125.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.10.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/log v0.11.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.31.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	gonum.org/v1/gonum v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.72.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.32.1 // indirect
	k8s.io/code-generator v0.32.1 // indirect
	k8s.io/gengo v0.0.0-20250130153323-76c5745d3511 // indirect
	k8s.io/gengo/v2 v2.0.0-20240911193312-2b36238f13e9 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250318190949-c8a335a9a2ff // indirect
	k8s.io/utils v0.0.0-20241104100929-3ea5e8cea738 // indirect
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
)
