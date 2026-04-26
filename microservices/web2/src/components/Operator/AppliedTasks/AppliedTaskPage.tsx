import React, { useEffect, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";
import TasksBlock from "./TasksBlock";
import { useAuth } from "../../../context/AuthContext";

type Task = {
  id: number;
  shift: number;
  number: number;
  amount: number;
  seed_ru?: string;
};

const AppliedTasksPage: React.FC = () => {
  usePageHeader("Задания", "Ваши активные задания");

  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const { user_id } = useAuth();

  const loadData = async () => {
    setLoading(true);
    setError(false);

    try {
      console.log(user_id);
      const res = await api.getList("tasks", user_id);
      setTasks(res);
    } catch (e) {
      console.error(e);
      toast.error("Не удалось загрузить задания");
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

  return <TasksBlock tasks={tasks} />;
};

export default AppliedTasksPage;