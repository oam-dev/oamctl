package convert

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/oam-dev/oamctl/pkg/util"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	deploymentName string
	serviceName    string
	ingressName    string
)

// migrateCmd represents the create command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "create oam ApplicationConfiguration yaml from exist k8s resource",
	Long:  `create oam ApplicationConfiguration yaml from exist k8s resource, e.g: deployment, service, job....`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(genOamYaml(deploymentName, serviceName, ingressName, args))
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	migrateCmd.Flags().StringVar(&deploymentName, "deployment", "", "name of Kubernetes Deployment")
	migrateCmd.Flags().StringVar(&serviceName, "service", "", "name of Kubernetes Service")
	migrateCmd.Flags().StringVar(&ingressName, "ingress", "", "name of Kubernetes Ingress")
	_ = migrateCmd.MarkFlagRequired("deployment")
}

func getDeploy(name string) *v1.Deployment {
	var deployment v1.Deployment
	_ = mgr.GetAPIReader().Get(context.TODO(), types.NamespacedName{Namespace: nameSpace, Name: name}, &deployment)
	return &deployment
}

func getService(name string) *v12.Service {
	var service v12.Service
	_ = mgr.GetAPIReader().Get(context.TODO(), types.NamespacedName{Namespace: nameSpace, Name: name}, &service)
	return &service
}

func getIngress(name string) *v1beta1.Ingress {
	var ingress v1beta1.Ingress
	_ = mgr.GetAPIReader().Get(context.TODO(), types.NamespacedName{Namespace: nameSpace, Name: name}, &ingress)
	return &ingress
}

func genOamYaml(deployment, service, ingress string, args []string) string {
	deploy := getDeploy(deployment)
	if deploy == nil {
		return fmt.Sprintf("can not get deployment %s", deployment)
	}
	var appName = deploy.Name
	if len(args) > 0 {
		appName = args[0]
	}

	svc := getService(service)
	ing := getIngress(ingress)
	if svc == nil {
		fmt.Println("start to render worker")
		res, _ := util.RenderWorker(struct {
			Deployment *v1.Deployment
			Name       string
		}{
			deploy,
			appName,
		})
		return res
	}
	fmt.Println("start to render server")
	res, err := util.RenderServer(struct {
		Deployment *v1.Deployment
		Service    *v12.Service
		Ingress    *v1beta1.Ingress
		Name       string
	}{
		deploy,
		svc,
		ing,
		appName,
	})
	if err != nil {
		logrus.Error(err)
	}
	return res
}
