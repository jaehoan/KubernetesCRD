package controller

import (
	"context"
	"fmt"
	"log"
	"time"
	"strconv"

	crd "v1"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// Run starts a CRD resource controller.
func (c *ExampleController) Run(ctx context.Context) error {
	fmt.Println("Watch CRD objects...")

	// Watch CRD objects
	_, err := c.watch(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Example resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *ExampleController) watch(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.ExampleClient,
		crd.ExampleResourcePlural,
		corev1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		source,
		&crd.Example{},
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,
		// CRD event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		},
	)

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *ExampleController) onAdd(obj interface{}) {
	test := obj.(*crd.Example)
	fmt.Println("[CONTROLLER] OnAdd " + test.ObjectMeta.Name + ", Created Time : " + time.Now().String())
	log.Println("[CONTROLLER] OnAdd " + test.ObjectMeta.Name + ", Created Time : " + time.Now().String())

	// Use DeepCopy() to make a deep copy of original object and modify this copy
	// or create a copy manually for better performance.
	testCopy := test.DeepCopy()
	testCopy.Status = crd.ExampleStatus{
		State:   crd.StateProcessed,
		Message: "Successfully processed by controller",
	}

	err := c.ExampleClient.Put().
		Name(test.ObjectMeta.Name).
		Namespace(test.ObjectMeta.Namespace).
		Resource(crd.ExampleResourcePlural).
		Body(testCopy).
		Do().
		Error()

	if err != nil {
		fmt.Println("ERROR updating status: " + err.Error())
		log.Println("ERROR updating status: " + err.Error())
	} else {
		fmt.Println("UPDATED status: " + testCopy.SelfLink)
		log.Println("UPDATED status: " + testCopy.SelfLink)
	}
}

func (c *ExampleController) onUpdate(oldObj, newObj interface{}) {
	old := oldObj.(*crd.Example)
	new := newObj.(*crd.Example)
	fmt.Println("[CONTROLLER] OnUpdate old - Spec.Foo: " + old.Spec.Foo + ", Spec.Bar: " + strconv.FormatBool(old.Spec.Bar))
	log.Println("[CONTROLLER] OnUpdate old - Spec.Foo: " + old.Spec.Foo + ", Spec.Bar: " + strconv.FormatBool(old.Spec.Bar))
	fmt.Println("[CONTROLLER] OnUpdate new - Spec.Foo: " + new.Spec.Foo + ", Spec.Bar: " + strconv.FormatBool(new.Spec.Bar))
	log.Println("[CONTROLLER] OnUpdate new - Spec.Foo: " + new.Spec.Foo + ", Spec.Bar: " + strconv.FormatBool(new.Spec.Bar))
}

func (c *ExampleController) onDelete(obj interface{}) {
	test := obj.(*crd.Example)
	fmt.Println("[CONTROLLER] OnDelete Namespace: " + test.ObjectMeta.Namespace + ", Name: " + test.ObjectMeta.Name)
	log.Println("[CONTROLLER] OnDelete Namespace: " + test.ObjectMeta.Namespace + ", Name: " + test.ObjectMeta.Name)
}
