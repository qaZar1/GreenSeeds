import React, { useState } from "react";
import {
  Box,
  Chip,
  Card,
  CardContent,
  Typography,
  Collapse,
  LinearProgress,
} from "@mui/material";

const MobileLogList = ({ data = [], loading = false, isFetchingMore = false, hasMore = true }) => {
  const [openId, setOpenId] = useState(null);

  const toggle = (id) => setOpenId(openId === id ? null : id);

  if (loading) return <LinearProgress />;

  if (!data.length)
    return <Box sx={{ textAlign: "center", p: 3 }}>Данные не найдены</Box>;

  return (
    <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
      {data.map((row) => {
        const key = row.id; // ← используем только id

        return (
          <Card key={key} onClick={() => toggle(key)} sx={{ cursor: "pointer" }}>
            <CardContent>
              <Typography variant="body2" color="text.secondary">
                {new Date(row.dt).toLocaleString()}
              </Typography>

              <Box sx={{ display: "flex", gap: 1, mt: 1 }}>
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
              </Box>

              <Typography variant="body1" sx={{ mt: 1 }}>
                {row.msg}
              </Typography>

              <Collapse in={openId === key}>
                <Box sx={{ mt: 2, p: 2 }}>
                  <Typography variant="body2">
                    <b>Функция:</b> {row.caller}
                  </Typography>
                  <Typography variant="body2">
                    <b>ID:</b> {row.request_id}
                  </Typography>
                  <Typography variant="body2">
                    <b>Пользователь:</b> {row.username}
                  </Typography>
                  <Typography variant="body2" sx={{ whiteSpace: "pre-wrap", mt: 1 }}>
                    <b>Сообщение:</b> {row.msg}
                  </Typography>

                  {row.raw && (
                    <Box
                      sx={{
                        mt: 1,
                        fontFamily: "monospace",
                        fontSize: 12,
                        whiteSpace: "pre-wrap",
                      }}
                    >
                      {typeof row.raw === "string"
                        ? row.raw
                        : JSON.stringify(row.raw, null, 2)}
                    </Box>
                  )}
                </Box>
              </Collapse>
            </CardContent>
          </Card>
        );
      })}

      {isFetchingMore && <LinearProgress />}

      {!hasMore && (
        <Box sx={{ textAlign: "center", p: 2 }}>
          Конец списка
        </Box>
      )}
    </Box>
  );
};

export default MobileLogList;
