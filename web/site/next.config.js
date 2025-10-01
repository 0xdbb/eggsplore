/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  eslint: {
    ignoreDuringBuilds: true, // optional if ESLint breaks builds
  },
  images: {
    unoptimized: true, // Netlify doesn't fully support Next image optimization without plugin
  },
};

export default nextConfig;
