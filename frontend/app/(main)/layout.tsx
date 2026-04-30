import BottomNavBar from "@/components/BottomNavBar";
import Navbar from "@/components/Navbar";
import Sidebar from "@/components/Sidebar";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="h-screen flex flex-col overflow-hidden pt-18">
      <Navbar />
      <div className="flex-1 flex flex-row overflow-hidden  ">
        <Sidebar />
        <div className="flex-1 overflow-y-auto md:pb-0 md:pl-50">
          {children}
        </div>
      </div>
      <BottomNavBar />
    </div>
  );
}
