import type { NextConfig } from "next";
import config from "./tailwind.config";

const nextConfig: NextConfig = {
  /* config options here */
  iwebpack: (config, { dev }) => {
    if (dev) {
      config.watchOptions = {
        poll: 1000,
        aggregateTimeout: 300,
      };
    }
    return config; // <-- This must be inside the webpack function
  },
};

export default nextConfig;
