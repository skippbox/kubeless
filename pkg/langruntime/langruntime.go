package langruntime

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/pkg/api/v1"
)

// Langruntime struct
type Langruntimes struct {
}

var availableRuntimes []RuntimeInfo

// ReadConfigMap reads the configmap
func (l Langruntimes) ReadConfigMap(c kubernetes.Interface) {

	cfgm, err := c.CoreV1().ConfigMaps("kubeless").Get("kubeless-config", metav1.GetOptions{})

	if err != nil {
		logrus.Info("ERROR!!!! ------ ")
		return
	}

	// logrus.Info("ConfigMap Data", cfgm.Data)
	// logrus.Info("ConfigMap Data runtime images", cfgm.Data["runtime-images"])

	err = yaml.Unmarshal([]byte(cfgm.Data["runtime-images"]), &availableRuntimes)

	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("configmap ri is Versions: ", availableRuntimes[0])
	logrus.Info("configmap ri is Versions with ID: ", availableRuntimes[0].ID)
	logrus.Info("configmap ri is Versions with depname: ", availableRuntimes[0].DepName)
}

// // Runtimes struct
// type runtimes struct {
// 	ID             string            `json:"ID"`
// 	Versions       []runtimeVersions `json:"versions"`
// 	DepName        string            `json:"depname,omitempty"`
// 	FileNameSuffix string            `json:"filenamesuffix,omitempty"`
// }

// type runtimeVersions struct {
// 	Name        string `json:"name"`
// 	Version     string `json:"version"`
// 	HTTPImage   string `json:"httpImage"`
// 	PubSubImage string `json:"pubsubImage"`
// 	InitImage   string `json:"initImage"`
// }

const (
	python27Http    = "kubeless/python@sha256:0f3b64b654df5326198e481cd26e73ecccd905aae60810fc9baea4dcbb61f697"
	python27Pubsub  = "kubeless/python-event-consumer@sha256:1aeb6cef151222201abed6406694081db26fa2235d7ac128113dcebd8d73a6cb"
	python27Init    = "tuna/python-pillow:2.7.11-alpine" // TODO: Migrate the image for python 2.7 to an official source (not alpine-based)
	python34Http    = "kubeless/python@sha256:e502078dc9580bb73f823504a6765dfc98f000979445cdf071900350b938c292"
	python34Pubsub  = "kubeless/python-event-consumer@sha256:d963e4cd58229d662188d618cd87503b3c749b126b359ce724a19a375e4b3040"
	python34Init    = "python:3.4"
	python36Http    = "kubeless/python@sha256:6300c2513ca51653ae698a31eacf6b2b8a16d2737dd3e244a8c9c11f6408fd35"
	python36Pubsub  = "kubeless/python-event-consumer@sha256:0a2f9162de56b7966b02b70a5a0bcff03badfd9d87b8ae3d13e5381abd00220f"
	python36Init    = "python:3.6"
	node6Http       = "kubeless/nodejs@sha256:2b25d7380d6ed06ad817f4ee1e177340a282788596b34464173bb8a967d83c02"
	node6Pubsub     = "kubeless/nodejs-event-consumer@sha256:1861c32d6a46b2fdfc3e3996daf690ff2c3d5ca19a605abd2af503011d68e221"
	node6Init       = "node:6.10"
	node8Http       = "kubeless/nodejs@sha256:f1426efe274ea8480d95270c98f6007ac64645e36291dbfa36d759b5c8b7b733"
	node8Pubsub     = "kubeless/nodejs-event-consumer@sha256:b301b02e463b586d9a32d5c1cb5a68c2a11e4fba9514e28d900fc50a78759af9"
	node8Init       = "node:8"
	ruby24Http      = "kubeless/ruby@sha256:738e4cdeb5f5feece236bbf4e46902024e4b9fc16db4f3791404fa27e8b0db15"
	ruby24Pubsub    = "kubeless/ruby-event-consumer@sha256:f9f50be51d93a98ae30689d87b067c181905a8757d339fb0fa9a81c6268c4eea"
	ruby24Init      = "bitnami/ruby:2.4"
	dotnetcore2Http = "allantargino/kubeless-dotnetcore@sha256:d321dc4b2c420988d98cdaa22c733743e423f57d1153c89c2b99ff0d944e8a63"
	dotnetcore2Init = "microsoft/aspnetcore-build:2.0"
	pubsubFunc      = "PubSub"
)

type runtimeVersion struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	HTTPImage   string `yaml:"httpImage"`
	PubSubImage string `yaml:"pubsubImage"`
	InitImage   string `yaml:"initImage"`
	ImageSecret string `yaml:"imageSecret,omitempty"`
}

// RuntimeInfo describe the runtime specifics (typical file suffix and dependency file name)
// and the supported versions
type RuntimeInfo struct {
	ID             string           `yaml:"ID"`
	Versions       []runtimeVersion `yaml:"versions"`
	DepName        string           `yaml:"depName"`
	FileNameSuffix string           `yaml:"fileNameSuffix"`
}

// var pythonVersions, nodeVersions, rubyVersions, dotnetcoreVersions []runtimeVersion

// func init() {
// 	python27 := runtimeVersion{version: "2.7", httpImage: python27Http, pubsubImage: python27Pubsub, initImage: python27Init}
// 	python34 := runtimeVersion{version: "3.4", httpImage: python34Http, pubsubImage: python34Pubsub, initImage: python34Init}
// 	python36 := runtimeVersion{version: "3.6", httpImage: python36Http, pubsubImage: python36Pubsub, initImage: python36Init}
// 	pythonVersions = []runtimeVersion{python27, python34, python36}

// 	node6 := runtimeVersion{version: "6", httpImage: node6Http, pubsubImage: node6Pubsub, initImage: node6Init}
// 	node8 := runtimeVersion{version: "8", httpImage: node8Http, pubsubImage: node8Pubsub, initImage: node8Init}
// 	nodeVersions = []runtimeVersion{node6, node8}

// 	ruby24 := runtimeVersion{version: "2.4", httpImage: ruby24Http, pubsubImage: ruby24Pubsub, initImage: ruby24Init}
// 	rubyVersions = []runtimeVersion{ruby24}

// 	dotnetcore2 := runtimeVersion{version: "2.0", httpImage: dotnetcore2Http, pubsubImage: "", initImage: dotnetcore2Init}
// 	dotnetcoreVersions = []runtimeVersion{dotnetcore2}

// 	availableRuntimes = []RuntimeInfo{
// 		{ID: "python", versions: pythonVersions, DepName: "requirements.txt", FileNameSuffix: ".py"},
// 		{ID: "nodejs", versions: nodeVersions, DepName: "package.json", FileNameSuffix: ".js"},
// 		{ID: "ruby", versions: rubyVersions, DepName: "Gemfile", FileNameSuffix: ".rb"},
// 		{ID: "dotnetcore", versions: dotnetcoreVersions, DepName: "requirements.xml", FileNameSuffix: ".cs"},
// 	}
// }

// GetRuntimes returns the list of available runtimes as strings
func GetRuntimes() []string {
	result := []string{}
	for _, runtimeInf := range availableRuntimes {
		for _, runtime := range runtimeInf.Versions {
			result = append(result, runtimeInf.ID+runtime.Version)
		}
	}
	return result
}

// IsValidRuntime returns true if passed runtime name is valid runtime
func IsValidRuntime(runtime string) bool {
	for _, validRuntime := range GetRuntimes() {
		if runtime == validRuntime {
			return true
		}
	}
	return false
}

func getAvailableRuntimesPerTrigger(imageType string) []string {
	var runtimeList []string
	for i := range availableRuntimes {
		for j := range availableRuntimes[i].Versions {
			if (imageType == "PubSub" && availableRuntimes[i].Versions[j].PubSubImage != "") || (imageType == "HTTP" && availableRuntimes[i].Versions[j].HTTPImage != "") {
				runtimeList = append(runtimeList, availableRuntimes[i].ID+availableRuntimes[i].Versions[j].Version)
			}
		}
	}
	return runtimeList
}

// extract the branch number from the runtime string
func getVersionFromRuntime(runtime string) string {
	re := regexp.MustCompile("[0-9.]+$")
	return re.FindString(runtime)
}

// GetRuntimeInfo returns all the info regarding a runtime
func GetRuntimeInfo(runtime string) (RuntimeInfo, error) {
	runtimeID := regexp.MustCompile("^[a-zA-Z]+").FindString(runtime)
	logrus.Info("availableruntim GetRuntimeInfo: ", availableRuntimes)
	for _, runtimeInf := range availableRuntimes {
		logrus.Info("Runtim ID is %v and expected is %v", runtimeInf, runtimeID)
		if runtimeInf.ID == runtimeID {
			return runtimeInf, nil
		}
	}
	return RuntimeInfo{}, fmt.Errorf("Unable to find %s as runtime", runtime)
}

func findRuntimeVersion(runtimeWithVersion string) (runtimeVersion, error) {
	version := getVersionFromRuntime(runtimeWithVersion)
	runtimeInf, err := GetRuntimeInfo(runtimeWithVersion)
	if err != nil {
		return runtimeVersion{}, err
	}
	for _, versionInf := range runtimeInf.Versions {
		if versionInf.Version == version {
			return versionInf, nil
		}
	}
	return runtimeVersion{}, fmt.Errorf("The given runtime and version %s is not valid", runtimeWithVersion)
}

// GetFunctionImage returns the image ID depending on the runtime, its version and function type
func GetFunctionImage(runtime, ftype string) (string, error) {
	runtimeInf, err := GetRuntimeInfo(runtime)
	if err != nil {
		return "", err
	}

	imageNameEnvVar := ""
	if ftype == pubsubFunc {
		imageNameEnvVar = strings.ToUpper(runtimeInf.ID) + getVersionFromRuntime(runtime) + "_PUBSUB_RUNTIME"
	} else {
		imageNameEnvVar = strings.ToUpper(runtimeInf.ID) + getVersionFromRuntime(runtime) + "_RUNTIME"
	}
	imageName := os.Getenv(imageNameEnvVar)
	if imageName == "" {
		versionInf, err := findRuntimeVersion(runtime)
		if err != nil {
			return "", err
		}
		if ftype == pubsubFunc {
			if versionInf.PubSubImage == "" {
				err = fmt.Errorf("The given runtime and version '%s' does not have a valid image for event based functions. Available runtimes are: %s", runtime, strings.Join(getAvailableRuntimesPerTrigger("PubSub")[:], ", "))
			} else {
				imageName = versionInf.PubSubImage
			}
		} else {
			if versionInf.HTTPImage == "" {
				err = fmt.Errorf("The given runtime and version '%s' does not have a valid image for HTTP based functions. Available runtimes are: %s", runtime, strings.Join(getAvailableRuntimesPerTrigger("HTTP")[:], ", "))
			} else {
				imageName = versionInf.HTTPImage
			}
		}
	}
	return imageName, nil
}

// GetBuildContainer returns a Container definition based on a runtime
func GetBuildContainer(runtime string, env []v1.EnvVar, installVolume v1.VolumeMount) (v1.Container, error) {
	runtimeInf, err := GetRuntimeInfo(runtime)
	if err != nil {
		return v1.Container{}, err
	}
	depsFile := path.Join(installVolume.MountPath, runtimeInf.DepName)
	versionInf, err := findRuntimeVersion(runtime)
	if err != nil {
		return v1.Container{}, err
	}

	var command string
	switch {
	case strings.Contains(runtime, "python"):
		command = "pip install --prefix=" + installVolume.MountPath + " -r " + depsFile
	case strings.Contains(runtime, "nodejs"):
		registry := "https://registry.npmjs.org"
		scope := ""
		for _, v := range env {
			if v.Name == "NPM_REGISTRY" {
				registry = v.Value
			}
			if v.Name == "NPM_SCOPE" {
				scope = v.Value + ":"
			}
		}
		command = "npm config set " + scope + "registry " + registry +
			" && npm install --prefix=" + installVolume.MountPath
	case strings.Contains(runtime, "ruby"):
		command = "bundle install --gemfile=" + depsFile + " --path=" + installVolume.MountPath
	}

	return v1.Container{
		Name:            "install",
		Image:           versionInf.InitImage,
		Command:         []string{"sh", "-c"},
		Args:            []string{command},
		VolumeMounts:    []v1.VolumeMount{installVolume},
		ImagePullPolicy: v1.PullIfNotPresent,
		WorkingDir:      installVolume.MountPath,
		Env:             env,
	}, nil
}

// UpdateDeployment object in case of custom runtime
func UpdateDeployment(dpm *v1beta2.Deployment, depsPath, runtime string) {
	switch {
	case strings.Contains(runtime, "python"):
		dpm.Spec.Template.Spec.Containers[0].Env = append(dpm.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
			Name:  "PYTHONPATH",
			Value: path.Join(depsPath, "lib/python"+getVersionFromRuntime(runtime)+"/site-packages"),
		})
	case strings.Contains(runtime, "nodejs"):
		dpm.Spec.Template.Spec.Containers[0].Env = append(dpm.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
			Name:  "NODE_PATH",
			Value: path.Join(depsPath, "node_modules"),
		})
	case strings.Contains(runtime, "ruby"):
		dpm.Spec.Template.Spec.Containers[0].Env = append(dpm.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
			Name:  "GEM_HOME",
			Value: path.Join(depsPath, "ruby/2.4.0"),
		})
	case strings.Contains(runtime, "dotnetcore"):
		dpm.Spec.Template.Spec.Containers[0].Env = append(dpm.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
			Name:  "DOTNETCORE_HOME",
			Value: "/usr/bin/",
		})
	}
}
