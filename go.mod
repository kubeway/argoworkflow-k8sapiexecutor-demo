module github.com/kubeway/argoworkflow-k8sapiexecutor-demo

go 1.15

require (
	github.com/argoproj/argo-workflows/v3 v3.1.9
	github.com/elazarl/goproxy v0.0.0-20201021153353-00ad82a08272 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sys v0.0.0-20210611083646-a4fc73990273 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
	k8s.io/api v0.19.6
	k8s.io/apimachinery v0.19.6
	k8s.io/client-go v0.19.6
)

replace (
	k8s.io/api v0.18.8 => k8s.io/api v0.0.0-20200831051839-f197499901bd
	k8s.io/apimachinery v0.18.8 => k8s.io/apimachinery v0.16.16-rc.0
	k8s.io/client-go v0.18.8 => k8s.io/client-go v0.16.4
)
