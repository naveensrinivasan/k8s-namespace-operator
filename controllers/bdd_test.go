package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	timeout  = time.Second * 20
	interval = time.Millisecond * 250
)

var _ = Describe("Notebook controller valid scenario", func() {
	const MyPodNamespace = "default"
	newNamespace := ""

	Context("When validating the Namespace controller", func() {
		It("The downward API of the namespace should have been set to ENV variable MY_POD_NAMESPACE", func() {
			// Set the namespace where the secret would be deployed.
			err := os.Setenv("MY_POD_NAMESPACE", MyPodNamespace)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("Should create namespace", func() {
			By("By creating a new namespace")
			newNamespace = fmt.Sprintf("test-%s", randStringRunes(6))
			ctx := context.Background()
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:      newNamespace,
					Namespace: newNamespace,
					Labels:    namespaceLabels,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())
		})
		It("It should have a secret by TargetSecretName of dockersecret", func() {
			By("Checking for existing of the secret being present")
			sec := corev1.Secret{}
			secretKey := types.NamespacedName{Name: dockersecret, Namespace: MyPodNamespace}
			Eventually(func() bool {
				err := k8sClient.Get(context.TODO(), secretKey, &sec)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})
		It(fmt.Sprintf("It should create a secret called dockerpullsecret in the namespace %s", newNamespace), func() {
			By("Checking for existing of the secret being present")
			sec := corev1.Secret{}
			secretKey := types.NamespacedName{Name: "dockerpullsecret", Namespace: newNamespace}
			Eventually(func() bool {
				err := k8sClient.Get(context.TODO(), secretKey, &sec)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})
	})
})
