// src/components/layout/CustomMenu.jsx
import React, { useState } from "react";
import { Menu, MenuItemLink } from "react-admin";
import {
  ExpandLess,
  ExpandMore,
  Warehouse as WarehouseIcon,
  Grass as GrassIcon,
  Receipt as ReceiptIcon,
  DragIndicator as DragIndicatorIcon,
  Assignment as AssignmentIcon,
  Link as LinkIcon,
  People as PersonIcon,
} from "@mui/icons-material";
import { Collapse, ListItemButton, ListItemText, Box } from "@mui/material";
import SettingsApplicationsIcon from "@mui/icons-material/SettingsApplications";

const CustomMenu = () => {
  const [openRefs, setOpenRefs] = useState(false);

  const handleToggleRefs = () => setOpenRefs(!openRefs);

  return (
    <Menu>
      {/* ✅ Остальные отдельные пункты */}
      <MenuItemLink to="/users" primaryText="Пользователи" leftIcon={<PersonIcon />} />
      <MenuItemLink to="/tasks" primaryText="Задания" leftIcon={<AssignmentIcon />} />

      <MenuItemLink to="/shifts" primaryText="Смены" leftIcon={<DragIndicatorIcon />} />
      <MenuItemLink to="/assignments" primaryText="Задания" leftIcon={<AssignmentIcon />} />
      <MenuItemLink to="/reports" primaryText="Отчеты" leftIcon={<AssignmentIcon />} />

      {/* ✅ Секция справочников */}
      <ListItemButton onClick={handleToggleRefs} sx={{pl: 2, pr: 2, height: 36}}>
        <Box sx={{ display: "flex", alignItems: "center", flexGrow: 1 }}>
          <SettingsApplicationsIcon sx={{ marginRight: 2 }} />
          <ListItemText primary="Настройки" />
        </Box>
        {openRefs ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>

      <Collapse in={openRefs} timeout="auto" unmountOnExit>
        <Box sx={{ pl: 4 }}>
            <MenuItemLink to="/bunkers" primaryText="Бункеры" leftIcon={<WarehouseIcon />} />
            <MenuItemLink to="/seeds" primaryText="Семена" leftIcon={<GrassIcon />} />
            <MenuItemLink to="/placements" primaryText="Расположение" leftIcon={<LinkIcon />} />
            <MenuItemLink to="/receipts" primaryText="Рецепты" leftIcon={<ReceiptIcon />} />
        </Box>
      </Collapse>
    </Menu>
  );
};

export default CustomMenu;
