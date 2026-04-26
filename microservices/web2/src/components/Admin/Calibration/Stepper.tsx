type Props = {
  steps: string[];
  current: number;
};

export function Stepper({ steps, current }: Props) {

  return (
    <div className="w-full select-none">

      <div className="flex items-start w-full">

        {steps.map((step, i) => {

          const done = i < current;
          const active = i === current;

          return (
            <div key={i} className="flex-1 flex flex-col items-center relative">

              {/* линия */}
              {i !== steps.length - 1 && (
                <div
                  className="absolute top-[17px] left-1/2 w-full h-[2px]"
                  style={{ transform: "translateX(17px)" }}
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
                  done
                    ? "bg-[var(--color-primary)] text-[var(--text-inverse)] border-[var(--color-primary)]"
                    : active
                    ? "border-[var(--color-primary)] text-[var(--color-primary)] bg-[var(--bg-card)]"
                    : "border-[var(--border-color)] text-[var(--text-secondary)] bg-[var(--bg-card)]"
                }
                `}
              >
                {done ? <i className="fas fa-check text-[11px]" /> : i + 1}
              </div>

              {/* текст */}
              <div
                className={`
                mt-[8px]
                text-[12px]
                text-center
                max-w-[120px]
                ${
                  active
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

    </div>
  );
}