import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";

export function SetupPage() {
  return (
    <div className="max-w-sm mx-auto px-6">
      <header className="pt-16 pb-8">
        <h1 className="text-4xl font-extrabold tracking-tight">Setup</h1>
        <p className="text-slate-700 mt-2">
          To begin using Openticket, you must first configure the initial admin
          user.
        </p>
      </header>
      <form action="">
        <fieldset className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input id="name" type="text" autoFocus />
          </div>

          <div className="space-y-2">
            <Label htmlFor="username">Username</Label>
            <Input id="username" type="text" />
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" />
          </div>

          <div className="space-y-2">
            <Label htmlFor="confirm-password">Confirm password</Label>
            <Input id="confirm-password" type="password" />
          </div>

          <Button className="w-full" size="lg">
            Setup first user
          </Button>
        </fieldset>
      </form>
    </div>
  );
}
