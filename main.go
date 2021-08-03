package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	servicebinding "github.com/redhat-developer/service-binding-operator/apis/binding/v1alpha1"

	sbBuilder "github.com/redhat-developer/service-binding-operator/pkg/reconcile/pipeline/builder"
	sbContext "github.com/redhat-developer/service-binding-operator/pkg/reconcile/pipeline/context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	ctrlr "sigs.k8s.io/controller-runtime"
)

func main() {

	var sbFile = flag.String("sb", "servicebinding.yaml", "path to Service Binding resource (has to be a yaml file)")
	flag.Parse()

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	mgr, err := ctrlr.NewManager(restConfig, ctrlr.Options{
		Scheme:                 runtime.NewScheme(),
		HealthProbeBindAddress: "0",
		MetricsBindAddress:     "0",
	})
	if err != nil {
		panic(err)
	}

	pipeline := sbBuilder.DefaultBuilder.WithContextProvider(
		sbContext.Provider(
			dynamicClient,
			sbContext.ResourceLookup(mgr.GetRESTMapper()),
		),
	).Build()

	data, err := ioutil.ReadFile(*sbFile)
	if err != nil {
		panic(err)
	}
	var obj servicebinding.ServiceBinding
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Service Binding: %s\n", obj.Name)

	repeat, err := pipeline.Process(&obj)
	fmt.Printf("should repeat: %t\n", repeat)
	fmt.Printf("err: %s", err)
}
