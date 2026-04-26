type ErrorStateProps = {
  onRetry: () => void;
};

const ErrorState = ({ onRetry }: ErrorStateProps) => {
  return (
    <div className="flex items-center justify-center h-full">
      <div className="text-center max-w-sm">

        <div className="text-[32px] mb-[12px] text-[var(--status-danger-text)]">
          ⚠️
        </div>

        <h2 className="text-[18px] font-semibold text-[var(--text-primary)]">
          Сервис недоступен
        </h2>

        <p className="mt-[8px] text-[13px] text-[var(--text-secondary)]">
          Не удалось загрузить данные. Попробуйте снова.
        </p>

        <button
          onClick={onRetry}
          className="mt-[16px] px-[16px] py-[8px] rounded-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] hover:bg-[var(--color-primary-hover)] transition-colors"
        >
          Повторить
        </button>
      </div>
    </div>
  );
};

export default ErrorState;