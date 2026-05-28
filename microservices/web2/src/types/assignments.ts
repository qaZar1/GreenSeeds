export interface Assignment {
  id?: number;
  shift: number;
  number: number;
  recipe: number;
  description: string;
  amount: number;
}

export interface ActiveTask {
  id: number;
  shift: number;
  number: number;
  recipe: number;
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
  recipe: number;
  required_amount: number;
  reports?: Report[];
}