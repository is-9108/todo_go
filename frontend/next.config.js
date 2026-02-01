/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "standalone", // Docker 用にスタンドアロン出力
};

module.exports = nextConfig;
