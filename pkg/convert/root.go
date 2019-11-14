package convert

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	cfgFile       string
	nameSpace     string
	runtimeScheme = runtime.NewScheme()
	mgr           ctrl.Manager
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oamctl",
	Short: "oamctl is a tiny tool help oam user manage oam app.",
	Long:  `oamctl is a tiny tool help oam user manage oam app, includes: migrate existed k8s resource to oam app, create/update/delete oam app...`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	_ = appsv1.AddToScheme(runtimeScheme)
	_ = v1beta1.AddToScheme(runtimeScheme)
	_ = corev1.AddToScheme(runtimeScheme)
	_ = v2beta2.AddToScheme(runtimeScheme)
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":52014", "The address the metric endpoint binds to.")
	flag.Parse()
	m, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: runtimeScheme, MetricsBindAddress: metricsAddr})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mgr = m

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oamctl.yaml)")

	rootCmd.PersistentFlags().StringVarP(&nameSpace, "namespace", "n", "default", "operate namespace")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".oamctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".oamctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
