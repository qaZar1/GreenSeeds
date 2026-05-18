type Props = {
  steps: string[];
  current: number;
};

export function Stepper({
  steps,
  current,
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

          const done = i < current;
          const active = i === current;

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
                    done
                      ? "border-[var(--color-primary)] bg-[var(--bg-card)] text-[var(--color-primary)] border-[var(--color-primary)]"
                      : active
                        ? "border-[var(--color-primary)] text-[var(--color-primary)] bg-[var(--bg-card)]"
                        : "rounded-full border border-[var(--border-color)] text-[var(--text-secondary)] bg-[var(--bg-page)] bg-[var(--bg-card)]"
                  }
                `}
              >
                {done
                  ? (
                    <i className="fas fa-check text-[11px]" />
                  )
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
                {prev.index + 1}
              </div>

              <div className="w-[36px] h-[2px] bg-[var(--color-primary)]/50" />
            </>
          )}

          {/* current */}
          <div
            className="
              w-[42px]
              h-[42px]
              rounded-full

              border-2
              border-[var(--color-primary)]

              bg-[var(--bg-card)]
              text-[var(--color-primary)]

              flex
              items-center
              justify-center

              text-[15px]
              font-semibold

              shrink-0
              z-10
            "
          >
            {current + 1}
          </div>

          {/* next */}
          {next && (
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

          <div className="text-[14px] font-medium text-[var(--text-primary)]">
            {steps[current]}
          </div>

          <div className="text-[12px] text-[var(--text-secondary)] mt-[2px]">
            Шаг {current + 1} из {steps.length}
          </div>

        </div>

      </div>

    </div>
  );
}