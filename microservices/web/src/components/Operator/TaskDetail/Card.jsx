import React from "react";
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

const TaskCard = () => {
  const record = useRecordContext();

  const progress =
    record && record.required_amount
      ? (record.completed_amount / record.required_amount) * 100
      : 0;

  const handleAction1 = () => {
    console.log("Действие 1 выполнено");
  };

  const handleAction2 = () => {
    console.log("Действие 2 выполнено");
  };

  return (
    <Card
      sx={{
        borderRadius: 3,
        boxShadow: 4,
        p: 3,
      }}
    >
      <CardContent>
        <Typography
          variant="h5"
          fontWeight="bold"
          gutterBottom
          textAlign="center"
        >
          Детали задания
        </Typography>

        <Divider sx={{ my: 2 }} />

        <SimpleShowLayout>
          <Box
            display="grid"
            gridTemplateColumns={{ xs: "1fr", sm: "1fr 1fr" }}
            gap={3}
            width="100%"
          >
            {/* Первая колонка */}
            <Box>
              <Typography variant="subtitle2" color="textSecondary">
                Номер
              </Typography>
              <Typography variant="body1">
                <TextField source="number" />
              </Typography>

              <Typography variant="subtitle2" color="textSecondary" mt={2}>
                Смена
              </Typography>
              <Typography variant="body1">
                <TextField source="shift" />
              </Typography>

              <Typography variant="subtitle2" color="textSecondary" mt={2}>
                Семена
              </Typography>
              <Typography variant="body1">
                <TextField source="seed" />
              </Typography>
            </Box>

            {/* Вторая колонка */}
            <Box>
              <Typography variant="subtitle2" color="textSecondary">
                Требуемое количество
              </Typography>
              <Typography variant="body1">
                <TextField source="required_amount" />
              </Typography>

              <Typography variant="subtitle2" color="textSecondary" mt={2}>
                Выполнено
              </Typography>
              <Typography variant="body1">
                <TextField source="completed_amount" />
              </Typography>
            </Box>

            {/* Прогресс бар, растянутый на обе колонки */}
            <Box gridColumn="1 / -1" mt={3}>
              <Typography variant="subtitle2" color="textSecondary">
                Прогресс выполнения
              </Typography>
              <LinearProgress
                variant="determinate"
                value={progress}
                sx={{ height: 10, borderRadius: 5, mt: 1 }}
              />
            </Box>

            {/* Кнопки, растянутые на обе колонки */}
            <Box
              gridColumn="1 / -1"
              mt={3}
              display="flex"
              justifyContent="space-between"
              gap={2}
            >
              <Button
                variant="contained"
                color="primary"
                fullWidth
                onClick={handleAction1}
              >
                Действие 1
              </Button>
              <Button
                variant="outlined"
                color="secondary"
                fullWidth
                onClick={handleAction2}
              >
                Действие 2
              </Button>
            </Box>
          </Box>
        </SimpleShowLayout>
      </CardContent>
    </Card>
  );
};

export default TaskCard;
