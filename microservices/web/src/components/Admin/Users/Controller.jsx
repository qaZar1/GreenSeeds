import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EmptyUser } from "./EmptyUser";
import MobileActionsMenu from "./Mobile";
import AdminToggleField from "./AdminToggleField";
import { DeleteButton } from "react-admin";
import ResetPasswordButton from "./ResetPwd";
import { FunctionField } from "react-admin";

const UserListContent = ({ isSmall, currentUsername }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay/>;
    if (error) return <EmptyUser />;

    return isSmall ? (
        <SimpleList
            primaryText={(record) => record.full_name}
            secondaryText={(record) =>
                record.is_admin ? "Администратор" : "Пользователь"
            }
            tertiaryText={record => (
                <>
                    <MobileActionsMenu record={record} currentUsername={currentUsername} />
                </>
            )}
            rowClick={false}
            empty={<EmptyUser />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptyUser />}
        >
            <TextField source="full_name" label="ФИО" />
            <FunctionField
                label="Администратор"
                render={record => <AdminToggleField record={record} currentUsername={currentUsername}/>}
            />
            <FunctionField
                label="Сброс пароля"
                render={record => <ResetPasswordButton record={record} />}
            />
            <FunctionField
                label="Удалить"
                render={record => (
                    <DeleteButton
                        label="Удалить"
                        record={record}
                        mutationMode="pessimistic"
                        disabled={record.username === currentUsername}
                    />
                )}
            />
        </Datagrid>
    )
};

export default UserListContent;
