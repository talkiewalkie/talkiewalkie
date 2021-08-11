import Head from "next/head";
import Link from "next/link";
import { ReactNode } from "react";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div>
      <Head>
        <title>TalkieWalkie</title>
        <meta name="description" content="TalkieWalkie webapp" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main
        style={{
          minHeight: "100vh",
          height: "100vh",
        }}
      >
        <div className="fixed z-10 w-screen bg-white py-4 border-b border-b-2 flex justify-center">
          <Link href="/">
            <a>TalkieWalkie</a>
          </Link>
        </div>
        <div className="bg-gray-100 min-h-full">
          <div
            style={{ maxWidth: "40rem" }}
            className="mx-auto px-2 py-20 h-full flex flex-col space-y-8"
          >
            {children}
          </div>
        </div>
      </main>
    </div>
  );
}
