export type Column<T> = {
  header: string;
  render: (row: T) => React.ReactNode;
  className?: string;
  headerClassName?: string;
  width?: string;
};