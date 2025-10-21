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
  Divider,
  useMediaQuery,
} from "@mui/material";
import { ExpandLess, ExpandMore } from "@mui/icons-material";
import {
  useListContext,
  DateField,
  BooleanField,
  ShowButton,
} from "react-admin";
import { useTheme } from "@mui/material/styles";

const GroupedDatagrid = () => {
  const { data, isLoading } = useListContext();
  const [openGroups, setOpenGroups] = useState({});
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  if (isLoading) return <p>Загрузка...</p>;
  if (!data || data.length === 0) return <p>Нет данных</p>;

  // Группировка по сменам
  const grouped = data.reduce((acc, item) => {
    acc[item.shift] = acc[item.shift] ? [...acc[item.shift], item] : [item];
    return acc;
  }, {});

  const toggleGroup = (shift) =>
    setOpenGroups((prev) => ({ ...prev, [shift]: !prev[shift] }));

  return (
    <Box sx={{ p: 2 }}>
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
            {isMobile ? (
              // 📱 Мобильный вид — карточки
              <Box sx={{ p: 2, bgcolor: "background.paper" }}>
                {reports
                .sort((a, b) => a.number - b.number)
                .map((rec, i) => (
                  <Card
                    key={rec.id || i}
                    variant="outlined"
                    sx={{
                      mb: 1.5,
                      borderRadius: 2,
                      boxShadow: 1,
                    }}
                  >
                    <CardContent>
                      <Typography variant="subtitle2" color="text.secondary">
                        Задание №{rec.number}
                      </Typography>
                      <Typography>Рецепт: {rec.receipt}</Typography>
                      <Typography variant="body2" sx={{ mt: 0.5 }}>
                        <DateField
                          record={rec}
                          source="dt"
                          showTime
                          locales="ru-RU"
                        />
                      </Typography>
                      <Typography sx={{ mt: 1 }}>
                        Успешно:{" "}
                        <BooleanField record={rec} source="success" />
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
                      <TableCell>
                        <DateField
                          record={rec}
                          source="dt"
                          showTime
                          locales="ru-RU"
                        />
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
