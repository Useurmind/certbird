// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/useurmind/certbird/k8s-controller/pkg/apis/cr/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CertRequestLister helps list CertRequests.
type CertRequestLister interface {
	// List lists all CertRequests in the indexer.
	List(selector labels.Selector) (ret []*v1.CertRequest, err error)
	// CertRequests returns an object that can list and get CertRequests.
	CertRequests(namespace string) CertRequestNamespaceLister
	CertRequestListerExpansion
}

// certRequestLister implements the CertRequestLister interface.
type certRequestLister struct {
	indexer cache.Indexer
}

// NewCertRequestLister returns a new CertRequestLister.
func NewCertRequestLister(indexer cache.Indexer) CertRequestLister {
	return &certRequestLister{indexer: indexer}
}

// List lists all CertRequests in the indexer.
func (s *certRequestLister) List(selector labels.Selector) (ret []*v1.CertRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CertRequest))
	})
	return ret, err
}

// CertRequests returns an object that can list and get CertRequests.
func (s *certRequestLister) CertRequests(namespace string) CertRequestNamespaceLister {
	return certRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CertRequestNamespaceLister helps list and get CertRequests.
type CertRequestNamespaceLister interface {
	// List lists all CertRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.CertRequest, err error)
	// Get retrieves the CertRequest from the indexer for a given namespace and name.
	Get(name string) (*v1.CertRequest, error)
	CertRequestNamespaceListerExpansion
}

// certRequestNamespaceLister implements the CertRequestNamespaceLister
// interface.
type certRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CertRequests in the indexer for a given namespace.
func (s certRequestNamespaceLister) List(selector labels.Selector) (ret []*v1.CertRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CertRequest))
	})
	return ret, err
}

// Get retrieves the CertRequest from the indexer for a given namespace and name.
func (s certRequestNamespaceLister) Get(name string) (*v1.CertRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("certrequest"), name)
	}
	return obj.(*v1.CertRequest), nil
}
