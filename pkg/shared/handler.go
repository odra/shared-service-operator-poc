package shared

import (
	"context"

	"github.com/aerogear/shared-service-operator-poc/pkg/apis/aerogear/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
)

func NewHandler(k8sClient kubernetes.Interface, sharedServiceClient dynamic.ResourceInterface, operatorNS string) sdk.Handler {
	return &Handler{
		k8client:k8sClient,
		operatorNS : operatorNS,
		sharedServiceClient: sharedServiceClient,
	}
}

type Handler struct {
	// Fill me
	k8client kubernetes.Interface
	operatorNS string
	sharedServiceClient dynamic.ResourceInterface
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {

	switch o := event.Object.(type) {
	case *v1alpha1.SharedService:
		fmt.Println("shared service recieved ", o.Namespace, o.Name, o.Status, event.Deleted)
		if event.Deleted{
			return h.handleSharedServiceDelete(o)
		}
		return h.handleSharedServiceCreateUpdate(o)
	case *v1alpha1.SharedServiceSlice:
		fmt.Println("shared service slice recieved ", o.Namespace, o.Name, o.Status, event.Deleted)
		if event.Deleted{
			return h.handleSharedServiceSliceDelete(o)
		}
		return h.handleSharedServiceSliceCreateUpdate(o)

	case *v1alpha1.SharedServiceClient:
		fmt.Println("shared service slice recieved ", o.Namespace, o.Name, o.Status, event.Deleted)
		if event.Deleted{
			return h.handleSharedServiceClientDelete(o)
		}
		return h.handleSharedServiceClientCreateUpdate(o)
	}
	return nil
}


func (h *Handler)handleSharedServiceCreateUpdate(service *v1alpha1.SharedService)error{
	fmt.Println("called handleSharedServiceCreateUpdate ")
	if service.Status.Ready{
		//delete the pod
	}
	if service.Status.Status == "" {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName:service.Name + "-",
				//Labels: extContext.Metadata,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  service.Name + "-create",
						Image: service.Spec.Image,
						Args: []string{
							"provision",
							"--extra-vars",
							"", // need to figure out how to get and pass the needed params
						},
						//Env:             createPodEnv(extContext),
						ImagePullPolicy: "Always",
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
				//ServiceAccountName: extContext.Account,
				//Volumes:            volumes,
			},
		}

		logrus.Infof(fmt.Sprintf("Creating pod %q in the %s namespace", pod.Name, h.operatorNS))
		_, err := h.k8client.CoreV1().Pods(h.operatorNS).Create(pod)
		if err != nil {
			return err
		}
		//watch pod until complete and update the status of the shared service
	}
	service.Status.Status = "provisioning"
	unstructObj := k8sutil.UnstructuredFromRuntimeObject(service)
	if _, err := h.sharedServiceClient.Update(unstructObj); err != nil{
		return err
	}

	return nil
}

func (h *Handler)handleSharedServiceDelete(service *v1alpha1.SharedService)error{
	fmt.Println("called handleSharedServiceDelete")
	return nil
}

func (h *Handler)handleSharedServiceSliceCreateUpdate(service *v1alpha1.SharedServiceSlice)error{
	fmt.Println("called handleSharedServiceSliceCreateUpdate")
	return nil
}

func (h *Handler)handleSharedServiceSliceDelete(service *v1alpha1.SharedServiceSlice)error{
	fmt.Println("called handleSharedServiceSliceDelete")
	return nil
}

func (h *Handler)handleSharedServiceClientCreateUpdate(serviceClient *v1alpha1.SharedServiceClient)error{
	fmt.Println("called handleSharedServiceClientCreateUpdate")
	return nil
}

func (h *Handler)handleSharedServiceClientDelete(serviceClient *v1alpha1.SharedServiceClient)error{
	fmt.Println("called handleSharedServiceClientDelete")
	return nil
}