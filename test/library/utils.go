package library

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (d DynamicResourceLoader) CreateTestingNS(baseName string) (*v1.Namespace, error) {
	name := fmt.Sprintf("%v", baseName)

	namespaceObj := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "",
		},
		Status: v1.NamespaceStatus{},
	}

	var got *v1.Namespace
	if err := wait.PollImmediate(1*time.Second, 30*time.Second, func() (bool, error) {
		var err error
		got, err = d.KubeClient.CoreV1().Namespaces().Create(context.Background(), namespaceObj, metav1.CreateOptions{})
		if err != nil {
			// t.Logf("Error creating namespace: %v", err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, err
	}

	return got, nil
}

func (d DynamicResourceLoader) DeleteTestingNS(baseName string) (bool, error) {
	name := fmt.Sprintf("%v", baseName)

	if err := wait.PollImmediate(1*time.Second, 30*time.Second, func() (bool, error) {

		// Poll until namespace is deleted
		if _, err := d.KubeClient.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{}); err != nil {
			if k8serrors.IsNotFound(err) {
				return true, err
			}
			return false, nil
		}
		return false, nil
	}); err != nil {
		return true, err
	}
	return false, nil
}
