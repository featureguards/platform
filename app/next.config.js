/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_APP_ENV: process.env.NEXT_PUBLIC_APP_ENV,
    NEXT_PUBLIC_MIXPANEL_ID: process.env.NEXT_PUBLIC_MIXPANEL_ID
  }
};

module.exports = nextConfig;
