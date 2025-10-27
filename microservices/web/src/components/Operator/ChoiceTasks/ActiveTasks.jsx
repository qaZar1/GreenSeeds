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
import { useNavigate } from "react-router-dom";

const Tasks = ({ tasks }) => {
  const navigate = useNavigate();

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
                </TableRow>
              </TableHead>
              <TableBody>
                {tasks.map((t, i) => (
                  <TableRow key={t.id} hover onClick={() => navigate(`/tasks/${t.id}`)}>
                    <TableCell>{t.number}</TableCell>
                    <TableCell>{t.amount}</TableCell>
                    <TableCell>{t.seed || "—"}</TableCell>
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
