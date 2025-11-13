import React, { useState } from "react";
import {
  Box,
  Typography,
  Collapse,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  IconButton,
  Card,
  CardContent,
  Button,
  useMediaQuery,
} from "@mui/material";
import { ExpandLess, ExpandMore } from "@mui/icons-material";
import { useListContext, DateField, BooleanField, ShowButton } from "react-admin";
import { useTheme } from "@mui/material/styles";

const GroupedDatagrid = () => {
  const { data, isLoading } = useListContext();
  const [openGroups, setOpenGroups] = useState({});
  const [showHistory, setShowHistory] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));
  const isMedium = useMediaQuery(theme.breakpoints.between("sm", "md"));

  if (isLoading) return <p>Загрузка...</p>;

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  // Разделяем сегодняшние и исторические данные
  const todaysData = data.filter((item) => new Date(item.dt) >= today);
  const historyData = data.filter((item) => new Date(item.dt) < today);

  // Группировка по сменам
  const grouped = (showHistory ? data : todaysData).reduce((acc, item) => {
    acc[item.shift] = acc[item.shift] ? [...acc[item.shift], item] : [item];
    return acc;
  }, {});

  const toggleGroup = (shift) => setOpenGroups((prev) => ({ ...prev, [shift]: !prev[shift] }));

  return (
    <Box sx={{ p: 2 }}>
      <Button
        variant="outlined"
        onClick={() => setShowHistory((prev) => !prev)}
        sx={{ mb: 2 }}
      >
        {showHistory ? "Скрыть историю" : "Показать историю"}
      </Button>

      {Object.entries(grouped)
        .sort(([a], [b]) => Number(b) - Number(a))
        .map(([shift, reports]) => (
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

            {/* Содержимое смены */}
            <Collapse in={openGroups[shift]} timeout="auto" unmountOnExit>
              {isMobile || isMedium ? (
                <Box
                  sx={{
                    p: 2,
                    display: "grid",
                    gridTemplateColumns: isMobile ? "1fr" : "1fr 1fr",
                    gap: 2,
                  }}
                >
                  {reports
                    .sort((a, b) => a.number - b.number)
                    .map((rec, i) => (
                      <Card
                        key={rec.id || i}
                        variant="outlined"
                        sx={{ borderRadius: 2, boxShadow: 1 }}
                      >
                        <CardContent>
                          <Typography variant="subtitle2" color="text.secondary">
                            Задание №{rec.number}
                          </Typography>
                          <Typography>Рецепт: {rec.receipt}</Typography>
                          <Typography>Номер выполнения: {rec.turn}</Typography>
                          <Typography>
                            Дата: <DateField record={rec} source="dt" showTime locales="ru-RU" />
                          </Typography>
                          <Typography>
                            Успешно: <BooleanField record={rec} source="success" />
                          </Typography>
                          <Box sx={{ mt: 1 }}>
                            <ShowButton record={rec} label="Показать" />
                          </Box>
                        </CardContent>
                      </Card>
                    ))}
                </Box>
              ) : (
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Номер задания</TableCell>
                      <TableCell>Рецепт</TableCell>
                      <TableCell>Номер выполнения</TableCell>
                      <TableCell>Дата</TableCell>
                      <TableCell>Успешно</TableCell>
                      <TableCell></TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {reports
                      .sort((a, b) => a.number - b.number)
                      .map((rec) => (
                        <TableRow key={rec.id}>
                          <TableCell>{rec.number}</TableCell>
                          <TableCell>{rec.receipt}</TableCell>
                          <TableCell>{rec.turn}</TableCell>
                          <TableCell>
                            <DateField record={rec} source="dt" showTime locales="ru-RU" />
                          </TableCell>
                          <TableCell>
                            <BooleanField record={rec} source="success" />
                          </TableCell>
                          <TableCell>
                            <ShowButton record={rec} label="Показать" />
                          </TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
              )}
            </Collapse>
          </Box>
        ))}
    </Box>
  );
};

export default GroupedDatagrid;
