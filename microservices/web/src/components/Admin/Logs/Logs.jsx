import React, { useState, useEffect, useCallback, useRef } from "react";
import {
    Box,
    TextField,
    MenuItem,
    Select,
    FormControl,
    InputLabel,
    useMediaQuery,
} from "@mui/material";
import { Title } from "react-admin";
import { useTheme } from "@mui/material/styles";
import dataProvider from "../../../dataProvider";
import MobileLogList from "./MobileLogList";
import SimpleLogTable from "./SimpleLogPage";

const LIMIT = 50;

const LogsPage = () => {
    const scrollRef = useRef();
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [more, setMore] = useState(true);

    const [search, setSearch] = useState("");
    const [level, setLevel] = useState("ALL");
    const [dateFrom, setDateFrom] = useState("");
    const [dateTo, setDateTo] = useState("");

    const [offset, setOffset] = useState(0);

    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

    const [query, setQuery] = useState("");
    useEffect(() => {
        const t = setTimeout(() => setQuery(search), 300);
        return () => clearTimeout(t);
    }, [search]);

    // reset filter
    useEffect(() => {
        setData([])
        setOffset(0);
        setMore(true);
    }, [query, level, dateFrom, dateTo]);

    const fetchLogs = useCallback(async () => {
        setLoading(true);

        const params = {
            filter: {
                search: query,
                level,
                date_from: dateFrom,
                date_to: dateTo,
            },
            pagination: { page: offset / LIMIT + 1, perPage: LIMIT },
            sort: { field: "dt", order: "DESC" },
        };

        try {
            const { data: newDataRaw } = await dataProvider.getList("logs", params);

            const newData = newDataRaw.map((item, index) => ({
                id: offset + index,
                ...item,
            }));

            setData((prev) => (offset === 0 ? newData : [...prev, ...newData]));
            setMore(newData.length === LIMIT);
        } catch {
            setData(offset === 0 ? [] : data);
            setMore(false);
        }

        setLoading(false);
    }, [offset, query, level, dateFrom, dateTo]);

    useEffect(() => {
        fetchLogs();
    }, [offset, fetchLogs]);

    const handleScroll = (e) => {
        if (!more || loading) return;

        const el = e.target;

        const bottom =
            el.scrollHeight - el.scrollTop - el.clientHeight < 150;

        if (bottom) {
            setOffset((prev) => prev + LIMIT);
        }
    };

    return (
        <Box sx={{ display: "flex", flexDirection: "column", height: "100%", p: 5 }}>
            <Title title="Логи системы" />

            <Box
                sx={{
                    display: "grid",
                    gap: 2,
                    gridTemplateColumns: {
                        xs: "1fr",
                        sm: "repeat(2, 1fr)",
                        md: "repeat(4, 1fr)",
                    },
                    mb: 2,
                }}
            >

                {/* 1 — поиск */}
                <TextField
                    label="Поиск"
                    size="small"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                />

                {/* 2 — уровень */}
                <FormControl size="small">
                    <InputLabel>Уровень</InputLabel>
                    <Select
                        value={level}
                        label="Уровень"
                        onChange={(e) => setLevel(e.target.value)}
                    >
                        <MenuItem value="ALL">Все</MenuItem>
                        <MenuItem value="INFO">INFO</MenuItem>
                        <MenuItem value="WARN">WARN</MenuItem>
                        <MenuItem value="ERROR">ERROR</MenuItem>
                    </Select>
                </FormControl>

                {/* 3 — дата с */}
                <TextField
                    label="Дата с"
                    type="date"
                    size="small"
                    value={dateFrom}
                    onChange={(e) => setDateFrom(e.target.value)}
                    slotProps={{ inputLabel: { shrink: true } }}
                    sx={{
                        gridColumn: {
                            xs: "span 1",
                            sm: "span 1",
                            md: "span 1",
                        },
                    }}
                />

                {/* 4 — дата до */}
                <TextField
                    label="Дата до"
                    type="date"
                    size="small"
                    value={dateTo}
                    onChange={(e) => setDateTo(e.target.value)}
                    slotProps={{ inputLabel: { shrink: true } }}
                    sx={{
                        gridColumn: {
                            xs: "span 1",
                            sm: "span 1",
                            md: "span 1",
                        },
                    }}
                />
            </Box>

            <Box
                ref={scrollRef}
                onScroll={handleScroll}
                sx={{
                    overflowY: "auto",
                    height: "calc(100vh - 200px)",
                }}
            >
                {isMobile ? (
                    <MobileLogList
                        data={data}
                        loading={loading && offset === 0}
                        isFetchingMore={loading && offset > 0}
                        hasMore={more}
                    />
                ) : (
                    <SimpleLogTable
                        data={data}
                        loading={loading && offset === 0}
                        isFetchingMore={loading && offset > 0}
                        hasMore={more}
                    />
                )}
            </Box>
        </Box>
    );
};

export default LogsPage;
