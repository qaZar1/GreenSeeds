import { useState } from "react";
import type { TaskRecord } from "../../../types/task";

export function useTaskController(record: TaskRecord) {
  const [isRunning, setIsRunning] = useState(false);
  const [isRetrying, setIsRetrying] = useState(false);
  const [completedAmount, setCompletedAmount] = useState(0);

  const requiredAmount = record.required_amount;

  const progress =
    requiredAmount > 0
      ? Math.round((completedAmount / requiredAmount) * 100)
      : 0;

  const handleSuccessStep = () => {
    setCompletedAmount((prev) => prev + 1);
    setIsRetrying(false);
  };

  const handleErrorStep = (turn: number, error: string) => {
    console.error("Step error:", turn, error);
    setIsRetrying(true);
  };

  return {
    isRunning,
    setIsRunning,
    isRetrying,
    completedAmount,
    requiredAmount,
    progress,
    handleSuccessStep,
    handleErrorStep,
  };
}