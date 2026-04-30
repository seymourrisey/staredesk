"use client";

import Image from "next/image";

export default function DashboardPage() {
  return (
    <div className="p-5 md:py-5 md:px-10 flex flex-col xl:flex-row items-center md:items-start gap-0 xl:gap-5">
      <div className="flex flex-col items-start px-5 md:px-0 w-100">
        <div className="bg-[#f4f4f4] border-2 rounded-[10px] w-full flex justify-between md:hidden py-3 pl-3 pr-5 font-bold">
          <div className="flex items-center gap-2">
            <div className="bg-[#FDB833] border h-5 w-5 rounded-[5px]" />
            <h1>Device Status</h1>
          </div>

          <h1>ONLINE</h1>
        </div>

        <div className="relative bg-white border-2 w-full max-w-100 p-5 mt-5 md:mt-0 flex flex-col gap-6 rounded-xs">
          <Image
            src="/optimal.png"
            alt="optimal"
            height={100}
            width={100}
            loading="eager"
            className="absolute right-3 bottom-0 h-32 w-32"
          />

          <p className="font-bold text-black/70">Current Condition</p>

          <div className="flex flex-col gap-2">
            <div className="flex items-center gap-3">
              <div className="h-5 w-5 rounded-full bg-[#FDB833] border" />
              <p className="text-4xl font-bold">Optimal</p>
            </div>

            <p>Maintaining healthy alignment</p>
          </div>
        </div>

        <div className="relative flex flex-col gap-6 w-full max-w-100 p-5 mt-5 md:mt-3 bg-black rounded-xs overflow-hidden">
          <Image
            src="/clock.png"
            alt="clock"
            width={100}
            height={100}
            loading="eager"
            className="absolute right-0 bottom-0 h-26 w-26"
          />

          <p className="font-bold text-[#FDB833]">Current Session</p>

          <div className="flex flex-col gap-2">
            <p className="text-white text-5xl tracking-tighter font-extrabold">
              12 : 14 : 25
            </p>
            <p className="text-white/70 text-sm font-semibold tracking-tight">
              Active focus time
            </p>
          </div>
        </div>

        <div className="border-2 bg-white rounded-xs flex flex-col w-full max-w-100 mt-5 md:mt-3 overflow-hidden">
          <p className="px-4 py-2 border-b-2 text-md font-bold bg-[#eeeeee] tracking-wide">
            TODAY&apos;S SUMMARY
          </p>

          <div className="p-4 flex flex-col gap-2 border-b-2 relative">
            <Image
              src="/study-table.png"
              alt="study-table"
              height={100}
              width={100}
              loading="eager"
              className="absolute right-0 bottom-0 h-20 w-auto"
            />

            <p className="text-xs font-bold">TOTAL FOCUS TIME</p>
            <p className="font-extrabold text-4xl tracking-tight">4h 25m</p>
          </div>

          <div className="p-4 flex flex-col gap-2 relative">
            <Image
              src="/hour-head.png"
              alt="hour-head"
              width={100}
              height={100}
              loading="eager"
              className="absolute right-0 bottom-0 h-20 w-auto"
            />

            <p className="text-xs font-bold">PEAK PRODUCTIVITY HOUR</p>
            <p className="font-extrabold text-4xl tracking-tight">11:00 AM</p>
          </div>
        </div>
      </div>

      <div className="border-2 bg-white mt-5 md:mt-3 xl:mt-0 w-full max-w-100 rounded-xs overflow-hidden">
        <h1 className="px-4 py-2 border-b-2 text-md font-bold bg-[#eeeeee] tracking-wide">
          LIVE TELEMETRY
        </h1>

        <div className="flex">
          <div className="flex flex-col border-r-2 flex-1">
            <div className="border-b-2 px-4 pb-3 pt-10 flex flex-col gap-1">
              <h1 className="font-extrabold text-6xl">45</h1>
              <p className="text-md font-bold text-black/45">Distance (cm)</p>
            </div>

            <div className="font-bold px-4 pb-3 pt-15">
              <div className="flex gap-2 items-center ">
                <div className="h-5 w-5 bg-[#FDB833] border rounded-full" />
                <p className="text-xl">Present</p>
              </div>

              <p className="text-sm font-bold text-black/45">PIR Active</p>
            </div>
          </div>

          <div className="relative flex-1 flex flex-col justify-end p-4">
            <Image
              src="/black-logo.png"
              alt="black-logo"
              loading="eager"
              width={100}
              height={100}
              className="absolute top-0 right-0 h-53 w-36"
            />
            <p className="text-sm font-bold text-black/45">Cahaya Cukup</p>
            <p className="text-6xl font-extrabold tracking-tight">720</p>
            <p className="text-md font-bold text-black/45">Light (LDR)</p>
          </div>
        </div>
      </div>
    </div>
  );
}
