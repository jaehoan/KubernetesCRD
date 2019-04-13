package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	ItemResourcePlural string = "examples"
	// GroupName is the group name used in this package.
	GroupName        string = "welcome.group"
	ExampleCRDName   string = ItemResourcePlural + "." + GroupName
	version          string = "v1"
	StateCreated     string = "Created"
	StateUpdated     string = "Updated"
	StateProcessed   string = "Processed"
	DefaultNamespace string = "default"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   GroupName,
		Version: version,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Item is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh deepcopy github.com/jinghzhu/k8scrd/client github.com/jinghzhu/k8scrd/apis "example:v1"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Item struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Attribute         ItemAttribute `json:"attribute"`
	Status            ItemStatus    `json:"status,omitempty"`
}

type ItemAttribute struct {
	WelcomeMsg string `json:"welcomeMsg"`
	SleepTime  int    `json:"sleepTime"`
}

type ItemStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Item `json:"items"`
}
