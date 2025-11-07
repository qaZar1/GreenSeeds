import React, { useState } from "react";
import {
  Card,
  CardContent,
  Typography,
  Divider,
  Box,
  LinearProgress,
  Button,
} from "@mui/material";
import { SimpleShowLayout, TextField, useRecordContext } from "react-admin";
import { useRobotWS } from "../../hooks/useRobotWs";
import TaskStartDialog from "./TaskStartDialog";

const TaskCard = () => {
  const record = useRecordContext();
  const [openDialog, setOpenDialog] = useState(false);

  const handleOpenStartConfirm = () => setOpenDialog(true);
  const handleCloseStartConfirm = () => setOpenDialog(false);

  const { rawState, displayState, sendCommand, sendGcode, isBoot } = useRobotWS(record);

  if (!record) return null;

  const progress = record.required_amount
    ? (record.completed_amount / record.required_amount) * 100
    : 0;

  const handleConfirmStart = (extraMode) => {
    setOpenDialog(false);
    startTask(extraMode);
  };

  const startTask = (extraMode) => {
    if (!record?.gcode || !record?.bunker || !record?.shift || !record?.number) return;
    let displayText = "";
    if (extraMode) {
      // Добавим спец заголовок
      displayText = `\x07Sorrel\x0A${record.shift}/${record.number}\x0A${record.completed_amount+1}/${record.required_amount}\x0A\x0D`;
    }
    sendGcode(
      record.gcode,
      record.bunker,
      record.shift,
      record.number,
      record.completed_amount + 1,
      displayText
    );
  };

  return (
    <Card
      sx={{
        borderRadius: 3,
        boxShadow: 4,
        p: 3,
        width: "100%",
        maxWidth: 600,
        margin: "auto",
        position: "relative",
        overflow: "hidden",
      }}
    >
      {/* Заголовок и индикатор */}
      <Box
        display="flex"
        alignItems="center"
        justifyContent="center"
        position="relative"
        mb={2}
      >
        <Box
          sx={{
            position: "absolute",
            left: 0,
            width: 16,
            height: 16,
            borderRadius: "50%",
            bgcolor: isBoot ? "green" : "red",
            border: "1px solid #000",
          }}
        />
        <Typography variant="h5" fontWeight="bold" textAlign="center">
          Детали задания
        </Typography>
      </Box>

      <Divider sx={{ mb: 2 }} />

      <SimpleShowLayout>
        <Box
          display="grid"
          gridTemplateColumns={{ xs: "1fr", sm: "1fr 1fr" }}
          gap={2}
          width="100%"
        >
          {/* Левая колонка */}
          <Box>
            <Typography variant="subtitle2" color="textSecondary">
              Номер
            </Typography>
            <Typography variant="body1">
              <TextField source="number" />
            </Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>
              Смена
            </Typography>
            <Typography variant="body1">
              <TextField source="shift" />
            </Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>
              Семена
            </Typography>
            <Typography variant="body1">
              <TextField source="seed" />
            </Typography>
          </Box>

          {/* Правая колонка */}
          <Box>
            <Typography variant="subtitle2" color="textSecondary">
              Требуемое количество
            </Typography>
            <Typography variant="body1">
              <TextField source="required_amount" />
            </Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>
              Выполнено
            </Typography>
            <Typography variant="body1">
              <TextField source="completed_amount" />
            </Typography>
          </Box>

          {/* Прогресс */}
          <Box gridColumn="1 / -1" mt={2}>
            <Typography variant="subtitle2" color="textSecondary" mb={0.5}>
              Прогресс выполнения
            </Typography>
            <LinearProgress
              variant="determinate"
              value={progress}
              sx={{ height: 10, borderRadius: 5 }}
            />
          </Box>

          {/* Кнопки */}
          <Box gridColumn="1 / -1" mt={2} display="flex" flexDirection="column" gap={1}>
          {rawState !== "MANUAL_MODE" && (
            <Box display="flex" gap={1}>
              <Button
                variant="contained"
                onClick={() => sendCommand("STATUS")}
                sx={{ flex: 1, minHeight: 40 }}
              >
                Узнать статус робота
              </Button>

              <Button
                variant="contained"
                disabled={rawState !== "READY"}
                onClick={handleOpenStartConfirm}
                sx={{
                  flex: 1,
                  minHeight: 40,
                  backgroundColor: "#4caf50",
                  color: "white",
                  "&:hover": { backgroundColor: "#45a049" },
                }}
              >
                Начать
              </Button>

              <Button
                variant="contained"
                color="error"
                disabled={rawState !== "END"}
                onClick={() => sendCommand("RETURN")}
                sx={{ flex: 1, minHeight: 40 }}
              >
                Вернуть каретку
              </Button>

              <TaskStartDialog
                open={openDialog}
                onClose={handleCloseStartConfirm}
                onConfirm={handleConfirmStart}
                task={record}
              />
            </Box>
          )}
            <Button
              variant="outlined"
              color="success"
              fullWidth
              onClick={() => sendCommand("SETSTATUS READY")}
            >
              DEV: Set READY
            </Button>
            <Button
              variant="outlined"
              color="success"
              fullWidth
              onClick={() => sendCommand("SETSTATUS MANUAL_MODE")}
            >
              DEV: Set MANUAL MODE
            </Button>

            <Typography align="center" sx={{ opacity: 0.7, mt: 0.5, mb: 0 }}>
              Состояние робота: {displayState}
            </Typography>
          </Box>
        </Box>
      </SimpleShowLayout>
    </Card>
  );
};

export default TaskCard;
