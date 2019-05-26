package controller

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/argoproj/argo-kube-notifier/pkg/apis/argoproj/v1alpha1"
	"strings"
	"testing"
)

var jsonStr = "{\"apiVersion\":\"argoproj.io/v1alpha1\",\"kind\":\"Workflow\",\"metadata\":{\"clusterName\":\"\",\"creationTimestamp\":\"2019-05-07T20:34:24Z\"," +
	"\"generateName\":\"output-artifact-s3-\",\"generation\":1,\"labels\":{\"workflows.argoproj.io/completed\":\"true\",\"workflows.argoproj.io/phase\":\"Error\"}," +
	"\"name\":\"output-artifact-s3-85lbf\",\"namespace\":\"default\",\"resourceVersion\":\"1876992\",\"selfLink\":\"/apis/argoproj.io/v1alpha1/namespaces/default/workflows/output-artifact-s3-85lbf\"," +
	"\"uid\":\"8085d32a-7107-11e9-9102-025000000001\"},\"spec\":{\"arguments\":{},\"entrypoint\":\"whalesay\",\"serviceAccountName\":\"argo-build\",\"templates\":[{\"container\":{\"args\":[\"cowsay " +
	"hello world | tee /tmp/hello_world.txt\"],\"command\":[\"sh\",\"-c\"],\"image\":\"docker/whalesay:latest\",\"name\":\"\",\"resources\":{}},\"inputs\":{},\"metadata\":{},\"name\":\"whalesay\"," +
	"\"outputs\":{\"artifacts\":[{\"name\":\"message\",\"path\":\"/tmp\",\"s3\":{\"accessKeySecret\":{\"key\":\"\"},\"bucket\":\"\",\"endpoint\":\"\",\"key\":\"hello_world.txt.tgz\",\"secretKeySecret\":" +
	"{\"key\":\"\"}}}]}}]},\"status\":{\"finishedAt\":\"2019-05-07T20:35:37Z\",\"message\":\"pods  not found\",\"nodes\":{\"output-artifact-s3-85lbf\":{\"displayName\":\"output-artifact-s3-85lbf\",\"finishedAt\":" +
	"\"2019-05-07T20:35:37Z\",\"id\":\"output-artifact-s3-85lbf\",\"message\":\"pods not found\",\"name\":\"output-artifact-s3-85lbf\",\"phase\":\"Error\",\"startedAt\":\"2019-05-07T20:35:37Z\",\"templateName\":" +
	"\"whalesay\",\"type\":\"Pod\"}},\"phase\":\"Error\",\"startedAt\":\"2019-05-07T20:34:24Z\"}}"

func TestRuleValidaion(t *testing.T) {

	condition := &v1alpha1.Condition{
		Jsonpath: "status/phase",
		Operator: "eq",
		Value:    "Error",
		ChildConditions: []v1alpha1.Condition{
			{
				Jsonpath: "status/nodes/*/phase",
				Operator: "eq",
				Value:    "Error",
			},
		},
	}
	doc, _ := jsonquery.Parse(strings.NewReader(jsonStr))
	status := ValidateCondition(condition, doc)

	fmt.Println(status)

	//assert.Nil(t, err)
	//assert.Equal()
}

func TestValidateRule(t *testing.T) {

	rule := v1alpha1.Rule{
		AllConditions: []v1alpha1.Condition{
			{
				Jsonpath: "status/phase",
				Operator: "eq",
				Value:    "Error",
			},
			{
				Jsonpath:     "status/nodes",
				JoinOperator: "and",
				ChildConditions: []v1alpha1.Condition{
					{
						Jsonpath:     "//id",
						Operator:     "eq",
						Value:        "output-artifact-s3-85lbf",
						JoinOperator: "and",
					},
					{
						Jsonpath: "//phase",
						Operator: "eq",
						Value:    "Error",
					},
				},
			},
		},
	}
	doc, _ := jsonquery.Parse(strings.NewReader(jsonStr))
	status := ValidateRule(&rule, doc)
	fmt.Println(status)
}

//func TestNestedRuleValidaion(t *testing.T) {
//
//	condition := &config.NestedCondition{
//		RootJsonpath: "status/nodes",
//		NotificationLevel: "good",
//		AllConditions: []config.Condition{
//			{
//				Jsonpath:"//phase",
//				Operator:"eq",
//				Value:"Error",
//
//			},
//			{
//				Jsonpath:"//id",
//				Operator:"eq",
//				Value:"output-artifact-s3-85lbf",
//
//			},
//		},
//	}
//	doc, _ := jsonquery.Parse(strings.NewReader(jsonStr))
//	status := ValidateNestedCondition(condition, doc )
//
//	fmt.Println(status)
//
//	//assert.Nil(t, err)
//	//assert.Equal()
//}
