package util

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const OAMServerTemplate = `
apiVersion: core.oam.dev/v1alpha1
kind: ComponentSchematic
metadata:
  name: {{.Deployment.Name}}
  namespace: {{.Deployment.Namespace}}
spec:
  workloadType: core.oam.dev/v1alpha1.Server
  containers:
{{- range $idx, $ele := .Deployment.Spec.Template.Spec.Containers}}
    - name: {{.Name}}
      image: {{.Image}}
{{- with .Command}}
      cmd:
{{- range $index, $element := .}}
        - {{$element}}
{{- end}}
{{- end}}
{{- with .Args}}
      args:
{{- range $index, $element := .}}
        - {{$element}}
{{- end}}
{{- end}}
{{- with .Env}}
      env:
{{- range $index, $element := .}}
        - name: {{$element.Name}}
          fromParam: {{$ele.Name}}-{{$element.Name}}
{{- end}}
{{- end}}
      ports:
{{- range $index, $element := .Ports}}
        - name: {{$element.Name}}
          containerPort: {{$element.ContainerPort}}
          protocol: {{$element.Protocol}}
{{- end}}
      {{- with .LivenessProbe}}
      livenessProbe:
        {{with .Exec -}}
        exec: 
         {{- range $index, $element := .Command}}
          - {{$element}}
         {{- end}}
        {{- end}}
        {{with .HTTPGet -}}
        httpGet:
          path: {{.Path}}
          port: {{json .Port | printf}}
          {{with .HTTPHeaders -}}
          httpHeaders:
            {{- range $index, $element := .}}
            name: {{$element.Name}}
            value: {{$element.Value}}
            {{- end}}
          {{- end}}
        {{- end}}
        {{with .TCPSocket -}}
        tcpSocket:
          port: {{json .Port | printf}}
        {{- end}}
        initialDelaySeconds: {{.InitialDelaySeconds}}
        periodSeconds: {{.TimeoutSeconds}}
        timeoutSeconds: {{.PeriodSeconds}}
        successThreshold: {{.SuccessThreshold}}
        failureThreshold: {{.FailureThreshold}}
      {{- end}}
      {{- with .ReadinessProbe}}
      readinessProbe:
        {{with .Exec -}}
        exec: 
         {{- range $index, $element := .Command}}
          - {{$element}}
         {{- end}}
        {{- end}}
        {{with .HTTPGet -}}
        httpGet:
          path: {{.Path}}
          port: {{json .Port | printf}}
          {{with .HTTPHeaders -}}
          httpHeaders:
            {{- range $index, $element := .}}
            name: {{$element.Name}}
            value: {{$element.Value}}
            {{- end}}
          {{- end}}
        {{- end}}
        {{with .TCPSocket -}}
        tcpSocket:
          port: {{json .Port | printf}}
        {{- end}}
        initialDelaySeconds: {{.InitialDelaySeconds}}
        periodSeconds: {{.TimeoutSeconds}}
        timeoutSeconds: {{.PeriodSeconds}}
        successThreshold: {{.SuccessThreshold}}
        failureThreshold: {{.FailureThreshold}}
      {{- end}}
{{- end}}
  parameters:
{{- range $index, $element := .Deployment.Spec.Template.Spec.Containers}}
{{- range $sub_index, $sub_element := $element.Env }}
    - name: {{$element.Name}}-{{$sub_element.Name}}
      value: {{$sub_element.Value}}
{{- end}}
{{- end}}
---
apiVersion: core.oam.dev/v1alpha1
kind: ApplicationConfiguration
metadata:
  name: {{.Name}}
  namespace: {{.Deployment.Namespace}}
spec:
  components:
    {{with .Deployment -}}
    - name: {{.Name}}
      instanceName: {{.Name}}_instance
      parameterValues:
		{{- range $index, $element := .Spec.Template.Spec.Containers}}
		{{- range $sub_index, $sub_element := $element.Env }}
        - name: {{$element.Name}}-{{$sub_element.Name}}
          value: {{$sub_element.Value}}
		{{- end}}
		{{- end}}
   {{- end}}
      traits:
        {{with .Deployment -}}
        - name: manual-scaler
          parameterValues:
            - name: replicaCount
              value: {{.Spec.Replicas}}
        {{- end}}
      {{- if and .Ingress.UID .Service.UID }}
        - name: ingress
          parameterValues:
         {{- range $index, $element := .Ingress.Spec.Rules }}
            - name: hostname
              value: {{.Host}}
            - name: path
              value: {{(index .HTTP.Paths 0).Path}}
            - name: service_port
              value: {{json (index .HTTP.Paths 0).Backend.ServicePort | print}}
      {{- end}}
      {{- end -}}
`

const OAMWorkerTemplate = `
apiVersion: core.oam.dev/v1alpha1
kind: ComponentSchematic
metadata:
  name: {{.Deployment.Name}}
  namespace: {{.Deployment.Namespace}}
spec:
  workloadType: core.oam.dev/v1alpha1.Server
  containers:
{{- range $idx, $ele := .Deployment.Spec.Template.Spec.Containers}}
    - name: {{.Name}}
      image: {{.Image}}
{{- with .Command}}
      cmd:
{{- range $index, $element := .}}
        - {{$element}}
{{- end}}
{{- end}}
{{- with .Args}}
      args:
{{- range $index, $element := .}}
        - {{$element}}
{{- end}}
{{- end}}
{{- with .Env}}
      env:
{{- range $index, $element := .}}
        - name: {{$element.Name}}
          fromParam: {{$ele.Name}}-{{$element.Name}}
{{- end}}
{{- end}}
      {{- with .LivenessProbe}}
      livenessProbe:
        {{with .Exec -}}
        exec: 
         {{- range $index, $element := .Command}}
          - {{$element}}
         {{- end}}
        {{- end}}
        {{with .HTTPGet -}}
        httpGet:
          path: {{.Path}}
          port: {{json .Port | printf}}
          {{with .HTTPHeaders -}}
          httpHeaders:
            {{- range $index, $element := .}}
            name: {{$element.Name}}
            value: {{$element.Value}}
            {{- end}}
          {{- end}}
        {{- end}}
        {{with .TCPSocket -}}
        tcpSocket:
          port: {{json .Port | printf}}
        {{- end}}
        initialDelaySeconds: {{.InitialDelaySeconds}}
        periodSeconds: {{.TimeoutSeconds}}
        timeoutSeconds: {{.PeriodSeconds}}
        successThreshold: {{.SuccessThreshold}}
        failureThreshold: {{.FailureThreshold}}
      {{- end}}
      {{- with .ReadinessProbe}}
      readinessProbe:
        {{with .Exec -}}
        exec: 
         {{- range $index, $element := .Command}}
          - {{$element}}
         {{- end}}
        {{- end}}
        {{with .HTTPGet -}}
        httpGet:
          path: {{.Path}}
          port: {{json .Port | printf}}
          {{with .HTTPHeaders -}}
          httpHeaders:
            {{- range $index, $element := .}}
            name: {{$element.Name}}
            value: {{$element.Value}}
            {{- end}}
          {{- end}}
        {{- end}}
        {{with .TCPSocket -}}
        tcpSocket:
          port: {{json .Port | printf}}
        {{- end}}
        initialDelaySeconds: {{.InitialDelaySeconds}}
        periodSeconds: {{.TimeoutSeconds}}
        timeoutSeconds: {{.PeriodSeconds}}
        successThreshold: {{.SuccessThreshold}}
        failureThreshold: {{.FailureThreshold}}
      {{- end}}
{{- end}}
    parameters:
{{- range $index, $element := .Deployment.Spec.Template.Spec.Containers}}
{{- range $sub_index, $sub_element := $element.Env }}
    - name: {{$element.Name}}-{{$sub_element.Name}}
      value: {{$sub_element.Value}}
{{- end}}
{{- end}}
---
apiVersion: core.oam.dev/v1alpha1
kind: ApplicationConfiguration
metadata:
  name: {{.Name}}
  namespace: {{.Deployment.Namespace}}
spec:
  components:
    {{with .Deployment -}}
    - name: {{.Name}}
      instanceName: {{.Name}}_instance
      parameterValues:
		{{- range $index, $element := .Spec.Template.Spec.Containers}}
		{{- range $sub_index, $sub_element := $element.Env }}
        - name: {{$element.Name}}-{{$sub_element.Name}}
          value: {{$sub_element.Value}}
		{{- end}}
		{{- end}}
   {{- end}}
      traits:
        {{with .Deployment -}}
        - name: manual-scaler
          parameterValues:
            - name: replicaCount
              value: {{.Spec.Replicas}}
        {{- end}}
`

var (
	serverTPL, workerTPL *template.Template
)

func init() {
	if t, err := template.New("OAMServer").
		Funcs(template.FuncMap{"json": genJson}).
		Parse(string(OAMServerTemplate)); err != nil {
		panic(err)
	} else {
		serverTPL = t
	}
	if t, err := template.New("OAMWorker").
		Funcs(template.FuncMap{"json": genJson}).
		Parse(string(OAMWorkerTemplate)); err != nil {
		panic(err)
	} else {
		workerTPL = t
	}
}

func genJson(v interface{}) string {
	jsonBytes, _ := json.Marshal(v)
	return string(jsonBytes)
}

func RenderServer(params interface{}) (string, error) {
	var buf bytes.Buffer
	err := serverTPL.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderWorker(params interface{}) (string, error) {
	var buf bytes.Buffer
	err := workerTPL.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
