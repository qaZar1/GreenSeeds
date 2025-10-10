import React, { useState, useEffect } from "react";
import {
    Edit,
    SimpleForm,
    TextInput,
    Toolbar,
    SaveButton,
    useNotify,
    PasswordInput,
    warning,
    useGetOne,
} from "react-admin";
import { useMediaQuery, Card, CardContent, Typography, Box, TextField } from "@mui/material";
import { jwtDecode } from "jwt-decode";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button } from "@mui/material";
import { BooleanInput } from "react-admin";
import { EmptyProfile } from "./EmptyProfile";
import { LoadingOverlay } from "../utils/Loading";

const Profile = () => {
    const [username, setUsername] = useState(null);
    const [openDialog, setOpenDialog] = useState(false);
    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const notify = useNotify();
    const { data, isLoading, error } = useGetOne("profile", { id: username });


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

    const handlePasswordReset = async () => {
        if (newPassword !== confirmPassword) {
            notify("Новые пароли не совпадают", "error");
            return;
        }

        const response = await fetch('/api/users/change-password', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                oldPassword: oldPassword,
                newPassword: newPassword,
            }),
        });

        if (response.status !== 204) {
            notify("Ошибка обновления пароля", "error");
            return;
        }

        notify("Пароль успешно обновлен", "success");

        setOpenDialog(false);
        setOldPassword("");
        setNewPassword("");
        setConfirmPassword("");
    };

    const ProfileToolbar = (props) => (
        <Toolbar {...props}>
            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    width: '100%',
                }}
            >
                <SaveButton label="Сохранить"/>
                <Button
                        variant="outlined"
                        color="secondary"
                        sx={{ ml: 2 }}
                        onClick={() => setOpenDialog(true)}
                    >
                        Сброс пароля
                </Button>
            </Box>
        </Toolbar>
    );

    if (!username || isLoading) return <LoadingOverlay text="Загрузка профиля..." />;

    if (error) return <EmptyProfile />;

    return (
        <>
            <Edit
                id={username}
                resource="profile"
                mutationMode="pessimistic"
                title="Профиль"
                sx={{
                    display: 'flex',
                    justifyContent: 'center',  // центрируем по горизонтали
                    padding: 2,
                }}
            >
                <SimpleForm toolbar={<ProfileToolbar />} empty={<EmptyProfile />}>
                    <Box
                        display="flex"
                        flexDirection="column"
                        alignItems="center"
                        gap={2}
                    >
                        <Box sx={{ width: '100%' }}>
                            <Typography variant="h6" gutterBottom>
                            Основная информация
                            </Typography>
                            <TextInput source="id" label="Имя пользователя" fullWidth disabled/>
                            <TextInput source="full_name" label="ФИО" fullWidth />
                            <BooleanInput source="is_admin" label="Администратор" disabled/>
                        </Box>
                    </Box>
                    </SimpleForm>
            </Edit>
            {/* Модалка смены пароля */}
            <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
                <DialogTitle>Смена пароля</DialogTitle>
                <DialogContent>
                    <TextField
                        label="Старый пароль"
                        type="password"
                        fullWidth
                        onChange={(e) => setOldPassword(e.target.value)}
                        value={oldPassword}
                    />
                    <TextField
                        label="Новый пароль"
                        type="password"
                        fullWidth
                        onChange={(e) => setNewPassword(e.target.value)}
                        value={newPassword}
                    />
                    <TextField
                        label="Повторите новый пароль"
                        type="password"
                        fullWidth
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        value={confirmPassword}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenDialog(false)}>Отмена</Button>
                    <Button onClick={handlePasswordReset} variant="contained">
                        Сохранить
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
};

export default Profile;
