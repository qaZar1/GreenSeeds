import React, { useState } from "react";
import {
  Box,
  Collapse,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Card,
  CardContent,
  Button,
  useMediaQuery,
} from "@mui/material";
import { useTheme } from "@mui/material/styles";
import { useListContext, DateField, ShowButton } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EmptyShift } from "./EmptyShift";

const ShiftListContent = () => {
  const { data, isLoading, error } = useListContext();
  const [showHistory, setShowHistory] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  if (isLoading) return <LoadingOverlay />;
  if (error) return <EmptyShift />;

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const todaysData = data.filter((item) => new Date(item.dt) >= today);
  const historyData = data
    .filter((item) => new Date(item.dt) < today)
    .sort((a, b) => new Date(b.dt) - new Date(a.dt));

  const renderTableRows = (items) =>
    items
      .sort((a, b) => a.number - b.number)
      .map((rec) => (
        <TableRow key={rec.id}>
          <TableCell>
            <DateField record={rec} source="dt" showTime locales="ru-RU" />
          </TableCell>
          <TableCell>
            <ShowButton record={rec} label="Показать" />
          </TableCell>
        </TableRow>
      ));

  const renderTable = (items) => {
    if (!items.length) return null;
    return (
      <Table size="small" sx={{ mb: 2 }}>
        <TableHead>
          <TableRow>
            <TableCell>Дата</TableCell>
            <TableCell></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{renderTableRows(items)}</TableBody>
      </Table>
    );
  };

  const renderCards = (items) =>
    items
      .sort((a, b) => a.number - b.number)
      .map((rec) => (
        <Card key={rec.id} variant="outlined" sx={{ mb: 2 }}>
          <CardContent>
            <div>
              Дата: <DateField record={rec} source="dt" showTime locales="ru-RU" />
            </div>
            <div>
              <ShowButton record={rec} label="Показать" />
            </div>
          </CardContent>
        </Card>
      ));

  return (
    <Box sx={{ p: 2 }}>
      {historyData.length > 0 && (
        <Button
          variant="outlined"
          onClick={() => setShowHistory((prev) => !prev)}
          sx={{ mb: 2 }}
        >
          {showHistory ? "Скрыть историю" : "Показать историю"}
        </Button>
      )}

      {isMobile ? renderCards(todaysData) : renderTable(todaysData)}

      <Collapse in={showHistory}>
        {isMobile ? renderCards(historyData) : renderTable(historyData)}
      </Collapse>
    </Box>
  );
};

export default ShiftListContent;
