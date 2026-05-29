import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "MarketLens BFF Gateway",
  description: "API Aggregator Gateway",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}