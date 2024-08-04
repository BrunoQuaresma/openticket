import { Avatar, AvatarFallback } from "./avatar";

type Size = "sm" | "md";

type UserAvatarProps = {
  name: string;
  size?: Size;
};

const classNameBySize: Record<Size, string> = {
  sm: "w-5 h-5 text-xs",
  md: "",
};

export function UserAvatar({ name, size = "md" }: UserAvatarProps) {
  return (
    <Avatar className={classNameBySize[size]}>
      <AvatarFallback>{name[0]}</AvatarFallback>
    </Avatar>
  );
}
