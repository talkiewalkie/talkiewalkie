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

// ---------------
// ----------- LOG
// ---------------

resource "aws_iam_policy" "logging_writer" {
  name        = "task-policy-logging-writer"
  description = "Policy that allows to create logging events"

  policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
               "logs:*"
           ],
           "Resource": "*"
       }
   ]
}
EOF
}


// -------------------
// ----------- SECRETS
// -------------------

resource "aws_iam_policy" "secrets_fetcher" {
  name        = "task-policy-secrets-fetcher"
  description = "Policy that allows to query secrets"

  // TODO: `secretsmanager:GetSecretValue` was not enough to pull secrets
  //        somehow, would be good to restrict to accurate set.
  policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
               "secretsmanager:*"
           ],
           "Resource": "*"
       }
   ]
}
EOF
}