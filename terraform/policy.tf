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
  for_each   = toset([
    aws_iam_role.nginx_task_role.name,
    aws_iam_role.nginx_task_execution_role.name
  ])
  role       = each.value
  policy_arn = aws_iam_policy.ecr.arn
}

// --------------
// ----------- S3
// --------------

resource "aws_iam_policy" "s3" {
  name        = "task-policy-s3-full"
  description = "Policy that allows full access to S3"

  policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
               "s3:*"
           ],
           "Resource": "*"
       }
   ]
}
EOF
}