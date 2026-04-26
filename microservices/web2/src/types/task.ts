import type { Report } from "./reports"

export type FailedReport = {
  turn: number;
  error?: string;
};

export type Bunker = {
  bunker: number;
  amount: number;
};

export type TaskRecord = {
  id: number;
  number: number;
  shift: string;
  seed: string;
  seed_ru: string;
  required_amount: number;
  gcode?: string;
  bunker: number;
  turn: number;
  receipt: string;
  reports?: Report[];
};

export type DecisionModalState = {
  open: boolean;
  reason: string;
  photo: string | null;
};

export type WSMessage =
  | {
      type: "BEGIN";
      params: {
        shift: string;
        number: number;
        receipt?: string;
        seed: string;
        turn: number;
        completed_amount: number;
        required_amount: number;
        bunker: number;
        gcode?: string;
        extraMode?: boolean;
      };
    }
  | {
      type: "DECISION";
      status: "OK" | "NOK";
      solution: string;
    };