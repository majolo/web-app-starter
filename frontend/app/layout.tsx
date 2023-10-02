import "./globals.css";
import React from "react";
import Navbar from "@/components/Navbar";

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
          {children}
        </main>
      </body>
    </html>
  );
}
