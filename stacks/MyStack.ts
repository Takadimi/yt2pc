import { StackContext, Api, Table } from "sst/constructs";

export function API({ stack }: StackContext) {
  const table = new Table(stack, "yt2pc", {
    fields: {
      PK: "string",
      SK: "string",
    },
    primaryIndex: { partitionKey: "PK", sortKey: "SK" },
  })

  const api = new Api(stack, "api", {
    defaults: {
      throttle: {
        burst: 50,
        rate: 5,
      },
      function: {
        environment: {
          "TABLE_NAME": table.tableName,
        },
        bind: [table],
      },
    },
    routes: {
      "GET /youtube/{id}": "functions/lambda/youtube/get/main.go",
      "GET /rss/youtube/{id}": "functions/lambda/rss/youtube/get/main.go",
    },
  });

  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}
