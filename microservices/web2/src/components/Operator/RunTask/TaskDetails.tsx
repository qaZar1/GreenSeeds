import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import TaskCard from "./Task";
import type { TaskRecord } from "../../../types/task";
import NotFoundPage from "../../pages/NotFoundPage";
import SproutLoader from "../../utils/Loader/SproutLoader";

const TaskDetails = () => {
  const { id } = useParams();

  const [task, setTask] = useState<TaskRecord | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;

    const load = async () => {
      try {
        const res = await api.getOne("task", id);
        setTask(res.data);
      } catch {
        toast.error("Ошибка загрузки задания");
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [id]);

  if (loading) return <SproutLoader />;
  if (!task) return <NotFoundPage />;

  return <TaskCard record={task} />;
};

export default TaskDetails;