import type { StepConfig, StepKey } from "../../../types/calibration";

export const calibrationFlow: Record<StepKey, StepConfig> = {
  prepare: {
    label: "Подготовка",
    next: "photo1"
  },

  photo1: {
    label: "Первое фото",
    action: async () => {
      const res = await fetch("/api/camera/photo");
      const blob = await res.blob();
      window.dispatchEvent(
        new CustomEvent("photo1", { detail: URL.createObjectURL(blob) })
      );
    },
    next: "move"
  },

  move: {
    label: "Перемещение платформы",
    action: async () => {
      await fetch("/api/calibration/move", { method: "POST" });
    },
    next: "photo2"
  },

  photo2: {
    label: "Второе фото",
    action: async () => {
      const res = await fetch("/api/camera/photo");
      const blob = await res.blob();
      window.dispatchEvent(
        new CustomEvent("photo2", { detail: URL.createObjectURL(blob) })
      );
    },
    next: "calculate"
  },

  calculate: {
    label: "Вычисления",
    action: async () => {
      const res = await fetch("/api/calibration/calc", {
        method: "POST"
      });

      const data = await res.json();

      window.dispatchEvent(
        new CustomEvent("calibrationResult", { detail: data })
      );
    },
    next: "save"
  },

  save: {
    label: "Сохранение результата",
    action: async () => {
      await fetch("/api/calibration/save", {
        method: "POST"
      });
    },
    next: null
  }
};