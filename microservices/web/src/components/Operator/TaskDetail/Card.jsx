import React, { useEffect, useMemo, useState } from "react";
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
import DecisionModal from "./DecisionModal";

const TaskCard = () => {
  const record = useRecordContext() || {};
  const {
    rawState,
    displayState,
    sendMessage,
    isBoot,
    control,
    amount,
    decisionModal,
    setDecisionModal
  } = useRobotWS(
    { record,
      onSuccessStep: () => {
        setFailedReports(prev =>
          prev.filter((_, i) => i !== completedAmount)
        );

        setCompletedAmount(prev => prev + 1);
      },
    onErrorStep: (turn, error) => {
      setFailedReports(prev => {
        const exists = prev.find(r => r.turn === turn);

        if (exists) {
          return prev.map(r =>
            r.turn === turn ? { ...r, error } : r
          );
        }

        return [...prev, { turn, error }];
      });
    },
  });

  const [openDialog, setOpenDialog] = useState(false);
  const [pendingExtraMode, setPendingExtraMode] = useState("");
  const [isRunning, setIsRunning] = useState(false);
  const [controlState, setControlState] = useState(null);
  const navigate = useNavigate();

  if (!record) return navigate("/tasks");

  const [failedReports, setFailedReports] = useState([]);


  const requiredAmount = record.required_amount || 0;
  const progress = requiredAmount > 0
  ? ((requiredAmount - failedReports.length) / requiredAmount) * 100
  : 0;

  const [completedAmount, setCompletedAmount] = useState(0);

  const currentTaskNumber = completedAmount + 1; 

  const isRetrying = useMemo(() => {
    if (!isRunning) return false;
    
    const failed = failedReports.find(
      r => Number(r.turn) === currentTaskNumber
    );

    return failed && failed.error && failed.error !== "";
  }, [isRunning, failedReports, currentTaskNumber]);


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

  const sendDecision = (ok, comment) => {
    if (!sendMessage) return; 
    sendMessage({
      type: "DECISION",
      status: ok ? "OK" : "NOK",
      solution: comment,
    });
    setDecisionModal({
      open: false,
      reason: "",
      photo: null,
    });
  };

  useEffect(() => {
    const failed = Array.isArray(record.reports)
      ? record.reports.filter(r => !r.success)
      : [];

    setFailedReports(failed);
  }, [record]);

  useEffect(() => {
    setControlState(control);
  }, [control]);

  useEffect(() => {
    if (!record) return;

    const failedReports = Array.isArray(record.reports)
      ? record.reports.filter(r => !r.success)
      : [];

    const realCompleted = record.required_amount - failedReports.length;
    setCompletedAmount(realCompleted);
  }, [record]);

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

    if (rawState === "READY" && isBoot === true) {
      const nextReport = failedReports.find(r => Number(r.turn) === currentTaskNumber);

      if (!nextReport) return;

      sendMessage({
        type: "BEGIN",
        params: {
          shift: record.shift,
          number: record.number,
          receipt: record.receipt,
          seed: record.seed,
          turn: currentTaskNumber,
          completed_amount: completedAmount,
          required_amount: record.required_amount,
          bunker: record.bunker,
          gcode: record.gcode,
          extraMode: pendingExtraMode,
        },
      });
    }
  }, [isRunning, rawState]);

  // --- UI СТИЛИЗАЦИЯ ---
  const statusColor = rawState === "READY" ? "success.main" :
                      rawState.includes("ERROR") || rawState.includes("FAULT") ? "error.main" :
                      rawState === "MANUAL_MODE" ? "warning.main" :
                      "info.main";

  const controlTextColor = controlState === true ? "success.dark" :
                           controlState === false ? "error.dark" :
                           "text.primary";


  const DataLabel = ({ value, color = 'text.primary'}) => (
    <Box>
        <Typography
            variant="body1"
            sx={{ color: color }}
        >
            {value}
        </Typography>
    </Box>
  );

  return (
    <Card
      sx={{
        borderRadius: 3,
        boxShadow: 4,
        p: { xs: 1, sm: 3 },
        maxWidth: 600,
        mx: "auto",
        boxSizing: "border-box"
      }}
    >
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
          Задание №{record.number}
        </Typography>
      </Box>

      <Divider sx={{ mb: 2 }} />

      <Box mb={3} p={1} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, textAlign: 'center' }}>
        <Typography variant="subtitle2" color="textSecondary">СЕМЕНА</Typography>
        <Typography variant="body1" textTransform={"uppercase"}>{record.seed_ru}</Typography>
      </Box>

      <SimpleShowLayout>
        <Box
          display="grid"
          gridTemplateColumns={{ xs: "1fr", sm: "1fr 1fr" }}
          width="100%"
          gap={2}
        >
          <Box
            justifySelf={{ xs: "stretch", sm: "start" }}
            textAlign={{ xs: "left", sm: "left" }}
          >
            <Typography variant="subtitle2" color="textSecondary" mt={1}>
              СМЕНА
            </Typography>
            <Typography variant="body1">{record.shift}</Typography>
          </Box>

          <Box
            justifySelf={{ xs: "stretch", sm: "end" }}
            textAlign={{ xs: "left", sm: "right" }}
          >
            <Typography variant="subtitle2" color="textSecondary" mt={1}>
              КОНТРОЛЬ КАЧЕСТВА
            </Typography>
            <DataLabel
              value={
                controlState === null
                  ? "Ожидание проверки"
                  : controlState
                    ? "Пройден"
                    : "Не пройден!"
              }
              color={controlTextColor}
            />
          </Box>
        </Box>

        {/* <Box display="grid" gridTemplateColumns={{ xs: "1fr", sm: "1fr 1fr" }} gap={2} width="100%">
          <Box>
            <Typography variant="subtitle2" color="textSecondary" mt={1}>СМЕНА</Typography>
            <Typography variant="body1">{record.shift}</Typography>
          </Box>

          <Box>
            <Typography variant="subtitle2" color="textSecondary" mt={1}>КОНТРОЛЬ КАЧЕСТВА</Typography>
            <DataLabel
              value={
                controlState === null
                  ? "Ожидание проверки"
                  : controlState
                    ? "Пройден"
                    : "Не пройден!"
              }
              color={controlTextColor}
            />
          </Box>

          {/* <Box>
            <Typography variant="subtitle2" color="textSecondary" mt={1}>ВЫПОЛНЕНО</Typography>
            <Typography variant="body1">
              {isNaN(completedAmount) ? 0 : completedAmount}
            </Typography>

            <Typography variant="subtitle2" color="textSecondary" mt={1}>ТРЕБУЕМОЕ КОЛИЧЕСТВО</Typography>
            <Typography variant="body1">{record.required_amount}</Typography>
          </Box>
        </Box> */}
        </SimpleShowLayout>
        <Box>
          <Divider sx={{ mb: 2 }} />

          <Box gridColumn="1 / -1" mt={2}>
            <Typography variant="subtitle2" color="textSecondary" mb={0.5} align="center" fontWeight="normal">
              Прогресс выполнения: {completedAmount} из {requiredAmount}
            </Typography>
            <LinearProgress variant="determinate" value={isNaN(progress) ? 0 : progress} sx={{ height: 10, borderRadius: 5 }} />
          </Box>

          <Box 
            mt={2} 
            p={2} 
            borderRadius={2} 
            sx={{ 
              border: `2px solid ${statusColor}`,
              bgcolor: 'background.paper',
            }}
          >
            <Typography
              variant="subtitle1"
              align="center"
              fontWeight="bold"
              color={statusColor}
            >
              СОСТОЯНИЕ: {displayState}
            </Typography>
            
            {isRunning && (
              <Box mt={0.5}>
                <Typography
                  variant="subtitle1"
                  align="center"
                  fontWeight="bold"
                  color="text.primary"
                >
                  ЛОТОК: {currentTaskNumber} из {requiredAmount}
                </Typography>
                
                {isRetrying && (
                  <Typography
                    variant="body2"
                    align="center"
                    mt={0.5}
                    color="error.dark"
                    fontWeight="bold"
                  >
                    ⚠️ ПОВТОРНЫЙ ПОСЕВ.
                  </Typography>
                )}
              </Box>
            )}
          </Box>
          

          <Box gridColumn="1 / -1" mt={2} display="flex" flexDirection="column" gap={1}>
            {rawState !== "MANUAL_MODE" && (
              <Box display="flex" gap={1} flexWrap="wrap">
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
          </Box>
        <DecisionModal
          open={decisionModal.open}
          reason={decisionModal.reason}
          photo={decisionModal.photo}
          onConfirm={(comment) => sendDecision(true, comment)}
          onReject={(comment) => sendDecision(false, comment)}
        />
      </Box>
    </Card>
  );
};

export default TaskCard;
