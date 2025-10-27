import React, { useState, useEffect } from "react";
import { Box, CircularProgress, Typography, Card, CardContent } from "@mui/material";
import { jwtDecode } from "jwt-decode";
import { Show, SimpleShowLayout, useShowContext } from "react-admin";
import Task from "./FreeTask";
import EmptyChoice from "./EmptyChoice";
import { List } from "react-admin";
import { useMediaQuery } from "@mui/material";
import ChoiceListContent from "./Controller";

const ChoiceList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));

    const [username, setUsername] = useState(null);

    useEffect(() => {
      try {
        const stored = localStorage.getItem("auth");
        if (stored) {
          const parsed = JSON.parse(stored);
          if (parsed?.token) {
            const decoded = jwtDecode(parsed.token);
            setUsername(decoded?.username);
          }
        }
      } catch (e) {
        console.warn("Ошибка получения профиля:", e);
      }

    }, []);

    if (!username) {
      return (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
          <CircularProgress />
        </Box>
      );
    }

    return (
        <List
            resource="choice"
            empty={<EmptyChoice />}
            {...props}
            pagination={false}
            perPage={false}
            sx={{ padding: 2 }}
            actions={false}
            title="Выбор задания"
            component={Box}
        >
            <ChoiceListContent isSmall={isSmall} isMedium={isMedium} username={username} />
        </List>
    );
};

export default ChoiceList;