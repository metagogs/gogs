apiVersion: v1
kind: Pod
metadata:
  name: {{.PodName}}
  namespace: {{.PodNamespace}}
  labels:
    app: {{.PodName}}
spec:
  containers:
    - name: gogs-main
      image: gogs-image
      ports:
      {{range .Deployments }}
        - name: {{.Name}}
          containerPort: {{.Port}}
          protocol: {{.Protocol}}
          {{if not .Svc}}hostPort: {{.Port}}{{end}}
      {{end}}
      resources:
        limits:
          cpu: 500m
          memory: 500Mi
        requests:
          cpu: 100m
          memory: 100Mi
      readinessProbe:
        httpGet:
          path: /health
          port: {{.HealthPort}}
          scheme: HTTP
        initialDelaySeconds: 5
        timeoutSeconds: 1
        periodSeconds: 5
        successThreshold: 1
        failureThreshold: 3
      imagePullPolicy: IfNotPresent
  restartPolicy: Always

{{if .Svc}}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: {{.PodName}}
  name: {{.PodName}}-svc
  namespace: {{.PodNamespace}}
spec:
  ports:
  {{range .Deployments }}
  - port: {{.Port}}
    targetPort: {{.Port}}
    protocol: {{.Protocol}}
  {{end}}
  selector:
    app: {{.PodName}}
  type: ClusterIP

{{end}}