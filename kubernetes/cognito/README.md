# Cognito + Ingress NGINX integration

This directory contains files related to integrating AWS Cognito authentication
with the Ingress NGINX (i.e. locking down our cluster to only authenticated
users).

The `cdk/` dir contains the CDK stack used to deploy the Cognito User Pool and
User Pool Client resources on AWS.

## Scripts

**NOTE: the CDK stack is configured to create a `cdk-outputs.json` file containing the `<user_pool_id>` and `<user_pool_client_id>` values.**

After creating the resources with CDK, you can try some of the following scripts
to test your configuration:

**_Create a new user using AWS CLI:_**

```bash
aws cognito-idp admin-create-user \
    --user-pool-id <user_pool_id> \
    --username jeff@example.com \
    --user-attributes \
        Name="given_name",Value="jeff" \
        Name="family_name",Value="smith"
```
