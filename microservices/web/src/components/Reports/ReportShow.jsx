import React from "react";
import { Show, SimpleShowLayout, TextField, DateField } from "react-admin";
import BackButton from "../utils/Back";
import { Card, CardContent, Typography, Divider, Box } from "@mui/material";
import { useMediaQuery } from "@mui/material";

const ReportShow = () => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    
    return (
        <Show actions={<BackButton />} sx={{ p: 3 }}>
            <Card sx={{ borderRadius: 3, boxShadow: 3, p: 2 }}>
                <CardContent sx={{ fontSize: 20 }}>
                    <Typography variant="h6" gutterBottom>
                        Детали отчета
                    </Typography>
                    <Divider sx={{ mb: 2 }} />

                    {isSmall || isMedium ? (
                        <SimpleShowLayout>
                            <TextField source="shift" emptyText="-" label="Смена" />
                            <TextField source="number" emptyText="-" label="Номер задания" />
                            <TextField source="receipt" emptyText="-" label="Рецепт" />
                            <TextField source="turn" emptyText="-" label="Номер выполнения" />
                            <DateField source="dt" showTime locales="ru-RU" label="Дата" />
                            <TextField source="success" emptyText="-" label="Успешно" />
                            <TextField source="error" emptyText="-" label="Ошибка" />
                            <TextField source="solution" emptyText="-" label="Решение" />
                            <TextField source="mark" emptyText="-" label="Маркировка" />
                        </SimpleShowLayout>
                    ) : (
                    <Box sx={{ display: "flex", flexDirection: "row", width: "100%"}}>
                        <Box sx={{ display: "flex", flexDirection: "column", flex: 1}}>
                            <SimpleShowLayout>
                                <TextField source="shift" emptyText="-" label="Смена" />
                                <TextField source="number" emptyText="-" label="Номер задания" />
                                <TextField source="receipt" emptyText="-" label="Рецепт" />
                                <TextField source="turn" emptyText="-" label="Номер выполнения" />
                                <DateField source="dt" showTime locales="ru-RU" label="Дата" />
                            </SimpleShowLayout>
                        </Box>
                        <Box sx={{ display: "flex", flexDirection: "column", flex: 1}}>
                            <SimpleShowLayout>
                                <TextField source="success" emptyText="-" label="Успешно" />
                                <TextField source="error" emptyText="-" label="Ошибка" />
                                <TextField source="solution" emptyText="-" label="Решение" />
                                <TextField source="mark" emptyText="-" label="Маркировка" />
                            </SimpleShowLayout>
                        </Box>
                    </Box>
                    )}
                </CardContent>
            </Card>
        </Show>
    );
};

export default ReportShow;
