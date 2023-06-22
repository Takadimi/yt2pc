// import { HostedZone } from "aws-cdk-lib/aws-route53";
import { StackContext } from "sst/constructs";

export function DNS(ctx: StackContext) {
  const name = ctx.stack.stage === "production" ? "production.yt2pc.ethanwoodward.net" : "dev.yt2pc.ethanwoodward.net";

  return {
    zone: "ethanwoodward.net",
    domain: ctx.stack.stage === "production" ? name : `${ctx.stack.stage}.${name}`,
  };
}
