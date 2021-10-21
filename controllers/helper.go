package controllers

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	myv1 "saha.com/mycrd/api/v1"
	"strings"
)

// newDeployment creates a new Deployment for a messi resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the messi resource that 'owns' it.
func newDeployment(jr10 *myv1.Neymar) *appsv1.Deployment {
	fmt.Println("newDeployment is called")
	labels := map[string]string{
		"app":        trimTheOwnerPartFromImageName(jr10.Spec.DeploymentImage),
		"controller": jr10.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jr10.Spec.DeploymentName,
			Namespace: jr10.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(jr10, myv1.GroupVersion.WithKind("Neymar")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: jr10.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "neymarcontainer",
							Image: jr10.Spec.DeploymentImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: jr10.Spec.ServiceTargetPort,
								},
							},
						},
					},
				},
			},
		},
	}
}
func newService(jr10 *myv1.Neymar) *corev1.Service {
	fmt.Println("newService is called")
	labels := map[string]string{
		"app":        trimTheOwnerPartFromImageName(jr10.Spec.DeploymentImage),
		"controller": jr10.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jr10.Spec.ServiceName,
			Namespace: jr10.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(jr10, myv1.GroupVersion.WithKind("Neymar")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     getTheServiceType(jr10.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					NodePort: int32(30011),
					Port:     jr10.Spec.ServicePort,
					TargetPort: intstr.IntOrString{
						IntVal: jr10.Spec.ServiceTargetPort,
					},
				},
			},
		},
	}
}

func trimTheOwnerPartFromImageName(s string) string {
	arr := strings.Split(s, "/")
	if len(arr) == 1 {
		return arr[0]
	}
	return arr[1]
}

func getTheServiceType(s string) corev1.ServiceType {
	if s == "NodePort" {
		return corev1.ServiceTypeNodePort
	} else if s == "ClusterIP" {
		return corev1.ServiceTypeClusterIP
	}
	return corev1.ServiceTypeClusterIP
}

/*
func stringToInt(s string) int32 {
	x, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Err in stringToInt().", err)
		return 0
	}
	return int32(x)
}*/
