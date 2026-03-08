import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  images: {
    dangerouslyAllowSVG: true,
    domains: ['static.usernames.app-backend.toolsforhumanity.com', 'api.dicebear.com'],
  },
  allowedDevOrigins: ['*'],
  reactStrictMode: false,
};

export default nextConfig;
