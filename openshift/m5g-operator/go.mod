module mosaic5g/m5g-operator

go 1.14

require (
	cloud.google.com/go v0.41.0
	contrib.go.opencensus.io/exporter/ocagent v0.4.12
	git.apache.org/thrift.git v0.12.0 // indirect
	github.com/Azure/go-autorest v12.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.1 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/PuerkitoBio/purell v1.1.1
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
	// github.com/appscode/jsonpatch v0.0.0-20190108182946-7c0e3b262f30
	github.com/appscode/jsonpatch v0.0.0-20190108182946-7c0e3b262f30 // indirect
	// gomodules.xyz/jsonpatch v2.0
	github.com/beorn7/perks v1.0.0
	github.com/census-instrumentation/opencensus-proto v0.2.0
	github.com/coreos/prometheus-operator v0.29.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emicklei/go-restful v2.9.6+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1
	github.com/go-openapi/jsonpointer v0.19.2
	github.com/go-openapi/jsonreference v0.19.2
	github.com/go-openapi/spec v0.19.2
	github.com/go-openapi/swag v0.19.4
	github.com/gobuffalo/envy v1.7.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/google/btree v1.0.0
	github.com/google/gofuzz v1.0.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.0
	github.com/gophercloud/gophercloud v0.2.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/grpc-ecosystem/grpc-gateway v1.9.3
	github.com/hashicorp/golang-lru v0.5.1
	github.com/imdario/mergo v0.3.7
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.6
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e
	github.com/markbates/inflect v1.0.4
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/operator-framework/operator-sdk v0.7.0
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/petar/GoLLRB v0.0.0-20190514000832-33fb24c13b99
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.6.0
	github.com/prometheus/procfs v0.0.3
	github.com/rogpeppe/go-internal v1.3.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/pflag v1.0.3
	go.opencensus.io v0.22.0
	go.uber.org/atomic v1.4.0
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20190626221950-04f50cda93cb
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools v0.0.0-20190703212419-2214986f1668
	google.golang.org/api v0.7.0
	google.golang.org/appengine v1.6.1
	google.golang.org/genproto v0.0.0-20190701230453-710ae3a149df
	google.golang.org/grpc v1.22.0
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/gengo v0.0.0-20190327210449-e17681d19d3a
	k8s.io/klog v0.3.3
	k8s.io/kube-openapi v0.0.0-20180711000925-0cf8f7e6ed1d
	sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools v0.1.8
	sigs.k8s.io/testing_frameworks v0.1.2 // indirect
	sigs.k8s.io/yaml v1.1.0
)
