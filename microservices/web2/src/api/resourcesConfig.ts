type UrlParams = Record<string, string | number>;

type Url = string | ((params: UrlParams) => string);

interface ResourceConfig {
  idKey: string;
  urls: Record<string, Url>;
  createPayload?: (data: any) => any;
  updatePayload?: (data: any) => any;
}

export const resourcesConfig: Record<string, ResourceConfig> = {
    bunkers: {
        idKey: "bunker",
        urls: {
            list: "/api/admin/bunkers/get",
            one: (id) => `/api/admin/bunkers/get/${id}`,
            create: "/api/admin/bunkers/add",
            update: "/api/admin/bunkers/update",
            delete: (id) => `/api/admin/bunkers/delete/${id}`,
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
            list: "/api/admin/seeds/get",
            one: (id) => `/api/admin/seeds/get/${id}`,
            create: "/api/admin/seeds/add",
            update: "/api/admin/seeds/update",
            delete: (id) => `/api/admin/seeds/delete/${id}`,
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
            list: "/api/admin/users/get",
            one: (id) => `/api/users/get/${id}`,
            create: "/api/admin/register",
            update: "/api/users/update",
            delete: (id) => `/api/admin/users/delete/${id}`,
        },
        createPayload: (d) => ({
            username: d.username,
            password: d.password,
            full_name: d.full_name,
            is_admin: Boolean(d.is_admin),
        }),
        updatePayload: (d) => ({
            username: d.username,
            password: d.password,
            full_name: d.full_name,
            is_admin: Boolean(d.is_admin),
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
            list: "/api/admin/placement/get",
            one: (id) => `/api/admin/placement/get/${id}`,
            create: "/api/admin/placement/add",
            update: "/api/admin/placement/update",
            delete: (id) => `/api/admin/placement/delete/${id}`,
        },
        createPayload: (d) => ({
            bunker: Number(d.bunker),
            seed: d.seed,
            amount: Number(d.amount),
        }),
        updatePayload: (d) => ({
            bunker: Number(d.id),
            seed: d.seed,
            seed_ru: d.seed_ru,
            amount: Number(d.amount),
        }),
    },

    receipts: {
        idKey: "receipt",
        urls: {
            list: "/api/admin/receipts/get",
            one: (id) => `/api/admin/receipts/get/${id}`,
            create: "/api/admin/receipts/add",
            update: "/api/admin/receipts/update",
            delete: (id) => `/api/admin/receipts/delete/${id}`,
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
            list: "/api/admin/shifts/get",
            one: (id) => `/api/admin/shifts/get/${id}`,
            create: "/api/admin/shifts/add",
            update: "/api/admin/shifts/update",
            delete: (id) => `/api/admin/shifts/delete/${id}`,
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
            list: "/api/admin/assignments/get",
            one: (id) => `/api/admin/assignments/get/${id}`,
            create: "/api/admin/assignments/add",
            update: "/api/admin/assignments/update",
            delete: (id) => `/api/admin/assignments/delete/${id}`,
        },
        createPayload: (d) => ({
            shift: Number(d.shift),
            number: Number(d.number),
            receipt: Number(d.receipt),
            amount: Number(d.amount),
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
            list: "/api/admin/reports/get",
            one: (id) => `/api/admin/reports/get/${id}`,
        },
    },

    choice: {
        idKey: "id",
        urls: {
            list: "/api/shifts/getWithoutUser",
            update: "/api/shifts/update",
        },
        updatePayload: (d) => ({
            user_id: d.user_id,
            shift: d.shift,
            dt: d.dt,
        }),
    },

    tasks: {
        idKey: "id",
        urls: {
            list: (id) => `/api/assignments/active-tasks/${id}`,
        },
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
            list: "/api/admin/logs/get",
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
    },
    takePhoto: {
        idKey: "numberPhoto",
        urls: {
            create: (p) => `/api/calibration/photo/${p.numberPhoto}`,
        },
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
            d_per_step: Number(d.dPerStep) || 0,
        })
    },
    calculate: {
        idKey: "session_id",
        urls: {
            create: "/api/calibration/calculate",
        },
        createPayload: (d) => ({
            steps: Number(d.steps),
        })
    },
    deviceSettings: {
        idKey: "key",
        urls: {
            list: "/api/admin/device-settings/get",
            one: (id) => `/api/admin/device-settings/get/${id}`,
            create: "/api/admin/device-settings/add",
            update: "/api/admin/device-settings/update",
            delete: (id) => `/api/admin/device-settings/delete/${id}`,
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
    changePass: {
        idKey: "key",
        urls: {
            update: "/api/users/change-password",
        },
        updatePayload: (d) => ({
            id: d.id,
            old_password: d.old_password,
            new_password: d.new_password,
        }),
    },
    auth: {
        idKey: "id",
        urls: {
            create: "/auth/login"
        },
        createPayload: (d) => ({
            username: d.username,
            password: d.password,
        })
    }
};
