import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  FormControlLabel,
  Checkbox,
} from "@mui/material";

const TaskStartDialog = ({ open, onClose, onConfirm, task }) => {
  const [extraMode, setExtraMode] = useState(false);

  if (!task) return null;

  const handleConfirm = () => {
    onConfirm(extraMode);
  };

  return (
    <Dialog 
      open={open} 
      onClose={onClose}
      fullWidth
      maxWidth="sm"
    >
      <DialogTitle>Подтверждение старта задания</DialogTitle>

      <DialogContent dividers>
        <Typography>Вы уверены, что хотите запустить выполнение задания?</Typography>

        <Box mt={2}>
          <Typography variant="subtitle2">Номер задания:</Typography>
          <Typography>{task.number}</Typography>

          <Typography variant="subtitle2" mt={2}>Смена:</Typography>
          <Typography>{task.shift}</Typography>

          <Typography variant="subtitle2" mt={2}>Семена:</Typography>
          <Typography>{task.seed_ru}</Typography>
        </Box>

        <Box mt={3}>
          <FormControlLabel
            control={
              <Checkbox
                checked={extraMode}
                onChange={(e) => setExtraMode(e.target.checked)}
              />
            }
            label="Включить отображение на дисплее?"
          />
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>Отмена</Button>
        <Button variant="contained" onClick={handleConfirm}>
          Подтвердить и начать
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TaskStartDialog;
