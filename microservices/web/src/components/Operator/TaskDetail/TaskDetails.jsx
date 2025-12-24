import React from "react";
import { useParams } from "react-router-dom";
import { Box } from "@mui/material";
import { Show } from "react-admin";
import TopToolbarWithBackButton from "../../utils/Back";
import TaskCard from "./Card";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";
import { useNotify } from "react-admin";

const TaskDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const notify = useNotify();

  useEffect(() => {
    if (!id || id === "undefined" || id === "null") {
      navigate("/tasks", { replace: true });
    }
  }, [id, navigate]);

  return (
    <Show
      resource="task"
      id={id}
      title="Детали задания"
      sx={{ padding: 2 }}
      actions={<TopToolbarWithBackButton to={`/tasks`} />}
      mutationMode="pessimistic"
      component={Box}
      queryOptions={{
        onError: () => {
          notify("Ошибка загрузки задания", { type: "error" });
        },
      }}
    >
      <TaskCard />
    </Show>
  );
};

export default TaskDetails;
