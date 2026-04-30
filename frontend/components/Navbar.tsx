import { LogOut } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="bg-black h-18 flex items-center justify-between pl-3 pr-7 md:px-10 fixed top-0 left-0 right-0 z-50">
      <Link href="/dashboard" className="flex items-center gap-0 md:gap-2">
        <Image
          src="/staredesk-logo-dua.png"
          alt="StareDesk"
          height={45}
          width={45}
          className="object-contain"
        />
        <h1 className="text-white font-bold text-2xl">StareDesk</h1>
      </Link>

      <LogOut className="flex md:hidden h-5 w-5 text-[#FDB833]" />

      <div className="bg-white hidden md:flex items-center gap-3 py-2 px-5 lg:px-3 rounded-md">
        <div className="h-5 w-5 rounded-full bg-amber-500"></div>
        <h1 className="font-bold tracking-wide">ONLINE</h1>
      </div>
    </nav>
  );
}
