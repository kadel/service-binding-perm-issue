A small example program to verify https://github.com/redhat-developer/service-binding-operator/issues/1003 



on OpenShift

1) create Deployments for dummy App and Service
```
$ oc apply -f -f deployment-app.yaml -f deployment-service.yaml
```

2) use SBO's Pipeline.Process to create binding
```
$ go run main.go 

Service Binding: test
should repeat: true
err: customresourcedefinitions.apiextensions.k8s.io "deployments.apps" is forbidden: User "developer" cannot get resource "customresourcedefinitions" in API group "apiextensions.k8s.io" at the cluster scope% 
``` 


