import { Badge } from "./badge";

type TicketLabelstProps = {
  labels: string[];
};

export function TicketLabels({ labels }: TicketLabelstProps) {
  return (
    <div className="space-x-1">
      {labels.map((l) => (
        <Badge variant="outline" key={l}>
          {l}
        </Badge>
      ))}
    </div>
  );
}
