# Deployment

## Live URL

http://pack-calculator-alb-1950116497.eu-central-1.elb.amazonaws.com

## How It Works

1. Push code to `main` branch
2. GitHub Actions runs tests
3. Docker image is built and pushed to AWS ECR
4. ECS Fargate service is updated with new image
5. Application Load Balancer routes traffic to healthy containers

## AWS Resources

| Resource | Name |
|----------|------|
| Region | eu-central-1 |
| ECR Repository | pack-calculator |
| ECS Cluster | pack-calculator |
| ECS Service | pack-calculator-service |
| Load Balancer | pack-calculator-alb |

## GitHub Secrets

Required secrets in repository settings:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
