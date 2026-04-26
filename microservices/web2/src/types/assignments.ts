export interface Assignment {
  id?: number;
  shift: number;
  number: number;
  receipt: number;
  description: string;
  amount: number;
}

export interface ActiveTask {
  id: number;
  shift: number;
  number: number;
  receipt: number;
  dt: string;
  amount: number;
  done_turns: number;
  seed: string;
  seed_ru: string;
}

export interface Task {
  id: number;
  shift: number;
  number: number;
  seed: string;
  seed_ru: string;
  bunker: number;
  gcode: string;
  receipt: number;
  required_amount: number;
  reports?: Report[];
}