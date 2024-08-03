import { Link, Navigate, Outlet } from "react-router-dom";
import { useStatus } from "../../status";
import { TicketCheckIcon } from "lucide-react";
import { UserAvatar } from "@/ui/user-avatar";

export function Dashboard() {
  const { data } = useStatus();

  if (!data.user) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="flex flex-col w-full h-screen">
      <div className="border-b px-6 h-14 flex justify-between items-center text-sm">
        <Link to="/" aria-label="Go to tickets">
          <Logo />
        </Link>

        <div className="justify-self-end">
          <div className="flex items-center gap-2">
            <UserAvatar name={data.user.name} />
            <span>{data.user.name}</span>
          </div>
        </div>
      </div>
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
