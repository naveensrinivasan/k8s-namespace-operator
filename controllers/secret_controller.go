/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	_ "strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	b64 "encoding/base64"

	serverv1alpha1 "github.com/naveensrinivasan/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// SecretReconciler reconciles a Secret object
type SecretReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

// +kubebuilder:rbac:groups=server.naveensrinivasan.dev,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=server.naveensrinivasan.dev,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secret,verbs=list;watch;get;patch;create
// +kubebuilder:rbac:groups=apps,resources=namespace,verbs=list;watch;get;update

func (r *SecretReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("secret", req.NamespacedName)
	const DockerPullSecret = "dockersecret"

	var sec serverv1alpha1.Secret
	if err := r.Get(ctx, req.NamespacedName, &sec); err != nil {
		if apierrors.IsNotFound(err) {

			// This is the case where the namespace is just created and watch event kicks in.
			if req.Namespace == "" {
				o := metav1.ObjectMeta{Namespace: req.Name, Name: DockerPullSecret}

				r.Create(ctx, &serverv1alpha1.Secret{
					ObjectMeta: o,
				})
				log.Info("scheduled for namespace", "namespace", req.Name)
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, nil
		}
		log.Info("error", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	secret, err := r.desiredSecret(sec)
	if err != nil {
		return ctrl.Result{}, err
	}
	foundSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: sec.Name, Namespace: sec.Namespace}, foundSecret)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info("Creating Secret", "namespace", sec.Namespace, "name", sec.Name)
		err = r.Create(ctx, &secret)
		if err != nil {
			log.Error(err, "unable to create Secret")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{Requeue: true}, nil
		}
	} else if err != nil {
		log.Error(err, "error getting Secret")
		return ctrl.Result{}, err
	}

	log.Info("reconciled Secret", "namespace", req.Namespace)
	return ctrl.Result{}, nil
}
func (r *SecretReconciler) desiredSecret(s serverv1alpha1.Secret) (corev1.Secret, error) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Namespace: s.Namespace, Name: s.Name},
		Data: map[string][]byte{
			"username": []byte(b64.StdEncoding.EncodeToString([]byte("foo"))),
			"password": []byte(b64.StdEncoding.EncodeToString([]byte("bar")))},
		StringData: nil,
	}

	ctrl.Log.Info("Secret", "s", s, "secret", secret, "scheme", r.Scheme)
	if err := controllerutil.SetControllerReference(&s, &secret, r.Scheme); err != nil {
		return secret, err
	}
	return secret, nil
}
func (r *SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&serverv1alpha1.Secret{}).
		Owns(&corev1.Secret{}).
		Watches(&source.Kind{Type: &corev1.Namespace{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(predicate.Funcs{CreateFunc: func(event event.CreateEvent) bool {
			// Look for only namespace with Kubeflow
			_, ok := event.Object.(*serverv1alpha1.Secret)
			if ok {
				return true
			}
			// Only for namespaces that have specific labels
			ns, e := event.Object.(*corev1.Namespace)
			if e {
				for k, _ := range ns.Labels {
					if k == "app.kubernetes.io/part-of" {
						return true
					}
				}
			}
			return false
		}}).Complete(r)
}
