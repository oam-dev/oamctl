package convert

import (
	"fmt"
	"testing"

	"k8s.io/api/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/oam-dev/oamctl/pkg/util"
	"github.com/stretchr/testify/assert"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRenderServer(t *testing.T) {
	var replica int32 = 3
	deploy := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deploy_name",
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replica,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:       "c1",
							Image:      "image1",
							Command:    []string{"cmd1", "cmd2"},
							Args:       []string{"arg1"},
							WorkingDir: "/app",
							Ports:      nil,
							Env: []corev1.EnvVar{
								{
									Name:  "env1",
									Value: "value1",
								},
							},
							Resources:     corev1.ResourceRequirements{},
							VolumeMounts:  nil,
							VolumeDevices: nil,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"ping"},
									},
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/v1/health",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 9999,
										},
										HTTPHeaders: []corev1.HTTPHeader{
											{Name: "H1", Value: "V1"},
										},
									},
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.String,
											StrVal: "9998",
										},
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      2,
								PeriodSeconds:       3,
								SuccessThreshold:    4,
								FailureThreshold:    5,
							},
							ReadinessProbe:  nil,
							ImagePullPolicy: "",
						},
						{
							Name:       "c2",
							Image:      "image2",
							Command:    []string{"cmd3"},
							Args:       []string{"arg2", "arg3"},
							WorkingDir: "/usr/bin",
							Ports: []corev1.ContainerPort{
								{
									Name:          "cp1",
									HostPort:      8800,
									ContainerPort: 8800,
									Protocol:      "http",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "env2",
									Value: "value2",
								},
							},
							Resources:     corev1.ResourceRequirements{},
							VolumeMounts:  nil,
							VolumeDevices: nil,
							LivenessProbe: nil,
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"ping"},
									},
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/v1/health",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 9999,
										},
										HTTPHeaders: []corev1.HTTPHeader{
											{Name: "H1", Value: "V1"},
										},
									},
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.String,
											StrVal: "9998",
										},
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      4,
								PeriodSeconds:       3,
								SuccessThreshold:    2,
								FailureThreshold:    1,
							},
							ImagePullPolicy: "",
						},
					},
				},
			},
		},
	}
	res, err := util.RenderServer(struct {
		Deployment *v1.Deployment
		Service    *corev1.Service
		Ingress    *v1beta1.Ingress
		Name       string
	}{Deployment: deploy, Name: "test_name", Service: nil, Ingress: nil})
	assert.NoError(t, err)
	fmt.Println(res)
}
