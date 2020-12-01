package controllers

import (
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	seleniumk8siov1 "selenosis-operator/api/v1"
)

const (
	appLabel                 = "app"
	selenosisDefaultAppName  = "selenosis"
	selenoidUIDefaultAppName = "selenoid-ui"
	browsersConfigVolumeName = "browsers-config"
)

var (
	rollingUpdateValues = struct {
		maxSurge, maxUnavailable intstr.IntOrString
	}{
		maxSurge:       intstr.FromInt(1),
		maxUnavailable: intstr.FromInt(1),
	}
)

func createSelenosisDeployment(hubspec *seleniumk8siov1.K8sSeleniumHub) (*appsv1.Deployment, error) {
	var errorMsg string
	if hubspec.Spec.SelenosisReplicas <= 0 && hubspec.Spec.BrowserLimit == "" {
		errorMsg = fmt.Sprintf("Unable to create deployment for selenosis with number of replicas %d with not set browser limit", hubspec.Spec.SelenosisReplicas)
		return nil, errors.New(errorMsg)
	}
	if hubspec.Spec.SelenosisPort <= 0 {
		errorMsg = fmt.Sprintf("Unable to create deployment for selenosis with selenoid port equal to %d", hubspec.Spec.SelenosisPort)
		return nil, errors.New(errorMsg)
	}
	selenosisDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: hubspec.Namespace,
			Name:      selenosisDefaultAppName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(hubspec.Spec.SelenosisReplicas),
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &rollingUpdateValues.maxUnavailable,
					MaxSurge:       &rollingUpdateValues.maxSurge,
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: createAppLabels(selenosisDefaultAppName),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: hubspec.Namespace,
					Labels:    createAppLabels(selenosisDefaultAppName),
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: browsersConfigVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "selenosis-config",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Args:            createSelenosisArgs(hubspec),
							Image:           hubspec.Spec.SelenosisImage,
							Name:            selenosisDefaultAppName,
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									Name:          "selenium",
									ContainerPort: 4444,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      browsersConfigVolumeName,
									MountPath: "/etc/selenosis",
								},
							},
							ReadinessProbe: createProbeForSelenosis(hubspec),
							LivenessProbe:  createProbeForSelenosis(hubspec),
						},
					},
					ImagePullSecrets: setupImagePullSecrets(hubspec),
				},
			},
		},
	}

	return selenosisDeployment, nil
}

func createBrowsersHeadlessService(hubspec *seleniumk8siov1.K8sSeleniumHub) (*corev1.Service, error) {
	var errorMsg string
	if hubspec.Spec.BrowsersServiceName == "" {
		errorMsg = fmt.Sprintf("Unable create browsers headless service with empty service name. Please provide browsers service name in K8sSeleniumHub specification")
		return nil, errors.New(errorMsg)
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hubspec.Spec.BrowsersServiceName,
			Namespace: hubspec.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"type": "browser",
			},
			ClusterIP:                "None",
			PublishNotReadyAddresses: true,
		},
	}
	return service, nil
}

func createSelenosisService(hubspec *seleniumk8siov1.K8sSeleniumHub) (*corev1.Service, error) {
	var errorMsg string
	if hubspec.Spec.SelenosisPort <= 0 {
		errorMsg = fmt.Sprintf("Error creating the selenosis service. Please specify the port number for selenosis service in K8sSeleniumHub specification")
		return nil, errors.New(errorMsg)
	}
	if hubspec.Spec.SelenosisServiceName == "" {
		errorMsg = fmt.Sprintf("Error creating the selenosis service. Please specify the selenosis service name in K8sSeleniumHub specification")
		return nil, errors.New(errorMsg)
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hubspec.Spec.SelenosisServiceName,
			Namespace: hubspec.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "selenium",
					Protocol: "TCP",
					Port:     hubspec.Spec.SelenosisPort,
					TargetPort: intstr.IntOrString{
						IntVal: hubspec.Spec.SelenosisPort,
					},
					NodePort: 31000,
				},
			},
			Selector:              createAppLabels(selenosisDefaultAppName),
			Type:                  corev1.ServiceTypeLoadBalancer,
			ExternalTrafficPolicy: "Cluster",
			SessionAffinity:       corev1.ServiceAffinityNone,
		},
	}
	return service, nil
}

//func createSelenidUIDeployment(hybspec *seleniumk8siov1.K8sSeleniumHub) (*appsv1.Deployment, error) {
//	var errorMsg string
//
//}

func createAppLabels(label string) map[string]string {
	return map[string]string{
		appLabel: label,
	}
}

func setupImagePullSecrets(hubspec *seleniumk8siov1.K8sSeleniumHub) []corev1.LocalObjectReference {
	var result []corev1.LocalObjectReference
	if hubspec.Spec.ImagePullSecretName != "" {
		result = append(result, corev1.LocalObjectReference{Name: hubspec.Spec.ImagePullSecretName})
		return result
	}
	return nil
}

func createSelenosisArgs(hubspec *seleniumk8siov1.K8sSeleniumHub) []string {
	var selenosisCmdArgs = []string{
		"/selenosis",
		"--browsers-config",
		"/etc/selenosis/browsers.yaml",
		"--namespace",
		hubspec.Namespace,
		"--service-name",
		hubspec.Spec.BrowsersServiceName,
		"--browser-limit",
		hubspec.Spec.BrowserLimit,
	}
	if hubspec.Spec.ImagePullSecretName != "" {
		selenosisCmdArgs = append(selenosisCmdArgs, hubspec.Spec.ImagePullSecretName)
	}
	return selenosisCmdArgs
}

func createProbeForSelenosis(hubspec *seleniumk8siov1.K8sSeleniumHub) *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthz",
				Port: intstr.IntOrString{
					IntVal: hubspec.Spec.SelenosisPort,
				},
			},
		},
		InitialDelaySeconds: 3,
		PeriodSeconds:       2,
	}
}
