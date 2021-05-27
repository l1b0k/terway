package cert

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AliyunContainerService/terway/pkg/utils"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	serverCertKey = "tls.crt"
	serverKeyKey  = "tls.key"

	caCertKey = "ca.crt"
)

// SyncCert sync cert for webhook
func SyncCert() error {
	certDir := viper.GetString("cert-dir")
	err := os.MkdirAll(certDir, os.ModeDir)
	if err != nil {
		return err
	}

	cs := utils.K8sClient
	// check secret
	var serverCertBytes, serverKeyBytes, caCertBytes []byte

	// get cert from secret or generate it
	existSecret, err := cs.CoreV1().Secrets(viper.GetString("controller-namespace")).Get(context.Background(), "terway-controlplane", metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("error get cert from secret, %w", err)
		}
		// create certs
		s, err := GenerateCerts(viper.GetString("controller-namespace"), "webhook", "cluster.local")
		if err != nil {
			return fmt.Errorf("error generate cert, %w", err)
		}

		serverCertBytes = s.Data[serverCertKey]
		serverKeyBytes = s.Data[serverKeyKey]
		caCertBytes = s.Data[caCertKey]

		// create secret this make sure one is the leader
		_, err = cs.CoreV1().Secrets(viper.GetString("controller-namespace")).Create(context.Background(), s, metav1.CreateOptions{})
		if err != nil {
			if !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error create cert to secret, %w", err)
			}
			secret, err := cs.CoreV1().Secrets(viper.GetString("controller-namespace")).Get(context.Background(), "terway-controlplane", metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error get cert from secret, %w", err)
			}
			_, err = base64.StdEncoding.Decode(serverCertBytes, secret.Data[serverCertKey])
			if err != nil {
				return fmt.Errorf("error decode from secret, %w", err)
			}
			_, err = base64.StdEncoding.Decode(serverKeyBytes, secret.Data[serverKeyKey])
			if err != nil {
				return fmt.Errorf("error decode from secret, %w", err)
			}
			_, err = base64.StdEncoding.Decode(caCertBytes, secret.Data[caCertKey])
			if err != nil {
				return fmt.Errorf("error decode from secret, %w", err)
			}
		}
	} else {
		_, err = base64.StdEncoding.Decode(serverCertBytes, existSecret.Data[serverCertKey])
		if err != nil {
			return fmt.Errorf("error decode from secret, %w", err)
		}
		_, err = base64.StdEncoding.Decode(serverKeyBytes, existSecret.Data[serverKeyKey])
		if err != nil {
			return fmt.Errorf("error decode from secret, %w", err)
		}
		_, err = base64.StdEncoding.Decode(caCertBytes, existSecret.Data[caCertKey])
		if err != nil {
			return fmt.Errorf("error decode from secret, %w", err)
		}
	}
	// write cert to file
	err = os.WriteFile(filepath.Join(certDir, serverCertKey), serverCertBytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error create secret file, %w", err)
	}
	err = os.WriteFile(filepath.Join(certDir, serverKeyKey), serverKeyBytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error create secret file, %w", err)
	}
	err = os.WriteFile(filepath.Join(certDir, caCertKey), caCertBytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error create secret file, %w", err)
	}

	// update webhook
	webhook, err := cs.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(context.Background(), "terway-controlplane", metav1.GetOptions{})
	if err != nil {
		return err
	}
	// only have one
	for _, hook := range webhook.Webhooks {
		if len(hook.ClientConfig.CABundle) != 0 {
			return nil
		}
		// patch ca
		webhook.Webhooks[0].ClientConfig.CABundle = caCertBytes
		patchBytes, err := json.Marshal(webhook)
		if err != nil {
			return err
		}
		err = wait.ExponentialBackoff(utils.DefaultPatchBackoff, func() (done bool, err error) {
			_, innerErr := cs.AdmissionregistrationV1().MutatingWebhookConfigurations().Patch(context.Background(), webhook.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
			if innerErr != nil {
				return false, err
			}
			return true, nil
		})
		return err
	}
	return nil
}

func GenerateCerts(serviceNamespace, serviceName, clusterDomain string) (*corev1.Secret, error) {
	var caPEM, serverCertPEM, serverPrivateKeyPEM *bytes.Buffer
	ca := &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{clusterDomain},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPK, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	caBytes, err := x509.CreateCertificate(cryptorand.Reader, ca, ca, &caPK.PublicKey, caPK)
	if err != nil {
		return nil, err
	}

	caPEM = new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, err
	}

	commonName := fmt.Sprintf("%s.%s.svc", serviceName, serviceNamespace)
	dnsNames := []string{serviceName,
		fmt.Sprintf("%s.%s", serviceName, serviceNamespace),
		commonName}

	cert := &x509.Certificate{
		DNSNames: dnsNames,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{clusterDomain},
		},
		NotBefore:   time.Now().Add(-time.Hour),
		NotAfter:    time.Now().AddDate(100, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	serverPrivateKey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	serverCertBytes, err := x509.CreateCertificate(cryptorand.Reader, cert, ca, &serverPrivateKey.PublicKey, caPK)
	if err != nil {
		return nil, err
	}

	serverCertPEM = new(bytes.Buffer)
	err = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})
	if err != nil {
		return nil, err
	}

	serverPrivateKeyPEM = new(bytes.Buffer)
	err = pem.Encode(serverPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivateKey),
	})
	if err != nil {
		return nil, err
	}

	return &corev1.Secret{Data: map[string][]byte{
		caCertKey:     caPEM.Bytes(),
		serverCertKey: serverCertPEM.Bytes(),
		serverKeyKey:  serverPrivateKeyPEM.Bytes(),
	}}, nil
}

// WriteFile writes data in the file at the given path
func WriteFile(filepath string, sCert *bytes.Buffer) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(sCert.Bytes())
	if err != nil {
		return err
	}
	return nil
}
