export const resourcesConfig = {
    bunkers: {
        idKey: "bunker",
        urls: {
            list: "/api/bunkers/get",
            one: (id) => `/api/bunkers/get/${id}`,
            create: "/api/bunkers/add",
            update: "/api/bunkers/update",
            delete: (id) => `/api/bunkers/delete/${id}`,
        },
        createPayload: (d) => ({
            bunker: d.bunker,
            distance: d.distance,
        }),
        updatePayload: (d) => ({
            bunker: d.id,
            distance: d.distance,
        }),
    },

    seeds: {
        idKey: "seed",
        urls: {
            list: "/api/seeds/get",
            one: (id) => `/api/seeds/get/${id}`,
            create: "/api/seeds/add",
            update: "/api/seeds/update",
            delete: (id) => `/api/seeds/delete/${id}`,
        },
        createPayload: (d) => ({
            seed: d.seed,
            seed_ru: d.seed_ru,
            min_density: d.min_density,
            max_density: d.max_density,
            tank_capacity: d.tank_capacity,
        }),
        updatePayload: (d) => ({
            seed: d.id,
            seed_ru: d.seed_ru,
            min_density: d.min_density,
            max_density: d.max_density,
            tank_capacity: d.tank_capacity,
            latency: d.latency,
        }),
    },

    users: {
        idKey: "username",
        urls: {
            list: "/api/users/get",
            one: (id) => `/api/users/get/${id}`,
            create: "/api/register",
            update: "/api/users/update",
            delete: (id) => `/api/users/delete/${id}`,
        },
        createPayload: (d) => ({
            username: d.username,
            password: d.password,
            full_name: d.full_name,
            is_admin: d.is_admin,
        }),
        updatePayload: (d) => ({
            username: d.username,
            password: d.password,
            full_name: d.full_name,
            is_admin: d.is_admin,
        }),
    },

    profile: {
        idKey: "username",
        urls: {
            one: (id) => `/api/users/get/${id}`,
            update: "/api/users/update",
        },
        updatePayload: (d) => ({
            username: d.username,
            full_name: d.full_name,
        }),
    },

    placements: {
        idKey: "bunker",
        urls: {
            list: "/api/placement/get",
            one: (id) => `/api/placement/get/${id}`,
            create: "/api/placement/add",
            update: "/api/placement/update",
            delete: (id) => `/api/placement/delete/${id}`,
        },
        createPayload: (d) => ({
            bunker: d.bunker,
            seed: d.seed,
            amount: d.amount,
        }),
        updatePayload: (d) => ({
            bunker: d.id,
            seed: d.seed,
            seed_ru: d.seed_ru,
            amount: d.amount,
        }),
    },

    receipts: {
        idKey: "receipt",
        urls: {
            list: "/api/receipts/get",
            one: (id) => `/api/receipts/get/${id}`,
            create: "/api/receipts/add",
            update: "/api/receipts/update",
            delete: (id) => `/api/receipts/delete/${id}`,
        },
        createPayload: (d) => ({
            seed: d.seed,
            gcode: d.gcode,
            description: d.description,
        }),
        updatePayload: (d) => ({
            receipt: d.id,
            seed: d.seed,
            gcode: d.gcode,
            description: d.description,
        }),
    },

    shifts: {
        idKey: "shift",
        urls: {
            list: "/api/shifts/get",
            one: (id) => `/api/shifts/get/${id}`,
            create: "/api/shifts/add",
            update: "/api/shifts/update",
            delete: (id) => `/api/shifts/delete/${id}`,
        },
        createPayload: (d) => ({
            dt: new Date(d.dt).toISOString(),
        }),
        updatePayload: (d) => ({
            shift: d.id,
            dt: new Date(d.dt).toISOString(),
        }),
    },

    assignments: {
        idKey: "id",
        urls: {
            list: "/api/assignments/get",
            one: (id) => `/api/assignments/get/${id}`,
            create: "/api/assignments/add",
            update: "/api/assignments/update",
            delete: (id) => `/api/assignments/delete/${id}`,
        },
        createPayload: (d) => ({
            shift: d.shift,
            number: d.number,
            receipt: d.receipt,
            amount: d.amount,
        }),
        updatePayload: (d) => ({
            id: d.id,
            shift: d.shift,
            number: d.number,
            receipt: d.receipt,
            amount: d.amount,
        }),
    },

    reports: {
        idKey: "id",
        urls: {
            list: "/api/reports/get",
            one: (id) => `/api/reports/get/${id}`,
        },
    },

    choice: {
        idKey: "shift",
        urls: {
            list: "/api/shifts/getWithoutUser",
            update: "/api/shifts/update",
        },
        updatePayload: (d) => ({
            shift: d.id,
            username: d.username,
            dt: d.dt,
        }),
    },

    tasks: {
        idKey: "username",
        urls: {
            list: (username) => `/api/assignments/active-tasks/${username}`,
        },
        updatePayload: (d) => ({
            shift: d.id,
            username: d.username,
            dt: d.dt,
        }),
    },

    task: {
        idKey: "id",
        urls: {
            one: (id) => `/api/assignments/task/${id}`,
        },
    },

    seedWithBunkers: {
        idKey: "id",
        urls: {
            list: (id) => `/api/seeds/getWithBunkers/${id}`,
        },
    },

    logs: {
        idKey: "id",
        urls: {
            list: "/api/logs/get",
        },
    },

    logs: {
        idKey: "id",
        urls: {
            list: "/api/logs/get",
        },
    },
    fill: {
        idKey: "seed",
        urls: {
            update: "/api/placement/fill",
        },
        updatePayload: (d) => ({
            seed: d.seed,
            percent: d.percent,
        }),
    },
    calibration: {
        idKey: "session_id",
        urls: {
            create: "/api/calibration/handshake",
            clear: "/api/calibration/clear",
            save: "/api/calibration/save",
        },
        createPayload: (d) => ({
            session_id: d.sessionId,
        }),
    },
    takePhoto: {
        idKey: "session_id",
        urls: {
            create: "/api/calibration/photo",
        },
        createPayload: (d) => ({
            session_id: d.sessionId,
            number_of_photo: d.numberOfPhoto,
        })
    },
    clearPhotos: {
        idKey: "session_id",
        urls: {
            create: "/api/calibration/clear",
        },
        createPayload: (d) => ({
            session_id: d.sessionId,
        })
    },
    saveCalibration: {
        idKey: "session_id",
        urls: {
            one: (id) => `/api/calibration/get-photos/${id}`,
            create: "/api/calibration/save",
        },
        createPayload: (d) => ({
            session_id: d.sessionId,
            d_per_step: Number(d.dPerStep),
        })
    },
    foundCalibration: {
        idKey: "session_id",
        urls: {
            create: "/api/calibration/found",
        },
        createPayload: (d) => ({
            session_id: d.sessionId,
            cir: Number(d.cir),
        })
    },
    "device-settings": {
        idKey: "key",
        urls: {
            list: "/api/device-settings/get",
            one: (id) => `/api/device-settings/get/${id}`,
            create: "/api/device-settings/add",
            update: "/api/device-settings/update",
            delete: (id) => `/api/device-settings/delete/${id}`,
        },
        createPayload: (d) => ({
            key: d.key,
            value: d.value,
        }),
        updatePayload: (d) => ({
            key: d.id,
            value: d.value,
        }),
    },
};
