import React, { useEffect, useState } from "react";
import {
  Box,
  Typography,
  Button,
  Stepper,
  Step,
  StepLabel,
  Paper,
  TextField,
  useMediaQuery,
  useTheme,
  Stack,
  Grid,
} from "@mui/material";
import { useNotify, useDataProvider } from "react-admin";

const steps = [
  "Подключение устройства",
  "Снимок №1",
  "Снимок №2",
  "Вычисление шага",
  "Проверка и сохранение",
];

const formatNum = (n) => Math.abs(Number(n)).toFixed(3);

export default function CalibrationPage() {
  const dataProvider = useDataProvider();
  const notify = useNotify();

  const [activeStep, setActiveStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isSaved, setIsSaved] = useState(false);
  const [firstPhoto, setFirstPhoto] = useState(null);
  const [secondPhoto, setSecondPhoto] = useState(null);

  const [calibration, setCalibration] = useState({
    deviceReady: false,
    distanceX: null,
    distanceY: null,
    cir: "",
    dPerStep: null,
  });

  const theme = useTheme();
  const isSmall = useMediaQuery(theme.breakpoints.down("sm"));

  /* ============================
     ШАГ 1 — HANDSHAKE
     ============================ */
  const handleHandshake = async () => {
    setLoading(true);
    try {
      const { data } = await dataProvider.create("calibration");

      localStorage.setItem("sessionId", data.session_id);

      setCalibration((p) => ({ ...p, deviceReady: true }));
      setActiveStep(1);
    } catch {
      notify("Устройство не ответило", { type: "error" });
    } finally {
      setLoading(false);
    }
  };

  /* ============================
     ШАГ 2 / 3 — ФОТО ЧЕРЕЗ BACKEND
     ============================ */
  const takeFirstPhoto = async () => {
    try {
      setLoading(true);
      const { data } = await dataProvider.create("takePhoto", {
        data: {
          sessionId: localStorage.getItem("sessionId"),
          numberOfPhoto: 1,
        },
      });

      setFirstPhoto(data.photo);
    } catch (e) {
      console.error("Ошибка:", e);
      notify("Ошибка при съёмке", { type: "error" });
    } finally {
      setLoading(false);
    }
  };

  const takeSecondPhoto = async () => {
    try {
      setLoading(true);
      const { data } = await dataProvider.create("takePhoto", {
        data: {
          sessionId: localStorage.getItem("sessionId"),
          numberOfPhoto: 2,
        },
      });

      setSecondPhoto(data.photo);
    } catch (e) {
      console.error("Ошибка:", e);
      notify("Ошибка при съёмке", { type: "error" });
    } finally {
      setLoading(false);
    }
  };

  /* ============================
     ОЧИСТКА ФОТО НА BACKEND
     ============================ */
  const clearPhotos = async () => {
    try {
      await dataProvider.custom("devices/calibration/clear", {
        method: "POST",
        body: { id: calibration.sessionId },
      });
    } catch {
      // намеренно игнорируем
    }
  };

  /* ============================
     НАЗАД С ОЧИСТКОЙ
     ============================ */
  const handleBack = async () => {
    if (activeStep > 1) {
      await clearPhotos();
      setCalibration((p) => ({ ...p, photo1: null, photo2: null, cir: "", dPerStep: null }));
      setFirstPhoto(null);
      setSecondPhoto(null);
    }

    setActiveStep((s) => Math.max(0, s - 1));
  };

  /* ============================
     ШАГ 4 — ПОДСЧЕТ
     ============================ */
  const handleCalc = async () => {
    setLoading(true);
    try {
      const { data } = await dataProvider.create("foundCalibration", {
        data: {
          sessionId: localStorage.getItem("sessionId"),
          cir: calibration.cir,
        },
      });

      setCalibration((p) => ({ ...p, distanceX: data.dx, distanceY: data.dy, dPerStep: data.d_per_step }));
      notify("Данные просчитаны", { type: "success" });
    } catch {
      notify("Ошибка просчета данных", { type: "error" });
    } finally {
      setLoading(false);
    }
  };

  /* ============================
     ШАГ 5 — СОХРАНЕНИЕ
     ============================ */
  const handleSave = async () => {
    setLoading(true);
    try {
      await dataProvider.create("saveCalibration", {
        data: {
          sessionId: localStorage.getItem("sessionId"),
          dPerStep: calibration.dPerStep,
        },
      });

      localStorage.removeItem("sessionId");
      setIsSaved(true);

      notify("Калибровка сохранена", { type: "success" });
    } catch {
      notify("Ошибка сохранения калибровки", { type: "error" });
    } finally {
      setLoading(false);
    }
  };

  /* ============================
     ЗАЩИТА ШАГОВ
     ============================ */
  useEffect(() => {
    if (activeStep === 1 && !calibration.deviceReady) setActiveStep(0);
  }, [activeStep, calibration]);

  /* ============================
     UI
     ============================ */
  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        overflow: "auto",
        p: { xs: 2, sm: 4 },
      }}
    >
      <Paper
        elevation={3}
        sx={{
          width: "100%",
          maxWidth: { xs: 480, sm: 720 },
          p: { xs: 2, sm: 4 },
          borderRadius: 2,
        }}
      >
        <Box sx={{ mb: 2 }}>
          <Typography variant={isSmall ? "h6" : "h5"} fontWeight={600} gutterBottom>
            Калибровка устройства
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Следуйте шагам для корректной калибровки
          </Typography>
        </Box>

        <Stepper
          activeStep={activeStep}
          alternativeLabel={!isSmall}
          orientation={isSmall ? "vertical" : "horizontal"}
          sx={{ mb: 2 }}
        >
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel
                sx={{
                  '& .MuiStepLabel-label': {
                    fontSize: isSmall ? '0.72rem' : undefined,
                    whiteSpace: 'normal',
                  },
                }}
              >
                {label}
              </StepLabel>
            </Step>
          ))}
        </Stepper>

        <Box
          sx={{
            minHeight: isSmall ? 240 : 180,
            mb: 2,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            border: "1px dashed",
            borderColor: "divider",
            borderRadius: 2,
            p: 2,
          }}
        >
          {/* Step 0 */}
          {activeStep === 0 && (
            <Button variant="contained" onClick={handleHandshake} disabled={loading} fullWidth={isSmall}>
              Подключить устройство
            </Button>
          )}

          {/* Step 1 - first photo */}
          {activeStep === 1 && !firstPhoto && (
            <Button variant="contained" onClick={takeFirstPhoto} disabled={loading} fullWidth={isSmall}>
              Сделать первый снимок
            </Button>
          )}

          {firstPhoto && activeStep === 1 && (
            <Stack spacing={2} alignItems="center" width="100%">
              <Box sx={{ display: "flex", justifyContent: "center", width: '100%' }}>
                <img
                  src={`data:image/png;base64,${firstPhoto}`}
                  alt="Первый снимок"
                  style={{ maxWidth: '100%', width: isSmall ? 160 : 200, height: 'auto', borderRadius: 8 }}
                />
              </Box>

              <Grid container spacing={1} justifyContent="center">
                <Grid item xs={12} sm="auto">
                  <Button
                    variant="outlined"
                    onClick={() => setFirstPhoto(null)}
                    fullWidth={isSmall}
                  >
                    Переснять
                  </Button>
                </Grid>

                <Grid item xs={12} sm="auto">
                  <Button
                    variant="contained"
                    onClick={() => setActiveStep(2)}
                    fullWidth={isSmall}
                  >
                    Далее
                  </Button>
                </Grid>
              </Grid>
            </Stack>
          )}

          {/* Step 2 - second photo */}
          {activeStep === 2 && !secondPhoto && (
            <Button variant="contained" onClick={takeSecondPhoto} disabled={loading} fullWidth={isSmall}>
              Сделать второй снимок
            </Button>
          )}

          {secondPhoto && activeStep === 2 && (
            <Stack spacing={2} alignItems="center" width="100%">
              <Box sx={{ display: "flex", justifyContent: "center", width: '100%' }}>
                <img
                  src={`data:image/png;base64,${secondPhoto}`}
                  alt="Второй снимок"
                  style={{ maxWidth: '100%', width: isSmall ? 160 : 200, height: 'auto', borderRadius: 8 }}
                />
              </Box>

              <Grid container spacing={1} justifyContent="center">
                <Grid item xs={12} sm="auto">
                  <Button
                    variant="outlined"
                    onClick={() => setSecondPhoto(null)}
                    fullWidth={isSmall}
                  >
                    Переснять
                  </Button>
                </Grid>

                <Grid item xs={12} sm="auto">
                  <Button
                    variant="contained"
                    onClick={() => setActiveStep(3)}
                    fullWidth={isSmall}
                  >
                    Далее
                  </Button>
                </Grid>
              </Grid>
            </Stack>
          )}

          {/* Step 3 - calc */}
          {activeStep === 3 && (
            <Box sx={{ width: isSmall ? '100%' : 250 }}>
              {calibration.distanceX === null ? (
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Введите количество пройденных шагов
                  </Typography>

                  <TextField
                    label="Количество шагов"
                    type="number"
                    value={calibration.cir}
                    onChange={(e) => setCalibration({ ...calibration, cir: e.target.value })}
                    fullWidth
                    inputProps={{ min: 0 }}
                  />

                  <Button
                    variant="contained"
                    onClick={handleCalc}
                    disabled={loading || !calibration.cir || Number(calibration.cir) < 1}
                    fullWidth={isSmall}
                  >
                    Подсчитать
                  </Button>
                </Stack>
              ) : (
                <Stack spacing={2} alignItems="center">
                  <Typography>Один шаг равен: {formatNum(calibration.dPerStep)} см</Typography>
                  <Button variant="contained" onClick={() => setActiveStep(4)} fullWidth={isSmall}>
                    Продолжить
                  </Button>
                </Stack>
              )}
            </Box>
          )}

          {/* Step 4 - save */}
          {activeStep === 4 && (
            <Box sx={{ textAlign: "center", width: '100%' }}>
              <Button variant="contained" onClick={handleSave} disabled={loading || isSaved} fullWidth={isSmall}>
                Сохранить калибровку
              </Button>
            </Box>
          )}
        </Box>

        <Box sx={{ display: "flex", justifyContent: "space-between", mt: 1 }}>
          <Button
            variant="outlined"
            disabled={activeStep === 0 || (activeStep === 4 && isSaved)}
            onClick={handleBack}
            fullWidth={isSmall}
          >
            Назад
          </Button>
        </Box>
      </Paper>
    </Box>
  );
}
