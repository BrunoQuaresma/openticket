import { Link, Navigate, Outlet } from "react-router-dom";
import { useStatus } from "../status";
import { CircleUserRoundIcon, TicketCheckIcon } from "lucide-react";
import { Tabs, TabsList, TabsTrigger } from "@/ui/tabs";

export function Dashboard() {
  const { data } = useStatus();

  if (!data.user) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="flex flex-col w-full h-screen">
      <header className="border-b px-6 h-14 grid grid-cols-3 justify-center items-center text-sm">
        <Logo />

        <Tabs value="open" className="justify-self-center">
          <TabsList>
            <TabsTrigger value="open" asChild>
              <Link to="/tickets">Open</Link>
            </TabsTrigger>
            <TabsTrigger value="close" asChild>
              <Link to="/tickets">Close</Link>
            </TabsTrigger>
          </TabsList>
        </Tabs>

        <div className="justify-self-end">
          <div className="flex items-center gap-2">
            <CircleUserRoundIcon className="w-6" strokeWidth={1.5} />
            <span>{data.user.name}</span>
          </div>
        </div>
      </header>
      <main className="h-full">
        <Outlet />
      </main>
    </div>
  );
}

function Logo() {
  return (
    <span className="flex gap-2 items-center">
      <TicketCheckIcon className="w-4" />
      <span className="font-medium">Openticket</span>
    </span>
  );
}
