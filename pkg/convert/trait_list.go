package convert

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/oam-dev/oamctl/pkg/apis/core.oam.dev/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/spf13/cobra"
)

// migrateCmd represents the create command
var listCmd = &cobra.Command{
	Use:   "trait-list",
	Short: "list all oam traits from k8s clusters",
	Long:  `list all oam traits from k8s clusters, with name, version and PRIMITIVES`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(listTraits(args))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func listTraits(args []string) string {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	v1alpha1.AddToScheme(scheme.Scheme)
	crdConfig := *cfg
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	traitsClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	traitsList := &v1alpha1.TraitList{}
	err = traitsClient.Get().Resource("traits").Do().Into(traitsList)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if len(traitsList.Items) < 1 {
		return "no traits found"
	}

	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', tabwriter.DiscardEmptyColumns)
	defer w.Flush()
	fmt.Fprintf(w, "%s\t%s\t%s\t", "NAME", "VERSION", "PRIMITIVES")
	fmt.Fprintf(w, "\n%s\t%s\t%s\t", "--------", "-------", "----------")

	for _, v := range traitsList.Items {
		fmt.Fprintf(w, "\n%s\t%s\t%s\t", v.Name, v.Annotations["version"], strings.Join(v.Spec.AppliesTo, ", "))
	}
	return ""
}
