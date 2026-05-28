import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyRecipe } from "./EmptyRecipe";
import RecipeListActions from "./Action";
import RecipeListContent from "./Controller";
import { useNotify } from "react-admin";

const RecipeList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    const notify = useNotify();

    return (
        <List
            resource="recipes"
            empty={<EmptyRecipe />}
            {...props}
            sx={{ padding: 2 }}
            actions={<RecipeListActions />}
            title="Рецепты"
            pagination={false}
            queryOptions={{
                onError: () => notify("Ошибка загрузки рецептов", { type: "error" }),
            }}
        >
            <RecipeListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default RecipeList;