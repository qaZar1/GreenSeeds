import React, { useState } from "react";
import {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Chip,
  LinearProgress,
  Collapse,
  Box,
} from "@mui/material";

const SimpleLogTable = ({ data = [], loading, isFetchingMore, hasMore }) => {
  const [openId, setOpenId] = useState(null);

  const toggle = (id) => setOpenId(openId === id ? null : id);

  return (
    <TableContainer sx={{ background: "transparent" }}>
      {loading && <LinearProgress />}

      <Table size="small" stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>Дата</TableCell>
            <TableCell>Уровень</TableCell>
            <TableCell>Сообщение</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {data.map((row) => {
            const rowKey = row.id;

            return (
                <React.Fragment key={rowKey}>
                    <TableRow hover onClick={() => toggle(rowKey)}>
                        <TableCell>{new Date(row.dt).toLocaleString()}</TableCell>

                        <TableCell>
                            <Chip
                                size="small"
                                label={row.lvl}
                                color={
                                    row.lvl === "ERROR"
                                        ? "error"
                                        : row.lvl === "WARN"
                                        ? "warning"
                                        : "success"
                                }
                            />
                        </TableCell>

                        <TableCell>{row.msg}</TableCell>
                    </TableRow>

                    <TableRow>
                        <TableCell colSpan={3} sx={{ p: 0 }}>
                            <Collapse in={openId === rowKey}>
                                <Box sx={{ p: 2}}>
                                    <div><b>Дата:</b> {row.dt}</div>
                                    <div><b>ID:</b> {row.request_id || '-'}</div>
                                    <div><b>Функция:</b> {row.caller || '-'}</div>
                                    <div><b>Пользователь:</b> {row.username || '-'}</div>
                                    <div><b>Сообщение:</b> {row.msg || '-'}</div>
                                </Box>
                            </Collapse>
                        </TableCell>
                    </TableRow>
                </React.Fragment>
            );
        })}


          {!loading && data.length === 0 && (
            <TableRow>
              <TableCell colSpan={3} align="center" sx={{ py: 3 }}>
                Данные не найдены
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {isFetchingMore && <LinearProgress sx={{ mt: 1 }} />}

      {!hasMore && !loading && data.length > 0 && (
        <Box sx={{ p: 2, textAlign: "center" }}>
          Конец списка
        </Box>
      )}
    </TableContainer>
  );
};

export default SimpleLogTable;
