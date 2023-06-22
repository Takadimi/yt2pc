import { StackContext, Api, use } from "sst/constructs";
import { DNS } from "./dns";

export function API({ stack }: StackContext) {
  const dns = use(DNS);

  const apiDomainName = "api."+dns.domain;

  const api = new Api(stack, "api", {
    defaults: {
      throttle: {
        burst: 10,
        rate: 1,
      },
      function: {
        environment: {},
        bind: [],
      },
    },
    routes: {
      "GET /youtube/{id}": "functions/lambda/youtube/get/main.go",
      "GET /rss/youtube/{id}": "functions/lambda/rss/youtube/get/main.go",
    },
    customDomain: {
      domainName: apiDomainName,
      hostedZone: dns.zone,
    },
  });

  const rssFeedFunc = api.getFunction("GET /rss/youtube/{id}");
  rssFeedFunc?.addEnvironment("API_ENDPOINT", api.customDomainUrl || api.url);

  stack.addOutputs({
    ApiEndpoint: api.customDomainUrl || api.url,
  });
}
