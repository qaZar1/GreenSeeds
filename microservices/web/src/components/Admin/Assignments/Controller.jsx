import React, { useState } from "react";
import {
  Box,
  Card,
  CardContent,
  Collapse,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  IconButton,
  useMediaQuery,
} from "@mui/material";
import { ExpandLess, ExpandMore } from "@mui/icons-material";
import { useTheme } from "@mui/material/styles";
import { useListContext, EditButton } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EmptyAssignments } from "./EmptyAssign";

const AssignmentsListContent = () => {
  const { data, isLoading, error } = useListContext();
  const [openGroups, setOpenGroups] = useState({});
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  if (isLoading) return <LoadingOverlay />;
  if (error || !data?.length) return <EmptyAssignments />;

  // группируем по shift
  const grouped = data.reduce((acc, item) => {
    if (!acc[item.shift]) acc[item.shift] = [];
    acc[item.shift].push(item);
    return acc;
  }, {});

  const toggleGroup = (shift) => {
    setOpenGroups((prev) => ({ ...prev, [shift]: !prev[shift] }));
  };

  const renderTableRows = (items) =>
    [...items]
      .sort((a, b) => a.number - b.number)
      .map((rec) => (
        <TableRow key={rec.id}>
          <TableCell>{rec.number}</TableCell>
          <TableCell>{rec.description}</TableCell>
          <TableCell>{rec.amount}</TableCell>
          <TableCell>
            <EditButton record={rec} label="Редактировать" />
          </TableCell>
        </TableRow>
      ));

  const renderCards = (items) =>
    [...items]
      .sort((a, b) => a.number - b.number)
      .map((rec) => (
        <Card key={rec.id} variant="outlined" sx={{ mb: 1 }}>
          <CardContent>
            <Typography variant="body2">Задание: {rec.number}</Typography>
            <Typography variant="body2">Рецепт: {rec.description}</Typography>
            <Typography variant="body2">Количество: {rec.amount}</Typography>
            <EditButton record={rec} label="Редактировать" />
          </CardContent>
        </Card>
      ));

  return (
    <Box sx={{ p: 2 }}>
      {Object.entries(grouped)
        .sort(([a], [b]) => Number(a) - Number(b))
        .map(([shift, items]) => (
          <Box
            key={shift}
            sx={{
              mb: 2,
              borderRadius: 2,
              boxShadow: 2,
              overflow: "hidden",
            }}
          >
            {/* Заголовок смены */}
            <Box
              onClick={() => toggleGroup(shift)}
              sx={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                p: 2,
                cursor: "pointer",
                bgcolor: "action.hover",
              }}
            >
              <Typography variant="subtitle1" fontWeight="bold">
                Смена №{shift}
              </Typography>
              <IconButton size="small">
                {openGroups[shift] ? <ExpandLess /> : <ExpandMore />}
              </IconButton>
            </Box>

            {/* Содержимое */}
            <Collapse in={openGroups[shift]}>
              <Box sx={{ p: 2 }}>
                {isMobile ? (
                  renderCards(items)
                ) : (
                  <Table size="small">
                    <TableHead>
                      <TableRow>
                        <TableCell>Задание</TableCell>
                        <TableCell>Рецепт</TableCell>
                        <TableCell>Количество</TableCell>
                        <TableCell />
                      </TableRow>
                    </TableHead>
                    <TableBody>{renderTableRows(items)}</TableBody>
                  </Table>
                )}
              </Box>
            </Collapse>
          </Box>
        ))}
    </Box>
  );
};

export default AssignmentsListContent;
