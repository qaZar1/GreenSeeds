import React from "react";
import { Card, Box, Typography, CardContent, Button, Divider } from "@mui/material";
import { DateField, useNotify } from "react-admin";
import AssignmentTurnedInIcon from "@mui/icons-material/AssignmentTurnedIn";
import { useUpdate } from "react-admin";

const Task = ({ task, username }) => {
  const [update] = useUpdate();
  const notify = useNotify();

  const handleTakeTask = async () => {
    try {
      await update(
        "choice", // ресурс
        {
          id: task.shift, // обязательно!
          data: {
            id: task.shift, // чтобы попало в bodyData.shift
            username: username, // логин текущего пользователя
            dt: task.dt,
          },
        },
        {
          onSuccess: () => {
            notify(`Задание ${task.id} успешно взято`);
            window.location.reload();
          },
          onError: (error) => {
            console.error("Ошибка обновления:", error);
          },
        }
      );
    } catch (e) {
      console.error("Ошибка при взятии задания:", e);
    }
  };


  return (
    <Card
      key={task.id}
      sx={{
        flex: "1 1 350px",       // базовая ширина карточки
        minWidth: "280px",       // не даёт стать слишком узкой
        maxWidth: "100%",        // растягивается, если остаётся место
        boxShadow: 4,
        borderRadius: 3,
        overflow: "hidden",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
      }}
    >
      {/* Верхняя часть — “шапка” */}
      <Box
        sx={{
          backgroundColor: "#2e7d32", // зелёный MUI success.dark
          color: "white",
          p: 2,
          display: "flex",
          alignItems: "center",
          gap: 1.5,
        }}
      >
        <AssignmentTurnedInIcon sx={{ fontSize: 30 }} />
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          Сменное задание №{task.shift}
        </Typography>
      </Box>

      {/* Основное содержимое */}
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography variant="body1" sx={{ mb: 1 }}>
          <strong>Дата:</strong>{" "}
          <DateField source="dt" record={task} showTime={true} />
        </Typography>

        <Typography variant="body1" sx={{ mb: 1 }}>
          <strong>Ответственный:</strong>{" "}
          {task.username || <em>не назначен</em>}
        </Typography>

        <Typography variant="body2" sx={{ color: "text.secondary" }}>
          Задания будут показаны после начала смены
        </Typography>
      </CardContent>

      <Divider />

      {/* Кнопка снизу */}
      <Box sx={{ display: "flex", justifyContent: "center", p: 2, pt: 1 }}>
        <Button
          variant="contained"
          color="success"
          sx={{
            textTransform: "none",
            borderRadius: 3,
            px: 5,
            py: 1.2,
            fontWeight: 600,
            fontSize: "1rem",
          }}
          onClick={() => handleTakeTask()}
        >
          Взять задание
        </Button>
      </Box>
    </Card>
  );
};

export default Task;
