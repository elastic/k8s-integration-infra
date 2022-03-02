package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	replicas := flag.Int("deployments", 2, "number of deployments replicas per namespace")
	namespaces := flag.Int("namespaces", 1, "number of namespaces")
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(*replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	nsName := "demo"
	for i := 1; i <= *namespaces; i++ {
		Nname := fmt.Sprintf("%s-%d", nsName, i)
		Dname := fmt.Sprintf("demo-deployment-%d", i)
		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: Nname}}
		nsClient := client.CoreV1().Namespaces()

		_, getErr := nsClient.Get(context.TODO(), Nname, metav1.GetOptions{})
		if getErr != nil {
			fmt.Printf("Namespace %s does not exist\n", Nname)
			_, err = nsClient.Create(context.TODO(), nsSpec, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Println("Namespace created")
		}
		fmt.Printf("Namespace %s exists\n", Nname)
		deploymentsClient := client.AppsV1().Deployments(Nname)

		result, getErr := deploymentsClient.Get(context.TODO(), Dname, metav1.GetOptions{})
		if getErr != nil {
			fmt.Printf("Deployment %s does not exist", Dname)
			// Create Deployment
			fmt.Println("Creating deployment...")
			deployment.Name = Dname
			result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
		} else {
			fmt.Printf("Deployment %s exists\n", Dname)
			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				// Retrieve the latest version of Deployment before attempting update
				// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

				result.Spec.Replicas = int32Ptr(int32(*replicas)) // update replica count
				_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
				return updateErr
			})
			if retryErr != nil {
				panic(fmt.Errorf("Update failed: %v", retryErr))
			}
			fmt.Printf("Updated Deployment %s\n", Dname)
		}
	}

}

func int32Ptr(i int32) *int32 { return &i }
