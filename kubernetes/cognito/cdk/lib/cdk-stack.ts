import {
  aws_cognito,
  CfnOutput,
  RemovalPolicy,
  Stack,
  StackProps,
} from "aws-cdk-lib";
import { Construct } from "constructs";

export class FlapFlapCognitoStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // Cognito User Pool
    const userPool = new aws_cognito.UserPool(this, "flapflap-userpool", {
      userPoolName: "flapflap-userpool",
      selfSignUpEnabled: true,
      signInAliases: { email: true },
      autoVerify: { email: true },
      standardAttributes: {
        givenName: {
          required: true,
          mutable: true,
        },
        familyName: {
          required: true,
          mutable: true,
        },
      },
      customAttributes: {
        isAdmin: new aws_cognito.StringAttribute({ mutable: true }),
      },
      passwordPolicy: {
        minLength: 6,
        requireLowercase: true,
        requireDigits: true,
        requireUppercase: false,
        requireSymbols: false,
      },
      accountRecovery: aws_cognito.AccountRecovery.EMAIL_ONLY,
      removalPolicy: RemovalPolicy.RETAIN,
    });

    const cognitoStandardAttributes: aws_cognito.StandardAttributesMask = {
      givenName: true,
      familyName: true,
      email: true,
      emailVerified: true,
      address: true,
      birthdate: true,
      gender: true,
      locale: true,
      middleName: true,
      fullname: true,
      nickname: true,
      phoneNumber: true,
      phoneNumberVerified: true,
      profilePicture: true,
      preferredUsername: true,
      profilePage: true,
      timezone: true,
      lastUpdateTime: true,
      website: true,
    };

    // Cognito User Pool Client
    const userPoolClient = new aws_cognito.UserPoolClient(
      this,
      "flapflap-userpool-client",
      {
        userPool,
        userPoolClientName: "flapflap-userpool-client",
        authFlows: {
          adminUserPassword: true,
          custom: true,

          // SRP = Secure Remote Password
          // See: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html
          userSrp: true,
        },
        supportedIdentityProviders: [
          aws_cognito.UserPoolClientIdentityProvider.COGNITO,
        ],
        readAttributes: new aws_cognito.ClientAttributes()
          .withStandardAttributes(cognitoStandardAttributes)
          .withCustomAttributes("isAdmin"),
        writeAttributes:
          new aws_cognito.ClientAttributes().withStandardAttributes({
            ...cognitoStandardAttributes,
            emailVerified: false,
            phoneNumberVerified: false,
          }),
      }
    );

    // OUTPUT: user pool id
    new CfnOutput(this, "flapflap-userpool-id", { value: userPool.userPoolId });

    // OUTPUT: user pool client id
    new CfnOutput(this, "flapflap-userpool-client-id", {
      value: userPoolClient.userPoolClientId,
    });
  }
}
