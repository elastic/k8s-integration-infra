package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

var deployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: "demo-deployment",
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(int32(1)),
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
				Annotations: map[string]string{
					"key": "val",
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

func main() {
	// parse flags
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	replicas := flag.Int("deployments", 2, "number of deployments replicas per namespace")
	namespaces := flag.Int("namespaces", 1, "number of namespaces")
	podLabels := flag.Int("podlabels", 1, "number of labels per pod")
	podAnnotations := flag.Int("podannotations", 1, "number of annotations per pod")
	delAll := flag.Bool("delete", false, "delete all")
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
	if *delAll {
		err = deleteNamespaces(client)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
	// Create namespaces and deployment per namespace
	for i := 1; i <= *namespaces; i++ {
		Nname := fmt.Sprintf("demo-%d", i)
		err := createNamespace(Nname, client)
		if err != nil {
			panic(err)
		}

		Dname := fmt.Sprintf("demo-deployment-%d", i)
		err = createOrUpdateDeployment(Dname, Nname, client, *replicas, *podLabels, *podAnnotations)
		if err != nil {
			panic(err.Error())
		}
	}

}

// createNamespace creates a new namespace if it does not exist
func createNamespace(name string, client *kubernetes.Clientset) error {
	nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	nsClient := client.CoreV1().Namespaces()

	_, getErr := nsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if getErr != nil {
		fmt.Printf("Namespace %s does not exist\n", name)
		_, err := nsClient.Create(context.TODO(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Namespace created")
	}
	fmt.Printf("Namespace %s exists\n", name)
	return nil
}

// deleteNamespaces creates a new namespace if it does not exist
func deleteNamespaces(client *kubernetes.Clientset) error {
	//nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	nsClient := client.CoreV1().Namespaces()
	res, getErr := nsClient.List(context.TODO(), metav1.ListOptions{})
	if getErr != nil {
		return getErr
	} else {
		for _, n := range res.Items {
			if strings.Contains(n.Name, "demo-") {
				err := nsClient.Delete(context.TODO(), n.Name, metav1.DeleteOptions{})
				if err != nil {
					return err
				}
				fmt.Printf("Namespace %s deleted\n", n.Name)
			}
		}

	}
	return nil
}

// createOrUpdateDeployment creates a new namespace if it does not exist or updates it if it does
func createOrUpdateDeployment(name, namespace string, client *kubernetes.Clientset, replicas, podLabels, podAnnotations int) error {
	deploymentsClient := client.AppsV1().Deployments(namespace)
	result, getErr := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if getErr != nil {
		fmt.Printf("Deployment %s does not exist\n", name)
		// Create Deployment
		fmt.Println("Creating deployment...")
		deployment.Name = name
		deployment.Spec.Replicas = int32Ptr(int32(replicas))
		for i := 1; i <= podLabels; i++ {
			lkey := fmt.Sprintf("app-%d", i)
			lval := fmt.Sprintf("demo-%d", i)
			deployment.Spec.Template.Labels[lkey] = lval
		}
		for i := 1; i <= podAnnotations; i++ {
			ankey := fmt.Sprintf("key-%d", i)
			anval := fmt.Sprintf("val-%d", i)
			deployment.Spec.Template.Annotations[ankey] = anval
		}
		result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	} else {
		fmt.Printf("Deployment %s exists\n", name)
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

			result.Spec.Replicas = int32Ptr(int32(replicas)) // update replica count
			for i := 1; i <= podLabels; i++ {
				lkey := fmt.Sprintf("app-%d", i)
				lval := fmt.Sprintf("demo-%d", i)
				result.Spec.Template.Labels[lkey] = lval
			}
			for i := 1; i <= podAnnotations; i++ {
				ankey := fmt.Sprintf("key-%d", i)
				anval := fmt.Sprintf("val-%d", i)
				result.Spec.Template.Annotations[ankey] = anval
			}
			_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return fmt.Errorf("Update failed: %v", retryErr)
		}
		fmt.Printf("Updated Deployment %s\n", name)
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
