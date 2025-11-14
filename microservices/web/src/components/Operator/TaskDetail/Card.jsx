import React, { useEffect, useState } from "react";
import {
  Card,
  Box,
  Typography,
  Divider,
  LinearProgress,
  Button,
} from "@mui/material";
import { SimpleShowLayout, TextField, useRecordContext } from "react-admin";
import { useRobotWS } from "../../hooks/useRobotWS";
import TaskStartDialog from "./TaskStartDialog";
import { useNavigate } from "react-router-dom";

const TaskCard = () => {
  const record = useRecordContext() || {};
  const { rawState, displayState, sendMessage, isBoot, control, amount } = useRobotWS({ record });

  const [openDialog, setOpenDialog] = useState(false);
  const [pendingExtraMode, setPendingExtraMode] = useState("");
  const [isRunning, setIsRunning] = useState(false);
  const [controlState, setControlState] = useState(false);
  const navigate = useNavigate();

  if (!record) return navigate("/tasks");

  const progress = record.required_amount
    ? (record.completed_amount / record.required_amount) * 100
    : 0;

  const [completedAmount, setCompletedAmount] = useState(
    record.completed_amount ?? 0
  );

  const handleOpenStartConfirm = () => setOpenDialog(true);
  const handleCloseStartConfirm = () => setOpenDialog(false);

  const handleConfirmStart = (extraMode) => {
    setOpenDialog(false);
    setPendingExtraMode(extraMode);
    setIsRunning(true);
  };

  const handleStop = () => {
    setIsRunning(false);
    setPendingExtraMode("");
  };

  useEffect(() => {
    if (control !== undefined) {
      setControlState(control);
    }
  }, [control]);

  useEffect(() => {
    if (record.completed_amount > completedAmount) {
      setCompletedAmount(record.completed_amount);
    }
  }, [record.completed_amount]);

  useEffect(() => {
    if (amount > completedAmount) {
      setCompletedAmount(amount);
    }
  }, [amount]);

  useEffect(() => {
    if (!isRunning) return;
    if (completedAmount >= record.required_amount) {
      handleStop();
      return;
    }

    if (rawState === "READY") {
      sendMessage({
        type: "BEGIN",
        params: {
          shift: record.shift,
          number: record.number,
          receipt: record.receipt,
          seed: record.seed,
          amount: completedAmount + 1,
          completed_amount: completedAmount,
          required_amount: record.required_amount,
          bunker: record.bunker,
          gcode: record.gcode,
          extraMode: pendingExtraMode,
        },
      });
    }
  }, [isRunning, rawState]);

  return (
    <Card sx={{ borderRadius: 3, boxShadow: 4, p: 3, width: "100%", maxWidth: 600, margin: "auto", position: "relative" }}>
      <Box display="flex" alignItems="center" justifyContent="center" position="relative" mb={2}>
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
        <Box display="grid" gridTemplateColumns={{ xs: "1fr", sm: "1fr 1fr" }} gap={2} width="100%">
          <Box>
            <Typography variant="subtitle2" color="textSecondary">Номер</Typography>
            <Typography variant="body1"><TextField source="number" /></Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>Смена</Typography>
            <Typography variant="body1"><TextField source="shift" /></Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>Семена</Typography>
            <Typography variant="body1"><TextField source="seed_ru" /></Typography>
          </Box>

          <Box>
            <Typography variant="subtitle2" color="textSecondary">Требуемое количество</Typography>
            <Typography variant="body1"><TextField source="required_amount" /></Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>Выполнено</Typography>
            <Typography variant="body1">{completedAmount}</Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>Контроль качества</Typography>
            <Typography variant="body1" component="div">
              {controlState ? "Пройден" : "Не пройден"}
            </Typography>
          </Box>

          <Box gridColumn="1 / -1" mt={2}>
            <Typography variant="subtitle2" color="textSecondary" mb={0.5}>Прогресс выполнения</Typography>
            <LinearProgress variant="determinate" value={progress} sx={{ height: 10, borderRadius: 5 }} />
          </Box>

          <Box gridColumn="1 / -1" mt={2} display="flex" flexDirection="column" gap={1}>
            {rawState !== "MANUAL_MODE" && (
              <Box display="flex" gap={1}>
                <Button
                  variant="contained"
                  disabled={rawState !== "READY" || isRunning}
                  onClick={handleOpenStartConfirm}
                  sx={{ flex: 1, minHeight: 40, backgroundColor: "#4caf50", color: "white", "&:hover": { backgroundColor: "#45a049" } }}
                >
                  Начать
                </Button>

                <Button
                  variant="contained"
                  onClick={handleStop}
                  disabled={!isRunning}
                  sx={{ flex: 1, minHeight: 40, backgroundColor: "#af5e4c", color: "white" }}
                >
                  Остановить
                </Button>

                <Button
                  variant="outlined"
                  onClick={() => sendMessage({ type: "SETSTATUS READY" })}
                  sx={{ flex: 1, minHeight: 40 }}
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
            <Typography align="center" sx={{ opacity: 0.7, mt: 0.5 }}>
              Состояние устройства: {displayState}
            </Typography>
          </Box>
        </Box>
      </SimpleShowLayout>
    </Card>
  );
};

export default TaskCard;
