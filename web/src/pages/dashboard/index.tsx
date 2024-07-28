import { Button } from "@/ui/button";
import { Link } from "react-router-dom";
import { Helmet } from "react-helmet-async";

export function IndexDashboardPage() {
  return (
    <>
      <Helmet>
        <title>Openticket - Tickets</title>
      </Helmet>
      <div className="w-full h-full flex items-center justify-center text-center">
        <div className="space-y-4">
          <hgroup className="space-y-1">
            <h3 className="text-2xl font-semibold tracking-tight">
              No tickets open
            </h3>
            <p className="text-muted-foreground">
              New open tickets will show up automatically
            </p>
          </hgroup>
          <Button asChild>
            <Link to="/tickets/new">Open a ticket</Link>
          </Button>
        </div>
      </div>
    </>
  );
}
