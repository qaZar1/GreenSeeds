export type Dynamic<T> = T | ((values: Record<string, any>) => T);

export interface FormField {
  name: string;
  label: Dynamic<string>;
  type: "text" | "number" | "select" | "datetime" | "textarea" | "password";
  required?: boolean;
  placeholder?: string;

  disabled?: Dynamic<boolean>;

  min?: Dynamic<number | undefined>;
  max?: Dynamic<number | undefined>;

  options?: Dynamic<{ label: string; value: any }[]>;
}