
# Argo-Kube-Notifier - Monitor kubernetes resources and get notification

Argo-kube-notifier will monitor the k8s resources changes and apply the rules for notification.



## Quickstart
```bash
kubectl create namespace argo
kubectl create -f https://raw.githubusercontent.com/sarabala1979/argo-kube-notifier/master/config/crds/argoproj_v1alpha1_notification.yaml
kubectl apply -n argo -f https://raw.githubusercontent.com/sarabala1979/argo-kube-notifier/master/config/manager/manager.yaml
```

#Create notification

example: [notification](https://raw.githubusercontent.com/sarabala1979/argo-kube-notifier/master/example/notification_deployment.yaml)

#### Simple Slack Notification with All Condition
```
apiVersion: argoproj.io/v1alpha1
kind: Notification
metadata:
  name: notification-1
spec:
  # Define resource need to be monitor.
  monitorResource:
    Group: argoproj.io
    Resource: workflows
    Version: v1alpha1
  #
  Namespace: default
  # Define the notifiers channel. Slack and email is supporting currently
  notifiers:
    -
      name: slack
      slack:
        channel: testargonotification
        hookurl: "https://hooks.slack.com"
  # Define the Rules to validate the resource spce
  rules:
    -
      # Define multiple condition with AND
      allConditions:
        -
          # Json Path for element need to validate
          jsonPath: status/phase
          # Operator
          operator: eq
          # value
          value: Succeeded
      # Define the Events which will trigger if above conditions met
      events:
        -
          emailSubject: "Argo Notification Condition Triggered {{.metadata.name}}"
          message: "Workflow = {{.metadata.name}} succesfully completed"
          notificationLevel: info
          notifierNames:
            - slack
      # Unique Rule Name per Notification
      name: rule1
```
#### Notification 2- Any Conditions with Slack and Email
```$xslt
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
        anyconditions:
          -
            jsonpath: "status/phase"
            operator: "eq"
            value: "Failed"
          -
            jsonpath: "status/phase"
            operator: "eq"
            value: "Error"
        events:
          -
            message: "Workflow ={{.metadata.name}} Failed."
            emailSubject: "[ALERT] Argo Notification Condition Triggered {{.metadata.name}}"
            notificationLevel: "critical"
            notifierNames:
              - "slack"
              - "email"

    notifiers:
      -
        name: "slack"
        slack:
          hookurl: "https://hooks.slack.com/services/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"
          channel: "demokubenotifier"
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