// src/dataProvider.js
const mockData = [
    { bunker: 1, distance: 100 },
    { bunker: 2, distance: 200 },
    { bunker: 3, distance: 300 },
    { bunker: 4, distance: 400 },
    { bunker: 5, distance: 500 },
];

const dataProvider = {
    getList: (resource, params) => {
        const dataWithId = mockData.map(item => ({ ...item, id: item.bunker }));
        return Promise.resolve({
            data: dataWithId,
            total: dataWithId.length,
        });
    },
    getOne: (resource, params) => {
        const record = mockData.find(item => String(item.bunker) === String(params.id));
        if (!record) {
            return Promise.reject(new Error(`Record with bunker=${params.id} not found`));
        }
        return Promise.resolve({ data: { ...record, id: record.bunker } });
    },
    create: (resource, params) => Promise.resolve({ data: { ...params.data, bunker: Date.now() } }),
    update: (resource, params) => Promise.resolve({ data: params.data }),
    delete: (resource, params) => Promise.resolve({ data: params.previousData }),
};

export default dataProvider;
