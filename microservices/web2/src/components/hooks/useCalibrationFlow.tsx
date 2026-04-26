import { useRef, useState } from "react";
import toast from "react-hot-toast";
import { calibrationApi } from "../Admin/Calibration/CalibrationApi";

export type Step =
  | "prepare"
  | "photo1"
  | "move"
  | "photo2"
  | "calculate"
  | "save";

const order: Step[] = [
  "prepare",
  "photo1",
  "move",
  "photo2",
  "calculate",
  "save"
];

export const useCalibrationFlow = () => {

  const [step, setStep] = useState<Step>("prepare");
  const [loading, setLoading] = useState(false);

  const sessionRef = useRef<string | null>(null);

  const prepareHandle = async () => {
    try {

      const res = await calibrationApi.handshake();

      const sessionId = res.headers.get("X-Calibration-Session");

      if (!sessionId) {
        toast.error("Не удалось получить сессию");
        return;
      }

      sessionRef.current = sessionId;

      setStep("photo1");

    } catch {
      toast.error("Не удалось начать калибровку");
    }
  };

  const photo1Handle = async () => {

    try {

      if (!sessionRef.current) {
        toast.error("Сессия отсутствует");
        return;
      }

      const res = await calibrationApi.photo(1, sessionRef.current);

      const blob = res.data;

      const url = URL.createObjectURL(blob);

      window.dispatchEvent(
        new CustomEvent("photo1", { detail: url })
      );

      setStep("move");

    } catch {
      toast.error("Не удалось сделать фото");
    }

  };

  const moveHandle = async () => {
    setStep("photo2");
  };

  const photo2Handle = async () => {

    try {

      if (!sessionRef.current) {
        toast.error("Сессия отсутствует");
        return;
      }

      const res = await calibrationApi.photo(2, sessionRef.current);

      const blob = res.data;

      const url = URL.createObjectURL(blob);

      window.dispatchEvent(
        new CustomEvent("photo2", { detail: url })
      );

      setStep("calculate");

    } catch {
      toast.error("Не удалось сделать второе фото");
    }

  };

  const calculateHandle = async (stepsCount?: number) => {

    try {

      if (!sessionRef.current) {
        toast.error("Сессия отсутствует");
        return;
      }

      const res = await calibrationApi.calculate(
        stepsCount ?? 10,
        sessionRef.current
      );

      const result = res.data;

      window.dispatchEvent(
        new CustomEvent("calibrationResult", {
          detail: result
        })
      );

      setStep("save");

    } catch {
      toast.error("Ошибка расчёта");
    }

  };

  const saveHandle = async () => {

    try {

      if (!sessionRef.current) {
        toast.error("Сессия отсутствует");
        return;
      }

      await calibrationApi.save(sessionRef.current);

      toast.success("Калибровка сохранена");
      
      window.dispatchEvent(new Event("calibrationReset"));

      sessionRef.current = null;

      setStep("prepare");

    } catch {
      toast.error("Ошибка сохранения");
    }

  };

  const run = async (stepsCount?: number) => {

    setLoading(true);

    try {

      switch (step) {

        case "prepare":
          await prepareHandle();
          break;

        case "photo1":
          await photo1Handle();
          break;

        case "move":
          await moveHandle();
          break;

        case "photo2":
          await photo2Handle();
          break;

        case "calculate":
          await calculateHandle(stepsCount);
          break;

        case "save":
          await saveHandle();
          break;

      }

    } finally {
      setLoading(false);
    }

  };

  const back = () => {

    const i = order.indexOf(step);

    if (i > 0) {
      setStep(order[i - 1]);
    }

  };

  return {
    step,
    run,
    back,
    loading
  };
};