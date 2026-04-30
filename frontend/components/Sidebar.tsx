"use client";

import { ChartPie, House, LogOut, Settings } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function Sidebar() {
  const pathname = usePathname();

  const links = [
    { href: "/dashboard", label: "HOME", icon: House },
    { href: "/analytics", label: "ANALYTICS", icon: ChartPie },
    { href: "/settings", label: "SETTINGS", icon: Settings },
  ];
  return (
    <aside className="w-50 hidden md:flex flex-col fixed items-center border-r h-full bg-white">
      <div className="flex flex-col w-full h-full justify-between ">
        <div>
          {links.map(({ href, label, icon: Icon }) => {
            const active = pathname === href;

            return (
              <Link
                key={href}
                href={href}
                className={`flex items-center gap-3 w-full px-5 py-4 border-b border-t border-black ${
                  active ? "bg-[#4d4949] text-[#FDB833]" : "bg-white text-black"
                }`}
              >
                <Icon size={18} />
                <span className="font-bold text-sm">{label}</span>
              </Link>
            );
          })}
        </div>

        <div className="bg-black text-[#FDB833] flex items-center pl-7 gap-3 h-13.75">
          <LogOut />
          <p className="font-bold">Log Out</p>
        </div>
      </div>
    </aside>
  );
}
