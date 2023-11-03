{{- define "controlPlaneTolerations" }}
- effect: NoSchedule
  key: node-role.kubernetes.io/control-plane
{{- end }}

{{- define "preferWorkerNodes" }}
- weight: 1
  preference:
    matchExpressions:
    - key: node-role.kubernetes.io/control-plane
      operator: DoesNotExist
{{- end }}