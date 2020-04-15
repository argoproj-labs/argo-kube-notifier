
# Argo-Kube-Notifier - Monitor kubernetes resources and get notification
## What is Argo Kube Notifier?
Argo Kube Notifier controller monitors Kubernetes resource updates and validates the configured rules to trigger events.  

## Why use Argo Kube Notifier?
Argo Kube Notifier is a fully customizable monitoring tool based on all kinds of Kubernetes resources. 
Argo Kube Notifier will monitor each resource's updates and apply rules on updates. Events will be trigged if a rule is met.
 
 


## Quickstart
```bash
kubectl create namespace argo
kubectl create -f https://raw.githubusercontent.com/argoproj-labs/argo-kube-notifier/master/config/crds/argoproj_v1alpha1_notification.yaml
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj-labs/argo-kube-notifier/master/config/manager/manager.yaml
```

## Create Notification

example: [notification](https://raw.githubusercontent.com/sarabala1979/argo-kube-notifier/master/example/notification_deployment.yaml)

Create secrets for Slack and Email
```
apiVersion: v1
kind: Secret
metadata:
  name: my-slack-secret
type: Opaque
stringData:
  hookURL: "https://hooks.slack.com/" #update your slack hookurl
```


#### Simple Slack Notification with All Conditions
```
apiVersion: argoproj.io/v1alpha1
kind: Notification
metadata:
  name: notification-1
spec:
  # Define resource need to be monitor.
  monitorResource:
    Group:    "apps"
    Resource: "deployments"
    Version:  "v1"
  #
  Namespace: default
  # Define the notifiers channel. Slack and email is supporting currently
  notifiers:
    -
      name: slack
      slack:
        hookUrlSecret:
          name: my-slack-secret
          key: hookURL
        channel: testargonotification
        
        hookurl: "https://hooks.slack.com"
  # Define the Rules to validate the resource spce
  rules:
    -
      allConditions:
        -
          jsonPath: "status/readyReplicas"
          operator: "ne"
          valueJsonPath: "status/replicas"

      events:
        -
          message: "Condition Triggered : Deployment = {{.metadata.name}} replicaset doesnt match. Required Replicas = {{.status.replicas}}, Current Replicas={{.status.readyReplicas}}"
          emailSubject: "Argo Notification Condition Triggered {{.metadata.name}}"
          notificationLevel: "warning"
          notifierNames:
            - "slack"

      name: rule1
      initialDelaySec: 60
      throttleMinutes: 5
    -
      allConditions:
        -
          jsonPath: "metadata/name"
          operator: "eq"
          value: "nginx-deployment"
        -
          jsonPath: "{{status/readyReplicas}} / {{ }}"
          operator: "lt"
          value: "30%"

      events:
        -
          message: "Condition Triggered : Deployment ={{.metadata.name}} replicaset doesnt match. Required Replicas = {{.status.replicas}}, Current Replicas={{.status.readyReplicas}}"
          emailSubject: "[ALERT] Argo Notification Condition Triggered {{.metadata.name}}"
          notificationLevel: "critical"
          notifierNames:
            - "slack"

      name: rule2
      initialDelaySec: 60
      throttleMinutes: 5
```
#### Notification 2 - Any Conditions with Slack and Email
```
apiVersion: argoproj.io/v1alpha1
kind: Notification
metadata:
  name: notification-2
spec:
  # Define resource need to be monitor.
  monitorResource:
    Group: argoproj.io
    Resource: workflows
    Version: v1alpha1
  #
  Namespace: default
  # Define the notifiers channel. Slack and email is supporting currently
    rules:
      -
        name: rule1
        initialDelaySec: 60
        throttleMinutes: 5
        anyConditions:
          -
            jsonPath: "status/phase"
            operator: "eq"
            value: "Failed"
          -
            jsonPath: "status/phase"
            operator: "eq"
            value: "Error"
        events:
          -
            message: "Workflow ={{.metadata.name}} Failed."
            emailSubject: "[ALERT] Argo Notification Condition Triggered {{.metadata.name}}"
            notificationLevel: "warning"
            notifierNames:
              - "slack"
              - "email"

    notifiers:
      -
       name: slack
       slack:
         hookUrlSecret:
           name: my-slack-secret
           key: hookURL
         channel: testargonotification
      -
        name: "email"
        email:
          smtphost: "smtp.gmail.com"
          smtpport: 465
          # Save the smtp service credential in secret in monitoring namespace
          usernameSecret:
            name: "my-email-secret"
            key: "username"
          passwordSecret:
            name: "my-email-secret"
            key: "password"
          FromEmailId: <EmailID>
          SenderList:
            - <EmailID>
            - <EmailID>
```
