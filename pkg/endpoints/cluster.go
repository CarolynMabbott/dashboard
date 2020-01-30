/*
Copyright 2019 The Tekton Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package endpoints

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/tektoncd/dashboard/pkg/logging"
	"github.com/tektoncd/dashboard/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Properties : properties we want to be able to retrieve via REST
type Properties struct {
	InstallNamespace string
}

const (
	tektonDashboardIngressName string = "tekton-dashboard"
	tektonDashboardRouteName   string = "tekton-dashboard"
)

// ProxyRequest does as the name suggests: proxies requests and logs what's going on
func (r Resource) ProxyRequest(request *restful.Request, response *restful.Response) {
	parsedURL, err := url.Parse(request.Request.URL.String())
	if err != nil {
		utils.RespondError(response, err, http.StatusNotFound)
		return
	}
	uri := request.PathParameter("subpath") + "?" + parsedURL.RawQuery
	forwardRequest := r.K8sClient.CoreV1().RESTClient().Verb(request.Request.Method).RequestURI(uri).Body(request.Request.Body)
	forwardRequest.SetHeader("Content-Type", request.HeaderParameter("Content-Type"))
	forwardResponse := forwardRequest.Do()

	responseBody, requestError := forwardResponse.Raw()
	if requestError != nil {
		utils.RespondError(response, requestError, http.StatusNotFound)
		return
	}
	response.Header().Add("Content-Type", utils.GetContentType(responseBody))
	response.Write(responseBody)
}

// GetIngress returns the Ingress endpoint called "tektonDashboardIngressName" in the requested namespace
func (r Resource) GetIngress(request *restful.Request, response *restful.Response) {
	requestNamespace := utils.GetNamespace(request)

	ingress, err := r.K8sClient.ExtensionsV1beta1().Ingresses(requestNamespace).Get(tektonDashboardIngressName, metav1.GetOptions{})

	if err != nil || ingress == nil {
		logging.Log.Errorf("Unable to retrieve any ingresses: %s", err)
		utils.RespondError(response, err, http.StatusInternalServerError)
		return
	}

	noRuleError := "no Ingress rules found labelled " + tektonDashboardIngressName

	// Harden this block to avoid Go panics (array index out of range)
	if len(ingress.Spec.Rules) > 0 { // Got more than zero entries?
		if ingress.Spec.Rules[0].Host != "" { // For that rule, is there actually a host?
			ingressHost := ingress.Spec.Rules[0].Host
			response.WriteEntity(ingressHost)
			return
		}
		logging.Log.Errorf("found an empty Ingress rule labelled %s", tektonDashboardIngressName)
	} else {
		logging.Log.Error(noRuleError)
	}

	logging.Log.Error("Unable to retrieve any Ingresses")
	utils.RespondError(response, err, http.StatusInternalServerError)
	return
}

// GetIngress returns the Ingress endpoint called "tektonDashboardIngressName" in the requested namespace
func (r Resource) GetEndpoints(request *restful.Request, response *restful.Response) {
	type element struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	}
	var responses []element
	requestNamespace := utils.GetNamespace(request)

	route, err := r.RouteClient.RouteV1().Routes(requestNamespace).Get(tektonDashboardIngressName, metav1.GetOptions{})
	noRuleError := "no Route found labelled " + tektonDashboardRouteName
	if err != nil || route == nil {
		logging.Log.Infof("Unable to retrieve any routes: %s", err)
	} else {
		if route.Spec.Host != "" { // For that rule, is there actually a host?
			routeHost := route.Spec.Host
			responses = append(responses, element{"Route", routeHost})
		} else {
			logging.Log.Error(noRuleError)
		}
	}

	ingress, err := r.K8sClient.ExtensionsV1beta1().Ingresses(requestNamespace).Get(tektonDashboardIngressName, metav1.GetOptions{})
	noRuleError = "no Ingress rules found labelled " + tektonDashboardIngressName
	if err != nil || ingress == nil {
		logging.Log.Infof("Unable to retrieve any ingresses: %s", err)
	} else {
		if len(ingress.Spec.Rules) > 0 { // Got more than zero entries?
			if ingress.Spec.Rules[0].Host != "" { // For that rule, is there actually a host?
				ingressHost := ingress.Spec.Rules[0].Host
				responses = append(responses, element{"Ingress", ingressHost})
			}
		} else {
			logging.Log.Error(noRuleError)
		}
	}

	if len(responses) != 0 {
		response.WriteEntity(responses)
	} else {
		logging.Log.Error("Unable to retrieve any Ingresses or Routes")
		utils.RespondError(response, err, http.StatusInternalServerError)
	}
}

// GetProperties is used to get the installed namespace only so far
func (r Resource) GetProperties(request *restful.Request, response *restful.Response) {
	properties := Properties{InstallNamespace: os.Getenv("INSTALLED_NAMESPACE")}
	response.WriteEntity(properties)
}

// Get dashboard version
func (r Resource) GetDashboardVersion(request *restful.Request, response *restful.Response) {
	config, err := rest.InClusterConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	properties := Properties{InstallNamespace: os.Getenv("INSTALLED_NAMESPACE")}
	version := ""

	api1 := clientset.AppsV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "app=tekton-dashboard",
	}
	deployments, err := api1.Deployments(properties.InstallNamespace).List(listOptions)

	for _, deployment := range deployments.Items {
		deploymentLabels := deployment.GetLabels()
		version = deploymentLabels["version"]
	}

	if version == "" {
		response.WriteEntity("Unknown")
		return
	}
	response.WriteEntity(version)
	return
}

// Get pipelines version
func (r Resource) GetPipelineVersion(request *restful.Request, response *restful.Response) {
	config, err := rest.InClusterConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	version := ""

	api1 := clientset.AppsV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/component=controller,app.kubernetes.io/name=tekton-pipelines",
	}
	deployments, err := api1.Deployments("tekton-pipelines").List(listOptions)

	for _, deployment := range deployments.Items {
		deploymentAnnotations := deployment.Spec.Template.GetAnnotations()
		version = deploymentAnnotations["tekton.dev/release"]

		if version == "" {
			deploymentImage := deployment.Spec.Template.Spec.Containers[0].Image
			if strings.Contains(deploymentImage, "pipeline/cmd/controller") {
				s := strings.SplitAfter(deploymentImage, ":")
				t := strings.Split(s[1], "@")
				version = t[0]
			}
		}
	}

	if version == "" {
		response.WriteEntity("Unknown")
		return
	}
	response.WriteEntity(version)
	return
}
