package util

import (
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// GetSecrets retrieves a secret value and memoizes the result
func GetSecrets(clientSet kubernetes.Interface, namespace, name, key string) ([]byte, error) {
	secretsIf := clientSet.CoreV1().Secrets(namespace)
	var secret *apiv1.Secret
	var err error
	_ = wait.ExponentialBackoff(retry.DefaultRetry, func() (bool, error) {
		secret, err = secretsIf.Get(name, metav1.GetOptions{})
		if err != nil {
			log.Warnf("Failed to get secret '%s': %v", name, err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return []byte{}, err
	}

	val, ok := secret.Data[key]
	if !ok {
		return []byte{}, nil
	}
	return val, nil
}
