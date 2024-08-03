import { Avatar, AvatarFallback } from "./avatar";

type UserAvatarProps = {
  name: string;
};

export function UserAvatar(props: UserAvatarProps) {
  return (
    <Avatar className="w-5 h-5 text-xs">
      <AvatarFallback>{props.name[0]}</AvatarFallback>
    </Avatar>
  );
}
