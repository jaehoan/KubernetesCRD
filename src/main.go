package main

import (
	"client"
	"context"
	k8scrdcontroller "controller"
	"fmt"
	"log"
	"os"
	"time"
	v1 "v1"

	"k8s.io/client-go/tools/clientcmd"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	kubeConfigPath := os.Getenv("KUBECONFIG")

	// Use kubeconfig to create client config.
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err)
	}

	apiextensionsclient.NewForConfig(clientConfig)

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	// logger init
	fpLog, err := os.OpenFile("crd.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	// Init a CRD.
	crd, err := v1.CreateCustomResourceDefinition(apiextensionsClientSet)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}
	// Just for cleanup.
	defer func() {
		if crd == nil {
			fmt.Println("Exit and clean ")
		} else {
			fmt.Println("Exit and clean " + crd.Name)
			apiextensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
		}
	}()

	// Make a new config for extension's API group and use the first one as the baseline.
	exampleClient, exampleScheme, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}

	// Start CRD controller.
	controller := k8scrdcontroller.Controller{
		Client: exampleClient,
		Scheme: exampleScheme,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	// Create a CRD client interface.
	crdClient := client.NewCrdClient(exampleClient, exampleScheme, v1.DefaultNamespace)

	// Code For CR Instance
	// Create an instance of CRD.
	instanceName := "example1"

	exampleInstance := &v1.Item{
		ObjectMeta: metav1.ObjectMeta{
			Name: instanceName,
		},
		Attribute: v1.ItemAttribute{
			WelcomeMsg: "Welcome to Kubernetes World",
			SleepTime:  10,
		},
		Status: v1.ItemStatus{
			State:   v1.StateCreated,
			Message: "Created but not processed yet",
		},
	}
	result, err := crdClient.Create(exampleInstance)
	if err == nil {
		fmt.Printf("CREATED: %#v", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v", result)
	} else {
		panic(err)
	}

	// Wait until the CRD object is handled by controller and its status is changed to Processed.
	err = client.WaitForInstanceProcessed(exampleClient, instanceName)
	if err != nil {
		panic(err)
	}
	fmt.Println("CREATE")

	// Get the list of CRs.
	exampleList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)

	// As there is a cleanup logic before, here it sleeps for a while for example view.
	sleepDuration := 10 * time.Second
	fmt.Printf("Sleep for %s...\n", sleepDuration.String())
	time.Sleep(sleepDuration)

	exampleGet, err := crdClient.Get(instanceName)
	if err != nil {
		panic(err)
	}

	resourceVersion := exampleGet.ObjectMeta.ResourceVersion

	// Update an instance of CRD
	exampleInstanceUpdate := &v1.Item{
		ObjectMeta: metav1.ObjectMeta{
			Name:            instanceName,
			ResourceVersion: resourceVersion,
		},
		Attribute: v1.ItemAttribute{
			WelcomeMsg: "Play with Kubernetes",
			SleepTime:  30,
		},
		Status: v1.ItemStatus{
			State:   v1.StateUpdated,
			Message: "Updated but not processed yet",
		},
	}
	resultUpdate, err := crdClient.Update(instanceName, exampleInstanceUpdate)
	if err == nil {
		fmt.Printf("UPDATED: %#v", resultUpdate)
	} else {
		panic(err)
	}

	fmt.Println("UPDATE")

	// Get the list of CRs.
	exampleUpdateList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleUpdateList)

	// As there is a cleanup logic before, here it sleeps for a while for example view.
	sleepDuration = 5 * time.Second
	fmt.Printf("Sleep for %s...\n", sleepDuration.String())
	time.Sleep(sleepDuration)

	// Delete an Instance of CRD
	err = crdClient.Delete(instanceName, &metav1.DeleteOptions{})
	if err == nil {
		fmt.Printf("DELETED Completely")
	} else {
		panic(err)
	}

	// Get the list of CRs.
	exampleFinalList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleFinalList)

	// As there is a cleanup logic before, here it sleeps for a while for example view.
	sleepDurationDelete := 5 * time.Second
	fmt.Printf("Sleep for %s...\n", sleepDurationDelete.String())
	time.Sleep(sleepDurationDelete)

}
