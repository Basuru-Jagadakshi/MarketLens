import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: 'standalone',

  env: {
    GO_API_BASE_URL: process.env.GO_API_BASE_URL ?? "http://localhost:8080/api/v1",
  },
};

export default nextConfig;
