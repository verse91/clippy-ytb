import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "@/styles/globals.css";
import { AuthProvider } from "@/lib/auth-context";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Clippy - Turn moments into clips",
  description: "Make viral clips with Clippy",
  other: {
    link: [
      {
        rel: "stylesheet",
        href: "https://cdn.boxicons.com/fonts/basic/boxicons.min.css",
      },
      {
        rel: "stylesheet",
        href: "https://cdn.boxicons.com/fonts/brands/boxicons-brands.min.css",
      },
    ],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased bg-black`}
      >
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
