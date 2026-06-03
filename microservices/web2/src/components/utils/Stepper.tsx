type Props = {
  steps: string[];
  current: number;
  completed?: boolean;
  errorStepIndex?: number | null;
};

export function Stepper({
  steps,
  current,
  completed = false,
  errorStepIndex = null,
}: Props) {

  // ================= MOBILE =================
  
  const prev =
    current > 0
      ? {
          index: current - 1,
          label: steps[current - 1],
        }
      : null;

  const next =
    current < steps.length - 1
      ? {
          index: current + 1,
          label: steps[current + 1],
        }
      : null;
      
  const isMobileError = errorStepIndex !== null && errorStepIndex === current;

  return (
    <div className="w-full select-none">

      {/* ================= DESKTOP ================= */}
      <div
        className="
          hidden lg:flex
          items-start
          w-full
        "
      >
        {steps.map((step, i) => {

          const isErrored = errorStepIndex !== null && i === errorStepIndex;
          const done = completed
            ? i <= current
            : i < current;

          const active = !completed && !isErrored && i === current;

          return (
            <div
              key={i}
              className="
                flex flex-col items-center relative
                min-w-[90px] lg:min-w-0 flex-1
              "
            >

              {/* линия */}
              {i !== steps.length - 1 && (
                <div
                  className="absolute top-[17px] left-1/2 w-full h-[2px]"
                  style={{
                    transform: "translateX(17px)",
                  }}
                >
                  <div className="w-full h-full bg-[var(--border-light)]">
                    {done && (
                      <div className="w-full h-full bg-[var(--color-primary)]" />
                    )}
                  </div>
                </div>
              )}

              {/* кружок */}
              <div
                className={`
                  z-10
                  w-[34px]
                  h-[34px]
                  rounded-full
                  flex
                  items-center
                  justify-center
                  text-[13px]
                  font-medium
                  border

                  ${
                    isErrored
                      ? "border-red-500 bg-red-50 text-red-500"
                      : done
                        ? "border-[var(--color-primary)] bg-[var(--bg-card)] text-[var(--color-primary)]"
                        : active
                          ? "border-[var(--color-primary)] text-[var(--color-primary)] bg-[var(--bg-card)]"
                          : "border border-[var(--border-color)] text-[var(--text-secondary)] bg-[var(--bg-card)]"
                  }
                `}
              >
                {isErrored
                  ? <i className="fas fa-exclamation text-[11px]" />
                  : done
                    ? <i className="fas fa-check text-[11px]" />
                    : i + 1}
              </div>

              {/* текст */}
              <div
                className={`
                  mt-[8px]
                  text-[11px] lg:text-[12px]
                  text-center
                  max-w-[80px] lg:max-w-[120px]
                  leading-[1.3]

                  ${
                    isErrored
                      ? "text-red-500 font-medium"
                      : active || done
                        ? "text-[var(--text-primary)]"
                        : "text-[var(--text-secondary)]"
                  }
                `}
              >
                {step}
              </div>
            </div>
          );
        })}
      </div>

      {/* ================= MOBILE ================= */}
      <div className="lg:hidden">

        {/* circles */}
        <div className="flex items-center justify-center">

          {/* prev */}
          {prev && (
            <>
              <div
                className="
                  w-[26px]
                  h-[26px]
                  rounded-full
                  border
                  border-[var(--color-primary)]

                  text-[var(--color-primary)]
                  bg-[var(--bg-card)]

                  flex
                  items-center
                  justify-center

                  text-[11px]
                  shrink-0
                "
              >
                <i className="fas fa-check text-[9px]" />
              </div>

              <div className="w-[36px] h-[2px] bg-[var(--color-primary)]" />
            </>
          )}

          {/* current */}
          <div
            className={`
              w-[42px]
              h-[42px]
              rounded-full

              flex
              items-center
              justify-center

              text-[15px]
              font-semibold

              shrink-0
              z-10

              ${
                isMobileError
                  ? "border-2 border-red-500 bg-red-50 text-red-500"
                  : completed
                    ? "border-2 border-[var(--color-primary)] bg-[var(--bg-card)] text-[var(--color-primary)]"
                    : "border-2 border-[var(--color-primary)] bg-[var(--bg-card)] text-[var(--color-primary)]"
              }
            `}
          >
            {isMobileError
              ? <i className="fas fa-exclamation text-[14px]" />
              : completed
                ? <i className="fas fa-check text-[14px]" />
                : current + 1}
          </div>

          {/* next */}
          {next && !completed && (
            <>
              <div className="w-[36px] h-[2px] bg-[var(--border-color)]" />

              <div
                className="
                  w-[26px]
                  h-[26px]
                  rounded-full
                  border
                  border-[var(--border-color)]

                  text-[var(--text-secondary)]
                  bg-[var(--bg-page)]

                  flex
                  items-center
                  justify-center

                  text-[11px]
                  opacity-80
                  shrink-0
                "
              >
                {next.index + 1}
              </div>
            </>
          )}

        </div>

        {/* current label */}
        <div className="mt-[12px] text-center">
          <div className={`text-[14px] font-medium ${isMobileError ? "text-red-500" : "text-[var(--text-primary)]"}`}>
            {isMobileError
              ? `Ошибка: ${steps[current]}`
              : completed
                ? "Завершено"
                : steps[current]}
          </div>

          <div className="text-[12px] text-[var(--text-secondary)] mt-[2px]">
            {isMobileError
              ? "Требуется внимание"
              : completed
                ? `Все ${steps.length} шагов выполнены`
                : `Шаг ${current + 1} из ${steps.length}`}
          </div>
        </div>

      </div>

    </div>
  );
}