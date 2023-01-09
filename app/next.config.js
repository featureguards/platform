/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: !!process.NEXT_PUBLIC_APP_ENV,
  env: {
    NEXT_PUBLIC_APP_ENV: process.env.NEXT_PUBLIC_APP_ENV,
    NEXT_PUBLIC_MIXPANEL_ID: process.env.NEXT_PUBLIC_MIXPANEL_ID
  }
};

module.exports = nextConfig;
