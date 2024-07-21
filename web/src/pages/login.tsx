import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { isSuccess, OpenticketSdk } from "../sdk";
import { Button } from "../ui/button";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  Form,
} from "../ui/form";
import { Input } from "../ui/input";
import { Navigate, useNavigate } from "react-router-dom";
import { useStatus } from "@/status";
import { useMutation } from "@tanstack/react-query";
import { useToast } from "@/ui/use-toast";
import { capitalize } from "@/utils";

const loginFormSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

type LoginFormValues = z.infer<typeof loginFormSchema>;

export function LoginPage() {
  const sdk = new OpenticketSdk();
  const { toast } = useToast();
  const navigate = useNavigate();
  const status = useStatus();
  const loginMutation = useMutation({
    mutationFn: sdk.login,
    onSuccess: async (res) => {
      if (!isSuccess(res)) {
        toast({
          title: "Error on login",
          description: capitalize(res.message),
        });
        return;
      }

      status.authenticate(res.data.user);
      navigate("/");
    },
    onError: (err) => {
      console.error(err);
    },
  });

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  if (status.data.user) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="max-w-sm mx-auto px-6">
      <header className="pt-16 pb-8">
        <h1 className="text-4xl font-extrabold tracking-tight">Login</h1>
      </header>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit((values) => {
            loginMutation.mutate(values);
          })}
        >
          <fieldset className="space-y-4" disabled={loginMutation.isPending}>
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="email" autoFocus {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full" size="lg">
              Login
            </Button>
          </fieldset>
        </form>
      </Form>
    </div>
  );
}
