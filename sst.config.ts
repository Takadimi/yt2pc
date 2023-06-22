import { SSTConfig } from "sst";
import { API } from "./stacks/api";
import { DNS } from "./stacks/dns";

export default {
  config(_input) {
    return {
      name: "yt2pc",
      region: "us-east-1",
    };
  },
  stacks(app) {
    if (app.stage !== "production") {
      app.setDefaultRemovalPolicy("destroy");
    }
    app.setDefaultFunctionProps({
      runtime: "go1.x",
    });
    app
      .stack(DNS)
      .stack(API);
  }
} satisfies SSTConfig;
