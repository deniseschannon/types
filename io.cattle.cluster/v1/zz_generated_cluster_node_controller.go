package v1

import (
	"sync"

	"context"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	ClusterNodeGroupVersionKind = schema.GroupVersionKind{
		Version: "v1",
		Group:   "io.cattle.cluster",
		Kind:    "ClusterNode",
	}
	ClusterNodeResource = metav1.APIResource{
		Name:         "clusternodes",
		SingularName: "clusternode",
		Namespaced:   false,
		Kind:         ClusterNodeGroupVersionKind.Kind,
	}
)

type ClusterNodeList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []ClusterNode
}

type ClusterNodeHandlerFunc func(key string, obj *ClusterNode) error

type ClusterNodeController interface {
	Informer() cache.SharedIndexInformer
	AddHandler(handler ClusterNodeHandlerFunc)
	Enqueue(namespace, name string)
	Start(threadiness int, ctx context.Context) error
}

type ClusterNodeInterface interface {
	Create(*ClusterNode) (*ClusterNode, error)
	Get(name string, opts metav1.GetOptions) (*ClusterNode, error)
	Update(*ClusterNode) (*ClusterNode, error)
	Delete(name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ClusterNodeList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() (ClusterNodeController, error)
}

type clusterNodeController struct {
	controller.GenericController
}

func (c *clusterNodeController) AddHandler(handler ClusterNodeHandlerFunc) {
	c.GenericController.AddHandler(func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*ClusterNode))
	})
}

type clusterNodeFactory struct {
}

func (c clusterNodeFactory) Object() runtime.Object {
	return &ClusterNode{}
}

func (c clusterNodeFactory) List() runtime.Object {
	return &ClusterNodeList{}
}

func NewClusterNodeClient(namespace string, config rest.Config) (ClusterNodeInterface, error) {
	objectClient, err := clientbase.NewObjectClient(namespace, config, &ClusterNodeResource, ClusterNodeGroupVersionKind, clusterNodeFactory{})
	return &clusterNodeClient{
		objectClient: objectClient,
	}, err
}

func (s *clusterNodeClient) Controller() (ClusterNodeController, error) {
	s.Lock()
	defer s.Unlock()

	if s.controller != nil {
		return s.controller, nil
	}

	controller, err := controller.NewGenericController(ClusterNodeGroupVersionKind.Kind+"Controller",
		s.objectClient)
	if err != nil {
		return nil, err
	}

	s.controller = &clusterNodeController{
		GenericController: controller,
	}
	return s.controller, nil
}

type clusterNodeClient struct {
	sync.Mutex
	objectClient *clientbase.ObjectClient
	controller   ClusterNodeController
}

func (s *clusterNodeClient) Create(o *ClusterNode) (*ClusterNode, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ClusterNode), err
}

func (s *clusterNodeClient) Get(name string, opts metav1.GetOptions) (*ClusterNode, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ClusterNode), err
}

func (s *clusterNodeClient) Update(o *ClusterNode) (*ClusterNode, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ClusterNode), err
}

func (s *clusterNodeClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *clusterNodeClient) List(opts metav1.ListOptions) (*ClusterNodeList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ClusterNodeList), err
}

func (s *clusterNodeClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

func (s *clusterNodeClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}