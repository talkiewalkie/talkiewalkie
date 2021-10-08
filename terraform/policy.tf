resource "aws_iam_policy" "ecr" {
  name        = "task-policy-ecr"
  description = "Policy that allows access to ECR"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ecr:*",
                "cloudtrail:LookupEvents"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "iam:CreateServiceLinkedRole"
            ],
            "Resource": "*",
            "Condition": {
                "StringEquals": {
                    "iam:AWSServiceName": [
                        "replication.ecr.amazonaws.com"
                    ]
                }
            }
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecr-task-role-policy-attachment" {
  role       = aws_iam_role.nginx_task_role.name
  policy_arn = aws_iam_policy.ecr.arn
}


// --------------
// ----------- S3
// --------------

resource "aws_iam_policy" "s3" {
  name        = "task-policy-s3"
  description = "Policy that allows access to S3"

  policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
               "dynamodb:CreateTable",
               "dynamodb:UpdateTimeToLive",
               "dynamodb:PutItem",
               "dynamodb:DescribeTable",
               "dynamodb:ListTables",
               "dynamodb:DeleteItem",
               "dynamodb:GetItem",
               "dynamodb:Scan",
               "dynamodb:Query",
               "dynamodb:UpdateItem",
               "dynamodb:UpdateTable"
           ],
           "Resource": "*"
       }
   ]
}
EOF
}