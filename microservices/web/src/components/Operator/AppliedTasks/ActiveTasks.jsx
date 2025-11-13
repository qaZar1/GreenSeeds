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
  Divider,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import EmptyTasks from "./EmptyTasks";

const Tasks = ({ tasks }) => {
  const navigate = useNavigate();

  if (!tasks || tasks.length === 0) {
    return (
      <EmptyTasks />
    );
  }

  // Группируем задачи по shift
  const grouped = tasks.reduce((acc, t) => {
    const shift = t.shift || "Без смены";
    if (!acc[shift]) acc[shift] = [];
    acc[shift].push(t);
    return acc;
  }, {});

  return (
    <Box display="flex" justifyContent="center" p={2} pl={5} pr={5}>
      <Card sx={{ width: "100%", borderRadius: 3, boxShadow: 4 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom align="center">
            Активные задания
          </Typography>

          {Object.entries(grouped).map(([shift, group]) => (
            <Box key={shift} mb={7}>
              <Typography
                variant="subtitle1"
                sx={{ fontWeight: "bold", mb: 1, color: "primary.main" }}
              >
                Смена: {shift}
              </Typography>

              <TableContainer component={Paper} sx={{ borderRadius: 2 }}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Задание</TableCell>
                      <TableCell>Кол-во</TableCell>
                      <TableCell>Культура</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {group.map((t) => (
                      <TableRow
                        key={t.id}
                        hover
                        sx={{ cursor: "pointer" }}
                        onClick={() => navigate(`/tasks/${t.id}`)}
                      >
                        <TableCell>{t.number}</TableCell>
                        <TableCell>{t.amount}</TableCell>
                        <TableCell>{t.seed_ru || "—"}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          ))}
        </CardContent>
      </Card>
    </Box>
  );
};

export default Tasks;
