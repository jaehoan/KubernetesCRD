package controller

import (
	"context"
	"fmt"
	"log"
	"strconv"

	crd "v1"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// Run starts a CRD resource controller.
func (c *Controller) Run(ctx context.Context) error {
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

func (c *Controller) watch(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.Client,
		crd.ItemResourcePlural,
		corev1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		source,
		&crd.Item{},
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

func (c *Controller) onAdd(obj interface{}) {
	test := obj.(*crd.Item)
	fmt.Println("[CONTROLLER] Create - Sleep Time: " + strconv.Itoa(test.Attribute.SleepTime) + ", Welcome Message: " + test.Attribute.WelcomeMsg)
	log.Println("[CONTROLLER] Create - Sleep Time: " + strconv.Itoa(test.Attribute.SleepTime) + ", Welcome Message: " + test.Attribute.WelcomeMsg)

	// Use DeepCopy() to make a deep copy of original object and modify this copy
	// or create a copy manually for better performance.
	testCopy := test.DeepCopy()
	testCopy.Status = crd.ItemStatus{
		State:   crd.StateProcessed,
		Message: "Successfully processed by controller",
	}

	err := c.Client.Put().
		Name(test.ObjectMeta.Name).
		Namespace(test.ObjectMeta.Namespace).
		Resource(crd.ItemResourcePlural).
		Body(testCopy).
		Do().
		Error()

	if err != nil {
		fmt.Println("ERROR updating status: " + err.Error())
		log.Println("ERROR updating status: " + err.Error())
	}
}

func (c *Controller) onUpdate(oldObj, newObj interface{}) {
	old := oldObj.(*crd.Item)
	new := newObj.(*crd.Item)
	if old.Attribute.SleepTime != new.Attribute.SleepTime || old.Attribute.WelcomeMsg != new.Attribute.WelcomeMsg {
		fmt.Println("[CONTROLLER] Before Update - Sleep Time: " + strconv.Itoa(old.Attribute.SleepTime) + ", Welcome Message: " + old.Attribute.WelcomeMsg)
		log.Println("[CONTROLLER] Before Update - Sleep Time: " + strconv.Itoa(old.Attribute.SleepTime) + ", Welcome Message: " + old.Attribute.WelcomeMsg)
		fmt.Println("[CONTROLLER] After  Update - Sleep Time: " + strconv.Itoa(new.Attribute.SleepTime) + ", Welcome Message: " + new.Attribute.WelcomeMsg)
		log.Println("[CONTROLLER] After  Update - Sleep Time: " + strconv.Itoa(new.Attribute.SleepTime) + ", Welcome Message: " + new.Attribute.WelcomeMsg)
	}
}

func (c *Controller) onDelete(obj interface{}) {
	test := obj.(*crd.Item)
	fmt.Println("[CONTROLLER] Delete Namespace: " + test.ObjectMeta.Namespace + ", Name: " + test.ObjectMeta.Name)
	log.Println("[CONTROLLER] Delete Namespace: " + test.ObjectMeta.Namespace + ", Name: " + test.ObjectMeta.Name)
}
