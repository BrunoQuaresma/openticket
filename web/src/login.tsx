import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";

export function LoginPage() {
  return (
    <div className="max-w-sm mx-auto px-6">
      <header className="pt-16 pb-8">
        <h1 className="text-4xl font-extrabold tracking-tight">Login</h1>
      </header>
      <form action="">
        <fieldset className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" autoFocus />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" />
          </div>

          <Button className="w-full" size="lg">
            Login
          </Button>
        </fieldset>
      </form>
    </div>
  );
}
