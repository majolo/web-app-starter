import "../globals.css";
import React from "react";
import Navbar from "@/app/components/Navbar";
import Providers from "@/app/(protected)/providers";
import Protected from "@/app/components/Protected";

export const metadata = {
  title: "Web App Starter",
  description: "Basic starter app.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <main className="min-h-screen bg-background flex flex-col items-center">
          <Navbar />
          <Providers>{children}</Providers>
          <Protected />
        </main>
      </body>
    </html>
  );
}
