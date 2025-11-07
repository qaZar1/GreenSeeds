import React, { useState, useEffect } from "react";
import { Box, CircularProgress, Typography, Card, CardContent } from "@mui/material";
import { jwtDecode } from "jwt-decode";
import EmptyTasks from "./EmptyTasks";
import { List } from "react-admin";
import { useMediaQuery } from "@mui/material";
import TasksListContent from "./Controller";

const AppliedTaskList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));

    const [username, setUsername] = useState(null);

    useEffect(() => {
      try {
        const stored = localStorage.getItem("auth");
        if (stored) {
          const parsed = JSON.parse(stored);
          if (parsed?.token) {
            const decoded = jwtDecode(parsed.token);
            setUsername(decoded?.username);
          }
        }
      } catch (e) {
        console.warn("Ошибка получения профиля:", e);
      }

    }, []);

    if (!username) {
      return (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
          <CircularProgress />
        </Box>
      );
    }

    return (
        <List
            resource="tasks"
            empty={<EmptyTasks />}
            {...props}
            pagination={false}
            perPage={false}
            sx={{ padding: 2, width: "100%" }}
            actions={false}
            title="Задания на смену"
            component={Box}
            filter={{ username }}
        >
            <TasksListContent isSmall={isSmall} isMedium={isMedium}/>
        </List>
    );
};

export default AppliedTaskList;