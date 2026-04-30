"use client";

import Image from "next/image";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex">
      {/* LEFT PANEL */}
      <div className="hidden xl:flex flex-60"></div>

      {/* RIGHT PANEL */}
      <div className="flex-40 flex flex-col items-center justify-center bg-white">
        <form
          action=""
          className="flex flex-col gap-2 items-center w-full max-w-sm px-8"
        >
          <Image
            src="/staredesk-logo.svg"
            alt="StareDesk"
            width={100}
            height={100}
            loading="eager"
            className="object-cover h-70 w-auto"
          />

          <div className="flex flex-col gap-1 w-full">
            <label className="font-bold tracking-wide">Email address</label>
            <input
              type="email"
              placeholder="user@domain.com"
              className="border-2 px-3 py-2 rounded-[5px] shadow-sm shadow-gray-800 placeholder:font-semibold"
            />
          </div>

          <div className="flex flex-col gap-1 w-full">
            <label className="font-bold tracking-wide">Password</label>
            <input
              type="password"
              placeholder="••••••••"
              className="border-2 px-3 py-2 rounded-[5px] shadow-sm shadow-gray-800 placeholder:font-semibold"
            />
          </div>

          <button
            type="submit"
            className="border-2 flex w-fit justify-center px-12 py-2 mt-4 items-center font-bold bg-[#FDB833] rounded-[5px] hover:bg-[#e8960f] transition-colors"
          >
            LOGIN
          </button>
        </form>
      </div>
    </div>
  );
}
