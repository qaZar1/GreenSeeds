export const steps = [
  "Подготовка",
  "Первое фото",
  "Перемещение платформы",
  "Второе фото",
  "Вычисления",
  "Сохранение результата"
] as const;

export type StepLabel = typeof steps[number];

export type StepKey =
  | "prepare"
  | "photo1"
  | "move"
  | "photo2"
  | "calculate"
  | "save";

export type StepConfig = {
  label: StepLabel;
  action?: () => Promise<void>;
  next: StepKey | null;
};