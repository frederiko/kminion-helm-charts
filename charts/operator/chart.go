// +gotohelm:ignore=true
package operator

import (
	_ "embed"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/redpanda-data/helm-charts/pkg/gotohelm/helmette"
	"github.com/redpanda-data/helm-charts/pkg/helm"
	"github.com/redpanda-data/helm-charts/pkg/kube"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

var (
	// Scheme is a [runtime.Scheme] with the appropriate extensions to load all
	// objects produced by the redpanda chart.
	Scheme = runtime.NewScheme()

	//go:embed Chart.yaml
	chartYAML []byte

	//go:embed values.yaml
	defaultValuesYAML []byte

	chartMeta helmette.Chart
)

func init() {
	must(scheme.AddToScheme(Scheme))
	must(certmanagerv1.AddToScheme(Scheme))
	must(monitoringv1.AddToScheme(Scheme))

	// NB: We can't directly unmarshal into a helmette.Chart as adding json
	// tags to it breaks gotohelm.
	var chart map[string]any
	must(yaml.Unmarshal(chartYAML, &chart))

	chartMeta = helmette.Chart{
		Name:       chart["name"].(string),
		Version:    chart["version"].(string),
		AppVersion: chart["appVersion"].(string),
	}
}

// ChartMeta returns a parsed version of redpanda's Chart.yaml.
func ChartMeta() helmette.Chart {
	return chartMeta
}

func Dot(release helmette.Release, values PartialValues, kubeConfig kube.Config) (*helmette.Dot, error) {
	valuesYaml, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	// NB: err1 is working around an issue in gotohelm's ASTs rewrites
	merged, err1 := helm.MergeYAMLValues("", defaultValuesYAML, valuesYaml)
	if err1 != nil {
		return nil, err1
	}

	return &helmette.Dot{
		Values:     merged,
		Chart:      ChartMeta(),
		Release:    release,
		KubeConfig: kubeConfig,
	}, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
