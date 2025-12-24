export const actionText = {
    getOne: {
        base: (resource) => `Ошибка при получении ресурса ${resource}`,
        400: (resource) => `Некорректные данные для получения ресурса ${resource}`,
        404: (resource) => `Ресурс ${resource} не найден`,
    },
    getList: {
        base: (resource) => `Ошибка при получении списка ${resource}`,
        400: (resource) => `Некорректные параметры запроса списка ${resource}`,
        404: (resource) => `Ресурс ${resource} не найден`,
    },
    create: {
        base: (resource) => `Ошибка при создании ресурса ${resource}`,
        200: (resource) => `Ресурс ${resource} создан`,
        400: (resource) => `Неверные данные для создания ресурса ${resource}`,
        409: (resource) => `Конфликт данных при создании ресурса ${resource}`,
    },
    update: {
        base: (resource) => `Ошибка при обновлении ресурса ${resource}`,
        200: (resource) => `Ресурс ${resource} обновлен`,
        400: (resource) => `Неверные данные для обновления ресурса ${resource}`,
        404: (resource) => `Ресурс ${resource} для обновления не найден`,
        409: (resource) => `Конфликт данных при обновлении ресурса ${resource}`,
    },
    remove: {
        base: (resource) => `Ошибка при удалении ресурса ${resource}`,
        200: (resource) => `Ресурс ${resource} удален`,
        400: (resource) => `Неверные данные для удаления ресурса ${resource}`,
        404: (resource) => `Ресурс ${resource} для удаления не найден`,
    },
};

export const handleApiError = async (resource, error, action) => {
    const status = error.status;

    const act = actionText[action];
    if (!act) throw new Error(`Unknown action: ${action}`);

    let message = act.base(resource);

    if (act[status]) {
        message = act[status](resource);
    } else {
        switch (status) {
            case 401:
                message = `Необходима авторизация`;
                break;
            case 403:
                message = `Недостаточно прав для операции над ресурсом ${resource}`;
                break;
            case 500:
                message = `Внутренняя ошибка сервера`;
                break;
        }
    }

    return message;
};
