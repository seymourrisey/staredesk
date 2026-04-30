"use client";

import { ChartPie, House, Settings } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function BottomNavBar() {
  const pathname = usePathname();

  const links = [
    { href: "/dashboard", label: "HOME", icon: House },
    { href: "/analytics", label: "ANALYTICS", icon: ChartPie },
    { href: "/settings", label: "SETTINGS", icon: Settings },
  ];

  return (
    <nav className="md:hidden h-18.75 bg-white fixed bottom-0 left-0 w-full z-50 flex items-center border-t">
      {links.map(({ href, label, icon: Icon }) => {
        const active = pathname === href;
        return (
          <Link
            key={href}
            href={href}
            className={`flex-1 h-full flex flex-col justify-center items-center gap-1 ${
              active ? "bg-black text-[#FDB833]" : "bg-white text-black"
            }`}
          >
            <Icon size={20} />
            <p className={"text-xs font-bold"}>{label}</p>
          </Link>
        );
      })}
    </nav>
  );
}
