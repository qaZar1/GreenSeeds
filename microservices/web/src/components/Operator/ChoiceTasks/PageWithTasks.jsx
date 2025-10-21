import React, { useState, useEffect } from "react";
import { Box, CircularProgress } from "@mui/material";
import { jwtDecode } from "jwt-decode";
import { ListContextProvider, useListController } from "react-admin";
import Task from "./Task";
import EmptyTask from "./EmptyTask";
import { Card, CardContent, List, ListItem, ListItemText, Typography } from "@mui/material";
import Tasks from "./Tasks";

const PageWithTasks = () => {
  const [username, setUsername] = useState(null);
  const [activeTasks, setActiveTasks] = useState([]);
  const [isLoadingActive, setIsLoadingActive] = useState(true);

  // react-admin контроллер (для списка всех заданий)
  const listContext = useListController({ resource: "tasks" });

  useEffect(() => {
    const loadActiveTasks = async () => {
      try {
        const stored = localStorage.getItem("auth");
        if (!stored) return;
        const parsed = JSON.parse(stored);
        if (!parsed?.token) return;

        const decoded = jwtDecode(parsed.token);
        const name = decoded?.username;
        setUsername(name);

        const response = await fetch(
          `/api/assignments/active-tasks/${name}`,
          {
            headers: {
              Authorization: `Bearer ${parsed.token}`,
            },
          }
        );
        if (response.status === 404) return;
        else if (!response.ok) throw new Error("Ошибка загрузки активных заданий");

        const data = await response.json();
        setActiveTasks(data || []);
      } catch (e) {
        console.warn("Ошибка при загрузке активных заданий:", e);
      } finally {
        setIsLoadingActive(false);
      }
    };

    loadActiveTasks();
  }, []);

  if (isLoadingActive || listContext.isLoading)
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
        <CircularProgress />
      </Box>
    );

  // 🔹 Если есть активные задания → показываем только их
  if (activeTasks && activeTasks.length > 0) {
    return (
        <Tasks tasks={activeTasks} />
      );
  }

  // 🔹 Иначе показываем список всех заданий (через react-admin dataProvider)
  return (
    <ListContextProvider value={listContext}>
      <TasksListFromContext username={username} />
    </ListContextProvider>
  );
};

const TasksListFromContext = ({ username }) => {
  const { data, isLoading } = useListController({ resource: "tasks" });

  if (isLoading)
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
        <CircularProgress />
      </Box>
    );

  if (!data || data.length === 0) return <EmptyTask />;

  return (
    <Box
      display="flex"
      flexWrap="wrap"
      justifyContent="center"
      gap={2}
      p={2}
    >
      {data.map((task) => (
        <Task key={task.id || task.shift} task={task} username={username} />
      ))}
    </Box>
  );
};

export default PageWithTasks;
