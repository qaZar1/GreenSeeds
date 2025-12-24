import React from "react";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Typography } from "@mui/material";
import { useState } from "react";

const BunkerCheckModal = ({ open, required, available, onReduce, onFill }) => {
  const [percent, setPercent] = useState("");
  console.log("req", required)

  return (
    <Dialog open={open}>
      <DialogTitle>Недостаточно семян</DialogTitle>
      <DialogContent>
        <Typography mb={2}>
          Требуется на задание: <b>{required}</b><br />
          Доступно в бункерах: <b>{available}</b>
        </Typography>

        <TextField
          fullWidth
          label="Заполнить бункер до N%"
          value={percent}
          type="number"
          onChange={(e) => setPercent(e.target.value)}
        />
      </DialogContent>

      <DialogActions>
        <Button
          variant="contained"
          color="primary"
          onClick={onReduce}
        >
          Уменьшить количество до {available}
        </Button>

        <Button
          variant="outlined"
          color="secondary"
          onClick={() => onFill(percent, required)}
        >
          Заполнить до {percent}%
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default BunkerCheckModal;
