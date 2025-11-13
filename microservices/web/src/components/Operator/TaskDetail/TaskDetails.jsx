import React from "react";
import { useParams } from "react-router-dom";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Divider,
  Stack,
  Chip,
} from "@mui/material";
import { Show, SimpleShowLayout, TextField } from "react-admin";
import TopToolbarWithBackButton from "../../utils/Back";
import TaskCard from "./Card";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";

const TaskDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();

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
    >
      <TaskCard />
    </Show>
  );
};

export default TaskDetails;
