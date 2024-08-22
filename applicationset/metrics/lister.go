package metrics

import (
	"context"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/pkg/client/listers/application/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Implements AppsetLister interface
type AppsetLister struct {
	Client ctrlclient.Client
}

func NewAppsetLister(Client ctrlclient.Client) ApplicationSetLister {
	return &AppsetLister{Client: Client}
}

func (l *AppsetLister) List(selector labels.Selector) (ret []*ApplicationSet, err error) {
	return clientListAppsets(l.Client, ctrlclient.ListOptions{})
}

// ApplicationSets returns an object that can list and get ApplicationSets.
func (l *AppsetLister) ApplicationSets(namespace string) ApplicationSetNamespaceLister {

	return &appsetNamespaceLister{
		Client:    l.Client,
		Namespace: namespace,
	}
}

// Implements ApplicationSetNamespaceLister
type appsetNamespaceLister struct {
	Client    ctrlclient.Client
	Namespace string
}

func (n *appsetNamespaceLister) List(selector labels.Selector) (ret []*ApplicationSet, err error) {
	return clientListAppsets(n.Client, ctrlclient.ListOptions{Namespace: n.Namespace})
}

func (n *appsetNamespaceLister) Get(name string) (*ApplicationSet, error) {
	appset := ApplicationSet{}
	err := n.Client.Get(context.TODO(), ctrlclient.ObjectKeyFromObject(&appset), &appset)
	return &appset, err
}

func clientListAppsets(client ctrlclient.Client, ListOptions ctrlclient.ListOptions) (ret []*ApplicationSet, err error) {

	var appsetlist ApplicationSetList
	var results []*ApplicationSet

	err = client.List(context.TODO(), &appsetlist, &ListOptions)

	if err == nil {
		for _, appset := range appsetlist.Items {
			results = append(results, appset.DeepCopy())
		}
	}

	return results, err
}
