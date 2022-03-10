#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { FlapFlapCognitoStack } from "../lib/cdk-stack";
import { Tags } from "aws-cdk-lib";

const app = new cdk.App();
const flapflapstack = new FlapFlapCognitoStack(app, "FlapFlapCognitoStack", {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});

// Tag all resources in this stack with "flapflap":
Tags.of(flapflapstack).add("flapflap", "flapflap");
