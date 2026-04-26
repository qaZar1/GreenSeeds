import React, { useEffect, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { Stepper } from "./Stepper";
import { steps } from "../../../types/calibration";
import { useCalibrationFlow } from "../../hooks/useCalibrationFlow";

type Result = {
  dx: number;
  dy: number;
};

const order = [
  "prepare",
  "photo1",
  "move",
  "photo2",
  "calculate",
  "save"
];

const nextLabel: Record<string, string> = {
  prepare: "Начать",
  photo1: "Сделать фото",
  move: "Далее",
  photo2: "Сделать фото",
  calculate: "Рассчитать",
  save: "Сохранить"
};

const CalibrationPage: React.FC = () => {

  usePageHeader("Калибровка", "Настройка позиционирования камеры");

  const { step, run, back, loading } = useCalibrationFlow();

  const [photo1, setPhoto1] = useState<string | null>(null);
  const [photo2, setPhoto2] = useState<string | null>(null);
  const [result, setResult] = useState<Result | null>(null);
  const [stepsCount, setStepsCount] = useState<number>(10);

  const stepIndex = order.indexOf(step);

  const handleRun = () => {
    run(stepsCount);
  };

  useEffect(() => {

    const p1 = (e: any) => setPhoto1(e.detail);
    const p2 = (e: any) => setPhoto2(e.detail);
    const res = (e: any) => setResult(e.detail);

    const reset = () => {
      setPhoto1(null);
      setPhoto2(null);
      setResult(null);
      setStepsCount(10);
    };

    window.addEventListener("photo1", p1);
    window.addEventListener("photo2", p2);
    window.addEventListener("calibrationResult", res);
    window.addEventListener("calibrationReset", reset);

    return () => {
      window.removeEventListener("photo1", p1);
      window.removeEventListener("photo2", p2);
      window.removeEventListener("calibrationResult", res);
      window.removeEventListener("calibrationReset", reset);
    };

  }, []);

  return (
    <div className="flex flex-col w-full space-y-[20px]">

      <div className="p-[20px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-card)]">
        <Stepper steps={steps as unknown as string[]} current={stepIndex} />
      </div>

      <div className="grid gap-[20px] grid-cols-1 lg:grid-cols-2">

        <div className="rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-card)] overflow-hidden">
          <div className="px-[20px] py-[14px] border-b border-[var(--border-light)] text-[14px] font-bold text-[var(--text-primary)] text-center">
            Камера
          </div>

          <div className="p-[20px]">
            <div className="h-[360px] bg-black rounded-[8px] flex items-center justify-center text-white/50 text-[13px]">
              live preview
            </div>
          </div>
        </div>

        <div className="rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-card)] overflow-hidden flex flex-col">

          <div className="px-[20px] py-[14px] border-b border-[var(--border-light)] text-[14px] font-bold text-[var(--text-primary)] text-center">
            Управление
          </div>

          <div className="p-[20px] space-y-[16px] flex flex-col flex-1">

            {photo1 && (
              <div className="space-y-[6px]">
                <div className="text-[13px] text-[var(--text-secondary)]">
                  Первое фото
                </div>

                <img
                  src={photo1}
                  className="max-h-[160px] w-full object-contain rounded-[8px] border border-[var(--border-color)]"
                />
              </div>
            )}

            {photo2 && (
              <div className="space-y-[6px]">
                <div className="text-[13px] text-[var(--text-secondary)]">
                  Второе фото
                </div>

                <img
                  src={photo2}
                  className="max-h-[160px] w-full object-contain rounded-[8px] border border-[var(--border-color)]"
                />
              </div>
            )}

            {step === "calculate" && (
              <div className="space-y-[6px]">
                <div className="text-[13px] text-[var(--text-secondary)]">
                  Количество шагов
                </div>

                <input
                  type="number"
                  min={1}
                  value={stepsCount}
                  onChange={(e) => setStepsCount(Number(e.target.value))}
                  className="w-full px-[12px] py-[8px] rounded-[8px] border border-[var(--border-color)] bg-[var(--bg-card)] text-[var(--text-primary)]"
                />
              </div>
            )}

            {result && (
              <div className="p-[14px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] space-y-[6px] text-[13px]">
                <div className="font-medium text-[var(--text-primary)]">
                  Результат калибровки
                </div>

                <div className="text-[var(--text-secondary)]">
                  X shift: <span className="text-[var(--text-primary)]">{result.dx}px</span>
                </div>

                <div className="text-[var(--text-secondary)]">
                  Y shift: <span className="text-[var(--text-primary)]">{result.dy}px</span>
                </div>
              </div>
            )}

            <div className="flex gap-[12px] mt-auto pt-[10px]">

              {step !== "prepare" && (
                <button
                  onClick={back}
                  disabled={loading}
                  className="flex-1 px-[16px] py-[10px] rounded-[10px] border border-[var(--border-color)] text-[var(--text-primary)] hover:bg-[var(--bg-page)] disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Назад
                </button>
              )}

              <button
                onClick={handleRun}
                disabled={loading}
                className="flex-1 px-[16px] py-[10px] rounded-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] hover:bg-[var(--color-primary-hover)] disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? "Выполнение..." : nextLabel[step]}
              </button>

            </div>

          </div>
        </div>

      </div>
    </div>
  );
};

export default CalibrationPage;