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
  Report as ReportIcon,
  Link as LinkIcon,
  People as PersonIcon,
  Task as TaskIcon,
  CalendarMonth as CalendarMonthIcon,
} from "@mui/icons-material";
import { Collapse, ListItemButton, ListItemText, Box } from "@mui/material";
import SettingsApplicationsIcon from "@mui/icons-material/SettingsApplications";
import { getRole } from "../../authProvider";
import ListAltIcon from '@mui/icons-material/ListAlt';

const CustomMenu = () => {
  const [openRefs, setOpenRefs] = useState(false);
  const role = getRole();

  const handleToggleRefs = () => setOpenRefs(!openRefs);

  return (
    <Menu>
      {role === 'admin' && [
      <MenuItemLink key="users" to="/users" primaryText="Пользователи" leftIcon={<PersonIcon />} />,
      ]}

      {role === 'operator' && [
          <MenuItemLink key="choice" to="/choice" primaryText="Выбор задания" leftIcon={<DragIndicatorIcon />} />,
          <MenuItemLink key="tasks" to="/tasks" primaryText="Задания на смену" leftIcon={<TaskIcon />}/>
      ]}

      {role === 'admin' && [
          <MenuItemLink key="shifts" to="/shifts" primaryText="План производства" leftIcon={<CalendarMonthIcon />} />,
          <MenuItemLink key="assignments" to="/assignments" primaryText="Сменные задания" leftIcon={<AssignmentIcon />} />,
          <MenuItemLink key="reports" to="/reports" primaryText="Отчеты" leftIcon={<ReportIcon />} />,
          <MenuItemLink key="logs" to="/logs" primaryText="Логи" leftIcon={<ListAltIcon />} />,
      ]}

      {/* Секция настроек */}
      {role === 'admin' && [
      <ListItemButton key="settings" onClick={handleToggleRefs} sx={{pl: 2, pr: 2, height: 36}}>
        <Box sx={{ display: "flex", alignItems: "center", flexGrow: 1 }}>
          <SettingsApplicationsIcon sx={{ marginRight: 2 }} />
          <ListItemText primary="Настройки" />
        </Box>
        {openRefs ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>,

      <Collapse key="collapse" in={openRefs} timeout="auto" unmountOnExit>
        <Box sx={{ pl: 4 }}>
            <MenuItemLink key="bunkers" to="/bunkers" primaryText="Бункеры" leftIcon={<WarehouseIcon />} />
            <MenuItemLink key="seeds" to="/seeds" primaryText="Семена" leftIcon={<GrassIcon />} />
            <MenuItemLink key="placements" to="/placements" primaryText="Расположение" leftIcon={<LinkIcon />} />
            <MenuItemLink key="receipts" to="/receipts" primaryText="Рецепты" leftIcon={<ReceiptIcon />} />
        </Box>
      </Collapse>
      ]}
    </Menu>
  );
};

export default CustomMenu;
