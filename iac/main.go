package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var projectTag = "aws-billing-lambda"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		resourceName := fmt.Sprintf("%s-iam-role-for-lambda", projectTag)
		iamLambda, err := iam.NewRole(
			ctx,
			resourceName,
			&iam.RoleArgs{
				AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
						"Effect": "Allow",
						"Principal": {
							"Service": "lambda.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}]
				}`),
				Tags: createNameTag(resourceName),
			})
		if err != nil {
			return fmt.Errorf("failed create iam role for lambda: %v", err)
		}
		ctx.Export(resourceName, iamLambda.ID())

		resourceName = fmt.Sprintf("%s-iam-policy-for-lambda", projectTag)
		iamLambdaPolicy, err := iam.NewRolePolicy(
			ctx,
			resourceName,
			&iam.RolePolicyArgs{
				Role: iamLambda.Name,
				Policy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "CloudWatchLogsPermissions",
						"Effect": "Allow",
						"Action": [
							"logs:CreateLogGroup",
							"logs:CreateLogStream",
							"logs:PutLogEvents"
						],
						"Resource": "*"
					},
					{
						"Sid": "CostExplorerPermissions",
						"Effect": "Allow",
						"Action": [
							"ce:GetCostAndUsage",
							"ce:GetDimensionValues"
						],
						"Resource": "*"
					}
				]
			}`),
			})
		if err != nil {
			return fmt.Errorf("failed create iam policy for lambda: %v", err)
		}
		ctx.Export(resourceName, iamLambdaPolicy.ID())

		resourceName = fmt.Sprintf("%s-lambda", projectTag)
		lambdaFunc, err := lambda.NewFunction(
			ctx,
			resourceName,
			&lambda.FunctionArgs{
				Runtime: pulumi.String("go1.x"),
				Handler: pulumi.String("handler"),
				Role:    iamLambda.Arn,
				Code:    pulumi.NewFileArchive("../src/handler.zip"), // TODO: 修正
				Environment: &lambda.FunctionEnvironmentArgs{
					Variables: pulumi.StringMap{
						"SLACK_WEBHOOK_URL": pulumi.String(""),
					},
				},
			})
		if err != nil {
			return fmt.Errorf("failed create iam policy for lambda: %v", err)
		}
		ctx.Export(resourceName, lambdaFunc.ID())

		// TODO: EventBridge

		return nil
	})
}

func createNameTag(tag string) pulumi.StringMap {
	return pulumi.StringMap{
		"Name": pulumi.String(tag),
	}
}
