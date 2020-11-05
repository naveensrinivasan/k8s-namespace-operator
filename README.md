# README

This is code is a `k8s` `controller` which will `watch` a `namespace` with a specific label for example `app.kubernetes.io/part-of` (a label for `kubeflow`) and create a `secret`.

This was built using `kubebuilder`.

## controller output

```
GOROOT=/usr/local/Cellar/go/1.15.2/libexec #gosetup
GOPATH=/Users/naveen/go #gosetup
/usr/local/Cellar/go/1.15.2/libexec/bin/go build -o /private/var/folders/rg/k3q045qd7_55mgdqwgqwcnt80000gn/T/___go_build_github_com_naveensrinivasan github.com/naveensrinivasan #gosetup
/private/var/folders/rg/k3q045qd7_55mgdqwgqwcnt80000gn/T/___go_build_github_com_naveensrinivasan
2020-11-05T18:31:22.869-0500    INFO    controller-runtime.metrics      metrics server is starting to listen    {"addr": ":8080"}
2020-11-05T18:31:22.869-0500    INFO    setup   starting manager
2020-11-05T18:31:22.869-0500    INFO    controller-runtime.manager      starting metrics server {"path": "/metrics"}
2020-11-05T18:31:22.870-0500    INFO    controller-runtime.controller   Starting EventSource    {"controller": "secret", "source": "kind source: /, Kind="}
2020-11-05T18:31:22.970-0500    INFO    controller-runtime.controller   Starting EventSource    {"controller": "secret", "source": "kind source: /, Kind="}
2020-11-05T18:31:23.071-0500    INFO    controller-runtime.controller   Starting EventSource    {"controller": "secret", "source": "kind source: /, Kind="}
2020-11-05T18:31:23.175-0500    INFO    controller-runtime.controller   Starting Controller     {"controller": "secret"}
2020-11-05T18:31:23.175-0500    INFO    controller-runtime.controller   Starting workers        {"controller": "secret", "worker count": 1}
2020-11-05T18:31:23.175-0500    INFO    Secret  {"s": {"kind":"Secret","apiVersion":"server.naveensrinivasan.dev/v1alpha1","metadata":{"name":"dockersecret","namespace":"naveen","selfLink":"/apis/server.naveensrinivasan.dev/v1alpha1/namespaces/naveen/secrets/dockersecret","uid":"1ca2b17a-e904-4981-b6d5-9202b8a06266","resourceVersion":"62464","generation":1,"creationTimestamp":"2020-11-05T23:01:26Z","managedFields":[{"manager":"___go_build_github_com_naveensrinivasan","operation":"Update","apiVersion":"server.naveensrinivasan.dev/v1alpha1","time":"2020-11-05T23:01:26Z","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{".":{},"f:name":{},"f:password":{},"f:userName":{}},"f:status":{}}}]},"spec":{"name":"","userName":"","password":""},"status":{}}, "secret": {"metadata":{"name":"dockersecret","namespace":"naveen","creationTimestamp":null},"data":{"password":"WW1GeQ==","username":"Wm05dg=="}}, "scheme": {}}
2020-11-05T18:31:23.176-0500    INFO    controllers.Secret      reconciled Secret       {"secret": "naveen/dockersecret", "namespace": "naveen"}
2020-11-05T18:31:23.176-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "secret", "request": "naveen/dockersecret"}
2020-11-05T18:31:23.176-0500    INFO    Secret  {"s": {"kind":"Secret","apiVersion":"server.naveensrinivasan.dev/v1alpha1","metadata":{"name":"dockersecret","namespace":"test1","selfLink":"/apis/server.naveensrinivasan.dev/v1alpha1/namespaces/test1/secrets/dockersecret","uid":"f839f81b-8856-4a99-9721-199f84452317","resourceVersion":"63806","generation":1,"creationTimestamp":"2020-11-05T23:21:32Z","managedFields":[{"manager":"___go_build_github_com_naveensrinivasan","operation":"Update","apiVersion":"server.naveensrinivasan.dev/v1alpha1","time":"2020-11-05T23:21:32Z","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{".":{},"f:name":{},"f:password":{},"f:userName":{}},"f:status":{}}}]},"spec":{"name":"","userName":"","password":""},"status":{}}, "secret": {"metadata":{"name":"dockersecret","namespace":"test1","creationTimestamp":null},"data":{"password":"WW1GeQ==","username":"Wm05dg=="}}, "scheme": {}}
2020-11-05T18:31:23.176-0500    INFO    controllers.Secret      reconciled Secret       {"secret": "test1/dockersecret", "namespace": "test1"}
2020-11-05T18:31:23.176-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "secret", "request": "test1/dockersecret"}
2020-11-05T18:31:23.185-0500    INFO    controllers.Secret      scheduled for namespace {"secret": "/test1", "namespace": "test1"}
2020-11-05T18:31:23.185-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "secret", "request": "/test1"}
2020-11-05T18:31:23.194-0500    INFO    controllers.Secret      scheduled for namespace {"secret": "/naveen", "namespace": "naveen"}
2020-11-05T18:31:23.194-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "secret", "request": "/naveen"}



```
