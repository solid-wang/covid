package main

import (
	"context"
	"fmt"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type name struct {
	a    *string
	amap map[string]*string
}

func main() {
	trigger := "goweb-c27b48e8-dev"
	config := &rest.Config{Host: "http://127.0.0.1:8080"}
	client, _ := clientset.NewForConfig(config)
	ci, err := client.BatchV1().ContinuousIntegrations("devops").Get(context.Background(), "goweb-c27b48e8", metav1.GetOptions{})
	continuousIntegration := ci.DeepCopy()
	continuousIntegration.Status.ContinuousDeploymentTrigger = &trigger
	continuousIntegration.Status.Phase = batchv1.DevOpsFailedStatus
	update, err := client.BatchV1().ContinuousIntegrations("devops").UpdateStatus(context.Background(), continuousIntegration, metav1.UpdateOptions{})
	//ci, err := clientSet.BatchV1().ContinuousIntegrations("default").Get(context.TODO(), "goweb", metav1.GetOptions{})
	//newci := ci.DeepCopy()
	//newci.Status.Phase = batchv1.DevOpsRunningStatus
	//////newci.Status.ObservedGeneration = newci.Generation
	//////newci.Spec.Image = "abcde"
	//update, err := clientSet.BatchV1().ContinuousIntegrations(ci.Namespace).UpdateStatus(context.Background(), newci, metav1.UpdateOptions{})
	//////update, err := clientSet.BatchV1().ContinuousIntegrations(ci.Namespace).Update(context.TODO(), newci, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	////fmt.Println(ci == nil)
	fmt.Println("update")
	fmt.Println(update)
	//watch, err := clientSet.DevopsV1().ContinuousIntegrations("default").Watch(context.Background(), metav1.ListOptions{})
	//if err != nil {
	//	fmt.Println("err")
	//	fmt.Println(err)
	//}
	//for {
	//	select {
	//	case <-watch.ResultChan():
	//		fmt.Println("enter")
	//		fmt.Println(watch.ResultChan())
	//		_ = <-watch.ResultChan()
	//	}
	//}

	//push := events.NewPush()
	//push.Ref = "refs/heads/master"
	//fmt.Println(push.GetBranch())
	//a := "123456789"
	//for _, c := range strings.Split(strings.TrimSpace(a[len(a):]), " ") {
	//	fmt.Printf("c %d", len(c))
	//	kv := strings.Split(c, "=")
	//	fmt.Printf("kv %s", kv[0])
	//}
	//fmt.Println(errors.IsNotFound(nil))
	//a := map[string]string{
	//	"x": "y",
	//}
	//delete(a, "x")
	//fmt.Println(a)
}
