import React, { useEffect, useState } from "react";
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

const control = "Пройден"

const TaskCard = () => {
  const record = useRecordContext();

  if (!record) return null;

  const [openDialog, setOpenDialog] = useState(false);
  const [pendingExtraMode, setPendingExtraMode] = useState("");
  const [isRunning, setIsRunning] = useState(false);

  const handleOpenStartConfirm = () => setOpenDialog(true);
  const handleCloseStartConfirm = () => setOpenDialog(false);

  const { rawState, displayState, sendCommand, sendGcode, isBoot } = useRobotWS(record);

  const progress = record.required_amount
    ? (record.completed_amount / record.required_amount) * 100
    : 0;

  const handleConfirmStart = (extraMode) => {
    setOpenDialog(false);
    setPendingExtraMode(extraMode);
    setIsRunning(true);
    set
  };

  const handleStop = () => {
    setIsRunning(false);
    setPendingExtraMode("");
  };

  useEffect(() => {
    if (!isRunning) return;
    if (record?.completed_amount >= record?.required_amount) {
      handleStop();
    }

    if (rawState === "READY") {
      startTask(pendingExtraMode);
    }
  }, [isRunning, rawState, record.completed_amount, record.required_amount]);

  const startTask = (extraMode) => {
    if (!record?.gcode || !record?.bunker || !record?.shift || !record?.number) return;
    let displayText = "";
    if (extraMode) {
      // Код отображения на дисплее
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
    setIsRunning(true);
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
              <TextField source="seed_ru" />
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
            <Typography variant="subtitle2" color="textSecondary" mt={1}>
              Контроль качества
            </Typography>
            <Typography variant="body1">
              <Typography variant="body2">{control}</Typography>
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
                disabled={rawState !== "READY" || isRunning}
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
                onClick={handleStop}
                disabled={!isRunning}
                sx={{
                  flex: 1,
                  minHeight: 40,
                  backgroundColor: "#af5e4c",
                  color: "white",
                  "&:hover": { backgroundColor: "#af5e4c" },
                }}
              >
                Остановить
              </Button>
              <Button
                variant="outlined"
                onClick={() => sendCommand("SETSTATUS READY")}
                sx={{
                  flex: 1,
                  minHeight: 40,
                }}
              >
                DEV: готов
              </Button>

              <TaskStartDialog
                open={openDialog}
                onClose={handleCloseStartConfirm}
                onConfirm={handleConfirmStart}
                task={record}
              />
            </Box>
          )}

            <Typography align="center" sx={{ opacity: 0.7, mt: 0.5, mb: 0 }}>
              Состояние устройства: {displayState}
            </Typography>
          </Box>
        </Box>
      </SimpleShowLayout>
    </Card>
  );
};

export default TaskCard;
