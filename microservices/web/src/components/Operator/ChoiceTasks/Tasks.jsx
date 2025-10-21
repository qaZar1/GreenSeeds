import React from "react";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from "@mui/material";

const Tasks = ({ tasks }) => {
  if (!tasks || tasks.length === 0) {
    return <Typography align="center">Нет заданий</Typography>;
  }

  return (
    <Box display="flex" justifyContent="center" p={2}>
      <Card sx={{ width: "100%", borderRadius: 3, boxShadow: 4 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom align="center">
            Активные задания
          </Typography>

          <TableContainer component={Paper} sx={{ borderRadius: 2 }}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Задание</TableCell>
                  <TableCell>Кол-во</TableCell>
                  <TableCell>Культура</TableCell>
                  <TableCell>Статус</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {tasks.map((t, i) => (
                  <TableRow key={t.id || i} hover>
                    <TableCell>{t.number}</TableCell>
                    <TableCell>{t.amount}</TableCell>
                    <TableCell>{t.seed || "—"}</TableCell>
                    <TableCell>{t.status || "Ожидает выполнения"}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Box>
  );
};

export default Tasks;
