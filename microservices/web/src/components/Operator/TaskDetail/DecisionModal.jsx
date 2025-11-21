import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  TextField,
  Box,
} from "@mui/material";

const DecisionModal = ({ open, onClose, onConfirm, onReject, reason, photo }) => {
  const [comment, setComment] = useState(""); // новое поле для текста

  const handleConfirm = () => {
    onConfirm(comment); // передаем текст вместе с подтверждением
    setComment(""); // очищаем поле
  };

  const handleReject = () => {
    onReject(comment); // передаем текст вместе с отклонением
    setComment("");
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Требуется подтверждение</DialogTitle>

      <DialogContent>
        <Box display="flex" gap={2}>
          {photo && (
            <img
              src={`data:image/jpeg;base64,${photo}`}
              alt="Проверка"
              style={{
                width: "40%",
                maxHeight: 300,
                objectFit: "contain",
                borderRadius: 8,
              }}
            />
          )}
          <Box flex={1} display="flex" flexDirection="column" gap={1}>
            <Typography variant="body1" sx={{ mb: 1 }}>
              {reason}
            </Typography>
            <TextField
              label="Комментарий"
              multiline
              rows={6}
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              variant="outlined"
              fullWidth
            />
          </Box>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button variant="contained" color="success" onClick={handleConfirm}>
          Считать успешным
        </Button>
        <Button variant="outlined" color="error" onClick={handleReject}>
          Считать неуспешным
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DecisionModal;
