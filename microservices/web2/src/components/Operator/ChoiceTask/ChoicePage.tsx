import React, { useEffect, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";
import TaskCard from "./TaskCard";

type ChoiceTask = {
  id: number;
  shift: number;
  dt: string;
};

const ChoicePage: React.FC = () => {
  usePageHeader("Выбор задания", "Доступные смены на сегодня");

  const [tasks, setTasks] = useState<ChoiceTask[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadData = async () => {
    setLoading(true);
    setError(false);

    try {
      const res = await api.getList("choice");
      setTasks(res);
    } catch (e) {
      console.error(e);
      toast.error("Не удалось загрузить смены");
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  if (loading) return <SproutLoader />;
  if (error) return <ErrorState onRetry={loadData} />;

  return (
    <div className="space-y-[16px]">

      {!tasks.length ? (
        <div className="text-center py-[40px] text-[var(--text-secondary)]">
          Смен на сегодня нет
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-[16px]">
          {tasks.map(task => (
            <TaskCard key={task.shift} task={task} onTaken={loadData} />
          ))}
        </div>
      )}

    </div>
  );
};

export default ChoicePage;